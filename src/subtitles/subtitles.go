package subtitles

import (
	// "fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
)

var Client *http.Client

type Subtitle struct {
	URL         string
	Description string
	Source      string
}

func (s *Subtitle) String() string {
	return s.Description
}

func sendGet(url string, params *url.Values) (string, error) {
	if params != nil {
		url = url + "?" + params.Encode()
	}
	resp, err := Client.Get(url)
	if err != nil {
		return "", err
	}

	text := readBody(resp)
	return text, nil
}
func readBody(resp *http.Response) string {
	bytes, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		log.Fatal(err)
	}

	dumpBytes, _ := httputil.DumpResponse(resp, true)
	log.Println(string(dumpBytes))

	text := string(bytes)
	return text
}

func concat(old1, old2 []Subtitle) []Subtitle {
	newslice := make([]Subtitle, len(old1)+len(old2))
	copy(newslice, old1)
	copy(newslice[len(old1):], old2)
	return newslice
}

func SearchSubtitles(name string) []Subtitle {
	// return yyetsSearchSubtitles(name)
	// return shooterSearch(name)
	yyetsSubs := make(chan []Subtitle)
	go func() {
		yyetsSubs <- yyetsSearchSubtitles(name)
	}()
	shooterSubs := make(chan []Subtitle)
	go func() {
		shooterSubs <- shooterSearch(name)
	}()

	return concat(<-yyetsSubs, <-shooterSubs)
}
func QuickDownload(url, path string) (bool, error) {
	resp, err := Client.Get(url)
	// bytes, err := httputil.DumpResponse(resp, false)
	// fmt.Println(string(bytes))
	if err != nil {
		return false, err
	}

	data, err := ioutil.ReadAll(resp.Body)
	// print(len(data))
	defer resp.Body.Close()

	if err != nil {
		return false, err
	}

	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, 0666)
	defer f.Close()
	if err != nil {
		return false, err
	}
	// fmt.Println(data)
	f.WriteAt(data, 0)
	return true, nil
}
