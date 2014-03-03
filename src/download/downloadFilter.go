package download

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"runtime"
	"strings"
	"time"
)

type downloadFilter struct {
	basicFilter
	url           string
	routineNumber int
}

func (df *downloadFilter) active() {
	url := df.url

	chFinishes := make([]chan bool, df.routineNumber)
	for i := 0; i < df.routineNumber; i++ {
		chFinishes[i] = make(chan bool)
		go func(ch chan bool) {
			downloadRoutine(url, df.input, df.output, df.quit)
			ch <- true
		}(chFinishes[i])
	}

	for _, ch := range chFinishes {
		<-ch
	}

	close(df.output)
}

func downloadRoutine(url string, input <-chan *block, output chan<- *block, quit chan bool) {
	// log.Print("download routine begin: ", url[:strings.Index(url, "?")])

	for {
		finalUrl, _, _, err := GetDownloadInfo(url)
		if err == nil {
			url = finalUrl
			break
		} else {
			log.Print(err)
		}

		select {
		case <-quit:
			return
		default:
			time.Sleep(time.Second * 2)
		}
	}
	// log.Print("final download url:", url)

	if strings.Contains(url, "192.168.1.1") {
		//AUSU router may redirect to error_page.html, download from this url will crap target file.
		return
	}

	for {
		select {
		case b, ok := <-input:
			if !ok {
				fmt.Println("downloadRoutine finish")
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

		// data, _ := httputil.DumpRequest(req, true)
		// println(string(data))

		// data, _ = httputil.DumpResponse(resp, false)
		// println(string(data))

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
			} else {
				bytes, _ := httputil.DumpResponse(resp, false)
				log.Print(string(bytes))
				log.Printf("download wrong data:%d,%d,%d", b.from, b.to, len(b.data))
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
		case <-time.After(NetworkTimeout): //cancelRequest if time.After before close(finish)
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
