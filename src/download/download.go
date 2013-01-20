package download

import (
	"bytes"
	// "fmt"
	"log"
	"net/http"
)

func doDownload(url string, path string, from, to int64) chan int64 {
	blockCnt := 5
	blockSize := 300 * 1024

	output := make(chan []byte, blockCnt)
	cntControl := make(chan bool, blockCnt)

	readyOutput := make(chan bool)
	go pipeDownload(url, from, blockSize, to, output, readyOutput, cntControl)
	readyOutput <- true

	progress := make(chan int64)
	go writeOutput(path, from, output, progress)

	return progress
}
func writeOutput(path string, from int64, output chan []byte, progress chan int64) {
	f := openOrCreateFileRW(path, from)
	defer f.Close()

	for b := range output {
		f.Write(b)

		progress <- int64(len(b))
	}

	defer close(progress)
}

func pipeDownload(url string, from int64, blockSize int, size int64, output chan []byte, readyOutput chan bool, cntControl chan bool) {
	for {
		if from >= size {
			<-readyOutput
			close(output)
			return
		}

		cntControl <- true // block if chan is full, not start new connection until one of current connections is complete.

		complete := make(chan bool)

		to := from + int64(blockSize)
		if to > size {
			to = size
		}

		go func(url string, from, to int64, output chan []byte, readyOutput, complete, cntControl chan bool) {
			//just block if network is down
			for {
				block, err := downloadBlock(url, from, to-1)
				if err == nil {
					<-readyOutput
					output <- block
					complete <- true
					<-cntControl
					return
				}
			}

		}(url, from, to, output, readyOutput, complete, cntControl)

		from = to
		readyOutput = complete
	}
}

func downloadBlock(url string, from, to int64) (data []byte, err error) {
	req := createDownloadRequest(url, from, to)

	resp, err := DownloadClient.Do(req)
	defer func() {
		if resp != nil && resp.Body != nil {
			resp.Body.Close()
		}
	}()
	if err != nil {
		return nil, err
	}

	buffer := bytes.NewBuffer(make([]byte, 0, to-from+1))
	_, err = buffer.ReadFrom(resp.Body)
	if err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func sampleDownload(url string, path string, from, to int64) chan int64 {
	output := make(chan []byte)
	go func(output chan []byte) {
		defer close(output)

		req := createDownloadRequest(url, from, -1)
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
	go writeOutput(path, from, output, progress)

	return progress
}

func getDownloadInfo(url string) (realURL string, name string, size int64) {
	req := createDownloadRequest(url, -1, -1)
	DownloadClient.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		realURL = req.URL.String()
		return nil
	}

	resp, err := DownloadClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	name, size = getFileInfo(resp.Header)
	if name == "" {
		name = getFileName(realURL)
	}

	return
}
