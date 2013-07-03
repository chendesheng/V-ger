package download

import (
	// "encoding/base64"
	"errors"
	"fmt"
	"io"
	// "io/ioutil"
	// "log"
	"native"
	"net/http"
	"os"
	// "sort"
	// "os/exec"
	// "path"
	// "sort"
	"strconv"
	"strings"
	"task"
	"time"
)

var DownloadClient *http.Client

var BaseDir string

// var task.taskDirName string

// func init() {
// 	taskDirName = "vger-tasks"
// }

// type Task struct {
// 	URL  string
// 	Size int64
// 	Name string //identifier (a little unsafe but more readable than url)
// 	// seconds from 1970-1-1
// 	StartTime int64

// 	DownloadedSize int64
// 	ElapsedTime    time.Duration
// 	isNew          bool

// 	LimitSpeed int64
// 	Speed      float64
// 	Status     string
// 	NameHash   string
// 	Est        time.Duration

// 	Autoshutdown bool
// }

// type taskSlice []*Task

// func (t taskSlice) Len() int           { return len(t) }
// func (t taskSlice) Less(i, j int) bool { return t[i].StartTime < t[j].StartTime }
// func (t taskSlice) Swap(i, j int)      { t[i], t[j] = t[j], t[i] }

func download(t *task.Task, control chan int, quit chan bool) {
	t.Status = "Downloading"
	task.SaveTask(t)

	if t.DownloadedSize < t.Size {
		progress := doDownload(t.URL, GetFilePath(t.Name), t.DownloadedSize, t.Size, t.LimitSpeed, control, quit)

		handleProgress(progress, t, quit)
	}

	t, _ = task.GetTask(t.Name)
	if t.DownloadedSize >= t.Size {
		// removeTask(t.Name)
		fmt.Printf("\nIt's done!\n\n")

		t.Status = "Finished"
		task.SaveTask(t)

		if t.Autoshutdown {
			go native.Shutdown(t.Name)
		} else {
			go native.SendNotification("V'ger Task Finished", t.Name)
			ResumeNextQueuedTask()
		}
	} else {
		fmt.Println(t.DownloadedSize, " ", t.Size)
		t.Status = "Stopped"
		task.SaveTask(t)
	}

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
	if n := getNumOfDownloadingTasks(); n > 0 {
		t.Status = "Queued"
		task.SaveTask(t)
		return t.Name, nil, nil
	}

	go download(t, control, quit)

	return t.Name, control, quit
}
func ResumeDownloadAsync(name string) (chan int, chan bool, error) {
	for _, t := range task.GetTasks() {
		if name == t.Name {
			t.IsNew = false

			control := make(chan int)
			quit := make(chan bool, 50)
			go func(t *task.Task, control chan int) {
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
func getOrNewTask(url string, name string) *task.Task {
	// fmt.Println("hello")
	url, filename, filesize := GetDownloadInfo(url)
	// fmt.Println("hello")
	// fmt.Println(url)
	if name == "" {
		name = filename
	}

	if name == "" {
		panic("File name must not be empty.")
	}

	if t, err := task.GetTask(name); err == nil {
		return t
	}

	t := new(task.Task)
	t.URL = url
	t.Name = name
	t.IsNew = true
	t.Size = filesize
	t.StartTime = time.Now().Unix()
	t.DownloadedSize = 0
	t.ElapsedTime = 0

	t.LimitSpeed = 0
	t.Speed = 0
	t.Status = "Stopped"

	// t.NameHash = hashName(t.Name)

	task.SaveTask(t)

	return t
}

// func hashName(name string) string {
// 	return strings.TrimRight(base64.URLEncoding.EncodeToString([]byte(name)), "=")
// }

// func task.GetTasks() []*task.Task {
// 	taskDir := path.Join(BaseDir, taskDirName)
// 	fileInfoes, err := ioutil.ReadDir(taskDir)
// 	if os.IsNotExist(err) {
// 		os.Mkdir(taskDir, 0666)
// 	} else if err != nil {
// 		log.Fatal(err)
// 	}

// 	tasks := make([]*task.Task, 0, len(fileInfoes))
// 	for _, f := range fileInfoes {
// 		if f.IsDir() {
// 			continue
// 		}

// 		name := f.Name()
// 		if t, err := getTask(name, taskDir); err == nil {
// 			tasks = append(tasks, t)
// 		}
// 	}

// 	// fmt.Printf("get tasks %v.\n", tasks)
// 	return tasks
// }

// func getTask(name string, taskDir string) (*task.Task, error) {
// 	if !strings.HasSuffix(name, ".vger-task.txt") {
// 		return nil, errors.New("Task file name error.")
// 	}

// 	t := new(task.Task)
// 	err := readJson(path.Join(taskDir, name), t)
// 	if err != nil {
// 		return nil, err
// 	}
// 	// if t.NameHash == "" {
// 	t.NameHash = hashName(t.Name)
// 	// }
// 	return t, nil
// }
func ResumeNextQueuedTask() {
	if t := GetNextQueuedTask(); t != nil {
		fmt.Println("Resume download ", t.Name)
		ResumeDownload(t.Name)
	}
}
func GetNextQueuedTask() *task.Task {
	tasks := task.GetTasks()

	var nextTask *task.Task
	startTime := time.Now().Unix()
	for _, t := range tasks {
		if t.Status == "Queued" && t.StartTime < startTime {
			startTime = t.StartTime
			nextTask = t
		}
	}

	return nextTask
}

// func GetTask(name string) (*task.Task, error) {
// 	name = fmt.Sprint(name, ".vger-task.txt")
// 	taskDir := path.Join(BaseDir, taskDirName)
// 	return getTask(name, taskDir)
// }
func SetAutoshutdown(name string, onOrOff bool) {
	if t, err := task.GetTask(name); err == nil {
		t.Autoshutdown = onOrOff
		task.SaveTask(t)
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
			if control != nil {
				taskControls[name] = taskControl{quit, control}
			}
			cmd.ack <- true
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
				t, err := task.GetTask(name)
				if err != nil {
					cmd.ack <- false
					cmd.result <- err.Error()
					return
				}
				t.Status = "Stopped"
				task.SaveTask(t)
				cmd.ack <- true
			}
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
		}
	}
}

func StartHandleCommands() {
	chanCommand = make(chan *command, 5)
	go handleCommands(chanCommand)
}
func LimitSpeed(name string, speed int) string {
	if t, err := task.GetTask(name); err == nil {
		t.LimitSpeed = int64(speed)
		task.SaveTask(t)

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
func QueueDownload(name string) error {
	t, err := task.GetTask(name)
	if err != nil {
		return err
	}

	t.Status = "Queued"
	task.SaveTask(t)
	return nil
}
func getNumOfDownloadingTasks() int {
	n := 0
	for _, t := range task.GetTasks() {
		if t.Status == "Downloading" {
			n++
		}
	}
	return n
}
func ResumeDownload(name string) string {
	cmd := newCommand("resume", name)
	chanCommand <- cmd
	if ok := <-cmd.ack; !ok {
		return <-cmd.result
	}
	return ""
}
func TryResumeDownload(name string) string {
	t, err := task.GetTask(name)
	if err != nil {
		return err.Error()
	}
	if n := getNumOfDownloadingTasks(); n > 0 {
		t.Status = "Queued"
		task.SaveTask(t)
		return ""
	}

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
