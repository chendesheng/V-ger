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
	"vger/nativejar"
	"vger/player/gui"
	"vger/player/movie"
	"vger/thunder"
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
		app.w = gui.NewWindow("V'ger", 390, 114) // default window size copy from QuickTime player
		app.w.SetSubFontSize(float64(util.ReadIntConfig("subtitle-font-size")))
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
					app.w.SendAlert(fmt.Sprintf("Couldn't open \"%s\".", filename))
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
}

func (app *appDelegate) OnCloseOpenPanel(filename string) {
	if len(filename) > 0 {
		app.OpenFile(filename)
	}
}

func (app *appDelegate) addVolume(offset int) {
	if volume := app.m.AddVolume(offset); volume >= 0 {
		app.w.SetVolume(volume)
		app.w.SetVolumeVisible(true)
	}
}

func (app *appDelegate) OnMouseWheel(deltaX, deltaY float64) {
	if deltaY != 0 {
		app.addVolume(int(deltaY * -10))
	}
}

func (app *appDelegate) OnWillSleep() {
	app.m.PausePlay()
}

func (app *appDelegate) OnDidWake() {
}

const (
	MENU_AUDIO = iota
	MENU_SUBTITLE
	MENU_SEARCH_SUBTITLE
	MENU_PLAY
	MENU_SEEK
	MENU_VOLUME
	MENU_SYNC_SUBTITLE
	MENU_SYNC_AUDIO
)

const (
	WILL_ENTER_FULL_SCREEN = iota
	DID_ENTER_FULL_SCREEN
	WILL_EXIT_FULL_SCREEN
	DID_EXIT_FULL_SCREEN
)

func (app *appDelegate) OnMenuClick(typ int, tag int) int {
	switch typ {
	case MENU_AUDIO:
		go app.m.SetAudioTrack(tag)
	case MENU_SUBTITLE:
		app.m.ToggleSubtitle(tag)
	case MENU_SEARCH_SUBTITLE:
		go app.m.ToggleSearchSubtitle()
	case MENU_PLAY:
		app.m.TogglePlay()
	case MENU_SEEK:
		app.onSeekMenuClick(tag)
	case MENU_VOLUME:
		app.addVolume(tag)
	case MENU_SYNC_SUBTITLE:
		app.onSyncSubtitleClick(tag)
	case MENU_SYNC_AUDIO:
		app.onSyncAudioClick(tag)
	}
	return 0
}

func (app *appDelegate) ImportSubtitle(filename string) {
	go app.m.ImportSubtitle(filename)
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

	app.w.SetControlsVisible(true, true)
}

func (app *appDelegate) onSyncSubtitleClick(typ int) {
	switch typ {
	case 0:
		if offset, err := app.m.SyncMainSubtitle(-200 * time.Millisecond); err == nil {
			app.w.ShowMessage(fmt.Sprint("Main Subtitle offset ", offset.String()), true)
		}
	case 1:
		if offset, err := app.m.SyncMainSubtitle(200 * time.Millisecond); err == nil {
			app.w.ShowMessage(fmt.Sprint("Main Subtitle offset ", offset.String()), true)
		}
	case 2:
		if offset, err := app.m.SyncSecondSubtitle(-200 * time.Millisecond); err == nil {
			app.w.ShowMessage(fmt.Sprint("Second Subtitle offset ", offset.String()), true)
		}
	case 3:
		if offset, err := app.m.SyncSecondSubtitle(200 * time.Millisecond); err == nil {
			app.w.ShowMessage(fmt.Sprint("Second Subtitle offset ", offset.String()), true)
		}
	}
}

func (app *appDelegate) onSyncAudioClick(tag int) {
	offset := app.m.SyncAudio(time.Duration(tag) * 100 * time.Millisecond)
	app.w.ShowMessage(fmt.Sprint("Audio offset ", offset.String()), true)
}

func (app *appDelegate) OnFullScreen(action int) {
	app.w.SetFullScreen(action)

	//switch action {
	//case WILL_ENTER_FULL_SCREEN, WILL_EXIT_FULL_SCREEN:
	//	app.m.ClockHold()
	//case DID_ENTER_FULL_SCREEN, DID_EXIT_FULL_SCREEN:
	//	app.m.ClockUnhold()
	//}
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

	movie.SearchSubtitleTimeout = time.Duration(util.ReadIntConfig("search-subtitle-timeout")) * time.Second

	thunder.UserName = util.ReadConfig("thunder-user")
	thunder.Password = util.ReadConfig("thunder-password")
	thunder.Gdriveid = util.ReadConfig("gdriveid")

	if http.DefaultClient.Jar == nil {
		http.DefaultClient.Jar, _ = cookiejar.New(nil)
		//http.DefaultClient.Jar, _ = nativejar.New()
	}

	log.Print("gdriveid:", thunder.Gdriveid)
	util.SetCookie("gdriveid", thunder.Gdriveid, "http://xunlei.com")

	networkTimeout := time.Duration(util.ReadIntConfig("network-timeout")) * time.Second
	transport := http.DefaultTransport.(*http.Transport)
	transport.ResponseHeaderTimeout = networkTimeout
	transport.MaxIdleConnsPerHost = 3

	app := &appDelegate{}
	gui.Run(app)
	return
}
