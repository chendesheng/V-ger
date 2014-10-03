package main

import (
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
	. "vger/player/movie"
	"vger/util"
)

type appDelegate struct {
	sync.Mutex
	w *gui.Window
	m *Movie
	t time.Duration
}

func (app *appDelegate) OpenFile(filename string) bool {
	log.Println("open file:", filename)

	if app.w == nil {
		app.w = gui.NewWindow("V'ger", 640, 360)
	}

	go func() {
		app.Lock()
		defer app.Unlock()

		if app.m != nil {
			app.m.SavePlaying()
			app.m.Close()
		}

		app.m = NewMovie()

		for i := 0; i < 3; i++ {
			err := app.m.Open(app.w, filename)

			if err == nil {
				app.m.PlayAsync()
				break
			} else {
				app.m = nil

				if i < 2 {
					log.Print(err)
				} else {
					log.Fatal(err)
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
func (app *appDelegate) ToggleSearchSubtitle() {
	log.Print("ToggleSearchSubtitle")

	if app.m != nil {
		go app.m.ToggleSearchSubtitle()
	}
}
func (app *appDelegate) OnOpenOpenPanel() {
	if app.m != nil {
		app.t = app.m.Hold()
	}
}
func (app *appDelegate) OnCloseOpenPanel(filename string) {
	if len(filename) > 0 {
		app.OpenFile(filename)
	} else {
		if app.m != nil {
			app.m.Unhold(app.t)
		}
	}
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
