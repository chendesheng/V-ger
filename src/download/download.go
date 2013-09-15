package download

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"native"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"time"
	"util"
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

func generateBlock(input chan *block, output chan<- *block, chBlockSize chan int64, from, size int64, quit <-chan bool) {
	log.Printf("generate block output: %v", output)
	blockSize := int64(400 * 1024)

	to := from + blockSize
	if to > size {
		to = size
	}

	//small blocksize after start,
	//change to a larger blocksize after 15 seconds
	changeBlockSize := time.NewTimer(time.Second * 15)
	startCnt := 5
	log.Printf("output %v", output)
	maxSpeed := int64(0)
	for {
		if startCnt < 0 {
			select {
			case _, ok := <-input:
				if !ok {
					return
				}
			case <-quit:
				return
			}
		} else {
			startCnt--
		}

		// b := time.Now()
		select {
		case maxSpeed = <-chBlockSize:
			// maxSpeed = 0
			if maxSpeed > 0 {
				blockSize = maxSpeed * 1024
			} else {
				blockSize = int64(100 * 1024)
				changeBlockSize.Reset(time.Second * 15)
			}
		case output <- &block{from, to, nil}:
			if to == size {
				fmt.Println("return generate block ", size)
				close(output)
				for {
					select {
					case _, ok := <-input:
						if !ok {
							return
						}
					case <-quit:
						return
					}
				}
			} else {
				from = to
				to = from + blockSize
				if to > size {
					to = size
				}
			}
		case <-changeBlockSize.C:
			if maxSpeed == 0 {
				blockSize = 400 * 1024
			}
			changeBlockSize.Stop()
		case <-quit:
			close(output)
			fmt.Println("quit generate block")
			return
		}
	}
}
func writeOutput(w io.Writer, input <-chan *block, output chan *block, quit chan bool) {
	defer func() {
		fmt.Println("close write output")
		close(output)
	}()

	pathErrNotifyTimes := 0
	for {
		select {
		case b, ok := <-input:
			if !ok {
				return
			}
			for {

				_, err := w.Write(b.data)
				b.data = nil

				if err == nil {
					select {
					case output <- b:
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
	for {
		select {
		case b, ok := <-input:
			if !ok {
				fmt.Println("downloadRoutine finish")
				close(output)
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

func sortOutput(input <-chan *block, output chan<- *block, quit <-chan bool, from int64) {
	dbmap := make(map[int64]*block)
	nextOutputFrom := from
	for {
		select {
		case db, ok := <-input:
			if db != nil {
				dbmap[db.from] = db
			}

			// log.Println(len(dbmap))
			for {
				if d, exist := dbmap[nextOutputFrom]; exist {
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

			if !ok {
				close(output)
				return
			}
		case <-quit:
			fmt.Println("sort output quit")
			return
		}
	}
}
func concurrentDownload(url string, input <-chan *block, output chan<- *block, quit chan bool) {
	log.Printf("concurrentDownload download input %v", input)
	output1 := make(chan *block)
	output2 := make(chan *block)
	output3 := make(chan *block)
	output4 := make(chan *block)
	output5 := make(chan *block)

	go downloadRoutine(url, input, output1, quit)
	go downloadRoutine(url, input, output2, quit)
	go downloadRoutine(url, input, output3, quit)
	go downloadRoutine(url, input, output4, quit)
	go downloadRoutine(url, input, output5, quit)

	for {
		select {
		case b, ok := <-output1:
			if !ok {
				output1 = nil
			} else {
				select {
				case output <- b:
				case <-quit:
				}
			}
		case b, ok := <-output2:
			if !ok {
				output2 = nil
			} else {
				select {
				case output <- b:
				case <-quit:
				}
			}
		case b, ok := <-output3:
			if !ok {
				output3 = nil
			} else {
				select {
				case output <- b:
				case <-quit:
				}
			}
		case b, ok := <-output4:
			if !ok {
				output4 = nil
			} else {
				select {
				case output <- b:
				case <-quit:
				}
			}
		case b, ok := <-output5:
			if !ok {
				output5 = nil
			} else {
				select {
				case output <- b:
				case <-quit:
				}
			}
		case <-time.After(taskRestartTimeout):
			fmt.Println("close quit")
			close(quit)
			return
		case <-quit:
			fmt.Println("currentDownload quit")
			return
		}
		if output1 == nil && output2 == nil && output3 == nil &&
			output4 == nil && output5 == nil {
			close(output)
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
