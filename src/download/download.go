package download

import (
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
	"util"
)

var networkTimeout time.Duration

func init() {
	networkTimeout = time.Duration(util.ReadIntConfig("network-timeout")) * time.Second

	// http.DefaultTransport.(*http.Transport).Dial = func(network, addr string) (net.Conn, error) {
	// 	c, err := net.Dial(network, addr)
	// 	log.Printf("%v", c)
	// 	return c, err
	// }
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

	name = strings.Replace(name, "/", "|", -1)
	name = strings.Replace(name, "\\", "|", -1)
	name = strings.TrimLeft(name, ".")

	if name == "" && size == 0 {
		err = fmt.Errorf("Broken resource")
	}
	return
}
