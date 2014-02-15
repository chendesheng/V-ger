package subtitles

import (
	// "fmt"
	// "io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	// "os"
)

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
	resp, err := http.Get(url)
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

func SearchSubtitles(name string, url string, result chan Subtitle) {
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

	if len(url) > 0 && strings.Contains(url, "gdl.lixian.vip.xunlei.com") {
		kankanFinish := make(chan bool)
		go func() {
			defer close(kankanFinish)
			defer func() {
				r := recover()
				if r != nil {
					log.Print(r.(error))
				}
			}()

			err := kankanSearch(url, result)
			if err != nil {
				log.Print(err)
			}

		}()

		<-kankanFinish
	}

	<-yyetsFinish
	<-shooterFinish

	println("search subtitle finish")
	close(result)
}

func SearchSubtitlesMaxCount(name string, result chan Subtitle, maxcnt int) {
	yyetsRes := make(chan Subtitle)
	yyetsCnt := 0
	go func() {
		err := yyetsSearchSubtitles(name, yyetsRes)
		if err != nil {
			log.Println(err)
		}
		close(yyetsRes)
	}()

	shooterRes := make(chan Subtitle)
	shooterCnt := 0
	go func() {
		err := shooterSearch(name, shooterRes)
		if err != nil {
			log.Println(err)
		}
		close(shooterRes)
	}()

	var yyetsFinish, shoooterFinish bool
	var s Subtitle
	for !(yyetsFinish && shoooterFinish) {
		select {
		case s, yyetsFinish = <-yyetsRes:
			yyetsCnt++
			if yyetsCnt < maxcnt {
				result <- s
			}
			break
		case s, shoooterFinish = <-shooterRes:
			shooterCnt++
			if shooterCnt < maxcnt {
				result <- s
			}
			break
		}

		if yyetsCnt >= maxcnt && shooterCnt >= maxcnt {
			close(result)
		}
	}

	close(result)
}
func QuickDownload(url string) ([]byte, error) {
	resp, err := http.Get(url)
	defer resp.Body.Close()
	// bytes, err := httputil.DumpResponse(resp, false)
	// fmt.Println(string(bytes))
	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadAll(resp.Body)
	// print(len(data))

	// if err != nil {
	// 	return err
	// }

	// fmt.Println(data)

	return data, err
}
