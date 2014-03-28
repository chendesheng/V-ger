package download

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
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
		case <-time.After(2 * time.Second):
			break
		}
	}
	// log.Print("final download url:", url)

	if strings.Contains(url, "192.168.") {
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

			b.data, err = readWithTimeout(req, resp, size, b.data, quit)

			if err == nil {
				select {
				case output <- b:
					return
				case <-quit:
					return
				}
			} else {
				log.Print(err)
				bytes, _ := httputil.DumpResponse(resp, false)
				log.Print(string(bytes))
				log.Printf("download wrong data:%d,%d,%d", b.from, b.to, len(b.data))
			}
		}

		select {
		case <-quit:
			return
		case <-time.After(time.Second):
			break
		}
	}
}

func readWithTimeout(req *http.Request, resp *http.Response, size int64, data []byte, quit chan bool) ([]byte, error) {
	// buffer := bytes.NewBuffer(data)
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

	_, err := io.ReadFull(resp.Body, data)
	defer resp.Body.Close()

	close(finish)

	return data, err
}
