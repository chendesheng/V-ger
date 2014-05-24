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

func fetchN(req *http.Request, n int) (resp *http.Response, err error) {
	for i := 0; i < n; i++ {
		resp, err = http.DefaultClient.Do(req)
		if err == nil {
			defer resp.Body.Close()
			return
		}
	}

	return
}

func GetDownloadInfo(url string) (finalUrl string, name string, size int64, err error) {
	req := createDownloadRequest(url, -1, -1)

	resp, err := fetchN(req, 3)
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
