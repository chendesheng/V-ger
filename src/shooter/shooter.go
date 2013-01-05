package shooter

import (
	"fmt"
	"github.com/peterbourgon/html"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"regexp"
	"strings"
	"time"
)

var Client *http.Client

type Subtitle struct {
	Name        string
	URL         string
	Description string
}

func (s *Subtitle) String() string {
	return s.Name + "\n" + s.Description
}

func attr2Map(attrs []html.Attribute) (m map[string]string) {
	m = make(map[string]string, len(attrs))
	for _, a := range attrs {
		m[a.Key] = a.Val
	}
	return
}
func hasClass(n *html.Node, class string) bool {
	return hasAttr(n, "class", class)
}
func hasId(n *html.Node, id string) bool {
	return hasAttr(n, "id", id)
}
func hasAttr(n *html.Node, key string, val string) bool {
	for _, a := range n.Attr {
		if a.Key == key && a.Val == val {
			return true
		}
	}

	return false
}
func getAttr(n *html.Node, key string) string {
	for _, a := range n.Attr {
		if a.Key == key {
			return a.Val
		}
	}

	return ""
}
func getText(n *html.Node) string {
	for _, c := range n.Child {
		if c.Type == html.TextNode {
			return strings.TrimSpace(c.Data)
		}
		if c.Type == html.ElementNode && c.Data == "a" {
			return getText(c)
		}
	}
	return ""
}
func getClass(n *html.Node, class string) []*html.Node {
	list := make([]*html.Node, 0)
	for _, c := range n.Child {
		if hasClass(c, class) {
			list = append(list, c)
		}
	}
	return list
}
func getClass1(n *html.Node, class string) *html.Node {
	for _, c := range n.Child {
		if hasClass(c, class) {
			return c
		}
	}
	return nil
}
func getId(n *html.Node, id string) *html.Node {
	for _, c := range n.Child {
		if hasId(c, id) {
			return c
		}
	}
	return nil
}
func getRidOfSpace(text string) string {
	r, _ := regexp.Compile(`[\s\r\n]+`)
	r1, _ := regexp.Compile(`：[\s\r\n]+`)
	return r1.ReplaceAllString(r.ReplaceAllString(text, " "), "：")
}
func getSubDesc(n *html.Node) string {
	desc := ""

	for _, c := range n.Child {
		if c.Data == "li" {
			temp := ""
			for _, cc := range c.Child {
				if cc.Type == html.TextNode {
					temp += strings.TrimSpace(cc.Data)
				}
				if cc.Type == html.ElementNode && cc.Data == "span" {
					temp += getText(cc)
				}
			}
			temp = strings.Replace(temp, "file,", "", 1)
			temp = strings.Replace(temp, "sub,", "", 1)
			desc += strings.TrimSpace(getRidOfSpace(temp)) + "\n"
		}
	}

	return desc
}
func getSub(n *html.Node) Subtitle {
	sub := Subtitle{}

	a := getClass1(getClass1(getClass1(n, "sublist_box_title"), "sublist_box_title_l"), "introtitle")
	sub.Name = getText(a)
	sub.URL = getAttr(a, "href")

	sub.Description = getSubDesc(getId(n, "sublist_ul"))

	log.Println(sub.Name)
	log.Println(sub.URL)
	log.Println(sub.Description)

	return sub
}
func SearchSubtitles(name string) []Subtitle {
	resp, err := Client.Get("http://www.shooter.cn/search/" + url.QueryEscape(name))
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	doc, err := html.Parse(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	subs := make([]Subtitle, 0)
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Data == "div" {
			if hasId(n, "resultsdiv") {
				for _, c := range getClass(n, "subitem") {
					subs = append(subs, getSub(c))
				}
				return
			}
		}

		for _, c := range n.Child {
			f(c)
		}
	}
	f(doc)

	return subs
}

//figure out file name from url while input name is empty
//return file name

func getFileId(text string) string {
	reg, _ := regexp.Compile(`var gFileidToBeDownlaod = ([^;]+);`)
	return reg.FindStringSubmatch(text)[1]
}
func getHash(text string) string {
	reg, _ := regexp.Compile(`shtg_filehash[+]?="([^"]+)"`)
	// reg1, _ := regexp.Compile(`shtg_filehash="([^"]+)"`)

	// hash := reg1.FindAllStringSubmatch(text)[1]
	hash := ""
	for _, s := range reg.FindAllStringSubmatch(text, -1) {
		hash += s[1]
	}
	return hash
}
func getSubId(webPageURL string) string {
	i := strings.LastIndex(webPageURL, "/") + 1
	return webPageURL[i : len(webPageURL)-4]
}
func setSubIdAndFileIdCookie(subId string, fileId string) {
	client := Client

	cookie := http.Cookie{
		Name:    "sub" + subId,
		Value:   "1",
		Domain:  "shooter.cn",
		Expires: time.Now().AddDate(100, 0, 0),
	}
	cookie2 := http.Cookie{
		Name:    "file" + fileId,
		Value:   "1",
		Domain:  "shooter.cn",
		Expires: time.Now().AddDate(100, 0, 0),
	}
	cookies := []*http.Cookie{&cookie, &cookie2}
	url, _ := url.Parse("http://www.shooter.com")
	client.Jar.SetCookies(url, cookies)
}
func decryptUrl(encryptedUrl string) string {
	a := encryptedUrl
	b := func(j string) string {
		g := ""

		for _, h := range j {
			if h+47 >= 126 {
				g += string(uint8(32 + (h+47)%126)) //32: space
			} else {
				g += string(uint8(h + 47))
			}
		}
		return g
	}

	d := func(g string) string {
		var j = len(g)
		j = j - 1
		h := ""
		for f := j; f >= 0; f-- {
			h += string(g[f])
		}
		return h
	}
	c := func(j string, h uint8, g uint8, f uint8) string {
		lj := uint8(len(j))
		return j[lj-f+g-h:lj-f+g] + j[lj-f:lj-f+g-h] + j[lj-f+g:] + j[0:lj-f]
	}

	if len(a) > 32 {
		switch string(a[0]) {
		case "o":
			return (b((c(a[1:], 8, 17, 27))))
			break
		case "n":
			return (b(d(c(a[1:], 6, 15, 17))))
			break
		case "m":
			return (d(c(a[1:], 6, 11, 17)))
			break
		case "l":
			return (d(b(c(a[1:], 6, 12, 17))))
			break
		case "k":
			return (c(a[1:], 14, 17, 24))
			break
		case "j":
			return (c(b(d(a[1:])), 11, 17, 27))
			break
		case "i":
			return (c(d(b(a[1:])), 5, 7, 24))
			break
		case "h":
			return (c(b(a[1:]), 12, 22, 30))
			break
		case "g":
			return (c(d(a[1:]), 11, 15, 21))
		case "f":
			return (c(a[1:], 14, 17, 24))
		case "e":
			return (c(a[1:], 4, 7, 22))
		case "d":
			return (d(b(a[1:])))
		case "c":
			return (b(d(a[1:])))
		case "b":
			return (d(a[1:]))
		case "a":
			return b(a[1:])
			break
		}
	}
	return a
}
func GetDownloadUrl(webPageURL string) (string, string) {
	webPageURL = "http://www.shooter.cn" + webPageURL

	subId := getSubId(webPageURL)
	log.Println(subId)

	pageHtml := sendGet(webPageURL, nil)
	fileId := getFileId(pageHtml)
	log.Println(fileId)

	loadmain := sendGet("http://www.shooter.cn/a/loadmain.js", nil)
	hash := getHash(loadmain)
	log.Println(hash)

	log.Println(fmt.Sprintf("http://www.shooter.cn/files/file3.php?hash=%s&fileid=%s", hash, fileId))
	encryptedUrl := sendGet(fmt.Sprintf("http://www.shooter.cn/files/file3.php?hash=%s&fileid=%s", hash, fileId), nil)
	log.Println(encryptedUrl)
	url := decryptUrl(encryptedUrl)
	log.Println(url)

	b := strings.LastIndex(url, "/") + 1
	e := strings.Index(url, "?")
	return "http://file0.shooter.cn" + url, url[b:e]
}

func sendGet(url string, params *url.Values) string {
	if params != nil {
		url = url + "?" + params.Encode()
	}
	resp, err := Client.Get(url)
	if err != nil {
		log.Fatal(err)
	}

	dumpBytes, _ := httputil.DumpResponse(resp, true)
	log.Println(string(dumpBytes))

	text := readBody(resp)
	return text
}
func readBody(resp *http.Response) string {
	bytes, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		log.Fatal(err)
	}

	dumpBytes, _ := httputil.DumpResponse(resp, true)
	log.Println(string(dumpBytes))

	text := string(bytes)
	return text
}
