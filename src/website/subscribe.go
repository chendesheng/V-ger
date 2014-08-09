package website

import (
	"download"
	"fmt"
	"io/ioutil"
	"log"
	"native"
	"net/http"
	"net/url"
	"runtime/debug"
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

	input, _ := ioutil.ReadAll(r.Body)
	url := string(input)

	log.Print("subscribeNewHandler:", url)
	s, tasks, err := subscribe.Parse(url)
	if err != nil {
		panic(err)
	}

	if s1 := subscribe.GetSubscribe(s.Name); s1 != nil {
		for _, t := range tasks {
			t1, err := task.GetTask(t.Name)
			if err != nil {
				if err == task.ErrNoTask {
					task.SaveTaskIgnoreErr(t)
				} else if exists, err := task.ExistsEpisode(t.Subscribe, t.Season, t.Episode); err == nil && !exists {
					task.SaveTaskIgnoreErr(t)
				} else {
					log.Println("episode exists")
				}
			} else {
				if t1.Subscribe != t.Subscribe || t1.Season != t.Season || t1.Episode != t.Episode {
					task.SaveTaskIgnoreErr(t)
				}
			}
		}

		writeJson(w, s1)
	} else {
		err := subscribe.SaveSubscribe(s)
		if err != nil {
			writeError(w, err)
			return
		}

		for _, t := range tasks {
			task.SaveTaskIgnoreErr(t)
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
				defer resp.Body.Close()
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

func unsubscribeHandler(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if re := recover(); re != nil {
			err := re.(error)

			writeError(w, err)
		}
	}()
	name := ""
	if len(r.URL.String()) > 13 {
		name, _ = url.QueryUnescape(r.URL.String()[13:])
	}

	log.Print("unsubscribe:", name)

	err := subscribe.RemoveSubscribe(name)
	if err != nil {
		writeError(w, err)
	}
}

func checkCache(s *subscribe.Subscribe, cachedlen int) (string, error) {
	resp, err := http.Get(s.URL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return "", fmt.Errorf("response status code: %d", resp.StatusCode)
	}

	bytes, err := ioutil.ReadAll(resp.Body)
	html := string(bytes)
	if err != nil {
		return "", err
	}

	if cachedlen > 0 && len(html) == cachedlen {
		// println(s.Name + " no change")
		return "", nil
	} else {
		return html, nil
	}
}
func updateOne(s *subscribe.Subscribe, cache map[string]int) {
	defer func() {
		if r := recover(); r != nil {
			log.Print("check " + s.Name)
			log.Print(r)
			log.Print(string(debug.Stack()))
		}
	}()

	html, err := checkCache(s, cache[s.Name])
	if err != nil {
		log.Print(err)
		return
	}
	if len(html) == 0 {
		return
	}

	cache[s.Name] = len(html)

	_, tasks, err := subscribe.Parse(s.URL)
	if err != nil {
		log.Print(err)
	} else {
		for _, t := range tasks {
			if exists, err := task.Exists(t.Name); err == nil && !exists {
				if exists, err := task.ExistsEpisode(t.Subscribe, t.Season, t.Episode); err == nil && !exists {
					log.Printf("subscribe new task: %v", t)

					if t.Season < 0 {
						task.SaveTaskIgnoreErr(t)
						continue
					}

					files, err := thunder.NewTask(t.Original, "")
					if err != nil {
						log.Print(err)
					}
					fmt.Printf("%v\n", files)
					if err == nil && len(files) == 1 && files[0].Percent == 100 {
						t.URL = files[0].DownloadURL
						_, _, size, _, err := download.GetDownloadInfo(t.URL, false)
						if err != nil {
							log.Print(err)
						} else {
							t.Size = size
							t.Status = "Stopped"
							task.SaveTaskIgnoreErr(t)
							native.SendNotification("A new episode is ready", t.Name)
							// task.StartNewTask2(t)
						}
					}
				}
			}
		}
	}
}
func UpdateAll(cache map[string]int) {
	subscribes := subscribe.GetSubscribes()
	for _, s := range subscribes {
		updateOne(s, cache)
	}
}

func Monitor() {
	time.Sleep(3 * time.Second)

	cache := make(map[string]int) //cache page length
	for {
		UpdateAll(cache)

		time.Sleep(3 * time.Minute)
	}
}
