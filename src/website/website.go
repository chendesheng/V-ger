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
	_ "net/http/pprof"
	"net/url"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"player/shared"
	"strconv"
	"strings"
	"subscribe"
	// "subtitles"
	"task"
	"thunder"
	"time"
	"util"
)

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
	t, err := task.GetTask(name)
	if err != nil {
		writeError(w, err)
		return
	}

	p := shared.GetPlaying(name)

	if s := subscribe.GetSubscribe(t.Subscribe); s != nil {
		if t.LastPlaying > time.Minute &&
			p != nil && p.Duration > 0 &&
			float64(t.LastPlaying)/float64(p.Duration) > 0.85 &&
			t.LastPlaying < s.Duration {
			subscribe.UpdateDuration(t.Subscribe, t.LastPlaying)
		}
	}

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
				t, _ := task.GetTask(name)
				t.URL = url
				task.SaveTask(t)
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
		writeJson(w, files)
	} else {
		writeError(w, err)
	}
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
		writeJson(w, files)
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
	writeJson(w, configs)
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
	cnt := 50
	tks := make([]*task.Task, 0)
	for _, t := range tasks {
		tks = append(tks, t)
		if len(tks) == cnt {
			err := writeJson(ws, tks)
			if err != nil {
				return
			}

			tks = tks[0:0]
		}
	}
	if len(tks) > 0 {
		err := writeJson(ws, tks)
		if err != nil {
			return
		}
	}

	ch := make(chan *task.Task, 20)
	// log.Println("website watch task change ", ch)
	task.WatchChange(ch)
	defer task.RemoveWatch(ch)

	// ws.SetDeadline(time.Now().Add(21 * time.Second))

	for {
		select {
		case t := <-ch:
			err := writeJson(ws, []*task.Task{t})
			if err != nil {
				return
			}
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
	if _, err := os.Stat(path); os.IsNotExist(err) {
		http.NotFound(w, r)
	} else {
		http.ServeFile(w, r, path)
	}
}

func playHandler(w http.ResponseWriter, r *http.Request) {
	name, _ := url.QueryUnescape(r.URL.String()[6:])
	fmt.Printf("play \"%s\".\n", name)

	playerPath := util.ReadConfig("video-player")
	t, err := task.GetTask(name)
	if err != nil {
		writeError(w, err)
		return
	}
	cmd := exec.Command("open", playerPath, "--args", t.URL)

	// config := util.ReadAllConfigs()
	// playerPath := config["video-player"]
	// util.KillProcess(playerPath)
	// cmd := exec.Command("open", playerPath, "--args", "http://"+config["server"]+"/video/"+name)

	err = cmd.Start()
	if err != nil {
		writeError(w, err)
	}
}

func writeError(w http.ResponseWriter, err error) {
	log.Print(err)
	w.Write([]byte(err.Error()))
}
func writeJson(w io.Writer, obj interface{}) error {
	text, err := json.Marshal(obj)
	if err != nil {
		return err
	} else {
		_, err := w.Write(text)
		return err
	}
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

	http.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "assets/favicon.png")
	})

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

	server := util.ReadConfig("server")

	log.Print("server ", server, " started.")
	err := http.ListenAndServe(server, nil)
	if err != nil {
		log.Fatal(err)
	}
}
