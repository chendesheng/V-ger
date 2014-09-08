package subtitles

import (
	"log"
	"net/http"
	"net/url"
	"runtime/debug"
	"vger/httpex"

	"github.com/peterbourgon/html"
)

// func Addic7edSubtitle(keywords string, quit chan struct{}) (name string, content string) {

// 	a := addic7ed{}
// 	a.quit = quit
// 	a.search()

// 	return a.downloadSubtitle()
// }

type addic7ed struct {
	name string
	quit chan struct{}
}

func (a *addic7ed) search(result chan Subtitle) error {
	log.Printf("Addic7ed search subtitle: %s", a.name)

	defer func() {
		r := recover()
		if r != nil {
			log.Print(r)
			log.Print(string(debug.Stack()))
		}
	}()

	params := url.Values{
		"search": {a.name},
		"Submit": {"Search"},
	}

	resp, err := httpex.Get("http://www.addic7ed.com/search.php?"+params.Encode(), a.quit)
	if err != nil {
		return nil
	}

	defer resp.Body.Close()

	searchResultPage := resp.Request.URL.String()

	doc, err := html.Parse(resp.Body)

	if err != nil {
		println(err)
	}
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Data == "body" {
			for _, c := range getTag(getTag(n, "center")[1], "div") {
				if hasId(c, "container95m") {
					if ok, href := a.parseItem(c); ok {
						u, _ := url.Parse(searchResultPage)
						uhref, _ := url.Parse(href)
						url := u.ResolveReference(uhref).String()
						header := http.Header{}
						header.Add("Referer", searchResultPage)
						result <- Subtitle{url, "", "Addic7ed", header}
						return
					}
				}
			}
			return
		}

		for _, c := range n.Child {
			f(c)
		}
	}
	f(doc)

	// name, content := downloadSubtitle(searchResultPage, subtitleHref)
	// result <- Subtitle{string(content), name, "Addic7ed"}
	return nil
}

func (a *addic7ed) parseItem(n *html.Node) (bool, string) {
	table := getTag1(getTag(getTag(getClass1(n, "tabel95").Child[0], "tr")[1], "td")[1], "table")
	if table != nil {
		tr := getTag(table.Child[0], "tr")[2]
		tds := getTag(tr, "td")
		if len(tds) >= 5 && getText(tds[2]) == "English" && getText(getTag1(tds[3], "b")) == "Completed" {
			href := getAttr(getClass1(tds[4], "buttonDownload"), "href")
			return true, href
		}
	}

	return false, ""
}
