package subtitles

import (
	// "fmt"
	// "http/httputil/url"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

var Client *http.Client

type Subtitle struct {
	URL         string
	Description string
}

func (s *Subtitle) String() string {
	return s.Description
}

func sendGet(url string, params *url.Values) string {
	if params != nil {
		url = url + "?" + params.Encode()
	}
	resp, err := Client.Get(url)
	if err != nil {
		log.Fatal(err)
	}

	dumpBytes, _ := httputil.DumpResponse(resp, true)
	log.Println(string(dumpBytes))

	text := readBody(resp)
	return text
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
	return concat(shooterSearch(name), yyetsSearchSubtitles(name))
}
