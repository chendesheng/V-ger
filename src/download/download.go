package download

import (
	"bytes"
	"fmt"
	// "runtime"
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

func doDownload(t *Task, url string, path string, from, to int64,
	maxSpeed int64, control chan int) chan int64 {

	input := make(chan block)
	output := make(chan *dataBlock)

	go generateBlock(input, from, to, maxSpeed, control)

	go pipeDownload(url, input, output)

	progress := make(chan int64)
	go writeOutput(path, from, output, progress)

	return progress
}
func generateBlock(input chan<- block, from, size int64, maxSpeed int64, control chan int) {
	blockSize := int64(400 * 1024)
	if maxSpeed > 0 {
		blockSize = maxSpeed * 1024
	}

	to := from + blockSize
	if to > size {
		to = size
	}
	// lastTime := time.Now()
	// lastTime2 := time.Now()
	r := time.Duration(0)
	for {
		// if maxSpeed > 0 {
		// 	now := time.Now()
		// 	dur := now.Sub(lastTime)
		// 	lastTime = now

		// 	fmt.Printf("input %d-%d\n", from, to)

		// 	if dur < time.Second {
		// 		fmt.Println("sleep dur ", dur)
		// 		time.Sleep(dur)
		// 	}
		// }
		b := time.Now()
		select {
		case cmd := <-control:
			if cmd == -1 {
				close(input)
				fmt.Println("input stopped")
				return
			} else {
				fmt.Println("set max speed")
				maxSpeed = int64(cmd)
				if maxSpeed > 0 {
					blockSize = maxSpeed * 1024
				} else {
					blockSize = int64(400 * 1024)
				}
			}
		case input <- block{from, to}:
			if to == size {
				// fmt.Println("input", from, "-", to)
				fmt.Println("return input")
				close(input)
				return
			} else {
				// fmt.Println("input", from, "-", to)
				from = to
				to = from + blockSize
				if to > size {
					to = size
				}
				if maxSpeed > 0 {
					d := time.Now().Sub(b)
					if d < time.Second {
						// fmt.Println("durs: d: ", d, " r: ", r)
						time.Sleep(time.Second - d - r)
						r -= time.Second
						if r < 0 {
							r = 0
						}
					} else {
						r = d - time.Second
					}
				}
				// if maxSpeed > 0 {
				// 	now := time.Now()
				// 	dur := now.Sub(lastTime2)
				// 	lastTime2 = now

				// 	fmt.Printf("input %d-%d\n", from, to)

				// 	if dur < time.Second {
				// 		fmt.Println("sleep dur 2", dur)
				// 		time.Sleep(dur)
				// 	}
				// }
			}
		}
	}
	fmt.Println("doDownload return")
}
func writeOutput(path string, from int64, output <-chan *dataBlock, progress chan int64) {
	f := openOrCreateFileRW(path, from)
	defer f.Close()

	defer func() {
		fmt.Println("close progress")
		close(progress)
	}()

	for db := range output {
		// fmt.Printf("writeOutput %d-%d\n", db.from, db.to)
		_, err := f.WriteAt(db.data, db.from)

		if err == nil {
			// fmt.Println("progress<-")
			progress <- db.to - db.from
			// fmt.Println("progress 111")
		} else {
			fmt.Printf("\n%s", err)
			log.Fatal(err)
		}
	}

	fmt.Println("writeOutput end")
}
func pipeDownload(url string, input <-chan block, output chan<- *dataBlock) {
	defer close(output)
	numOfConn := make(chan bool, 5)
	defer close(numOfConn)

	prevComplete := make(chan bool, 1)
	prevComplete <- true

	// for b := range input {
	// fmt.Printf("1   %v\n", b)
	for b := range input {
		// fmt.Println("14")
		numOfConn <- true
		// fmt.Println("7")
		complete := make(chan bool)
		go func(b block, output chan<- *dataBlock, numOfConn, prevComplete, complete chan bool) {
			// defer func() {
			// 	if r := recover(); r != nil {
			// 		fmt.Println("Recovered in f ", r)
			// 	}
			// }()
			//just block if network is down
			for {
				// fmt.Println("10")
				data, err := downloadBlock(url, b)
				// fmt.Println("11")
				// data, err := make([]byte, 1), error(nil)

				if err == nil {
					// fmt.Printf("%v 1\n", b)
					<-numOfConn
					// fmt.Printf("%v 2\n", b)
					<-prevComplete
					// fmt.Printf("%v 3\n", b)
					close(prevComplete)
					// fmt.Printf("%v 4\n", b)
					output <- &dataBlock{from: b.from, to: b.to, data: data}

					// fmt.Printf("%v 5\n", b)
					complete <- true
					// 
					// fmt.Printf("%v 6\n", b)
					return
				} else {
					// fmt.Println("12")
					log.Println(err)
				}
			}
		}(b, output, numOfConn, prevComplete, complete)
		// fmt.Println("9")
		prevComplete = complete
		// fmt.Println("13")

	}
	<-prevComplete

	fmt.Println("pipeDownload return")
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
		return nil, errors.New(fmt.Sprintf("not download whole block. download length: %d, need length: %d\n", n, to-from))
	}

	return buffer.Bytes(), nil
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
