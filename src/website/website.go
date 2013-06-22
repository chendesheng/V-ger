package website

import (
	"code.google.com/p/cookiejar"
	"download"
	"encoding/json"
	"fmt"
	"native"
	// "regexp"
	// "net/http/httputil"
	"path"
	// "html/template"
	"io/ioutil"
	"os/exec"
	"strconv"
	"strings"
	// "regexp"
	// "io"
	"runtime"
	// "strings"
	// "encoding/json"
	"b1"
	"errors"
	"log"
	"mime"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"subtitles"
	"thunder"
	"time"
)

func init() {
	config = readConfig()
	if logPath, ok := config["log"]; ok {
		f, err := os.OpenFile(logPath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			log.Fatal(err)
		}
		log.SetOutput(f)
	}

	client := &http.Client{
		Jar: cookiejar.NewJar(true),
	}
	cookie := http.Cookie{
		Name:    "gdriveid",
		Value:   config["gdriveid"],
		Domain:  "xunlei.com",
		Expires: time.Now().AddDate(100, 0, 0),
	}
	cookies := []*http.Cookie{&cookie}
	url, _ := url.Parse("http://vip.lixian.xunlei.com")
	client.Jar.SetCookies(url, cookies)

	download.DownloadClient = client
	thunder.Client = client
	subtitles.Client = client
	b1.Client = client

	download.BaseDir = config["dir"]

	runtime.GOMAXPROCS(runtime.NumCPU())

	native.WebSiteAddress = config["server"]
}

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
	cmd := exec.Command("open", path.Join(download.BaseDir, name))
	cmd.Start()

	w.Write([]byte(``))
}

func trashHandler(w http.ResponseWriter, r *http.Request) {
	name, _ := url.QueryUnescape(r.URL.String()[7:])
	fmt.Printf("trash \"%s\".\n", name)

	native.MoveFileToTrash(download.BaseDir, name)
	time.Sleep(time.Second)
	native.MoveFileToTrash(path.Join(download.BaseDir, "vger-tasks"), fmt.Sprint(name, ".vger-task.txt"))
	w.Write([]byte(``))
}

func resumeHandler(w http.ResponseWriter, r *http.Request) {
	name, _ := url.QueryUnescape(r.URL.String()[8:])
	fmt.Printf("resume download \"%s\".\n", name)

	w.Write([]byte(download.TryResumeDownload(name)))
}
func newTaskHandler(w http.ResponseWriter, r *http.Request) {
	var name string
	if len(r.URL.String()) > 4 {
		name, _ = url.QueryUnescape(r.URL.String()[5:])
	} else {
		name = ""
	}
	fmt.Println("newTaskhandler", name)
	input, _ := ioutil.ReadAll(r.Body)
	url := string(input)
	fmt.Printf("add download \"%s\".\n", url)

	w.Write([]byte(download.NewDownload(url, name)))
}
func thunderNewHandler(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if re := recover(); re != nil {
			err := re.(error)

			w.Write([]byte(err.Error()))
		}
	}()

	input, _ := ioutil.ReadAll(r.Body)
	url := string(input)

	thunder.Login(config["thunder-user"], config["thunder-password"])
	files, err := thunder.NewTask(url)
	if err == nil {
		text, _ := json.Marshal(files)
		w.Write([]byte(text))
	} else {
		w.Write([]byte(err.Error()))
	}
}
func thunderTorrentHandler(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if re := recover(); re != nil {
			err := re.(error)

			w.Write([]byte(err.Error()))
		}
	}()
	// res, _ := httputil.DumpRequest(r, true)
	// fmt.Println(string(res))
	fmt.Println("thunder torrent handler")
	f, _, err := r.FormFile("torrent")
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	input, _ := ioutil.ReadAll(f)

	thunder.Login(config["thunder-user"], config["thunder-password"])

	files, err := thunder.NewTaskWithTorrent(input)
	if err == nil {
		text, _ := json.Marshal(files)
		w.Write([]byte(text))
	} else {
		w.Write([]byte(err.Error()))
	}
}
func stopHandler(w http.ResponseWriter, r *http.Request) {
	name, _ := url.QueryUnescape(r.URL.String()[6:])
	fmt.Printf("stop download \"%s\".\n", name)

	w.Write([]byte(download.StopDownload(name)))
}
func limitHandler(w http.ResponseWriter, r *http.Request) {
	name, _ := url.QueryUnescape(r.URL.String()[7:])
	input, _ := ioutil.ReadAll(r.Body)
	speed, _ := strconv.Atoi(string(input))
	fmt.Printf("download \"%s\" limit speed %dKB/s.\n", name, speed)

	w.Write([]byte(download.LimitSpeed(name, speed)))
}
func setAutoShutdownHandler(w http.ResponseWriter, r *http.Request) {
	name, _ := url.QueryUnescape(r.URL.String()[14:])
	input, _ := ioutil.ReadAll(r.Body)
	autoshutdown := string(input)
	fmt.Printf("Autoshutdown task \"%s\" %s.", name, autoshutdown)

	download.SetAutoshutdown(name, autoshutdown == "on")
}

// func queueHandler(w http.ResponseWriter, r *http.Request) {
// 	name, _ := url.QueryUnescape(r.URL.String()[7:])

// 	if err := download.QueueDownload(name); err != nil {
// 		w.Write([]byte(err.Error()))
// 	}
// }
func progressHandler(w http.ResponseWriter, r *http.Request) {
	tasks := download.GetTasks()
	text, _ := json.Marshal(tasks)
	w.Write([]byte(text))
}

type command struct {
	ack    chan bool
	result chan string

	name string
	arg  string
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
func subtitlesSearchHandler(w http.ResponseWriter, r *http.Request) {
	print("subtitlesSearchHandler")
	movieName, _ := url.QueryUnescape(r.URL.String()[strings.LastIndex(r.URL.String(), "/")+1:])
	data := getSubList2(filterMovieName2(movieName))
	text, _ := json.Marshal(data)
	w.Write([]byte(text))
}

func subtitlesDownloadHandler(w http.ResponseWriter, r *http.Request) {
	movieName, _ := url.QueryUnescape(r.URL.String()[strings.LastIndex(r.URL.String(), "/")+1:])
	input, _ := ioutil.ReadAll(r.Body)
	url := string(input)
	// name := getFileName(url)
	url, name, _ := download.GetDownloadInfo(url)
	if ok, err := subtitles.QuickDownload(url, path.Join(download.BaseDir, name)); !ok {
		w.Write([]byte(err.Error()))
		return
	} else {
		extractSubtitle(name, movieName)
	}
}

func appStatusHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(fmt.Sprintf("# of goruntine: %d.", runtime.NumGoroutine())))
}
func appShutdownHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("bye"))
	go func() {
		time.Sleep(time.Second)
		os.Exit(1)
	}()
}
func appGCHandler(w http.ResponseWriter, r *http.Request) {
	runtime.GC()
}

func playHandler(w http.ResponseWriter, r *http.Request) {
	name, _ := url.QueryUnescape(r.URL.String()[6:])
	fmt.Printf("open \"%s\".\n", name)
	cmd := exec.Command("open", config["video-player"], "--args", "http://"+config["server"]+"/video/"+name)
	cmd.Start()
}

func writeError(w http.ResponseWriter, err error) {
	w.Write([]byte(err.Error()))
}
func videoHandler(w http.ResponseWriter, r *http.Request) {
	name, _ := url.QueryUnescape(r.URL.String()[7:])
	t, err := download.GetTask(name)
	if err != nil {
		writeError(w, err)
	}
	if t.Status == "Downloading" {
		download.StopDownload(name)
	}

	url, size := t.URL, t.Size

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

	// first request only need download by single request
	// if ra.start == 0 && sendSize == size {
	// 	err := download.SingleRoutineDownload(url, w, ra.start, ra.start+sendSize)
	// 	if err != nil {
	// 		w.Write([]byte(err.Error()))
	// 	}

	// 	return
	// }
	download.Play(url, w, ra.start, ra.start+sendSize)
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

func Run() {
	download.StartHandleCommands()

	http.Handle("/favicon.ico", http.NotFoundHandler())

	http.HandleFunc("/assets/", assetsHandler)

	http.HandleFunc("/", viewHandler)
	http.HandleFunc("/open/", openHandler)
	http.HandleFunc("/play/", playHandler)
	http.HandleFunc("/video/", videoHandler)
	http.HandleFunc("/resume/", resumeHandler)
	http.HandleFunc("/stop/", stopHandler)
	http.HandleFunc("/progress", progressHandler)
	http.HandleFunc("/new/", newTaskHandler)
	http.HandleFunc("/new", newTaskHandler)
	http.HandleFunc("/limit/", limitHandler)
	http.HandleFunc("/trash/", trashHandler)
	http.HandleFunc("/autoshutdown/", setAutoShutdownHandler)
	// http.HandleFunc("/queue/", queueHandler)

	http.HandleFunc("/thunder/new", thunderNewHandler)
	http.HandleFunc("/thunder/torrent", thunderTorrentHandler)

	http.HandleFunc("/subtitles/search/", subtitlesSearchHandler)
	http.HandleFunc("/subtitles/download/", subtitlesDownloadHandler)

	http.HandleFunc("/app/status", appStatusHandler)
	http.HandleFunc("/app/shutdown", appShutdownHandler)
	http.HandleFunc("/app/gc", appGCHandler)

	//resume downloading tasks
	tasks := download.GetTasks()
	hasDownloading := false
	for _, t := range tasks {
		if t.Status == "Downloading" {
			hasDownloading = true
			download.ResumeDownload(t.Name)
		}
	}
	if !hasDownloading {
		download.ResumeNextQueuedTask()
	}

	server := config["server"]

	log.Print("server ", server, " started.")
	err := http.ListenAndServe(server, nil)
	if err != nil {
		log.Fatal(err)
	}
}
