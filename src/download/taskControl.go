package download

import (
	"fmt"
	"log"
	"native"
	"net/http"
	"os"
	"runtime"
	"task"
	"time"
	"util"
)

type taskControl struct {
	quit     chan bool
	maxSpeed chan int
	t        *task.Task
}

func (tc *taskControl) stopDownload() {
	ensureQuit(tc.quit)
}

func (tc *taskControl) limitSpeed(speed int) error {
	select {
	case tc.maxSpeed <- speed:
		break
	case <-time.After(time.Second * 5):
		return fmt.Errorf("Limit speed operation timeout")
	}

	return nil
}

var DownloadClient *http.Client

var baseDir string
var taskControls map[string]taskControl

func init() {
	baseDir = util.ReadConfig("dir")
	taskControls = make(map[string]taskControl)
}

func monitorTask() {
	ch := make(chan []*task.Task)
	log.Println("task control watch task: ", ch)
	task.WatchChange(ch)

	for tks := range ch {
		for _, t := range tks {
			if tc, ok := taskControls[t.Name]; ok {
				// log.Printf("monitor task %v\n", t)

				if t.Status == "Stopped" {
					tc.stopDownload()
					delete(taskControls, t.Name)
				}
				if t.Status == "Deleted" {
					tc.stopDownload()
					delete(taskControls, t.Name)
				}
				if t.Status == "Finished" {
					delete(taskControls, t.Name)
				}
				if t.LimitSpeed != tc.t.LimitSpeed {
					tc.limitSpeed(t.LimitSpeed)
				}
			} else {
				if t.Status == "Downloading" {
					log.Printf("download task %v\n", t)

					control := make(chan int)
					quit := make(chan bool, 50)
					taskControls[t.Name] = taskControl{quit, control, t}

					go download(t, control, quit)
				}
			}

			if t.Status == "Deleted" {
				dir := util.ReadConfig("dir")
				native.MoveFileToTrash(dir, t.Name)
				native.MoveFileToTrash(task.TaskDir, fmt.Sprint(t.Name, ".vger-task.txt"))
				native.MoveFileToTrash(dir, fmt.Sprint(t.Name, ".zip"))
				native.MoveFileToTrash(dir, fmt.Sprint(t.Name, ".rar"))
			}
		}
	}
}
func Start() {
	go monitorTask()

	time.Sleep(time.Millisecond * 100)

	//resume downloading tasks
	tasks := task.GetTasks()
	hasDownloading := false
	for _, t := range tasks {
		if t.Status == "Downloading" {
			hasDownloading = true

			control := make(chan int)
			quit := make(chan bool, 50)
			taskControls[t.Name] = taskControl{quit, control, t}

			go download(t, control, quit)
		}
	}
	if !hasDownloading {
		task.ResumeNextTask()
	}
}

func download(t *task.Task, control chan int, quit chan bool) {
	if t.DownloadedSize < t.Size {
		f := openOrCreateFileRW(getFilePath(t.Name), t.DownloadedSize)
		defer f.Close()

		progress := doDownload(t.URL, f, t.DownloadedSize, t.Size, int64(t.LimitSpeed), control, quit)

		handleProgress(progress, t, quit)
	}

	t, err := task.GetTask(t.Name)
	if err != nil {
		return
	}

	if t.Status == "Deleted" {
		return
	}

	if t.DownloadedSize >= t.Size {
		fmt.Printf("\nIt's done!\n\n")

		t.Status = "Finished"
		task.SaveTask(t)

		if t.Autoshutdown {
			go native.Shutdown(t.Name)
		} else {
			go native.SendNotification("V'ger Task Finished", t.Name)
			task.ResumeNextTask()
		}
	}
}

func ensureQuit(quit chan bool) {
	log.Println("ensure quit")

	buf := make([]byte, 20000)
	log.Println(string(buf[:runtime.Stack(buf, false)]))

	for i := 0; i < 50; i++ {
		select {
		case quit <- true:
		case <-time.After(time.Second * 1):
			return
		}
	}
}

func getFilePath(name string) string {
	return fmt.Sprintf("%s%c%s", baseDir, os.PathSeparator, name)
}