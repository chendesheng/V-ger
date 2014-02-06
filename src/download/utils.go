package download

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func createDownloadRequest(urlString string, from, to int64) *http.Request {
	req := new(http.Request)
	req.Method = "GET"
	req.URL, _ = url.Parse(urlString)
	req.Header = make(http.Header)
	// fmt.Println(urlString)
	addRangeHeader(req, from, to)
	return req
}
func cancelRequest(req *http.Request) {
	http.DefaultTransport.(*http.Transport).CancelRequest(req)
}
func addRangeHeader(req *http.Request, from, to int64) {
	if from == to || (from <= 0 && to < 0) {
		return
	}
	if to < 0 {
		req.Header.Add("Range", fmt.Sprintf("bytes=%d-", from))
	} else {
		req.Header.Add("Range", fmt.Sprintf("bytes=%d-%d", from, to))
	}
}

func openOrCreateFileRW(path string, position int64) (*os.File, error) {
	// log.Print("open or create file " + path)

	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}

	offset, err := f.Seek(position, 0)
	if err != nil {
		return nil, err
	}
	if offset != position {
		return nil, fmt.Errorf("\nerror offset")
	}
	return f, nil
}
func getFileInfo(header http.Header) (name string, size int64) {
	// log.Printf("%v\n", header)
	if len(header["Content-Disposition"]) > 0 {
		contentDisposition := header["Content-Disposition"][0]
		regexFile := regexp.MustCompile(`filename="?([^"]+)"?`)

		if match := regexFile.FindStringSubmatch(contentDisposition); len(match) > 1 {
			name = match[1]
		} else {
			name = ""
		}
	}

	if cr := header["Content-Range"]; len(cr) > 0 {
		regexSize := regexp.MustCompile(`/(\d+)`)

		sizeStr := regexSize.FindStringSubmatch(cr[0])[1]
		size, _ = strconv.ParseInt(sizeStr, 10, 64)
	} else {
		if len(header["Content-Length"]) > 0 {
			size, _ = strconv.ParseInt(header["Content-Length"][0], 10, 64)
		} else {
			size = 0
		}
	}

	return
}
func getFileName(fullURL string) string {
	e := strings.Index(fullURL, "?")
	if e < 0 {
		e = len(fullURL)
	}
	s := strings.LastIndex(fullURL, `/`) + 1
	if s >= e {
		return ""
	}

	name, _ := url.QueryUnescape(fullURL[s:e])
	return name
}
