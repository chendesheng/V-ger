package download

import (
	"bytes"
	"fmt"
	// "strconv"
	// "io"
	"errors"
	"log"
	"net/http"
	"time"
)

type block struct {
	from, to int64
}
type dataBlock struct {
	from, to int64
	data     []byte
}

func doDownload(url string, path string, from, to int64, maxSpeed int64) chan int64 {
	blockCnt := 5

	control := make(chan block, blockCnt)
	output := make(chan *dataBlock, blockCnt)

	go func(control chan block, from, size int64) {
		blockSize := int64(400 * 1024)
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
func writeOutput(path string, from int64, output chan *dataBlock, progress chan int64) {
	f := openOrCreateFileRW(path, from)
	defer f.Close()

	for db := range output {
		_, err := f.WriteAt(db.data, db.from)
		if err == nil {
			progress <- db.to - db.from
		} else {
			fmt.Printf("\n%s", err)
			log.Fatal(err)
		}
	}

	defer close(progress)
}

func pipeDownload(url string, control chan block, output chan *dataBlock) {
	numOfConn := make(chan bool, 8)
	prevComplete := make(chan bool, 1)
	prevComplete <- true

	for b := range control {
		numOfConn <- true
		complete := make(chan bool)
		go func(b block, output chan *dataBlock, numOfConn, prevComplete, complete chan bool) {
			//just block if network is down
			for {
				data, err := downloadBlock(url, b)

				if err == nil {
					<-numOfConn
					<-prevComplete
					close(prevComplete)

					// log.Printf("write output %v\n", b)
					output <- &dataBlock{from: b.from, to: b.to, data: data}
					complete <- true
					return
				} else {
					log.Println(err)
				}
			}
		}(b, output, numOfConn, prevComplete, complete)
		prevComplete = complete
	}
	<-prevComplete

	close(output)
	close(numOfConn)
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
	n, err := buffer.ReadFrom(resp.Body)
	if err != nil {
		return nil, err
	}

	if n != (to - from) {
		data = nil
		err = errors.New(fmt.Sprintf("not download whole block. download length: %d, need length: %d\n", n, to-from))
	} else {
		data = buffer.Bytes()
		err = nil
	}

	return
}

// func sampleDownload(url string, path string) {
// 	output := make(chan []byte)
// 	go func(output chan []byte) {
// 		defer close(output)

// 		req := createDownloadRequest(url, 0, -1)
// 		resp, err := DownloadClient.Do(req)
// 		if err != nil {
// 			log.Fatal(err)
// 		}

// 		for {
// 			buffer := make([]byte, 40000)
// 			readLen, _ := resp.Body.Read(buffer)
// 			if readLen == 0 {
// 				break
// 			}
// 			output <- buffer[:readLen]
// 		}
// 	}(output)

// 	progress := make(chan int64)
// 	writeOutput(path, 0, output, progress)
// }

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
