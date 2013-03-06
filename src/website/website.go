package website

import (
	"code.google.com/p/cookiejar"
	"download"
	"encoding/json"
	"fmt"
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
	"log"
	"net/http"
	"net/url"
	"os"
	"subtitles"
	"thunder"
	"time"
)

func init() {
	f, err := os.OpenFile("vger.log", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		log.Fatal(err)
	}
	log.SetOutput(f)

	config = readConfig()

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

func playHandler(w http.ResponseWriter, r *http.Request) {
	name, _ := url.QueryUnescape(r.URL.String()[6:])
	fmt.Printf("play \"%s\".\n", name)
	cmd := exec.Command("open", fmt.Sprintf("%s%c%s", download.BaseDir, os.PathSeparator, name))
	cmd.Start()

	w.Write([]byte(``))
}

func resumeHandler(w http.ResponseWriter, r *http.Request) {
	name, _ := url.QueryUnescape(r.URL.String()[8:])
	fmt.Printf("resume download \"%s\".\n", name)

	w.Write([]byte(download.ResumeDownload(name)))
}
func newTaskHandler(w http.ResponseWriter, r *http.Request) {
	input, _ := ioutil.ReadAll(r.Body)
	url := string(input)

	fmt.Printf("add download \"%s\".\n", url)

	w.Write([]byte(download.NewDownload(url)))
}
func thunderNewHandler(w http.ResponseWriter, r *http.Request) {
	input, _ := ioutil.ReadAll(r.Body)
	url := string(input)

	thunder.Login(config["thunder-user"], config["thunder-password"])
	files := thunder.NewTask(url)

	text, _ := json.Marshal(files)
	w.Write([]byte(text))
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
func deleteHandler(w http.ResponseWriter, r *http.Request) {

}
func progressHandler(w http.ResponseWriter, r *http.Request) {
	tasks := download.GetTasks()
	// download.SortTasksByCreateTime(tasks)
	text, _ := json.Marshal(tasks)
	w.Write([]byte(text))
	// w.Write([]byte(fmt.Sprintf("<h3>Go routine numbers: %d</h3>", runtime.NumGoroutine())))
}

type command struct {
	ack    chan bool
	result chan string

	name string
	arg  string
}

func assetsHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.URL.Path)
	path := r.URL.Path[1:]
	if _, err := os.OpenFile(path, os.O_RDONLY, 0666); os.IsNotExist(err) {
		http.NotFound(w, r)
	} else {
		http.ServeFile(w, r, path)
	}
}
func subtitlesSearchHandler(w http.ResponseWriter, r *http.Request) {
	movieName, _ := url.QueryUnescape(r.URL.String()[strings.LastIndex(r.URL.String(), "/")+1:])
	data := getSubList(movieName, []filter{filterMovieName1, filterMovieName2})
	text, _ := json.Marshal(data)
	w.Write([]byte(text))
}

func subtitlesDownloadHandler(w http.ResponseWriter, r *http.Request) {
	movieName, _ := url.QueryUnescape(r.URL.String()[strings.LastIndex(r.URL.String(), "/")+1:])
	input, _ := ioutil.ReadAll(r.Body)
	url := string(input)
	name := getFileName(url)
	if ok, err := subtitles.QuickDownload(url, path.Join(download.BaseDir, name)); !ok {
		w.Write([]byte(err.Error()))
		return
	} else {
		extractSubtitle(name, movieName)
	}
}
func Run() {
	download.StartHandleCommands()

	http.Handle("/favicon.ico", http.NotFoundHandler())

	http.HandleFunc("/assets/", assetsHandler)

	http.HandleFunc("/", viewHandler)
	http.HandleFunc("/play/", playHandler)
	http.HandleFunc("/resume/", resumeHandler)
	http.HandleFunc("/stop/", stopHandler)
	http.HandleFunc("/progress", progressHandler)
	http.HandleFunc("/new", newTaskHandler)
	http.HandleFunc("/limit/", limitHandler)
	http.HandleFunc("/thunder/new", thunderNewHandler)

	http.HandleFunc("/subtitles/search/", subtitlesSearchHandler)
	http.HandleFunc("/subtitles/download/", subtitlesDownloadHandler)

	server := config["server"]
	fmt.Println("server ", server, " started.")
	err := http.ListenAndServe(server, nil)
	if err != nil {
		log.Fatal(err)
	}
}
