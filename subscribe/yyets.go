package subscribe

import (
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
	"strings"
	"time"
	"vger/task"
	"vger/util"

	"vger/html"
)

var YYetsFormats map[string]struct{} //[]string{"1080p", "720p", "bd-720p", "web-dl"}

func init() {
	//init default formats
	YYetsFormats = make(map[string]struct{})
	YYetsFormats["1080p"] = struct{}{}
	YYetsFormats["720p"] = struct{}{}
	YYetsFormats["bd-720p"] = struct{}{}
	YYetsFormats["web-dl"] = struct{}{}
}

func parseSubscribeInfo(r io.Reader) (s *Subscribe, err error) {
	s = &Subscribe{}
	var doc *html.Node
	doc, err = html.Parse(r)

	if err != nil {
		return nil, err
	}

	defer func() {
		r := recover()
		if r != nil {
			err = r.(error)
			return
		}
	}()
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Data == "div" {
			if hasClass(n, "resource-con") {
				s.Banner = getAttr(getTag1(getTag1(getClass1(n, "fl-img"), "p"), "a"), "href")
				props := getTag(getTag1(getClass1(n, "fl-info"), "ul"), "li")
				for _, c := range props {
					span := getTag1(c, "span")
					if span != nil {
						k := getText(span)
						if k == "英文：" {
							if len(c.Child) > 1 {
								s.Name = s.Name + getRidOfTags(c.Child[1])
							}
						}
					}
				}
				s.Name = strings.TrimLeft(s.Name, ".")
				s.Name = strings.TrimSpace(s.Name)
				s.Name = strings.Replace(s.Name, "/", "|", -1)
				s.Name = strings.Replace(s.Name, "\\", "|", -1)
				return
			}
		}

		for _, c := range n.Child {
			f(c)
		}
	}

	f(doc)

	return
}

func parseEpisodes(r io.Reader) (result []*task.Task, err error) {

	doc, err := html.Parse(r)

	if err != nil {
		return nil, err
	}

	defer func() {
		r := recover()
		if r != nil {
			err = r.(error)
			return
		}
	}()

	result = make([]*task.Task, 0)

	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Data == "div" {
			if hasClass(n, "media-list") {
				for _, li := range getTag(getTag1(n, "ul"), "li") {
					ft := strings.ToLower(getAttr(li, "format"))
					if _, ok := YYetsFormats[ft]; ok {
						season, _ := strconv.Atoi(getAttr(li, "season"))
						episode, _ := strconv.Atoi(getAttr(li, "episode"))

						t := parseSingle(li)
						t.Season = season
						t.Episode = episode

						result = append(result, t)
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

	return
}
func Parse(url string) (s *Subscribe, result []*task.Task, err error) {
	YYetsLogin(util.ReadConfig("yyets-user"), util.ReadConfig("yyets-password"))

	i := strings.LastIndex(url, "/")
	downloadPageUrl := fmt.Sprintf("%s/list%s", url[:i], url[i:])
	log.Println("downloadPageUrl:", downloadPageUrl)

	//parse episodes list
	resp, err := http.Get(downloadPageUrl)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()

	t, err := parseEpisodes(resp.Body)
	if err != nil {
		return nil, nil, err
	}
	log.Println("counts:", len(t))

	//parse subscirbe info
	resp, err = http.Get(url)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()

	s, err = parseSubscribeInfo(resp.Body)
	if err != nil {
		return nil, nil, err
	}

	s.URL = url
	s.Source = "YYets"
	s.Autodownload = true
	if len(s.Name) == 0 {
		respdata, _ := httputil.DumpResponse(resp, false)
		log.Print(string(respdata))
		log.Print(t)
		log.Print(s)
		return nil, nil, errors.New("No subscirbe found")
	}

	filter := make(map[int]map[int]struct{})
	var tks []*task.Task
	for _, tk := range t {
		if season, ok := filter[tk.Season]; ok {
			if _, ok := season[tk.Episode]; !ok {
				season[tk.Episode] = struct{}{}
				tks = append(tks, tk)
			}
		} else {
			filter[tk.Season] = make(map[int]struct{})
			filter[tk.Season][tk.Episode] = struct{}{}
			tks = append(tks, tk)
		}

		tk.Subscribe = s.Name
	}

	return s, tks, nil
}

func downloadBannerImage(url string) (string, error) {

	resp, err := http.Get(url)
	// bytes, err := httputil.DumpResponse(resp, false)
	// fmt.Println(string(bytes))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	// print(len(data))

	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(data), nil

	// return true, nil
}

func parseSingle(n *html.Node) *task.Task {
	t := &task.Task{}
	t.Status = "New"

	t.Name = getText(getTag1(getClass1(n, "fl"), "a"))
	t.StartTime = time.Now().Unix()

	c := getClass1(n, "fr")
	if ed2k := getChildAttr1(c, "type", "ed2k"); ed2k != nil {
		t.Original = getAttr(ed2k, "href")
	} else if magnet := getChildAttr1(c, "type", "magnet"); magnet != nil {
		t.Original = getAttr(magnet, "href")
	} else if thunder := getChildAttr1(c, "pvgniyjm", "*"); thunder != nil {
		t.Original = getAttr(thunder, "pvgniyjm")
	} else {
		t.Original = ""
	}

	return t
}

func YYetsLogin(name string, password string) {
	if YYetsIfLogined() {
		return
	}

	data := url.Values{
		"url_back": []string{""},
		"from":     []string{"loginpage"},
		"account":  []string{name},
		"password": []string{password},
		"remember": []string{"1"},
	}

	resp, err := http.PostForm("http://www.zimuzu.tv/user/login/ajaxlogin", data)
	if err != nil {
		log.Print(err)
	} else {
		log.Print("YYets login success")
	}
	bytes, _ := httputil.DumpResponse(resp, false)
	log.Print(string(bytes))

	resp.Body.Close()
}

func YYetsIfLogined() bool {
	u, _ := url.Parse("http://www.zimuzu.tv")
	cookies := http.DefaultClient.Jar.Cookies(u)
	for _, c := range cookies {
		if c.Name == "GINFO" && strings.HasPrefix(c.Value, "uid") {
			return true
		}
	}

	return false
}
