package website

import (
	"code.google.com/p/go.net/websocket"
	"download"
	"encoding/json"
	"errors"
	"fmt"
	"html"
	"io"
	"io/ioutil"
	"log"
	"mime"
	"native"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"subscribe"
	// "subtitles"
	"task"
	"thunder"
	"time"
	"util"
)

func pick(arr []string, emptyMessage string) (int, string) {
	if len(arr) == 0 {
		if emptyMessage != "" {
			fmt.Println(emptyMessage)
		}
		return -1, ""
	}

	for i, item := range arr {
		fmt.Printf("[%d] %s\n", i+1, item)
	}

	next := ""
	i := 0
	fmt.Scanf("%d%s", &i, &next)
	i--
	if i >= 0 && i < len(arr) {
		return i, next
	}
	fmt.Println("pick wrong number.")
	return -1, ""
}
func checkIfSubtitle(input string) bool {
	return !(strings.Contains(input, "://") || strings.HasSuffix(input, ".torrent") || strings.HasPrefix(input, "magnet:"))
}
func checkIfSpeed(input string) (int64, bool) {
	num, err := strconv.ParseUint(input, 10, 64)
	if err != nil {
		return 0, false
	}
	if num > 10*1024*1024 {
		num = 10 * 1024 * 1024
	}
	return int64(num), true
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "main.html")
}

func openHandler(w http.ResponseWriter, r *http.Request) {
	name, _ := url.QueryUnescape(r.URL.String()[6:])
	fmt.Printf("open \"%s\".\n", name)
	// cmd := exec.Command("./player", fmt.Sprintf("-task=%s", name))
	t, err := task.GetTask(name)
	p := util.ReadConfig("dir")
	if err == nil && t != nil {
		p = path.Join(p, t.Subscribe)
	}
	cmd := exec.Command("open", path.Join(p, name))

	cmd.Start()

	w.Write([]byte(``))
}

func trashHandler(w http.ResponseWriter, r *http.Request) {
	name, _ := url.QueryUnescape(r.URL.String()[7:])
	log.Printf("trash \"%s\".\n", name)

	task.DeleteTask(name)
}

func resumeHandler(w http.ResponseWriter, r *http.Request) {
	name, _ := url.QueryUnescape(r.URL.String()[8:])
	fmt.Printf("resume download \"%s\".\n", name)

	if err := task.ResumeTask(name); err != nil {
		writeError(w, err)
	}
}
func newTaskHandler(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if re := recover(); re != nil {
			err := re.(error)

			log.Print(err)

			w.Write([]byte(html.EscapeString(err.Error())))
		}
	}()

	name := ""
	if len(r.URL.String()) > 4 {
		name, _ = url.QueryUnescape(r.URL.String()[5:])
	}

	input, _ := ioutil.ReadAll(r.Body)

	if url := string(input); url != "" {

		_, name2, size, err := download.GetDownloadInfo(url)
		if err != nil {
			_, name2, size, err = download.GetDownloadInfo(url)

			if err != nil {
				writeError(w, err)
			}
			return
		}

		if name == "" {
			name = name2
		}

		fmt.Printf("add download \"%s\".\nname: %s\n", url, name)

		if t, err := task.GetTask(name); err == nil {
			if t.Status == "New" {
				t.URL = url
				t.Size = size
				task.StartNewTask2(t)
			} else if t.Status == "Finished" {
				w.Write([]byte("File has been downloaded."))
			} else {
				log.Print("task already exists")
				task.ResumeTask(name)
			}
		} else if err := task.StartNewTask(name, url, size); err != nil {
			writeError(w, err)
		} else {
			native.SendNotification("V'ger add task", name)
		}
	}
}
func thunderNewHandler(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if re := recover(); re != nil {
			err := re.(error)

			log.Print(err)

			w.Write([]byte(html.EscapeString(err.Error())))
		}
	}()

	input, _ := ioutil.ReadAll(r.Body)

	m := make(map[string]string)
	json.Unmarshal(input, &m)

	url := string(m["url"])
	verifycode := string(m["verifycode"])

	println("thunderNewHandler:", url, verifycode)

	files, err := thunder.NewTask(url, verifycode)
	if err == nil {
		text, _ := json.Marshal(files)
		w.Write([]byte(text))
	} else {
		writeError(w, err)
	}
}
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
			if t1, _ := task.GetTask(t.Name); t1 == nil {
				task.SaveTask(t)
			}
		}

		text, _ := json.Marshal(s1)
		w.Write([]byte(text))
	} else {
		subscribe.SaveSubscribe(s)

		for _, t := range tasks {
			task.SaveTask(t)
		}

		text, _ := json.Marshal(s)
		w.Write([]byte(text))
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

					w.Write(subscribe.GetBannerImage(name))
				}
			}
		}
	} else {
		w.WriteHeader(404)
		w.Write([]byte("Unknown subscribe"))
	}
}
func subscribeHandler(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if re := recover(); re != nil {
			err := re.(error)

			writeError(w, err)
		}
	}()
	text, _ := json.Marshal(subscribe.GetSubscribes())
	w.Write([]byte(text))
}

func thunderVerifyCodeHandler(w http.ResponseWriter, r *http.Request) {
	h := w.Header()
	h.Add("Content-Type", "image/jpeg")
	thunder.WriteValidationCode(w)
}

func thunderTorrentHandler(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if re := recover(); re != nil {
			err := re.(error)

			writeError(w, err)
		}
	}()
	// res, _ := httputil.DumpRequest(r, true)
	// fmt.Println(string(res))
	fmt.Println("thunder torrent handler")
	f, _, err := r.FormFile("torrent")
	if err != nil {
		writeError(w, err)
		return
	}
	input, _ := ioutil.ReadAll(f)

	// thunder.Login(config["thunder-user"], config["thunder-password"])

	files, err := thunder.NewTaskWithTorrent(input)
	if err == nil {
		text, _ := json.Marshal(files)
		w.Write([]byte(text))
	} else {
		writeError(w, err)
	}
}
func stopHandler(w http.ResponseWriter, r *http.Request) {
	name, _ := url.QueryUnescape(r.URL.String()[6:])
	fmt.Printf("stop download \"%s\".\n", name)

	if err := task.StopTask(name); err != nil {
		writeError(w, err)
	}

	fmt.Println("stop download finish")
}
func limitHandler(w http.ResponseWriter, r *http.Request) {
	input, _ := url.QueryUnescape(r.URL.String()[7:])
	speed, _ := strconv.Atoi(string(input))
	fmt.Printf("limit speed %dKB/s.\n", speed)

	util.SaveConfig("max-speed", input)

	if err := download.LimitSpeed(int64(speed)); err != nil {
		writeError(w, err)
	}
}
func configHandler(w http.ResponseWriter, r *http.Request) {
	configs := util.ReadAllConfigs()
	text, _ := json.Marshal(configs)
	w.Write([]byte(text))
}
func configSimultaneousHandler(w http.ResponseWriter, r *http.Request) {
	input, _ := ioutil.ReadAll(r.Body)
	cnt, _ := strconv.Atoi(string(input))
	if cnt > 0 {
		// oldcnt := util.ReadIntConfig("simultaneous-downloads")
		downloadingCnt := task.NumOfDownloadingTasks()

		for i := cnt; i < downloadingCnt; i++ {
			task.QueueDownloadingTask()
		}

		for i := downloadingCnt; i < cnt; i++ {
			task.ResumeNextTask()
		}

		util.SaveConfig("simultaneous-downloads", string(input))
	} else {
		writeError(w, fmt.Errorf("Simultaneous must greater than zero."))
	}
}
func setAutoShutdownHandler(w http.ResponseWriter, r *http.Request) {
	// name, _ := url.QueryUnescape(r.URL.String()[14:])
	input, _ := ioutil.ReadAll(r.Body)

	util.SaveConfig("shutdown-after-finish", string(input))
	// fmt.Printf("Autoshutdown task \"%s\" %s.", name, autoshutdown)
	// task.SetAutoshutdown(name, autoshutdown == "on")
}

func progressHandler(ws *websocket.Conn) {
	tasks := task.GetTasks()
	text, _ := json.Marshal(tasks)

	io.WriteString(ws, string(text))

	ch := make(chan *task.Task)
	// log.Println("website watch task change ", ch)
	task.WatchChange(ch)
	defer task.RemoveWatch(ch)

	for {
		select {
		case <-ch:
			text, _ := json.Marshal(task.GetTasks())
			io.WriteString(ws, string(text))
			break
		case <-time.After(time.Second * 20):
			//close connection every 20 seconds
			//if client is alive, it should reconnect to server
			//prevent socket connection & goroutine leak
			return
		}
	}
}

func assetsHandler(w http.ResponseWriter, r *http.Request) {
	// fmt.Println(r.URL.Path)
	path := r.URL.Path[1:]
	if _, err := os.OpenFile(path, os.O_RDONLY, 0666); os.IsNotExist(err) {
		http.NotFound(w, r)
	} else {
		http.ServeFile(w, r, path)
	}
}

func playHandler(w http.ResponseWriter, r *http.Request) {
	name, _ := url.QueryUnescape(r.URL.String()[6:])
	fmt.Printf("open \"%s\".\n", name)

	config := util.ReadAllConfigs()

	playerPath := config["video-player"]

	util.KillProcess(playerPath)

	cmd := exec.Command("open", playerPath, "--args", "http://"+config["server"]+"/video/"+name)
	cmd.Start()
}

func writeError(w http.ResponseWriter, err error) {
	log.Print(err)
	w.Write([]byte(err.Error()))
}
func videoHandler(w http.ResponseWriter, r *http.Request) {
	name, _ := url.QueryUnescape(r.URL.String()[7:])
	t, err := task.GetTask(name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if t.Status == "Downloading" {
		task.StopTask(name)
	}

	size := t.Size

	code := http.StatusOK

	// If Content-Type isn't set, use the file's extension to find it.
	ctype := w.Header().Get("Content-Type")
	if ctype == "" {
		ctype = mime.TypeByExtension(filepath.Ext(name))
		if ctype != "" {
			w.Header().Set("Content-Type", ctype)
		} else {
			w.Header().Set("Content-Type", "application/octet-stream")
		}
	}

	sendSize := size

	ranges, err := parseRange(r.Header.Get("Range"), size)
	if err != nil {
		http.Error(w, err.Error(), http.StatusRequestedRangeNotSatisfiable)
		return
	}
	ra := ranges[0]
	sendSize = ra.length
	code = http.StatusPartialContent
	w.Header().Set("Content-Range", ra.contentRange(size))
	w.Header().Set("Accept-Ranges", "bytes")
	if w.Header().Get("Content-Encoding") == "" {
		w.Header().Set("Content-Length", strconv.FormatInt(sendSize, 10))
	}
	w.WriteHeader(code)

	download.Play(t, w, ra.start, ra.start+sendSize)
}

type httpRange struct {
	start, length int64
}

func (r httpRange) contentRange(size int64) string {
	return fmt.Sprintf("bytes %d-%d/%d", r.start, r.start+r.length-1, size)
}

// parseRange parses a Range header string as per RFC 2616.
func parseRange(s string, size int64) ([]httpRange, error) {
	if s == "" {
		return nil, nil // header not present
	}
	const b = "bytes="
	if !strings.HasPrefix(s, b) {
		return nil, errors.New("invalid range")
	}
	var ranges []httpRange
	for _, ra := range strings.Split(s[len(b):], ",") {
		ra = strings.TrimSpace(ra)
		if ra == "" {
			continue
		}
		i := strings.Index(ra, "-")
		if i < 0 {
			return nil, errors.New("invalid range")
		}
		start, end := strings.TrimSpace(ra[:i]), strings.TrimSpace(ra[i+1:])
		var r httpRange
		if start == "" {
			// If no start is specified, end specifies the
			// range start relative to the end of the file.
			i, err := strconv.ParseInt(end, 10, 64)
			if err != nil {
				return nil, errors.New("invalid range")
			}
			if i > size {
				i = size
			}
			r.start = size - i
			r.length = size - r.start
		} else {
			i, err := strconv.ParseInt(start, 10, 64)
			if err != nil || i > size || i < 0 {
				return nil, errors.New("invalid range")
			}
			r.start = i
			if end == "" {
				// If no end is specified, range extends to end of the file.
				r.length = size - r.start
			} else {
				i, err := strconv.ParseInt(end, 10, 64)
				if err != nil || r.start > i {
					return nil, errors.New("invalid range")
				}
				if i >= size {
					i = size - 1
				}
				r.length = i - r.start + 1
			}
		}
		ranges = append(ranges, r)
	}
	return ranges, nil
}

func cocoaTestHandler(w http.ResponseWriter, r *http.Request) {
	action, _ := url.QueryUnescape(r.URL.String()[11:])
	log.Printf("test %s", action)
	switch action {
	case "notification":
		native.SendNotification("title", "infoText")
		break
	}
}
func Run() {
	go Monitor()

	http.Handle("/favicon.ico", http.NotFoundHandler())

	http.HandleFunc("/assets/", assetsHandler)

	http.HandleFunc("/", viewHandler)
	http.HandleFunc("/open/", openHandler)
	http.HandleFunc("/play/", playHandler)
	http.HandleFunc("/video/", videoHandler)
	http.HandleFunc("/resume/", resumeHandler)
	http.HandleFunc("/stop/", stopHandler)

	http.Handle("/progress", websocket.Handler(progressHandler))

	http.HandleFunc("/new/", newTaskHandler)
	http.HandleFunc("/limit/", limitHandler)
	http.HandleFunc("/config", configHandler)
	http.HandleFunc("/config/simultaneous", configSimultaneousHandler)
	http.HandleFunc("/trash/", trashHandler)
	http.HandleFunc("/autoshutdown", setAutoShutdownHandler)
	// http.HandleFunc("/queue/", queueHandler)

	http.HandleFunc("/subscribe/new", subscribeNewHandler)
	http.HandleFunc("/subscribe", subscribeHandler)
	http.HandleFunc("/subscribe/banner/", subscribeBannerHandler)

	http.HandleFunc("/thunder/new", thunderNewHandler)
	http.HandleFunc("/thunder/torrent", thunderTorrentHandler)
	http.HandleFunc("/thunder/verifycode", thunderVerifyCodeHandler)
	http.HandleFunc("/thunder/verifycode/", thunderVerifyCodeHandler)

	http.Handle("/subtitles/search/", websocket.Handler(subtitlesSearchHandler))
	http.HandleFunc("/subtitles/download/", subtitlesDownloadHandler)

	http.HandleFunc("/app/status", appStatusHandler)
	http.HandleFunc("/app/shutdown", appShutdownHandler)
	http.HandleFunc("/app/gc", appGCHandler)

	http.HandleFunc("/cocoatest/", cocoaTestHandler)

	server := util.ReadConfig("server")

	log.Print("server ", server, " started.")
	err := http.ListenAndServe(server, nil)
	if err != nil {
		log.Fatal(err)
	}
}

func UpdateAll() {
	subscribes := subscribe.GetSubscribes()
	for _, s := range subscribes {
		_, tasks, err := subscribe.Parse(s.URL)
		if err != nil {
			log.Print(err)
		} else {
			for _, t := range tasks {
				if b, err := task.Exists(t.Name); err == nil && !b {
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

func Monitor() {
	time.Sleep(3 * time.Second)

	for {
		UpdateAll()

		time.Sleep(30 * time.Second)
	}
}
