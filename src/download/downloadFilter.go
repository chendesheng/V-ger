package download

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"strings"
	"sync"
	"time"
)

var errStopFetch = errors.New("stop fetch")
var errReadTimeout = errors.New("read timeout")

type downloadFilter struct {
	basicFilter
	url           string
	routineNumber int
}

func (df *downloadFilter) active() {
	defer df.closeOutput()

	wg := sync.WaitGroup{}
	wg.Add(df.routineNumber)

	for i := 0; i < df.routineNumber; i++ {
		go func() {
			defer wg.Done()
			df.downloadRoutine()
		}()
	}

	wg.Wait()

	log.Print("downloadFilter return")
}
func (df *downloadFilter) downloadRoutine() {
	url := df.url

	url, _, _, err := GetDownloadInfoN(url, 10000000, df.quit)
	if err != nil {
		return
	}

	if strings.Contains(url, "192.168.") {
		//AUSU router may redirect to error_page.html, download from this url will crap target file.
		return
	}

	for {
		select {
		case b, ok := <-df.input:
			if !ok {
				fmt.Println("downloadRoutine finish")
				return
			}

			// trace(fmt.Sprint("download filter input:", b.from, b.to))

			df.downloadBlock(url, b)
		case <-df.quit:
			fmt.Println("downloadRoutine quit")
			return
		}
	}
}
func (df *downloadFilter) downloadBlock(url string, b block) {
	for {
		req := createDownloadRequest(url, b.from, b.from+int64(len(b.data))-1)
		err := requestWithTimeout(req, b.data, df.quit)

		if err == nil {
			df.writeOutput(b)
			// trace(fmt.Sprint("downloadFilter writeoutput:", b.from, b.to))
			return
		} else {
			select {
			case <-df.quit:
				return
			default:
			}
		}
		df.wait(100 * time.Millisecond)
	}
}

func requestWithTimeout(req *http.Request, data []byte, quit chan bool) (err error) {
	finish := make(chan error)
	var resp *http.Response
	go func() {
		defer close(finish)

		resp, err = http.DefaultClient.Do(req)
		if err != nil {
			return
		}
		defer resp.Body.Close()

		_, err = io.ReadFull(resp.Body, data)
	}()

	select {
	case <-time.After(NetworkTimeout): //cancelRequest if time.After before close(finish)
		cancelRequest(req)
		err = errReadTimeout //return not nil error is required
		break
	case <-quit:
		cancelRequest(req)
		err = errStopFetch
		break
	case <-finish:
		if err != nil {
			log.Print(err)
			if resp != nil {
				bytes, _ := httputil.DumpResponse(resp, false)
				log.Print(string(bytes))
			}
		}
		break
	}

	return
}
