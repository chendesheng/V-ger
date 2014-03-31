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

func SearchSubtitles(name string, url string, result chan Subtitle, quit chan bool) {
	yyetsFinish := make(chan bool)
	go func() {
		err := yyetsSearchSubtitles(name, result, quit)
		if err != nil {
			log.Println(err)
		}
		close(yyetsFinish)
	}()

	shooterFinish := make(chan bool)
	go func() {
		err := shooterSearch(name, result, quit)
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

			err := kankanSearch(url, result, quit)
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

func SearchSubtitlesMaxCount(name string, url string, result chan Subtitle, maxcnt int, quit chan bool) {
	yyetsRes := make(chan Subtitle)
	yyetsQuit := make(chan bool)
	yyetsCnt := 0
	go func() {
		err := yyetsSearchSubtitles(name, yyetsRes, yyetsQuit)
		if err != nil {
			log.Println(err)
		}
		close(yyetsRes)
	}()

	shooterRes := make(chan Subtitle)
	shooterQuit := make(chan bool)
	shooterCnt := 0
	go func() {
		err := shooterSearch(name, shooterRes, shooterQuit)
		if err != nil {
			log.Println(err)
		}
		close(shooterRes)
	}()

	kankanFinish := make(chan bool)
	if len(url) > 0 && strings.Contains(url, "gdl.lixian.vip.xunlei.com") {
		go func() {
			defer close(kankanFinish)
			defer func() {
				r := recover()
				if r != nil {
					log.Print(r.(error))
				}
			}()

			err := kankanSearch(url, result, quit)
			if err != nil {
				log.Print(err)
			}
			log.Print("kankan search finished")
		}()
	} else {
		close(kankanFinish)
	}

	yyetsOK, shoooterOK := true, true
	var s Subtitle
	for yyetsOK || shoooterOK {
		select {
		case s, yyetsOK = <-yyetsRes:
			if yyetsOK {
				yyetsCnt++
				if yyetsCnt < maxcnt {
					println("yyets:", s.Description)
					result <- s
				} else {
					close(yyetsQuit)
					yyetsOK = false
				}
			}
			break
		case s, shoooterOK = <-shooterRes:
			if shoooterOK {
				shooterCnt++
				if shooterCnt < maxcnt {
					println("shooter:", s.Description)
					result <- s
				} else {
					close(shooterQuit)
					shoooterOK = false
				}
			}
			break
		case <-quit:
			close(result)
			return
		}
	}

	<-kankanFinish
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
