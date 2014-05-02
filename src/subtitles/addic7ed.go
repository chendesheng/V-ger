package subtitles

import (
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"runtime/debug"
	"github.com/peterbourgon/html"
	// "net/http/cookiejar"
	// "net/http/httputil"
	"net/url"
)

func Addic7edSubtitle(keywords string) (name string, content string) {
	defer func() {
		r := recover()
		if r != nil {
			log.Print(r)
			log.Print(string(debug.Stack()))

			name = ""
			content = ""
		}
	}()
	a := addic7ed{}
	a.search(keywords)

	return a.downloadSubtitle()
}

type addic7ed struct {
	searchResultPage string
	subtitleHref     string
}

func (a *addic7ed) search(keywords string) {
	params := url.Values{
		"search": {keywords},
		"Submit": {"Search"},
	}

	log.Print(params.Encode())

	resp, err := http.Get("http://www.addic7ed.com/search.php?" + params.Encode())
	if err != nil {
		println(err)
		return
	}

	defer resp.Body.Close()

	a.searchResultPage = resp.Request.URL.String()
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
						a.subtitleHref = href
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

func (a *addic7ed) downloadSubtitle() (string, string) {
	u, _ := url.Parse(a.searchResultPage)
	println(u)
	uhref, _ := url.Parse(a.subtitleHref)
	subtitleUrl := u.ResolveReference(uhref)

	log.Print("subtitle url:", subtitleUrl.String())

	req, _ := http.NewRequest("GET", subtitleUrl.String(), nil)
	req.Header.Add("Referer", a.searchResultPage)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Print(err)
		return "", ""
	}
	defer resp.Body.Close()
	// data, _ := httputil.DumpResponse(resp, true)
	// println(string(data))
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Print(err)
		return "", ""
	}

	name := ""

	contentDisposition := resp.Header["Content-Disposition"][0]
	regexFile := regexp.MustCompile(`filename="?([^"]+)"?`)

	if match := regexFile.FindStringSubmatch(contentDisposition); len(match) > 1 {
		name = match[1]
	}

	return name, string(data)
}
