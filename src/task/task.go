package task

import (
	"encoding/base64"
	// "encoding/json"
	// "errors"
	"fmt"
	"util"
	// "io"
	"io/ioutil"
	"log"
	// "native"
	// "net/http"
	"os"
	"path"
	// "strconv"
	// "regexp"
	"strings"
	"time"
)

var TaskDir string

func init() {
	watchers = make([]chan []*Task, 0)
	TaskDir = path.Join(util.ReadConfig("dir"), "vger-tasks")
}

type Task struct {
	URL  string
	Size int64
	Name string //identifier (a little unsafe but more readable than url)
	// seconds from 1970-1-1
	StartTime int64

	DownloadedSize int64
	ElapsedTime    time.Duration
	IsNew          bool

	LimitSpeed int64
	Speed      float64
	Status     string
	NameHash   string
	Est        time.Duration

	Autoshutdown bool
}

func taskInfoFileName(name string) string {
	if !strings.HasSuffix(name, ".vger-task.txt") {
		name = fmt.Sprint(name, ".vger-task.txt")
	}
	return path.Join(TaskDir, name)
}

func hashName(name string) string {
	return strings.TrimRight(base64.URLEncoding.EncodeToString([]byte(name)), "=")
}
func NewTask(name string, url string, size int64) *Task {
	t := new(Task)
	t.URL = url
	t.Name = name
	t.IsNew = true
	t.Size = size
	t.StartTime = time.Now().Unix()
	t.DownloadedSize = 0
	t.ElapsedTime = 0

	t.LimitSpeed = 0
	t.Speed = 0
	t.Status = "Stopped"

	t.NameHash = hashName(t.Name)

	SaveTask(t)

	return t
}
func GetTask(name string) (*Task, error) {
	t := new(Task)
	err := util.ReadJson(taskInfoFileName(name), t)
	if err != nil {
		return nil, err
	}

	return t, nil
}

func GetTasks() []*Task {
	fileInfoes, err := ioutil.ReadDir(TaskDir)
	if err != nil {
		log.Print(err)
		return make([]*Task, 0)
	}

	tasks := make([]*Task, 0, len(fileInfoes))
	for _, f := range fileInfoes {
		name := f.Name()

		if strings.HasPrefix(name, ".") || f.IsDir() || !strings.HasSuffix(name, ".vger-task.txt") { //exculding hidden files
			continue
		}

		if t, err := GetTask(name); err == nil {
			tasks = append(tasks, t)
		}
	}

	return tasks
}

func GetDownloadingTask() (*Task, bool) {
	for _, t := range GetTasks() {
		if t.Status == "Downloading" {
			return t, true
		}
	}

	return nil, false
}

func SaveTask(t *Task) (err error) {
	_, err = ioutil.ReadDir(TaskDir)
	if os.IsNotExist(err) {
		os.Mkdir(TaskDir, 0777)
	}

	err = util.WriteJson(taskInfoFileName(t.Name), t)
	if err == nil {
		writeChangeEvent()
	}

	return
}

//never use, use move to trash instead.
func RemoveTask(name string) error {
	err := os.Remove(taskInfoFileName(name))
	if err != nil {
		fmt.Printf("Remove task [%s] failed: %s\n", name, err)
		return err
	}

	writeChangeEvent()

	return nil
}

var watchers []chan []*Task

func WatchChange(ch chan []*Task) {
	if ch == nil {
		panic("ch cannot be nil")
	}

	for _, w := range watchers {
		if w == ch {
			return
		}
	}

	watchers = append(watchers, ch)
	// chTaskChange = ch
}

func RemoveWatch(ch chan []*Task) {
	for i, w := range watchers {
		if w == ch {
			watchers = append(watchers[:i], watchers[i+1:]...)
		}
	}
}

//call this function after modify task file directly, like trash task.
func UpdateFiles() {
	writeChangeEvent()
}

func writeChangeEvent() {
	tks := GetTasks()
	for _, w := range watchers {
		select {
		case w <- tks:
			break
		case <-time.After(time.Second):
			break
		}
	}
}
