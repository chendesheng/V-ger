package website

import (
	"encoding/json"
	"errors"
	"fmt"
	"html"
	"io"
	"io/ioutil"
	"log"
	"mime"
	"net/http"
	_ "net/http/pprof"
	"os/exec"
	"path"
	"path/filepath"
	"runtime/debug"
	"strconv"
	"strings"
	"time"
	"vger/download"
	"vger/native"
	"vger/player/shared"
	"vger/task"
	"vger/thunder"
	"vger/util"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
} //use default buffer size

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
	http.ServeFile(w, r, "index.html")
}

func openHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]

	fmt.Printf("open \"%s\".\n", name)
	// cmd := exec.Command("./player", fmt.Sprintf("-task=%s", name))
	t, err := task.GetTask(name)
	p := util.ReadConfig("dir")
	if err == nil && t != nil {
		p = path.Join(p, t.Subscribe)
	}
	cmd := exec.Command("open", path.Join(p, name))

	err = cmd.Start()
	if err != nil {
		writeError(w, err)
	}
}

func trashHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]

	log.Printf("trash \"%s\".\n", name)

	err := task.DeleteTask(name)
	if err != nil {
		writeError(w, err)
		return
	} else {
		err = shared.DeleteSubtitle(name)
		if err != nil {
			log.Print(err)
		}

		err = shared.DeletePlaying(name)
		if err != nil {
			log.Print(err)
		}
	}
}

func resumeHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]

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

	vars := mux.Vars(r)
	name := vars["name"]

	input, _ := ioutil.ReadAll(r.Body)

	if url := string(input); url != "" {

		_, name2, size, _, err := download.GetDownloadInfo(url, false)

		if err != nil {
			writeError(w, err)
			return
		}

		if name == "" {
			name = name2
		}

		fmt.Printf("add download \"%s\".\nname: %s\n", url, name)

		if t, err := task.GetTask(name); err == nil {
			if t.Status == "Finished" {
				w.Write([]byte("File has already been downloaded."))
			} else if t.Status != "Downloading" && t.Status != "Stopped" {
				if t.Status == "Deleted" {
					log.Print("deleted task")
					t.DownloadedSize = 0
				}
				t.URL = url
				t.Size = size
				t.Status = "Stopped"
				if err := task.SaveTask(t); err != nil {
					writeError(w, err)
				}
			}
		} else if err := task.NewTask(name, url, size, "Stopped"); err != nil {
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
			log.Print(string(debug.Stack()))

			w.Write([]byte(html.EscapeString(err.Error())))
		}
	}()

	input, _ := ioutil.ReadAll(r.Body)

	m := make(map[string]string)
	err := json.Unmarshal(input, &m)
	if err != nil {
		writeError(w, err)
		return
	}

	url := string(m["url"])
	verifycode := string(m["verifycode"])

	log.Print("thunderNewHandler:", url, verifycode)

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
	vars := mux.Vars(r)
	name := vars["name"]

	fmt.Printf("stop download \"%s\".\n", name)

	if err := task.StopTask(name); err != nil {
		writeError(w, err)
	}

	fmt.Println("stop download finish")
}
func limitHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	input := vars["speed"]
	speed, _ := strconv.Atoi(string(input))
	fmt.Printf("limit speed %dKB/s.\n", speed)

	util.SaveConfig("max-speed", input)

	if err := download.LimitSpeed(speed); err != nil {
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
			err := task.QueueDownloadingTask()
			if err != nil {
				log.Print(err)
			}
		}

		for i := downloadingCnt; i < cnt; i++ {
			err, _ := task.ResumeNextTask()
			if err != nil {
				log.Print(err)
			}
		}

		util.SaveConfig("simultaneous-downloads", string(input))
	} else {
		writeError(w, fmt.Errorf("Simultaneous must greater than zero."))
	}
}
func setAutoShutdownHandler(w http.ResponseWriter, r *http.Request) {
	input, _ := ioutil.ReadAll(r.Body)

	util.SaveConfig("shutdown-after-finish", string(input))
	// fmt.Printf("Autoshutdown task \"%s\" %s.", name, autoshutdown)
	// task.SetAutoshutdown(name, autoshutdown == "on")
}

func progressHandler(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		writeError(w, err)
		return
	}

	tasks := task.GetTasks()
	cnt := 50
	tks := make([]*task.Task, 0)
	for _, t := range tasks {
		tks = append(tks, t)
		if len(tks) == cnt {
			err := ws.WriteJSON(tks) //writeJson(ws, tks)
			if err != nil {
				return
			}

			tks = tks[0:0]
		}
	}
	if len(tks) > 0 {
		err := ws.WriteJSON(tks) //writeJson(ws, tks)
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
			err := ws.WriteJSON([]*task.Task{t}) //writeJson(ws, []*task.Task{t})
			if err != nil {
				return
			}
			break
		case <-time.After(time.Second * 20):
			//close connection every 20 seconds
			//if client is alive, it should reconnect to server
			//prevent socket connection & goroutine leak
			ws.Close()
			return
		}
	}
}

func playHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]

	fmt.Printf("play \"%s\".\n", name)

	// playerPath := util.ReadConfig("video-player")
	t, err := task.GetTask(name)
	if err != nil {
		writeError(w, err)
		return
	}

	fmt.Printf("open %s", fmt.Sprintf("vgerplayer://%s", t.URL))
	cmd := exec.Command("open", fmt.Sprintf("vgerplayer://%s", t.URL)) //playerPath, "--args", t.URL)

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
	log.Print(string(debug.Stack()))

	w.Write([]byte(err.Error()))
}
func writeJson(w io.Writer, obj interface{}) {
	text, err := json.Marshal(obj)
	if err != nil {
		log.Print(err)
	} else {
		_, err := w.Write(text)
		if err != nil {
			log.Print(err)
		}
	}
}
func videoHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]

	t, err := task.GetTask(name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if t.Status == "Downloading" {
		err := task.StopTask(name)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
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

type MyServer struct {
	r *mux.Router
}

func (s MyServer) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	s.r.ServeHTTP(w, req)
}

func Run(isDebug bool) {
	if !isDebug {
		go Monitor()
	}

	err, _ := util.MakeSurePathExists(path.Join(util.ReadConfig("dir"), "subs"))
	if err != nil {
		log.Fatal(err)
	}

	r := mux.NewRouter()

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./index.html")
	})
	r.HandleFunc("/newclient", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./newindex.html")
	})
	r.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./assets/favicon.png")
	})

	r.HandleFunc("/open/{name}", openHandler)

	r.HandleFunc("/play/{name}", playHandler)
	r.HandleFunc("/video/{name}", videoHandler)
	r.HandleFunc("/resume/{name}", resumeHandler)
	r.HandleFunc("/stop/{name}", stopHandler)

	r.HandleFunc("/progress", progressHandler)

	r.HandleFunc("/new/{name}", newTaskHandler)
	r.HandleFunc("/new", newTaskHandler)

	r.HandleFunc("/limit/{speed:[0-9]+}", limitHandler)
	r.HandleFunc("/config", configHandler)
	r.HandleFunc("/config/simultaneous", configSimultaneousHandler)
	r.HandleFunc("/trash/{name}", trashHandler)
	r.HandleFunc("/autoshutdown", setAutoShutdownHandler)
	// http.HandleFunc("/queue/", queueHandler)

	r.HandleFunc("/subscribe/new", subscribeNewHandler)
	r.HandleFunc("/subscribe", subscribeHandler)
	r.HandleFunc("/subscribe/banner/{name}", subscribeBannerHandler)
	r.HandleFunc("/unsubscribe/{name}", unsubscribeHandler)

	r.HandleFunc("/thunder/new", thunderNewHandler)
	r.HandleFunc("/thunder/torrent", thunderTorrentHandler)
	r.HandleFunc("/thunder/verifycode", thunderVerifyCodeHandler)
	r.HandleFunc("/thunder/verifycode/", thunderVerifyCodeHandler)

	r.HandleFunc("/subtitles/search/{movie}", subtitlesSearchHandler)
	r.HandleFunc("/subtitles/download/{movie}", subtitlesDownloadHandler)

	r.HandleFunc("/app/status", appStatusHandler)
	r.HandleFunc("/app/shutdown", appShutdownHandler)
	r.HandleFunc("/app/gc", appGCHandler)

	r.PathPrefix("/assets/").Handler(http.FileServer(http.Dir(".")))

	http.Handle("/", MyServer{r})

	server := util.ReadConfig("server")

	log.Print("server ", server, " started.")
	err = http.ListenAndServe(server, nil)
	if err != nil {
		log.Fatal(err)
	}
}
