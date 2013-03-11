package download

import (
	"bytes"
	"fmt"
	// "runtime"
	// "strconv"
	// "errors"
	"io"
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

func doDownload(t *Task, url string, path string, from, to int64,
	maxSpeed int64, control chan int, quit chan bool) chan int64 {

	input := make(chan block)
	output := make(chan *dataBlock)

	go generateBlock(input, from, to, maxSpeed, control, quit)

	go pipeDownload(url, input, output, quit)

	progress := make(chan int64)
	go writeOutput(path, from, output, progress, quit)

	return progress
}
func generateBlock(input chan<- block, from, size int64, maxSpeed int64, control chan int, quit chan bool) {
	blockSize := int64(400 * 1024)
	if maxSpeed > 0 {
		blockSize = maxSpeed * 1024
	}

	to := from + blockSize
	if to > size {
		to = size
	}
	r := time.Duration(0)
	for {
		b := time.Now()
		select {
		case cmd := <-control:
			fmt.Println("set max speed")
			maxSpeed = int64(cmd)
			if maxSpeed > 0 {
				blockSize = maxSpeed * 1024
			} else {
				blockSize = int64(400 * 1024)
			}
		case input <- block{from, to}:
			if to == size {
				fmt.Println("return input")
				close(input)
				return
			} else {
				from = to
				to = from + blockSize
				if to > size {
					to = size
				}
				if maxSpeed > 0 {
					d := time.Now().Sub(b)
					if d < time.Second {
						time.Sleep(time.Second - d - r)
						r -= time.Second
						if r < 0 {
							r = 0
						}
					} else {
						r = d - time.Second
					}
				}
			}
		case <-quit:
			close(input)
			fmt.Println("input quit")
			return
		}
	}
}
func writeOutput(path string, from int64, output <-chan *dataBlock, progress chan int64, quit chan bool) {
	f := openOrCreateFileRW(path, from)
	defer f.Close()

	defer func() {
		fmt.Println("close progress")
		close(progress)
	}()

	for {
		select {
		case db, ok := <-output:
			if !ok {
				break
			}
			_, err := f.WriteAt(db.data, db.from)

			if err == nil {
				select {
				case progress <- db.to - db.from:
				case <-quit:
					return
				}
			} else {
				fmt.Printf("\n%s", err)
				log.Fatal(err)
			}
		case <-quit:
			fmt.Println("write output quit")
			return

		}
	}

	fmt.Println("writeOutput end")
}
func pipeDownload(url string, input <-chan block, output chan<- *dataBlock, quit chan bool) {
	numOfConn := make(chan bool, 5)
	defer close(numOfConn)

	prevComplete := make(chan bool, 1)
	prevComplete <- true

	for {
		select {
		case b, ok := <-input:
			if !ok {
				break
			}
			select {
			case numOfConn <- true:
			case <-quit:
				return
			}
			complete := make(chan bool)
			go func(b block, output chan<- *dataBlock, numOfConn, prevComplete, complete chan bool) {
				for {
					chanRes, closer, err := downloadBlock(url, b, quit)

					select {
					case data := <-chanRes:

						if err == nil {

							select {
							case <-numOfConn:
							case <-quit:
								return
							}

							select {
							case <-prevComplete:
							case <-quit:
								return
							}

							close(prevComplete)

							select {
							case output <- &dataBlock{from: b.from, to: b.to, data: data}:
							case <-quit:
								return
							}

							select {
							case complete <- true:
							case <-quit:
								return
							}
							return
						} else {
							log.Println(err)
						}
					case <-quit:
						fmt.Println("download block quit")
						closer.Close()
						return
					}
				}
			}(b, output, numOfConn, prevComplete, complete)
			prevComplete = complete
		case <-quit:
		}
	}

	select {
	case <-prevComplete:
		close(output)
		fmt.Println("pipdownload return")
		return
	case <-quit:
		fmt.Println("pipdownload quit")
		return
	}
}

func downloadBlock(url string, b block, quit chan bool) (chan []byte, io.Closer, error) {
	from, to := b.from, b.to
	req := createDownloadRequest(url, from, to-1)

	result := make(chan []byte)

	resp, err := DownloadClient.Do(req)
	if err != nil {
		go func() { result <- make([]byte, 0) }()
		return result, nil, err
	}

	go func() {
		defer func() {
			if resp != nil && resp.Body != nil {
				resp.Body.Close()
			}
		}()

		defer func() {
			if r := recover(); r != nil {
				fmt.Println(r)
			}
		}()
		buffer := bytes.NewBuffer(make([]byte, 0, to-from))
		buffer.ReadFrom(resp.Body)

		select {
		case result <- buffer.Bytes():
		case <-quit:
			fmt.Println("downloadblock inner quit")
		}
	}()
	return result, resp.Body, nil
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
