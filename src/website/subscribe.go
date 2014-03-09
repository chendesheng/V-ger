package website

import (
	"download"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"subscribe"
	"task"
	"thunder"
	"time"
)

func subscribeNewHandler(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if re := recover(); re != nil {
			err := re.(error)

			writeError(w, err)
		}
	}()
	log.Print("subscribeHandler")

	input, _ := ioutil.ReadAll(r.Body)
	url := string(input)

	println(url)
	s, tasks, err := subscribe.Parse(url)
	if err != nil {
		panic(err)
	}

	if s1 := subscribe.GetSubscribe(s.Name); s1 != nil {
		for _, t := range tasks {
			if exists, err := task.Exists(t.Name); err == nil && !exists {
				if exists, err := task.ExistsEpisode(t.Subscribe, t.Season, t.Episode); err == nil && !exists {
					task.SaveTask(t)
				} else {
					println("episode exists")
				}
			}
		}

		writeJson(w, s1)
	} else {
		subscribe.SaveSubscribe(s)

		for _, t := range tasks {
			task.SaveTask(t)
		}
		writeJson(w, s)
	}
}
func subscribeBannerHandler(w http.ResponseWriter, r *http.Request) {
	name := ""
	if len(r.URL.String()) > 17 {
		name, _ = url.QueryUnescape(r.URL.String()[18:])
	}
	println(name)
	s := subscribe.GetSubscribe(name)
	if s != nil {
		bytes := subscribe.GetBannerImage(name)
		if len(bytes) > 0 {
			h := w.Header()
			h.Add("Cache-Control", "max-age=3153600000") //100 years
			w.Write(bytes)
		} else {
			resp, err := http.Get(s.Banner)
			if err != nil {
				writeError(w, err)
			} else {
				bytes, err = ioutil.ReadAll(resp.Body)
				if err != nil {
					writeError(w, err)
				} else {
					subscribe.SaveBannerImage(name, bytes)

					h := w.Header()
					h.Add("Cache-Control", "max-age=3153600000") //100 years
					w.Write(bytes)
				}
			}
		}
	} else {
		if name == "Downloads" {
			// ioutil.ReadFile("filename")
			http.ServeFile(w, r, "assets/vger.png")
		} else {
			w.WriteHeader(404)
			w.Write([]byte("Unknown subscribe"))
		}
	}
}
func subscribeHandler(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if re := recover(); re != nil {
			err := re.(error)

			writeError(w, err)
		}
	}()
	writeJson(w, subscribe.GetSubscribes())
}

func checkCache(s *subscribe.Subscribe, cachedlen int) (string, error) {
	resp, err := http.Get(s.URL)
	defer resp.Body.Close()

	if err != nil {
		return "", err
	}

	bytes, err := ioutil.ReadAll(resp.Body)
	html := string(bytes)
	if err != nil {
		return "", err
	}

	if cachedlen > 0 && len(html) == cachedlen {
		println(s.Name + " no change")
		return "", nil
	} else {
		return html, nil
	}
}

func UpdateAll(cache map[string]int) {
	subscribes := subscribe.GetSubscribes()
	for _, s := range subscribes {
		println("check " + s.Name)

		html, err := checkCache(s, cache[s.Name])
		if err != nil {
			log.Print(err)
			continue
		}
		if len(html) == 0 {
			continue
		}

		cache[s.Name] = len(html)

		subscribe.ParseReader(strings.NewReader(html))

		_, tasks, err := subscribe.Parse(s.URL)
		if err != nil {
			log.Print(err)
		} else {
			for _, t := range tasks {
				if exists, err := task.Exists(t.Name); err == nil && !exists {
					if exists, err := task.ExistsEpisode(t.Subscribe, t.Season, t.Episode); err == nil && !exists {
						log.Printf("subscribe new task: %v", t)

						if t.Season < 0 {
							task.SaveTask(t)
							continue
						}

						files, err := thunder.NewTask(t.Original, "")
						if err != nil {
							log.Print(err)
						}
						fmt.Printf("%v\n", files)
						if err == nil && len(files) == 1 && files[0].Percent == 100 {
							t.URL = files[0].DownloadURL
							_, _, size, err := download.GetDownloadInfo(t.URL)
							if err != nil {
								log.Print(err)
							} else {
								t.Size = size
								task.SaveTask(t)
								task.StartNewTask2(t)
							}
						}
					}
				}
			}
		}
	}
}

func Monitor() {
	time.Sleep(3 * time.Second)

	cache := make(map[string]int) //cache page length
	for {
		UpdateAll(cache)

		time.Sleep(30 * time.Second)
	}
}
