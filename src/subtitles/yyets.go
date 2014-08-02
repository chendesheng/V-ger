package subtitles

import (
	"net/http"
	"net/url"
	"regexp"
	"strings"

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
func yyetsSearchSubtitles(name string, result chan Subtitle, quit chan struct{}) error {
	resp, err := http.Get("http://www.yyets.com/search/index?type=subtitle&order=uptime&keyword=" + url.QueryEscape(name))
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
					// log.Printf("%v", s)
					select {
					case result <- s:
						break
					case <-quit:
						return
					}

					if count++; count > 10 {
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
