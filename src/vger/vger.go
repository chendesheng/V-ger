package main

import (
	"code.google.com/p/cookiejar"
	"download"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"thunder"
	"time"
)

func init() {
	client := &http.Client{
		Jar: cookiejar.NewJar(true),
	}
	cookie := http.Cookie{
		Name:    "gdriveid",
		Value:   "5120E7CE422D1E3F34D7ED1501A1C86A",
		Domain:  "xunlei.com",
		Expires: time.Now().AddDate(100, 0, 0),
	}
	cookies := [1]*http.Cookie{&cookie}
	url, _ := url.Parse("http://vip.lixian.xunlei.com")
	client.Jar.SetCookies(url, cookies[0:])

	download.DownloadClient = client

	thunder.Client = client

	thunder.Login("129697884", "057764593828")
}

func main() {
	// var url, name string
	var url string

	if len(os.Args) > 1 {
		url = os.Args[1]
		tasks := thunder.NewTask(url)

		var selectedTask thunder.ThunderTask

		for i, t := range tasks {
			fmt.Printf("[%d] %s  %s %d%%\n", i+1, t.Name, t.Size, t.Percent)
		}
		i := 0
		_, err := fmt.Scanf("%d", &i)
		if err != nil {
			log.Fatal(err)
		}
		i--
		if i >= 0 && i < len(tasks) {
			selectedTask = tasks[i]
		} else {
			fmt.Println("pick wrong number.")
			return
		}

		if selectedTask.Percent < 100 {
			fmt.Println("task is not ready.")
			return
		}
		download.BeginDownload(selectedTask.DownloadURL, selectedTask.Name)
		return
	} else {
		tasks := download.GetTasks()
		if len(tasks) == 0 {
			fmt.Println("no unfinished task.")
			return
		}
		for i, t := range tasks {
			fmt.Printf("[%d] %s  %s\n", i+1, t.Name, t.StartDate)
		}
		i := 0
		_, err := fmt.Scanf("%d", &i)
		if err != nil {
			log.Fatal(err)
		}
		i--
		if i >= 0 && i < len(tasks) {
			t := tasks[i]
			download.BeginDownload(t.URL, t.Name)
		} else {
			fmt.Println("pick wrong number.")
			return
		}
	}
}
