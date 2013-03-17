package download

import (
	"bytes"
	"fmt"
	// "sort"
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
	data     []byte
}

func newDataBlock(from, to int64) *block {
	return &block{from, to, make([]byte, 0)}
}

func doDownload(url string, path string, from, to int64,
	maxSpeed int64, control chan int, quit chan bool) chan int64 {

	input := make(chan *block)
	output := make(chan *block)

	go generateBlock(input, from, to, maxSpeed, control, quit)

	go concurrentDownload(url, input, output, quit, from, to)

	progress := make(chan int64)
	go writeOutput(path, from, output, progress, quit)

	return progress
}
func generateBlock(input chan<- *block, from, size int64, maxSpeed int64, control chan int, quit chan bool) {
	blockSize := int64(60 * 1024)
	if maxSpeed > 0 {
		blockSize = maxSpeed * 1024
	}

	to := from + blockSize
	if to > size {
		to = size
	}
	r := time.Duration(0)

	//small blocksize after start,
	//change to a larger blocksize after 30 seconds
	changeBlockSize := time.Tick(time.Second * 30)
	for {
		b := time.Now()
		select {
		case cmd := <-control:
			fmt.Println("set max speed")
			maxSpeed = int64(cmd)
			if maxSpeed > 0 {
				blockSize = maxSpeed * 1024
			} else {
				blockSize = int64(300 * 1024)
			}
		case input <- newDataBlock(from, to):
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
		case <-changeBlockSize:
			blockSize = 300 * 1024
		case <-quit:
			close(input)
			fmt.Println("input quit")
			return
		}
	}
}
func writeOutput(path string, from int64, output <-chan *block, progress chan int64, quit chan bool) {
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
				return
			}
			_, err := f.WriteAt(db.data, db.from)
			db.data = nil

			if err == nil {
				select {
				case progress <- db.to - db.from:
					break
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

func downloadRoutine(url string, input <-chan *block, output chan<- *block, quit chan bool) {

	defer func() {
		if r := recover(); r != nil {
			fmt.Println(r)
		}
	}()
	for {
		select {
		case b, ok := <-input:
			if !ok {
				fmt.Println("downloadRoutine finish")
				// close(output)
				return
			}

		tryDownloadBlock:
			for {
				//fmt.Printf("downloadBlock %s %v\n", url, b)
				chanRes, closer, err := downloadBlock(url, b, quit)
				select {
				case data := <-chanRes:
					if err != nil {
						break
					}
					b.data = data
					// fmt.Printf("write downloadBlock %v\n", b)
					select {
					case output <- b:
						break tryDownloadBlock
					case <-quit:
						return
					}
				case <-quit:
					fmt.Println("downloadRoutine quit close")
					closer.Close()
					return
				}
			}
		case <-quit:
			fmt.Println("downloadRoutine quit")
			return
		}
	}
}
func createDownloadRoutine(url string, output chan<- *block, quit chan bool) chan<- *block {
	input := make(chan *block)
	go func(url string, input <-chan *block, output chan<- *block, quit chan bool) {
		downloadRoutine(url, input, output, quit)
	}(url, input, output, quit)
	return input
}
func sortOutput(input <-chan *block, output chan<- *block, quit chan bool, from int64, to int64) {
	dbmap := make(map[int64]*block)
	var nextOutputFrom = from
	for {
		select {
		case db, _ := <-input:
			if db == nil {
				break
			}

			dbmap[db.from] = db
			for {
				if d, ok := dbmap[nextOutputFrom]; ok {
					fmt.Printf("sort output %d-%d\n", d.from, d.to)
					select {
					case output <- d:
					case <-quit:
						return
					}
					nextOutputFrom = d.to
					delete(dbmap, db.from)
				} else {
					break
				}
			}
			if nextOutputFrom == to {
				close(output)
				fmt.Println("sortOutput finish")
				return
			}
		case <-quit:
			fmt.Println("sort output quit")
			return
		}
	}
}
func concurrentDownload(url string, input <-chan *block, output chan<- *block, quit chan bool, from, to int64) {
	disorderOutput := make(chan *block)
	chan1 := createDownloadRoutine(url, disorderOutput, quit)
	chan2 := createDownloadRoutine(url, disorderOutput, quit)
	chan3 := createDownloadRoutine(url, disorderOutput, quit)
	chan4 := createDownloadRoutine(url, disorderOutput, quit)
	chan5 := createDownloadRoutine(url, disorderOutput, quit)
	chan6 := createDownloadRoutine(url, disorderOutput, quit)

	go func(input <-chan *block, output chan<- *block, quit chan bool, from, to int64) {
		sortOutput(input, output, quit, from, to)
	}(disorderOutput, output, quit, from, to)

	for {
		select {
		case b, ok := <-input:
			if !ok {
				fmt.Println("concurrentDownload finish")
				close(chan1)
				close(chan2)
				close(chan3)
				close(chan4)
				close(chan5)
				close(chan6)
				return
			}
			select {
			case chan1 <- b:
			case chan2 <- b:
			case chan3 <- b:
			case chan4 <- b:
			case chan5 <- b:
			case chan6 <- b:
			case <-quit:
				fmt.Println("currentDownload quit")
				return
			}
			// fmt.Println("write to downloadRoutine")
			// chan1 <- b
		case <-quit:
			fmt.Println("currentDownload quit2")
			return
		}
	}

}
func downloadBlock(url string, b *block, quit chan bool) (chan []byte, io.Closer, error) {
	from, to := b.from, b.to
	req := createDownloadRequest(url, from, to-1)

	result := make(chan []byte)
	resp, err := DownloadClient.Do(req)
	if err != nil {
		fmt.Println(err)
		go func() { result <- make([]byte, 0) }()
		return result, nil, err
	}

	go func(quit chan bool, resp *http.Response) {
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
		// fmt.Println("read from buffer")
		buffer.ReadFrom(resp.Body)
		// fmt.Println("read from buffer end")

		select {
		case result <- buffer.Bytes():
			// fmt.Println("write result end")
			break
		case <-quit:
			fmt.Println("downloadblock inner quit")
			return
		}
	}(quit, resp)
	return result, resp.Body, nil
}

func GetDownloadInfo(url string) (realURL string, name string, size int64) {
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

func SingleRoutineDownload(url string, w io.Writer, from, to int64) error {
	req := createDownloadRequest(url, from, to-1)

	resp, err := DownloadClient.Do(req)
	if err != nil {
		return err
	}

	defer func() {
		if r := recover(); r != nil {
			fmt.Println(r)
		}
	}()

	io.Copy(w, resp.Body)

	return nil
}
