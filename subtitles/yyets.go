package subtitles

import (
	"log"
	"net/url"
	"regexp"
	"strings"
	"vger/httpex"

	"github.com/peterbourgon/html"
)

func yyetsParseSub(n *html.Node) Subtitle {
	sub := Subtitle{}

	a := getTag1(getClass1(getClass1(n, "search_info_ls"), "all_search_li2"), "a")

	pageUrl := getAttr(a, "href")
	id := pageUrl[strings.LastIndex(pageUrl, "/")+1:]

	sub.URL = "http://www.yyets.com/subtitle/index/download?id=" + id

	text := getRidOfTags(a)
	regClean := regexp.MustCompile("([[][^]]*[]])")

	sub.Description = regClean.ReplaceAllString(text, "")

	sub.Source = "YYets"
	return sub
}

type yyetsSearch struct {
	name   string
	maxcnt int
	quit   chan struct{}
}

func (y *yyetsSearch) search(result chan Subtitle) error {
	log.Printf("YYets search subtitle: %s %d", y.name, y.maxcnt)

	resp, err := httpex.Get("http://www.yyets.com/search/index?type=subtitle&order=uptime&keyword="+url.QueryEscape(y.name), y.quit)
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
		if n.Data == "ul" {
			if hasClass(n, "allsearch") {
				for _, c := range getTag(n, "li")[1:] { //skip first item, first item is ad
					s := yyetsParseSub(c)

					select {
					case <-y.quit:
						return
					case result <- s:
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
