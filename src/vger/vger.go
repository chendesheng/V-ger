package vger

import (
	"code.google.com/p/cookiejar"
	"download"
	"fmt"
	// "regexp"
	"runtime"
	"strings"
	// "io"
	// "encoding/json"
	"log"
	"net/http"
	"net/url"
	"os"
	"shooter"
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
	shooter.Client = client

	thunder.Login(config["thunder-user"], config["thunder-password"])

	download.BaseDir = config["dir"]

	runtime.GOMAXPROCS(runtime.NumCPU())
}

func pick(arr []string, emptyMessage string) int {
	if len(arr) == 0 {
		if emptyMessage != "" {
			fmt.Println(emptyMessage)
		}
		return -1
	}

	for i, item := range arr {
		fmt.Printf("[%d] %s\n", i+1, item)
	}

	i := 0
	_, err := fmt.Scanf("%d", &i)
	if err != nil {
		log.Fatal(err)
	}
	i--
	if i >= 0 && i < len(arr) {
		return i
	}
	fmt.Println("pick wrong number.")
	return -1
}
func checkIfDownload(input string) bool {
	return strings.Contains(input, "://") || strings.HasSuffix(input, ".torrent")
}

func main() {
	if len(os.Args) > 1 {
		input := os.Args[1]
		if !checkIfDownload(input) {
			getMovieSub(input)
			return
		}

		tasks := thunder.NewTask(input)

		arr := make([]string, len(tasks))
		for i, s := range tasks {
			arr[i] = s.String()
		}
		i := pick(arr, "")
		if i != -1 {
			selectedTask := tasks[i]
			if selectedTask.Percent < 100 {
				fmt.Println("the task is not ready.")
				return
			}

			getMovieSub(selectedTask.Name)

			download.BeginDownload(selectedTask.DownloadURL, selectedTask.Name)
		}
	} else {
		tasks := download.GetTasks()

		arr := make([]string, len(tasks))
		for i, s := range tasks {
			arr[i] = s.String()
		}
		i := pick(arr, "no unfinished task.")
		if i != -1 {
			selectedTask := tasks[i]
			download.BeginDownload(selectedTask.URL, selectedTask.Name)
		}
	}
}
