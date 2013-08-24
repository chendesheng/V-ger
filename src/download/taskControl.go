package download

import (
	"fmt"
	"log"
	"native"
	"net/http"
	"os"
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

func ensureQuit(quit chan bool) {
	select {
	case <-quit:
		// Since no one write to quit channel,
		// the channel must be closed when pass through receive operation.
		break
	case <-time.After(time.Millisecond):
		close(quit)
	}
}

var DownloadClient *http.Client
var baseDir string = util.ReadConfig("dir")
var taskControls map[string]taskControl = make(map[string]taskControl)

func monitorTask() {
	ch := make(chan *task.Task)
	log.Println("task control watch task: ", ch)
	task.WatchChange(ch)

	for t := range ch {
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

				if t.Autoshutdown {
					go native.Shutdown(t.Name)
				} else {
					go native.SendNotification("V'ger Task Finished", t.Name)
					task.ResumeNextTask()
				}
			}
			if t.LimitSpeed != tc.t.LimitSpeed {
				tc.limitSpeed(t.LimitSpeed)
				tc.t = t
			}
		} else {
			if t.Status == "Downloading" {
				log.Printf("download task %v\n", t)

				control := make(chan int)
				quit := make(chan bool, 50)
				taskControls[t.Name] = taskControl{quit, control, t}

				if t.DownloadedSize == 0 {
					native.SendNotification("V'ger task begin", t.Name)
				}
				go download(t, control, quit)
			}
		}

		if t.Status == "Deleted" {
			dir := util.ReadConfig("dir")
			err := native.MoveFileToTrash(dir, t.Name)
			if err != nil {
				log.Println(err)
			}
			native.MoveFileToTrash(task.TaskDir, fmt.Sprint(t.Name, ".vger-task.txt"))
			native.MoveFileToTrash(dir, fmt.Sprint(t.Name, ".zip"))
			native.MoveFileToTrash(dir, fmt.Sprint(t.Name, ".rar"))
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
			quit := make(chan bool)
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
		log.Println(t.Name, " Finished")

		t.Status = "Finished"
		task.SaveTask(t)
	}
}

func getFilePath(name string) string {
	return fmt.Sprintf("%s%c%s", baseDir, os.PathSeparator, name)
}
