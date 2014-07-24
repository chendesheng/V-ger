package download

import (
	"block"
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
	isFinalUrl    bool
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

	if !df.isFinalUrl {
		req, _ := http.NewRequest("GET", url, nil)
		resp, err := fetchN(req, 1000000, df.quit)
		if err != nil {
			log.Print(err)
			return
		}
		url = resp.Request.URL.String()
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

			// trace(fmt.Sprint("download filter input:", b.From, b.to))

			df.downloadBlock(url, b)
		case <-df.quit:
			// fmt.Println("downloadRoutine quit")
			return
		}
	}
}
func (df *downloadFilter) downloadBlock(url string, b block.Block) {
	for {
		req := createDownloadRequest(url, b.From, b.From+int64(len(b.Data))-1)
		err := requestWithTimeout(req, b.Data, df.quit)

		if err == nil {
			// println("download routine write output:", b.From)
			df.writeOutput(b)
			// trace(fmt.Sprint("downloadFilter writeoutput:", b.From, b.to))
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

func requestWithTimeout(req *http.Request, data []byte, quit chan struct{}) (err error) {
	finish := make(chan error)
	var resp *http.Response
	go func() {
		defer close(finish)
		var err error
		resp, err = http.DefaultClient.Do(req)
		if err != nil {
			return
		}
		defer resp.Body.Close()

		_, err = io.ReadFull(resp.Body, data)
		// if err != nil {
		// 	println("download routine ReadFull:", err.Error())
		// }
	}()

	// println("download routine NetworkTimeout:", NetworkTimeout)
	select {
	case <-time.After(NetworkTimeout): //cancelRequest if time.After before close(finish)
		// println("download routine timeout")
		cancelRequest(req)
		err = errReadTimeout //return not nil error is required
		break
	case <-quit:
		cancelRequest(req)
		err = errStopFetch
		break
	case <-finish:
		// println("download routine finish")
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
