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
var taskControls map[string]*taskControl = make(map[string]*taskControl)

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
				if t.DownloadedSize == 0 {
					native.SendNotification("V'ger task begin", t.Name)
				}

				control := make(chan int64)
				tc := &taskControl{nil, control, t}
				taskControls[t.Name] = tc
				go download(tc)
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
			tc := &taskControl{nil, control, t}
			taskControls[t.Name] = tc
			go download(tc)
		}
	}
	if !hasDownloading {
		task.ResumeNextTask()
	}
}

func download(tc *taskControl) {
	t := tc.t
	if t.DownloadedSize >= t.Size {
		return
	}

	f, err := openOrCreateFileRW(path.Join(baseDir, t.Name), t.DownloadedSize)
	if err != nil {
		return
	}
	defer f.Close()

	for t.Status == "Downloading" {
		tc.quit = make(chan bool)

		if t.DownloadedSize < t.Size {

			progress := doDownload(t.URL, f, t.DownloadedSize, t.Size, int64(t.LimitSpeed), tc.maxSpeed, tc.quit)
			if progress != nil {
				handleProgress(progress, t, tc.quit)
			}

		}

		var err error
		t, err = task.GetTask(t.Name)
		tc.t = t

		if err != nil {
			log.Print(err)
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
	}
}
