package download

import (
	"fmt"
	"log"
	"native"
	// "net/http"
	"path"
	"task"
	"time"
	"util"
)

type taskControl struct {
	quit     chan bool
	maxSpeed chan int64
	t        *task.Task
}

func (tc *taskControl) stopDownload() {
	ensureQuit(tc.quit)
}

func (tc *taskControl) limitSpeed(speed int64) error {
	select {
	case tc.maxSpeed <- speed:
		break
	case <-time.After(time.Second * 5):
		return fmt.Errorf("Limit speed operation timeout")
	}

	return nil
}

func ensureQuit(quit chan bool) {
	defer func() {
		recover()
	}()

	select {
	case <-quit:
		// Since no one write to quit channel,
		// the channel must be closed when pass through receive operation.
		break
	default:
		close(quit)
	}
}

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
					native.Shutdown(t.Name)
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
				log.Printf("start download: %v\n", t.Name)

				control := make(chan int64)
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

			control := make(chan int64)
			quit := make(chan bool)
			taskControls[t.Name] = taskControl{quit, control, t}

			go download(t, control, quit)
		}
	}
	if !hasDownloading {
		task.ResumeNextTask()
	}
}

func download(t *task.Task, control chan int64, quit chan bool) {
	if t.DownloadedSize < t.Size {
		f, err := openOrCreateFileRW(path.Join(baseDir, t.Name), t.DownloadedSize)
		if err != nil {
			return
		}

		defer f.Close()

		progress := doDownload(t.URL, f, t.DownloadedSize, t.Size, int64(t.LimitSpeed), control, quit)
		if progress != nil {
			handleProgress(progress, t, quit)
		}
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

		return
	}

	if t.Status == "Downloading" {
		log.Println("restart downloading: ", t.Name)

		t.Status = "Stopped"
		task.SaveTask(t)

		t.Status = "Downloading"
		t.Speed = 0
		task.SaveTask(t)

		return
	}

}
