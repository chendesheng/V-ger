package subtitles

import (
	"github.com/peterbourgon/html"
	// "log"
	"net/url"
	"regexp"
	"strings"
)

func yyetsParseSub(n *html.Node) Subtitle {
	sub := Subtitle{}

	a := getTag1(getClass1(getClass1(n, "search_info_ls"), "all_search_li2"), "a")

	pageUrl := getAttr(a, "href")
	id := pageUrl[strings.LastIndex(pageUrl, "/")+1:]

	sub.URL = "http://www.yyets.com/php/subtitle/index/download?id=" + id

	text := getRigOfTags(a)
	regClean := regexp.MustCompile("([[][^]]*[]])")

	sub.Description = regClean.ReplaceAllString(text, "")

	sub.Source = "YYets"
	return sub
}
func yyetsSearchSubtitles(name string, result chan Subtitle) error {
	resp, err := Client.Get("http://www.yyets.com/php/search/index?type=subtitle&order=uptime&keyword=" + url.QueryEscape(name))
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
				for _, c := range getTag(n, "li") {
					s := yyetsParseSub(c)
					// log.Printf("%v", s)
					result <- s

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
