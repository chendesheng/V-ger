package download

import (
	// "fmt"
	"log"
	"native"
	// "net/http"
	"path"
	"task"
	"time"
	"util"
)

type taskControl struct {
	quit       chan bool
	t          *task.Task
	chMaxSpeed chan int64
}

func (tc *taskControl) stopDownload() {
	ensureQuit(tc.quit)
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

			if t.Status == "Stopped" || t.Status == "Queued" {
				log.Print("queued task ", t.Name)
				tc.stopDownload()
				delete(taskControls, t.Name)
			}
			if t.Status == "Deleted" || t.Status == "New" {
				tc.stopDownload()
				delete(taskControls, t.Name)

				dir := util.ReadConfig("dir")
				err := native.MoveFileToTrash(dir, path.Join(t.Subscribe, t.Name))
				if err != nil {
					log.Println(err)
				}
			}
			if t.Status == "Finished" {
				delete(taskControls, t.Name)

				native.SendNotification("V'ger Task Finished", t.Name)
				if _, ok := task.ResumeNextTask(); !ok {
					if !task.HasDownloadingOrPlaying() {
						if util.ReadBoolConfig("shutdown-after-finish") {
							native.Shutdown(t.Name)
						}
					}
				}
			}

		} else {
			if t.Status == "Downloading" {
				if t.DownloadedSize == 0 {
					native.SendNotification("V'ger task begin", t.Name)
				}

				tc := &taskControl{nil, t, nil}
				taskControls[t.Name] = tc
				go download(tc)
			}
			if t.Status == "Deleted" {
				dir := util.ReadConfig("dir")
				err := native.MoveFileToTrash(dir, path.Join(t.Subscribe, t.Name))
				if err != nil {
					log.Println(err)
				}
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

			tc := &taskControl{nil, t, nil}
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

	dir := path.Join(baseDir, t.Subscribe)
	util.MakeSurePathExists(dir)

	f, err := openOrCreateFileRW(path.Join(dir, t.Name), t.DownloadedSize)
	if err != nil {
		return
	}
	defer f.Close()

	for t.Status == "Downloading" {
		tc.quit = make(chan bool)
		tc.chMaxSpeed = make(chan int64)

		if t.DownloadedSize < t.Size {
			doDownload(t, f, t.DownloadedSize, t.Size, int64(util.ReadIntConfig("max-speed")), tc.chMaxSpeed, tc.quit)
			log.Print("download return")
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
