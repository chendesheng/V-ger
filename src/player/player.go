package main

import (
	"dbHelper"
	"filelock"
	"fmt"
	"log"
	. "logger"
	"path"
	. "player/gui"
	// . "player/libav"
	. "player/movie"
	. "player/shared"
	"runtime"
	"util"
)

func init() {
	InitLog(util.ReadConfig("playerlog"))

	log.Print("log initialized.")

	runtime.GOMAXPROCS(runtime.NumCPU() - 1)
}

func formatSubtitleName(filename string, lang1, lang2 string) string {
	if len(lang2) > 0 {
		return fmt.Sprintf("[%s] %s", lang1, filename)
	} else {
		return fmt.Sprintf("[%s%s] %s", lang1, lang2, filename)
	}
}

type appDelegate struct {
	w *Window
	m *Movie
}

func (app *appDelegate) OpenFile(filename string) bool {
	log.Println("open file:", filename)
	name := path.Base(filename)

	subs := make([]string, 0)
	local := GetSubtitles(name)
	if len(local) > 0 {
		for _, s := range local {
			subs = append(subs, s.Name)
		}
	}
	log.Printf("%v", subs)

	if app.m != nil {
		app.m.Close()
		app.m = nil
	}

	app.m = NewMovie()
	app.m.Open(app.w, filename, subs)
	app.m.PlayAsync()

	return len(filename) > 0
}

func (app *appDelegate) WillTerminate() {
	if app.m != nil {
		app.m.Close()
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

	// NetworkInit()
	app := &appDelegate{}
	Initialize(app)
	app.w = NewWindow("V'ger", 1024, 576)
	PollEvents()
	return
}
