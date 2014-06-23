package download

import (
	"log"
	"toutf8"
	// "bytes"
	"fmt"
	// "io"
	// "log"
	// "native"b
	"net/http"
	// "os"
	// "path/filepath"
	// "regexp"
	// "runtime"
	"strings"
	"time"
)

var NetworkTimeout time.Duration

func fetchN(req *http.Request, n int, quit chan struct{}) (resp *http.Response, err error) {
	finish := make(chan struct{})
	go func() {
		defer close(finish)

		for i := 0; i < n; i++ {
			resp, err = http.DefaultClient.Do(req)
			if err == nil {
				return
			}
			time.Sleep(100 * time.Millisecond)
		}
	}()

	select {
	case <-quit:
		cancelRequest(req)
		err = errStopFetch
		break
	case <-finish:
		break
	}

	return
}

func GetDownloadInfo(url string) (finalUrl string, name string, size int64, err error) {
	return GetDownloadInfoN(url, 3, nil)
}

func GetDownloadInfoN(url string, retryTimes int, quit chan struct{}) (finalUrl string, name string, size int64, err error) {
	req := createDownloadRequest(url, -1, -1)

	resp, err := fetchN(req, retryTimes, quit)
	if err != nil {
		return "", "", 0, err
	}
	defer resp.Body.Close()

	finalUrl = resp.Request.URL.String()

	name, size = getFileInfo(resp.Header)
	if name == "" {
		name = getFileName(finalUrl)
	}

	if name == "" || size == 0 {
		err = fmt.Errorf("Broken resource")
	}

	encoding, err := toutf8.GuessEncoding([]byte(name))
	if err != nil {
		log.Print(err)
	}

	if encoding != "utf-8" && encoding != "ascii" {
		log.Print("file name encoding:", encoding)
		utf8name, err := toutf8.ConvertToUTF8From(name, "gb18030")
		if err != nil {
			log.Print(err)
		} else {
			name = utf8name
		}
	}

	name = strings.Replace(name, "/", "|", -1)
	name = strings.Replace(name, "\\", "|", -1)
	name = strings.TrimLeft(name, ".")

	return
}
