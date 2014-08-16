package subtitles

import (
	"sync"

	// "fmt"
	// "io"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"

	// "os"
)

type subSearcher interface {
	search(chan Subtitle) error
}

type Subtitle struct {
	URL         string
	Description string
	Source      string
	Context     http.Header
}

func (s *Subtitle) String() string {
	return s.Description
}

func httpGet(url string, quit chan struct{}) (*http.Response, error) {
	finish := make(chan error)
	defer close(finish)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	go func() {
		select {
		case <-quit:
			cancelRequest(req)
		case <-finish:
		}
	}()

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func sendGet(url string, params *url.Values, quit chan struct{}) (string, error) {
	if params != nil {
		url = url + "?" + params.Encode()
	}
	resp, err := httpGet(url, quit)
	if err != nil {
		return "", err
	}

	return readBody(resp.Body), nil
}

func readBody(body io.ReadCloser) string {
	defer body.Close()
	bytes, err := ioutil.ReadAll(body)
	if err != nil {
		log.Fatal(err)
	}

	text := string(bytes)
	return text
}

func SearchSubtitles(name string, url string, result chan Subtitle, quit chan struct{}) {
	searchers := []subSearcher{
		&yyetsSearch{name, 1, quit},
		&shooterSearch{name, 2, quit},
		&kankanSearch{url, quit},
		&addic7ed{name, quit},
	}

	w := sync.WaitGroup{}
	w.Add(len(searchers))
	for _, s := range searchers {
		go func(s subSearcher) {
			if err := s.search(result); err != nil {
				log.Print(err)
			}
			w.Done()
		}(s)
	}
	w.Wait()
	close(result)
}
