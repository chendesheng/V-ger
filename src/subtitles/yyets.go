package subtitles

import (
	"github.com/peterbourgon/html"
	"log"
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

	log.Printf("%v\n", sub)
	sub.Source = "YYets"
	return sub
}
func yyetsSearchSubtitles(name string) []Subtitle {
	resp, err := Client.Get("http://www.yyets.com/php/search/index?type=subtitle&order=uptime&keyword=" + url.QueryEscape(name))
	if err != nil {
		return make([]Subtitle, 0)
	}
	defer resp.Body.Close()

	doc, err := html.Parse(resp.Body)

	if err != nil {
		log.Fatal(err)
	}
	subs := make([]Subtitle, 0)
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Data == "ul" {
			if hasClass(n, "allsearch") {
				for _, c := range getTag(n, "li") {
					subs = append(subs, yyetsParseSub(c))
				}
				return
			}
		}

		for _, c := range n.Child {
			f(c)
		}
	}
	f(doc)

	if len(subs) > 10 {
		subs = subs[:10]
	}

	return subs
}
