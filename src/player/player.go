package main

import (
	"dbHelper"
	"filelock"
	. "player/shared"
	"time"

	// "fmt"
	"log"
	. "logger"
	"path"
	. "player/gui"
	. "player/movie"
	"runtime"
	"util"
)

func init() {
	InitLog("VgerPlayer", util.ReadConfig("log"))

	log.Print("log initialized.")

	runtime.GOMAXPROCS(runtime.NumCPU() - 1)
}

type appDelegate struct {
	w *Window
	m *Movie
}

func (app *appDelegate) OpenFile(filename string) bool {
	log.Println("open file:", filename)

	if app.m != nil {
		app.m.Close()
		app.m = nil
	}

	app.m = NewMovie()

	go func() {
		app.m.Open(app.w, filename)
		app.m.PlayAsync()
	}()

	return len(filename) > 0
}

func (app *appDelegate) WillTerminate() {
	app.m.SavePlaying()

	done := make(chan bool)
	go func() {
		if app.m != nil {

			app.m.Close()
		}
		done <- true
	}()
	select {
	case <-done:
		return
	case <-time.After(100 * time.Millisecond):
		log.Print("WillTerminate timeout")
		return
	}
}
func (app *appDelegate) SearchSubtitleMenuItemClick() {
	log.Print("SearchSubtitleMenuItemClick")

	go app.m.SearchDownloadSubtitle()
}
func (app *appDelegate) OnOpenOpenPanel() {
	if app.m != nil {
		app.m.Pause()
	}
}
func (app *appDelegate) OnCloseOpenPanel(filename string) {
	if app.m != nil {
		app.m.Resume()
	}

	if len(filename) > 0 {
		app.OpenFile(filename)
	}
}
func main() {
	dbHelper.Init("sqlite3", path.Join(util.ReadConfig("dir"), "vger.db"))

	filelock.DefaultLock, _ = filelock.New("/tmp/vger.db.lock.txt")

	runtime.LockOSThread()

	util.SetCookie("gdriveid", util.ReadConfig("gdriveid"), "http://xunlei.com")

	app := &appDelegate{}
	Initialize(app)
	app.w = NewWindow("V'ger", 1024, 576)
	PollEvents()
	return
}
