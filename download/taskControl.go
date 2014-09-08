package download

import (
	"log"
	"native"
	"net/http"
	"path"
	"subscribe"
	"task"
	"time"
	"util"
)

type taskControl struct {
	quit       chan struct{}
	t          *task.Task
	chMaxSpeed chan int
}

func (tc *taskControl) stopDownload() {
	ensureQuit(tc.quit)
}

func ensureQuit(quit chan struct{}) {
	if quit != nil {
		defer func() {
			err := recover()
			if err != nil {
				log.Print(err)
			}
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
}

var BaseDir string
var taskControls map[string]*taskControl = make(map[string]*taskControl)

func monitorTask() {
	ch := make(chan *task.Task, 20)
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

				go trashTask(t)
			}
			if t.Status == "Finished" {
				delete(taskControls, t.Name)

				native.SendNotification("V'ger Task Finished", t.Name)
				if _, ok := task.ResumeNextTask(); !ok {
					if !task.HasDownloadingOrPlaying() {
						if util.ReadBoolConfig("shutdown-after-finish") {
							err := native.Shutdown(t.Name)
							if err != nil {
								log.Print(err)
							}
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
			if t.Status == "Deleted" || t.Status == "New" {
				go trashTask(t)
			}
		}
	}
}
func trashTask(t *task.Task) {
	err := native.MoveFileToTrash(path.Join(BaseDir, t.Subscribe), t.Name)
	if err != nil {
		log.Println(err)
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

			filesize := util.GetFileSize(path.Join(BaseDir, t.Subscribe, t.Name))
			if t.DownloadedSize > filesize {
				t.DownloadedSize = filesize
				if t.DownloadedSize > 1000 {
					t.DownloadedSize -= 1000
				}
				task.SaveTaskIgnoreErr(t)
			}

			go download(tc)
		}
	}
	if !hasDownloading {
		err, _ := task.ResumeNextTask()
		if err != nil {
			log.Print(err)
		}
	}
}

func download(tc *taskControl) {
	t := tc.t
	if t.DownloadedSize >= t.Size {
		if t.Status == "Downloading" {
			t.Status = "Finished"
			task.SaveTaskIgnoreErr(t)
		}
		return
	}

	dir := path.Join(BaseDir, t.Subscribe)
	if _, exists := util.MakeSurePathExists(dir); !exists {
		s := subscribe.GetSubscribe(t.Subscribe)
		if s != nil {
			resp, err := http.Get(s.Banner)
			if err != nil {
				log.Print(err)
			} else {
				defer resp.Body.Close()
				native.DefaultNativeAPI.SetIcon(dir, resp.Body)
			}
		}
	}

	f, err := openOrCreateFileRW(path.Join(dir, t.Name), t.DownloadedSize)
	if err != nil {
		return
	}
	defer f.Close()

	for t.Status == "Downloading" {
		tc.quit = make(chan struct{})
		tc.chMaxSpeed = make(chan int)

		if t.DownloadedSize < t.Size {
			doDownload(t, f, t.DownloadedSize, t.Size,
				util.ReadIntConfig("max-speed"),
				tc.chMaxSpeed,
				util.ReadSecondsConfig("task-restart-timeout"),
				tc.quit)
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
			task.SaveTaskIgnoreErr(t)

			return
		}
	}
}
