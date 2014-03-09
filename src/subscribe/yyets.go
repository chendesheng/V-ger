package subscribe

import (
	"encoding/base64"
	"io"
	// "fmt"
	"github.com/peterbourgon/html"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"task"
	"time"
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
						if len(c.Child) > 1 {
							s.Name = s.Name + c.Child[1].Data
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

	err = nil
	return
}
func Parse(url string) (s *Subscribe, result []*task.Task, err error) {
	resp, err := http.Get(url)
	defer resp.Body.Close()

	if err != nil {
		return nil, nil, err
	}

	s, t, err := parse(resp.Body)
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

	data, err := ioutil.ReadAll(resp.Body)
	// print(len(data))
	defer resp.Body.Close()

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
