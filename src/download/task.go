package download

import (
	"fmt"
	"io"
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
func BeginDownload(url string, name string, maxSpeed int64) string {
	if DownloadClient == nil {
		DownloadClient = http.DefaultClient
	}

	t := getOrNewTask(url, name)

	progress := doDownload(t.URL, t.Path, t.DownloadedSize, t.Size, maxSpeed)
	// progress := sampleDownload(t.URL, t.Path, t.DownloadedSize, t.Size)
	handleProgress(progress, t)

	removeTask(t.Name)
	fmt.Printf("\nIt's done!\n\n")
	return t.Name
}

func DownloadSmallFile(url string, name string) (filename string, err error) {
	if name == "" {
		url, name, _ = getDownloadInfo(url)
	}

	resp, err := DownloadClient.Get(url)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	f, err := os.OpenFile(fmt.Sprintf("%s%c%s", BaseDir, os.PathSeparator, name),
		os.O_CREATE|os.O_RDWR, 0666)
	defer f.Close()
	if err != nil {
		return "", err
	}

	f.Seek(0, 0)
	for {
		b := make([]byte, 5000)
		n, err := resp.Body.Read(b)
		if n > 0 {
			f.Write(b[:n])
		}

		if err == io.EOF {
			break
		}
	}

	return name, nil
}

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
