package subscribe

import (
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"net/http/httputil"
	"net/url"
	"util"
	// "fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"task"
	"time"
	"github.com/peterbourgon/html"
)

func parseEpisodes(n *html.Node, season int, subscribeName string, format string, result *map[int]*task.Task) {
	for _, c := range getTag(n, "li") {
		if strings.ToLower(getAttr(c, "format")) == format {
			episode, _ := strconv.Atoi(getAttr(c, "episode"))
			if _, ok := (*result)[episode]; !ok {
				t := parseSingle(c)
				t.Subscribe = subscribeName
				t.Season = season
				t.Episode = episode
				// println(t)
				(*result)[t.Episode] = t
			}
		}
	}
}
func parse(r io.Reader) (s *Subscribe, result []*task.Task, err error) {

	doc, err := html.Parse(r)

	if err != nil {
		return nil, nil, err
	}

	defer func() {
		r := recover()
		if r != nil {
			err = r.(error)
			return
		}
	}()

	result = make([]*task.Task, 0)

	s = &Subscribe{}
	s.Source = "YYets"
	s.Autodownload = true

	var f func(*html.Node)
	f = func(n *html.Node) {

		if n.Data == "ul" {
			if hasClass(n, "resod_list") {
				season, _ := strconv.Atoi(getAttr(n, "season"))

				if season > 100 { // means it's not normal show episodes, may be trailers.
					season = -season // put it on bottom
				}

				res := make(map[int]*task.Task)
				parseEpisodes(n, season, s.Name, "720p", &res)
				parseEpisodes(n, season, s.Name, "web-dl", &res)
				parseEpisodes(n, season, s.Name, "bd-720p", &res)
				parseEpisodes(n, season, s.Name, "1080p", &res)

				for _, t := range res {
					result = append(result, t)
				}

				return
			}
		}

		if n.Data == "div" {
			if hasClass(n, "res_infobox") {
				s.Banner = getAttr(getTag1(getClass1(n, "f_l_img"), "a"), "href")
				props := getTag(getClass1(getClass1(n, "f_r_info"), "r_d_info"), "li")
				for _, c := range props {
					k := getRigOfTags(getTag1(c, "span"))
					if k == "英文：" {
						log.Print("get name")
						if len(c.Child) > 1 {
							s.Name = s.Name + getRigOfTags(c.Child[1])
							log.Print("get name:", s.Name)
						}
					}
					// if k == "播出：" {
					// 	s.Name = fmt.Sprintf("[%s] %s", getRigOfTags(getTag1(c, "strong")), s.Name)
					// }
				}
				s.Name = strings.TrimLeft(s.Name, ".")
				s.Name = strings.TrimSpace(s.Name)
				s.Name = strings.Replace(s.Name, "/", "|", -1)
				s.Name = strings.Replace(s.Name, "\\", "|", -1)
			}
		}

		for _, c := range n.Child {
			f(c)
		}
	}
	f(doc)

	// fmt.Printf("%v\n", s)
	// encoded, err := downloadBannerImage(s.Banner)
	// if err == nil {
	// 	s.Banner = encoded
	// }
	if len(s.Name) == 0 {
		err = fmt.Errorf("Can't find name of the subscribe.")
	} else {
		err = nil
	}

	return
}
func Parse(url string) (s *Subscribe, result []*task.Task, err error) {
	YYetsLogin(util.ReadConfig("yyets-user"), util.ReadConfig("yyets-password"))

	resp, err := http.Get(url)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()

	s, t, err := parse(resp.Body)
	if len(s.Name) == 0 {
		respdata, _ := httputil.DumpResponse(resp, false)
		log.Print(string(respdata))
	}

	if err != nil {
		return nil, nil, err
	}
	s.URL = url

	return s, t, err
}
func ParseReader(r io.Reader) (s *Subscribe, result []*task.Task, err error) {
	return parse(r)
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

	c := getClass1(getTag1(getClass1(getClass1(n, "lks"), "lks-1"), "a"), "a")
	t.Name = getRigOfTags(c)
	t.StartTime = time.Now().Unix()

	c = getClass1(getClass1(n, "pks"), "download")
	if ed2k := getChildAttr1(c, "type", "ed2k"); ed2k != nil {
		t.Original = getAttr(ed2k, "href")
	} else if magnet := getChildAttr1(c, "type", "magnet"); magnet != nil {
		t.Original = getAttr(magnet, "href")
	} else if thunder := getChildAttr1(c, "thunderhref", "*"); thunder != nil {
		t.Original = getAttr(thunder, "thunderhref")
	} else if a := getChildAttr1(c, "href", "*"); a != nil {
		t.Original = getAttr(a, "href")
	} else {
		t.Original = ""
	}

	return t
}

func YYetsLogin(name string, password string) {
	if YYetsIfLogined() {
		log.Println("logined")
		return
	} else {
		log.Println("not logined")
	}
	//http://www.yyets.com/user/login/ajaxLogin
	data := url.Values{
		"type":     []string{"nickname"},
		"account":  []string{name},
		"password": []string{password},
		"remember": []string{"1"},
	}

	resp, err := http.PostForm("http://www.yyets.com/user/login/ajaxLogin", data)
	if err != nil {
		log.Print(err)
	}
	defer resp.Body.Close()

	respData, err := httputil.DumpResponse(resp, true)
	println(string(respData))
}

func YYetsIfLogined() bool {
	u, _ := url.Parse("http://www.yyets.com")
	cookies := http.DefaultClient.Jar.Cookies(u)
	for _, c := range cookies {
		fmt.Printf("%#v\n", c)
		if c.Name == "GINFO" && strings.HasPrefix(c.Value, "uid") {
			return true
		}
	}

	return false
}
