package download

import (
	"bytes"
	"fmt"
	// "io"
	"log"
	"net/http"
	"time"
)

func doDownload(url string, path string, from, to int64, maxSpeed int64) chan int64 {
	blockCnt := 5
	// blockSize := 300 * 1024

	control := make(chan block, blockCnt)
	output := make(chan []byte, blockCnt)

	go func(control chan block, from, size int64) {
		blockSize := int64(300 * 1024)
		if maxSpeed > 0 {
			blockSize = maxSpeed / 10 * 1024
		}

		for {
			to := from + blockSize
			if to <= size {
				control <- block{from, to}
				from += blockSize
			} else {
				control <- block{from, size}
				close(control)
				break
			}
			if maxSpeed > 0 {
				time.Sleep(time.Millisecond * 100)
			}
		}
	}(control, from, to)

	go pipeDownload(url, control, output)

	progress := make(chan int64)
	go writeOutput(path, from, output, progress)

	return progress
}
func writeOutput(path string, from int64, output chan []byte, progress chan int64) {
	f := openOrCreateFileRW(path, from)
	defer f.Close()

	for b := range output {
		_, err := f.Write(b)
		if err == nil {
			progress <- int64(len(b))
		} else {
			fmt.Printf("\n%s", err)
			log.Fatal(err)
		}
	}

	defer close(progress)
}

func pipeDownload(url string, control chan block, output chan []byte) {
	numOfConn := make(chan bool, 5)
	var prevComplete chan bool

	for b := range control {
		complete := make(chan bool)
		go func(b block, output chan []byte, numOfConn, prevComplete, complete chan bool) {
			//just block if network is down
			for {
				numOfConn <- true
				block, err := downloadBlock(url, b)
				<-numOfConn
				if err == nil {
					if prevComplete != nil {
						<-prevComplete
					}
					output <- block
					complete <- true
					return
				}
			}
		}(b, output, numOfConn, prevComplete, complete)
		prevComplete = complete
	}
	if prevComplete != nil {
		<-prevComplete
	}
	close(output)
}

func downloadBlock(url string, b block) (data []byte, err error) {
	from, to := b.from, b.to
	req := createDownloadRequest(url, from, to-1)

	resp, err := DownloadClient.Do(req)
	defer func() {
		if resp != nil && resp.Body != nil {
			resp.Body.Close()
		}
	}()
	if err != nil {
		return nil, err
	}

	buffer := bytes.NewBuffer(make([]byte, 0, to-from))
	_, err = buffer.ReadFrom(resp.Body)
	if err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

func sampleDownload(url string, path string) {
	output := make(chan []byte)
	go func(output chan []byte) {
		defer close(output)

		req := createDownloadRequest(url, 0, -1)
		resp, err := DownloadClient.Do(req)
		if err != nil {
			log.Fatal(err)
		}

		for {
			buffer := make([]byte, 40000)
			readLen, _ := resp.Body.Read(buffer)
			if readLen == 0 {
				break
			}
			output <- buffer[:readLen]
		}
	}(output)

	progress := make(chan int64)
	writeOutput(path, 0, output, progress)
}

func getDownloadInfo(url string) (realURL string, name string, size int64) {
	req := createDownloadRequest(url, -1, -1)
	DownloadClient.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		temp := req.URL.String()
		if temp != "" {
			url = temp
		}
		return nil
	}

	resp, err := DownloadClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	name, size = getFileInfo(resp.Header)
	if name == "" {
		name = getFileName(url)
	}
	realURL = url
	return
}
