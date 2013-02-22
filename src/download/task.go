package download

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

var DownloadClient *http.Client
var BaseDir string

type Task struct {
	URL       string
	Size      int64
	Name      string //identifier (a little unsafe but more readable than url)
	Path      string
	StartDate string

	DownloadedSize int64
	ElapsedTime    time.Duration
	isNew          bool
}

func (t *Task) String() string {
	return fmt.Sprintf("%s %s %.2f%%", t.Name, t.StartDate, float32(t.DownloadedSize)/float32(t.Size)*100)
}
func GetTasks() []*Task {
	return getTasks()
}
func BeginDownload(url string, name string, maxSpeed int64) {
	if DownloadClient == nil {
		DownloadClient = http.DefaultClient
	}

	t := getOrNewTask(url, name)
	// fmt.Printf("%v", *t)
	// if t.DownloadedSize > 120*1024 {
	// 	t.DownloadedSize -= 100 * 1024
	// }
	progress := doDownload(t.URL, t.Path, t.DownloadedSize, t.Size, maxSpeed)
	// progress := sampleDownload(t.URL, t.Path, t.DownloadedSize, t.Size)
	handleProgress(progress, t)

	removeTask(t.Name)
	fmt.Printf("\nIt's done!\n\n")
}

// func FastDownload(url string, name string) string {
// 	path := GetFilePath(name)
// 	sampleDownload(url, path)
// 	return path
// }
func GetFilePath(name string) string {
	return fmt.Sprintf("%s%c%s", BaseDir, os.PathSeparator, name)
}
func taskInfoFileName(taskName string) string {
	return fmt.Sprintf("tasks%c%s.vger-task.txt", os.PathSeparator, taskName)
}
func saveTask(t *Task) {
	writeJson(taskInfoFileName(t.Name), *t)
}
func removeTask(name string) {
	err := os.Remove(taskInfoFileName(name))
	if err != nil {
		fmt.Printf("Remove task [%s] failed: %s\n", name, err)
	}
}
func getOrNewTask(url string, name string) *Task {
	// fmt.Println("hello")
	url, filename, filesize := getDownloadInfo(url)
	// fmt.Println("hello")
	// fmt.Println(url)
	if name == "" {
		name = filename
	}

	for _, t := range getTasks() {
		if name == t.Name {
			t.isNew = false
			return t
		}
	}

	t := new(Task)
	t.URL = url
	t.Name = name
	t.isNew = true
	t.Size = filesize
	t.StartDate = time.Now().String()
	t.Path = GetFilePath(name)
	t.DownloadedSize = 0
	t.ElapsedTime = 0

	saveTask(t)

	return t
}

func getTasks() []*Task {
	fileInfoes, err := ioutil.ReadDir("tasks")
	if err != nil {
		log.Fatal(err)
	}

	tasks := make([]*Task, 0, len(fileInfoes))
	for _, f := range fileInfoes {
		name := f.Name()
		if f.IsDir() || !strings.HasSuffix(name, ".vger-task.txt") {
			continue
		}

		t := new(Task)
		readJson("tasks/"+name, t)
		tasks = append(tasks, t)
	}

	return tasks
}
