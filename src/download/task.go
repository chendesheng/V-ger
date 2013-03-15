package download

import (
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"native"
	"net/http"
	"os"
	// "os/exec"
	"path"
	// "sort"
	"strconv"
	"strings"
	"time"
)

var DownloadClient *http.Client
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
	isNew          bool

	LimitSpeed int64
	Speed      float64
	Status     string
	NameHash   string
	Est        time.Duration

	Autoshutdown bool
}

// func BeginDownload(url string, name string, maxSpeed int64) string {
// 	if DownloadClient == nil {
// 		DownloadClient = http.DefaultClient
// 	}

// 	t := getOrNewTask(url, name)

// 	control := make(chan int)
// 	progress := doDownload(t, t.URL, GetFilePath(t.Name), t.DownloadedSize, t.Size, maxSpeed, control, quit)
// 	handleProgress(progress, t)

// 	removeTask(t.Name)
// 	fmt.Printf("\nIt's done!\n\n")
// 	return t.Name
// }
func download(t *Task, control chan int, quit chan bool) {
	t.Status = "Downloading"
	saveTask(t)

	if t.DownloadedSize < t.Size {
		progress := doDownload(t.URL, GetFilePath(t.Name), t.DownloadedSize, t.Size, t.LimitSpeed, control, quit)

		handleProgress(progress, t, quit)
	}

	t, _ = GetTask(t.Name)
	if t.DownloadedSize >= t.Size {
		// removeTask(t.Name)
		fmt.Printf("\nIt's done!\n\n")

		t.Status = "Finished"

		if t.Autoshutdown {
			go native.Shutdown(t.Name)
		} else {
			go native.SendNotification("V'ger Task Finished", t.Name)
		}

	} else {
		fmt.Println(t.DownloadedSize, " ", t.Size)
		t.Status = "Stopped"
	}
	saveTask(t)

	go func() {
		timeout := time.After(time.Second * 1)
		for i := 0; i < 50; i++ {
			select {
			case quit <- true:
			case <-timeout:
				return
			}
		}
	}()
}
func DownloadAsync(url string, name string) (string, chan int, chan bool) {
	control := make(chan int)
	quit := make(chan bool, 50)
	t := getOrNewTask(url, name)

	go download(t, control, quit)

	return t.Name, control, quit
}
func ResumeDownloadAsync(name string) (chan int, chan bool, error) {
	for _, t := range GetTasks() {
		if name == t.Name {
			t.isNew = false

			control := make(chan int)
			quit := make(chan bool, 50)
			go func(t *Task, control chan int) {
				download(t, control, quit)
			}(t, control)

			return control, quit, nil
		}
	}

	return nil, nil, errors.New("task not exists")
}

func DownloadSmallFile(url string, name string) (filename string, err error) {
	if name == "" {
		url, name, _ = GetDownloadInfo(url)
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
	return path.Join(BaseDir, taskDirName, fmt.Sprint(taskName, ".vger-task.txt"))
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
	url, filename, filesize := GetDownloadInfo(url)
	// fmt.Println("hello")
	// fmt.Println(url)
	if name == "" {
		name = filename
	}

	if t, err := GetTask(name); err == nil {
		return t
	}

	t := new(Task)
	t.URL = url
	t.Name = name
	t.isNew = true
	t.Size = filesize
	t.StartTime = time.Now().Unix()
	t.DownloadedSize = 0
	t.ElapsedTime = 0

	t.LimitSpeed = 0
	t.Speed = 0
	t.Status = "Stopped"

	t.NameHash = hashName(t.Name)

	saveTask(t)

	return t
}

func hashName(name string) string {
	return strings.TrimRight(base64.URLEncoding.EncodeToString([]byte(name)), "=")
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
		if f.IsDir() {
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

func getTask(name string, taskDir string) (*Task, error) {
	if !strings.HasSuffix(name, ".vger-task.txt") {
		return nil, errors.New("Task file name error.")
	}

	t := new(Task)
	err := readJson(path.Join(taskDir, name), t)
	if err != nil {
		return nil, err
	}
	// if t.NameHash == "" {
	t.NameHash = hashName(t.Name)
	// }
	return t, nil
}
func GetTask(name string) (*Task, error) {
	name = fmt.Sprint(name, ".vger-task.txt")
	taskDir := path.Join(BaseDir, taskDirName)
	return getTask(name, taskDir)
}
func SetAutoshutdown(name string, onOrOff bool) {
	if t, err := GetTask(name); err == nil {
		t.Autoshutdown = onOrOff
		saveTask(t)
	}
}

type command struct {
	ack    chan bool
	result chan string

	name string
	arg  string
}

var chanCommand chan *command

func newCommand(name, arg string) *command {
	return &command{make(chan bool), make(chan string), name, arg}
}

type taskControl struct {
	quit     chan bool
	maxSpeed chan int
}

func handleCommands(chanCommand chan *command) {
	taskControls := make(map[string]taskControl)

	for cmd := range chanCommand {
		switch cmd.name {
		case "new":
			args := strings.Split(cmd.arg, "####")
			name, url := args[0], args[1]
			name, control, quit := DownloadAsync(url, name)
			taskControls[name] = taskControl{quit, control}
			cmd.ack <- true
			break
		case "resume":
			name := cmd.arg
			if _, ok := taskControls[name]; !ok {
				control, quit, err := ResumeDownloadAsync(name)
				if err != nil {
					cmd.ack <- false
					cmd.result <- err.Error()
				} else {
					taskControls[name] = taskControl{quit, control}
					cmd.ack <- true
					close(cmd.result)
				}
			} else {
				cmd.ack <- false
				cmd.result <- "task_not_exists"
			}
			break
		case "stop":
			name := cmd.arg
			fmt.Println("handle stopped")

			if control, ok := taskControls[name]; ok {
				delete(taskControls, name)
				for i := 0; i < 50; i++ {
					control.quit <- true
				}
				cmd.ack <- true
				close(cmd.result)
			} else {
				cmd.ack <- false
				cmd.result <- "task_not_exists"
			}
			break
		case "limit":
			args := strings.Split(cmd.arg, ":::")
			name := args[0]
			fmt.Println(name)
			if control, ok := taskControls[name]; ok {
				speed, _ := strconv.Atoi(args[1])
				fmt.Println("up to ", speed)
				control.maxSpeed <- speed
				cmd.ack <- true
				close(cmd.result)
			} else {
				cmd.ack <- false
				cmd.result <- "task_not_exists"
			}
			break
		}
	}
}

func StartHandleCommands() {
	chanCommand = make(chan *command, 5)
	go handleCommands(chanCommand)
}
func LimitSpeed(name string, speed int) string {
	if t, err := GetTask(name); err == nil {
		t.LimitSpeed = int64(speed)
		saveTask(t)

		if t.Status != "Downloading" {
			return ""
		}
	} else {
		return "task has been deleted."
	}

	cmd := newCommand("limit", fmt.Sprintf("%s:::%d", name, speed))
	chanCommand <- cmd
	if ok := <-cmd.ack; !ok {
		return <-cmd.result
	}
	return ""
}
func StopDownload(name string) string {
	cmd := newCommand("stop", name)
	chanCommand <- cmd
	if ok := <-cmd.ack; !ok {
		return <-cmd.result
	}

	return ""
}
func ResumeDownload(name string) string {
	cmd := newCommand("resume", name)
	chanCommand <- cmd
	if ok := <-cmd.ack; !ok {
		return <-cmd.result
	}
	return ""
}
func NewDownload(url string, name string) string {
	cmd := newCommand("new", fmt.Sprint(name, "####", url))
	chanCommand <- cmd
	if ok := <-cmd.ack; !ok {
		return <-cmd.result
	}
	return ""
}
