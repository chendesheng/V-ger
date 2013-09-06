package download

import (
	"bytes"
	"fmt"
	"native"
	// "net"
	"path/filepath"
	"regexp"
	"strings"
	"util"
	// "sort"
	"runtime"
	// "strconv"
	// "errors"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

var taskRestartTimeout time.Duration
var networkTimeout time.Duration

func init() {
	networkTimeout = time.Duration(util.ReadIntConfig("network-timeout")) * time.Second
	taskRestartTimeout = time.Duration(util.ReadIntConfig("task-restart-timeout")) * time.Second

	// http.DefaultTransport.(*http.Transport).Dial = func(network, addr string) (net.Conn, error) {
	// 	c, err := net.Dial(network, addr)
	// 	log.Printf("%v", c)
	// 	return c, err
	// }
}

type block struct {
	from, to int64
	data     []byte
}

func doDownload(url string, w io.Writer, from, to int64,
	maxSpeed int64, control chan int64, quit chan bool) chan int64 {

	for {
		finalUrl, _, _, err := GetDownloadInfo(url)
		if err == nil {
			url = finalUrl
			break
		}

		select {
		case <-quit:
			return nil
		default:
			time.Sleep(time.Second * 2)
		}
	}

	input := make(chan *block)
	output := make(chan *block)

	go generateBlock(input, from, to, maxSpeed, control, quit)

	go concurrentDownload(url, input, output, quit, from, to)

	progress := make(chan int64)

	go writeOutput(w, from, output, progress, quit)

	return progress
}
func generateBlock(input chan<- *block, from, size int64, maxSpeed int64, control chan int64, quit <-chan bool) {
	blockSize := int64(100 * 1024)
	if maxSpeed > 0 {
		blockSize = maxSpeed * 1024
	}

	to := from + blockSize
	if to > size {
		to = size
	}

	//small blocksize after start,
	//change to a larger blocksize after 15 seconds
	changeBlockSize := time.NewTimer(time.Second * 15)
	for {
		b := time.Now()
		select {
		case maxSpeed := <-control:
			fmt.Println("set max speed")
			if maxSpeed > 0 {
				blockSize = maxSpeed * 1024
			} else {
				blockSize = int64(100 * 1024)
				changeBlockSize.Reset(time.Second * 15)
			}
		case input <- &block{from, to, nil}:
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
						time.Sleep(time.Second - d)
					}
				}
			}
		case <-changeBlockSize.C:
			if maxSpeed == 0 {
				blockSize = 200 * 1024
			}
			changeBlockSize.Stop()
		case <-quit:
			close(input)
			fmt.Println("input quit")
			return
		}
	}
}
func writeOutput(w io.Writer, from int64, output <-chan *block, progress chan int64, quit chan bool) {
	defer func() {
		fmt.Println("close progress")
		close(progress)
	}()

	pathErrNotifyTimes := 0
	for {
		select {
		case db, ok := <-output:
			if !ok {
				return
			}
			for {

				_, err := w.Write(db.data)
				db.data = nil

				if err == nil {
					select {
					case progress <- db.to - db.from:
						break
					case <-quit:
						return
					}
					break
				} else if perr, ok := err.(*os.PathError); ok {
					log.Print(err)

					if pathErrNotifyTimes == 0 { //only report once
						native.SendNotification("Error write "+filepath.Base(perr.Path), perr.Err.Error())
					}
					pathErrNotifyTimes++
					if pathErrNotifyTimes > 100 {
						log.Fatal(err)
						return
					}

					select {
					case <-quit:
						return
					case <-time.After(time.Second * 2):
						break
					}
				} else {
					log.Print(err)
					ensureQuit(quit)
					return
				}
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
			downloadBlock(url, b, output, quit)
		case <-quit:
			fmt.Println("downloadRoutine quit")
			return
		}
	}
}
func downloadBlock(url string, b *block, output chan<- *block, quit chan bool) {
	for {
		req := createDownloadRequest(url, b.from, b.to-1)

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Println(err)
		} else {
			size := b.to - b.from

			data, err := readWithTimeout(req, resp, size, quit)
			if err != nil {
				log.Print(err)
			}
			if err == nil && int64(len(data)) == size {
				b.data = data
				select {
				case output <- b:
					return
				case <-quit:
					return
				}
			}
		}

		select {
		case <-quit:
			return
		default:
			runtime.Gosched()
		}
	}
}

func readWithTimeout(req *http.Request, resp *http.Response, size int64, quit chan bool) ([]byte, error) {
	buffer := bytes.NewBuffer(make([]byte, 0, size+bytes.MinRead))
	finish := make(chan error)
	go func() {
		select {
		case <-time.After(networkTimeout): //cancelRequest if time.After before close(finish)
			cancelRequest(req)
		case <-quit:
			cancelRequest(req)
			return
		case <-finish: //close(finish) before time.After
			return
		}
	}()

	_, err := buffer.ReadFrom(resp.Body)
	close(finish)

	if err != nil {
		return nil, err
	} else {
		return buffer.Bytes(), nil
	}
}

func createDownloadRoutine(url string, output chan<- *block, quit chan bool) chan<- *block {
	input := make(chan *block)
	go downloadRoutine(url, input, output, quit)
	return input
}
func sortOutput(input <-chan *block, output chan<- *block, quit <-chan bool, from int64, to int64) {
	dbmap := make(map[int64]*block)
	nextOutputFrom := from
	for {
		select {
		case db, _ := <-input:
			if db == nil {
				break
			}

			dbmap[db.from] = db

			// log.Println(len(dbmap))
			for {
				if d, ok := dbmap[nextOutputFrom]; ok {
					// fmt.Printf("sort output %d-%d\n", d.from, d.to)
					select {
					case output <- d:
						nextOutputFrom = d.to
						delete(dbmap, d.from)
						break
					case <-quit:
						return
					}
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
	// chan6 := createDownloadRoutine(url, disorderOutput, quit)

	go sortOutput(disorderOutput, output, quit, from, to)

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
				// close(chan6)
				return
			}
			select {
			case chan1 <- b:
			case chan2 <- b:
			case chan3 <- b:
			case chan4 <- b:
			case chan5 <- b:
			// case chan6 <- b:
			case <-time.After(taskRestartTimeout):
				fmt.Println("close quit")
				close(quit)
				return
			case <-quit:
				fmt.Println("currentDownload quit")
				return
			}
		case <-quit:
			fmt.Println("currentDownload quit2")
			return
		}
	}

}

func GetDownloadInfo(url string) (finalUrl string, name string, size int64, err error) {
	req := createDownloadRequest(url, -1, -1)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", "", 0, err
	}
	defer resp.Body.Close()

	finalUrl = resp.Request.URL.String()

	name, size = getFileInfo(resp.Header)
	if name == "" {
		name = getFileName(finalUrl)
	}

	reg := regexp.MustCompile("[/]")
	name = reg.ReplaceAllString(name, "")
	name = strings.TrimLeft(name, ".")

	if name == "" && size == 0 {
		err = fmt.Errorf("Broken resource")
	}

	return
}
