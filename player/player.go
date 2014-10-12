package main

import (
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"path"
	"runtime"
	"sync"
	"time"
	"vger/dbHelper"
	"vger/filelock"
	"vger/logger"
	"vger/player/gui"
	"vger/player/movie"
	"vger/util"
)

type appDelegate struct {
	sync.Mutex
	w *gui.Window
	m *movie.Movie
	t time.Duration
}

func (app *appDelegate) OpenFile(filename string) bool {
	log.Println("open file:", filename)

	if app.w == nil {
		app.w = gui.NewWindow("V'ger", 390, 120) // default window size copy from QuickTime player
	}

	go func() {
		app.Lock()
		defer app.Unlock()

		if app.m != nil {
			app.m.SavePlaying()
			app.m.Close()
		}

		app.m = movie.New()
		gui.SetPlayer(app.m)

		for i := 0; i < 3; i++ {
			err := app.m.Open(app.w, filename)

			if err == nil {
				app.m.PlayAsync()

				gui.SendAddRecentOpenedFile(filename)
				break
			} else {
				app.m.Reset()

				if i >= 2 {
					log.Print(err)
					if len(app.m.Filename) > 0 {
						filename = app.m.Filename
					}
					app.w.SendAlert(fmt.Sprintf("Coundn't open \"%s\".", filename))
					break
				}
			}
		}
	}()

	return len(filename) > 0
}

func (app *appDelegate) WillTerminate() {
	m := app.m
	if m != nil {
		m.SavePlaying()
		app.w.DestoryRender()
	}
}
func (app *appDelegate) OnOpenOpenPanel() {
	go func() {
		if app.m != nil {
			app.t = app.m.Hold()
		}
	}()
}
func (app *appDelegate) OnCloseOpenPanel(filename string) {
	if len(filename) > 0 {
		app.OpenFile(filename)
	} else {
		go func() {
			if app.m != nil {
				app.m.Unhold(app.t)
			}
		}()
	}
}

func (app *appDelegate) OnMenuClick(typ int, tag int) int {
	switch typ {
	case 0:
		go app.m.SetAudioTrack(tag)
	case 1:
		app.m.ToggleSubtitle(tag)
	case 2:
		go app.m.ToggleSearchSubtitle()
	case 3:
		app.m.TogglePlay()
	case 4:
		app.onSeekMenuClick(tag)
	case 5:
		app.m.AddVolume(tag)
	case 6:
		app.onSyncSubtitleClick(tag)
	case 7:
		app.onSyncAudioClick(tag)
	}
	return 0
}

func (app *appDelegate) onSeekMenuClick(typ int) {
	switch typ {
	case 0:
		app.m.SeekBySubtitle(false) //backward
	case 1:
		app.m.SeekBySubtitle(true) //forward
	case 2:
		app.m.SeekOffset(-10 * time.Second)
	case 3:
		app.m.SeekOffset(10 * time.Second)
	}
}

func (app *appDelegate) onSyncSubtitleClick(typ int) {
	switch typ {
	case 0:
		app.m.SyncMainSubtitle(-200 * time.Millisecond)
	case 1:
		app.m.SyncMainSubtitle(200 * time.Millisecond)
	case 2:
		app.m.SyncSecondSubtitle(-200 * time.Millisecond)
	case 3:
		app.m.SyncSecondSubtitle(200 * time.Millisecond)
	}
}

func (app *appDelegate) onSyncAudioClick(tag int) {
	app.m.SyncAudio(time.Duration(tag) * 100 * time.Millisecond)
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU() - 1)
	runtime.LockOSThread()

	logbase := util.ReadConfig("log")
	logger.InitLog("[Player]", path.Join(logbase, "player.log"))

	go func() {
		err := http.ListenAndServe("localhost:8080", nil)
		if err != nil {
			log.Print(err)
		}
	}()

	dbHelper.Init("sqlite3", path.Join(util.ReadConfig("dir"), "vger.db"))

	filelock.DefaultLock, _ = filelock.New("/tmp/vger.db.lock.txt")

	util.SetCookie("gdriveid", util.ReadConfig("gdriveid"), "http://xunlei.com")

	networkTimeout := time.Duration(util.ReadIntConfig("network-timeout")) * time.Second
	transport := http.DefaultTransport.(*http.Transport)
	transport.ResponseHeaderTimeout = networkTimeout
	transport.MaxIdleConnsPerHost = 3

	app := &appDelegate{}
	gui.Run(app)
	return
}
