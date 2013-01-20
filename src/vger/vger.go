package main

import (
	"code.google.com/p/cookiejar"
	"download"
	"fmt"
	"runtime"
	"strings"
	// "io"
	"log"
	"net/http"
	"net/url"
	"os"
	"shooter"
	"thunder"
	"time"
)

func init() {
	// f, err := os.OpenFile("vger.log", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0666)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// log.SetOutput(f)

	client := &http.Client{
		Jar: cookiejar.NewJar(true),
	}
	cookie := http.Cookie{
		Name:    "gdriveid",
		Value:   "5120E7CE422D1E3F34D7ED1501A1C86A",
		Domain:  "xunlei.com",
		Expires: time.Now().AddDate(100, 0, 0),
	}
	cookies := []*http.Cookie{&cookie}
	url, _ := url.Parse("http://vip.lixian.xunlei.com")
	client.Jar.SetCookies(url, cookies)

	download.DownloadClient = client

	thunder.Client = client
	shooter.Client = client

	thunder.Login("129697884", "057764593828")

	runtime.GOMAXPROCS(runtime.NumCPU())
}
func pick(arr []interface{}, emptyMessage string) interface{} {
	if len(arr) == 0 {
		if emptyMessage != "" {
			fmt.Println(emptyMessage)
		}
		return nil
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
		return arr[i]
	}
	fmt.Println("pick wrong number.")
	return nil
}
func getMovieSub(movieName string) {
	subs := shooter.SearchSubtitles(movieName)

	arr := make([]interface{}, len(subs))
	for i, s := range subs {
		arr[i] = s
	}
	selected := pick(arr, ":( no subtitle.")
	if selected != nil {
		selectedSub := selected.(shooter.Subtitle)
		url, name := shooter.GetDownloadUrl(selectedSub.URL)
		download.BeginDownload(url, name)
	}

}
func main() {
	var url string

	if len(os.Args) > 1 {
		if os.Args[1] == "s" {
			name := os.Args[2]
			getMovieSub(name)
			return
		}

		url = os.Args[1]
		tasks := thunder.NewTask(url)

		arr := make([]interface{}, len(tasks))
		for i, s := range tasks {
			arr[i] = s
		}
		selected := pick(arr, "")
		if selected != nil {
			selectedTask := selected.(thunder.ThunderTask)
			if selectedTask.Percent < 100 {
				fmt.Println("the task is not ready.")
				return
			}

			fmt.Println("choose a subtitle:")
			getMovieSub(selectedTask.Name[:strings.LastIndex(selectedTask.Name, ".")])

			download.BeginDownload(selectedTask.DownloadURL, selectedTask.Name)
		}
	} else {
		tasks := download.GetTasks()

		arr := make([]interface{}, len(tasks))
		for i, s := range tasks {
			arr[i] = s
		}
		selected := pick(arr, "no unfinished task.")
		if selected != nil {
			selectedTask := selected.(*download.Task)
			download.BeginDownload(selectedTask.URL, selectedTask.Name)
		}

	}
}
