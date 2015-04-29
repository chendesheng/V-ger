package subtitles

import (
	"errors"
	"log"
	"net/url"
	"vger/httpex"

	"vger/html"
)

func yyetsParseSub(n *html.Node, quit chan struct{}) *Subtitle {
	sub := &Subtitle{}

	a := getTag1(getTag1(getTag1(getTag1(n, "div"), "div"), "div"), "a")

	pageURL := "http://www.zimuzu.tv" + getAttr(a, "href")

	var err error
	sub.URL, sub.Description, err = getDownloadLink(pageURL, quit)
	if err != nil {
		log.Print(err)
		return nil
	}

	sub.Source = "YYets"
	return sub
}

func getDownloadLink(url string, quit chan struct{}) (string, string, error) {
	log.Println("YYets: getDownloadLink", url)
	resp, err := httpex.Get(url, quit)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()

	doc, err := html.Parse(resp.Body)

	if err != nil {
		return "", "", err
	}

	var f func(*html.Node) (string, string, error)
	f = func(n *html.Node) (string, string, error) {
		if n.Data == "div" {
			if hasClass(n, "subtitle-links") {
				a := getTag1(getTag1(n, "h3"), "a")
				return getAttr(a, "href"), getRidOfTags(a), nil
			}
		}

		for _, c := range n.Child {
			if lnk, des, err := f(c); err == nil {
				return lnk, des, nil
			}
		}

		return "", "", errors.New("YYest: Can't find subtitle link")
	}

	return f(doc)
}

type yyetsSearch struct {
	name   string
	maxcnt int
	quit   chan struct{}
}

func (y *yyetsSearch) search(result chan Subtitle) error {
	log.Printf("YYets search subtitle: %s %d", y.name, y.maxcnt)
	resp, err := httpex.Get("http://www.zimuzu.tv/search/index?type=subtitle&order=uptime&keyword="+url.QueryEscape(y.name), y.quit)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	doc, err := html.Parse(resp.Body)

	if err != nil {
		return err
	}

	count := 0
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Data == "div" {
			if hasClass2(n, "search-result") {
				n1 := getTag1(n, "ul")
				for _, c := range getTag(n1, "li") {
					s := yyetsParseSub(c, y.quit)
					if s == nil {
						continue
					}

					select {
					case <-y.quit:
						return
					case result <- *s:
					}
					if count++; count >= y.maxcnt {
						return
					}
				}
				return
			}
		}

		for _, c := range n.Child {
			f(c)
		}
	}
	f(doc)

	return nil
}
