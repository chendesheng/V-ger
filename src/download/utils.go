package download

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
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
func openOrCreateFileRW(path string, position int64) *os.File {
	log.Print("open or create file " + path)

	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		log.Fatal(err)
	}

	offset, err := f.Seek(position, 0)
	if err != nil {
		log.Fatal(err)
	}
	if offset != position {
		fmt.Println("\nerror offset")
		os.Exit(1)
	}
	return f
}
func getFileInfo(header http.Header) (name string, size int64) {
	if len(header["Content-Disposition"]) > 0 {
		contentDisposition := header["Content-Disposition"][0]
		regexFile := regexp.MustCompile(`filename="([^"]+)"`)

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
		size, _ = strconv.ParseInt(header["Content-Length"][0], 10, 64)
	}

	return
}
func getFileName(fullURL string) string {
	e := strings.Index(fullURL, "?")
	if e < 0 {
		e = len(fullURL)
	}
	name, _ := url.QueryUnescape(fullURL[strings.LastIndex(fullURL, `/`)+1 : e])
	return name
}
func writeJson(path string, object interface{}) {
	data, err := json.Marshal(object)
	if err != nil {
		log.Fatal(err)
	}

	ioutil.WriteFile(path, data, 0666)
}
func readJson(path string, object interface{}) error {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	// log.Println("read json: ")
	// log.Println(path)
	// log.Println(string(data))

	return json.Unmarshal(data, &object)
}
