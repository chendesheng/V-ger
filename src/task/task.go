package task

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	// "io"
	"io/ioutil"
	"log"
	// "native"
	// "net/http"
	"os"
	"path"
	// "strconv"
	"regexp"
	"strings"
	"time"
)

var BaseDir string
var taskDirName string

func init() {
	taskDirName = "vger-tasks"
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

func taskInfoFileName(taskName string) string {
	return path.Join(BaseDir, taskDirName, fmt.Sprint(taskName, ".vger-task.txt"))
}
func SaveTask(t *Task) {
	writeJson(taskInfoFileName(t.Name), *t)
}
func RemoveTask(name string) {
	err := os.Remove(taskInfoFileName(name))
	if err != nil {
		fmt.Printf("Remove task [%s] failed: %s\n", name, err)
	}
}

func CleanName(name string) string {
	return filterMovieName2(name)
}

func filterMovieName2(name string) string {
	name = filterMovieName1(name)
	reg, _ := regexp.Compile("(?i)720p|x[.]264|BluRay|DTS|x264|1080p|H[.]264|AC3|[.]ENG|[.]BD|Rip|H264|HDTV|-IMMERSE|-DIMENSION|xvid|[[]PublicHD[]]|[.]Rus|Chi_Eng|DD5[.]1|HR-HDTV|[.]AAC|[0-9]+x[0-9]+|blu-ray|Remux|dxva|dvdscr")
	name = string(reg.ReplaceAll([]byte(name), []byte("")))
	name = strings.Replace(name, ".", " ", -1)
	name = strings.TrimSpace(name)

	return name
}
func filterMovieName1(name string) string {
	index := strings.LastIndex(name, ".")
	if index > 0 {
		name = name[:index]
	}
	index = strings.LastIndex(name, "-")
	if index > 0 {
		name = name[:index]
	}
	return name
}

func GetTasks() []*Task {
	taskDir := path.Join(BaseDir, taskDirName)
	fileInfoes, err := ioutil.ReadDir(taskDir)
	if os.IsNotExist(err) {
		os.Mkdir(taskDir, 0666)
	} else if err != nil {
		log.Fatal(err)
	}

	tasks := make([]*Task, 0, len(fileInfoes))
	for _, f := range fileInfoes {
		// log.Print(f.Name())
		if strings.HasPrefix(f.Name(), ".") || f.IsDir() { //exculding hidden files
			continue
		}

		name := f.Name()
		if t, err := getTask(name, taskDir); err == nil {
			tasks = append(tasks, t)
		}
	}

	// fmt.Printf("get tasks %v.\n", tasks)
	return tasks
}

func hashName(name string) string {
	return strings.TrimRight(base64.URLEncoding.EncodeToString([]byte(name)), "=")
}
func GetDownloadingTask() (*Task, bool) {
	for _, t := range GetTasks() {
		if t.Status == "Downloading" {
			return t, true
		}
	}

	return nil, false
}
func GetTask(name string) (*Task, error) {
	name = fmt.Sprint(name, ".vger-task.txt")
	taskDir := path.Join(BaseDir, taskDirName)
	return getTask(name, taskDir)
}
func getTask(name string, taskDir string) (*Task, error) {
	if !strings.HasSuffix(name, ".vger-task.txt") {
		return nil, errors.New("Task file name error.")
	}

	t := new(Task)
	err := readJson(path.Join(taskDir, name), t)
	if err != nil {
		return nil, err
	}

	t.NameHash = hashName(t.Name)

	return t, nil
}

func writeJson(path string, object interface{}) {
	data, err := json.Marshal(object)
	if err != nil {
		log.Fatal(err)
	}

	ioutil.WriteFile(path, data, 0666)
}
func readJson(path string, object interface{}) error {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	// log.Println("read json: ")
	// log.Println(path)
	// log.Println(string(data))

	json.Unmarshal(data, &object)
	return nil
}
