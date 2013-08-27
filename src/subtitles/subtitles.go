package subtitles

import (
	// "fmt"
	"io/ioutil"
	"log"
	"net/http"
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

	text := string(bytes)
	return text
}

func SearchSubtitles(name string, result chan Subtitle) {
	yyetsFinish := make(chan bool)
	go func() {
		err := yyetsSearchSubtitles(name, result)
		if err != nil {
			log.Println(err)
		}
		close(yyetsFinish)
	}()

	shooterFinish := make(chan bool)
	go func() {
		err := shooterSearch(name, result)
		if err != nil {
			log.Println(err)
		}
		close(shooterFinish)
	}()

	<-yyetsFinish
	<-shooterFinish
	close(result)
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
