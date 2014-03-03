package main

import (
	"bytes"
	"dbHelper"
	"download"
	"filelock"
	"fmt"
	"io/ioutil"
	"log"
	. "logger"
	"os"
	"path"
	"player/gui"
	. "player/shared"
	"runtime"
	// "strings"
	"subtitles"
	"task"
	"thunder"
	"time"
	"toutf8"
	"unicode/utf8"
	"util"
	// "website"
)

func init() {
	InitLog(util.ReadConfig("playerlog"))

	log.Print("log initialized.")

	runtime.GOMAXPROCS(runtime.NumCPU() - 1)
}

func downloadSubs(movieName, url string, search string, quit chan bool) []string {
	chSubs := make(chan subtitles.Subtitle)
	thunder.Login()
	go subtitles.SearchSubtitlesMaxCount(search, url, chSubs, 2, quit)

	subFileDir := path.Join(util.ReadConfig("dir"), "subs", movieName)
	util.MakeSurePathExists(subFileDir)

readSubs:
	for {
		select {
		case s, ok := <-chSubs:
			if !ok {
				break readSubs
			}

			log.Printf("%v", s)
			// text, _ := json.Marshal(s)
			// io.WriteString(ws, string(text))
			url, subname, _, err := download.GetDownloadInfo(s.URL)
			if err != nil {
				return nil
			}

			if subname == "content" {
				subname = s.Description + ".srt" //always use srt
			}

			subFile := path.Join(subFileDir, subname)

			println("subfile:", subFile)
			data, err := subtitles.QuickDownload(url)
			if err != nil {
				log.Print(err)
			} else {
				if util.CheckExt(subname, "rar", "zip") {
					ioutil.WriteFile(subFile, data, 0666)

					unar := path.Join(path.Dir(os.Args[0]), "unar")
					if os.Args[0] == "." {
						unar = "./unar"
					}
					log.Print(unar)
					util.Extract(unar, subFile)
				} else {
					data = bytes.Replace(data, []byte{'+'}, []byte{' '}, -1)

					spaceBytes := make([]byte, 4)
					n := utf8.EncodeRune(spaceBytes, '　')
					spaceBytes = spaceBytes[:n]
					data = bytes.Replace(data, spaceBytes, []byte{' '}, -1)

					data = bytes.Replace(data, []byte{'\\', 'N'}, []byte{'\n'}, -1)

					ioutil.WriteFile(subFile, data, 0666)
				}
			}
		case <-quit:
			return nil
		}
	}

	subs := make([]string, 0)
	util.EmulateFiles(subFileDir, func(filename string) {
		log.Print("try convert to utf8:", filename)

		utf8Text, _, err := toutf8.ConverToUTF8(filename)
		if err == nil {
			log.Print("convert to utf8 success")
			ioutil.WriteFile(filename, []byte(utf8Text), 0666)
			name := path.Base(filename)
			InsertSubtitle(&Sub{movieName, name, 0, utf8Text, path.Ext(filename)[1:], "", ""})
			subs = append(subs, name)
		}
	}, "srt", "ass")

	log.Printf("%v", subs)
	return subs
}

type appDelegate struct {
	w *gui.Window
	m *movie
}

func SearchDownloadSubtitle(m *movie) {
	name := m.p.Movie
	m.w.SendShowMessage("Downloading subtitles...", false)
	defer m.w.SendHideMessage()
	tk, _ := task.GetTask(name)
	var search = util.CleanMovieName(name)
	if tk != nil && len(tk.Subscribe) != 0 && tk.Season > 0 {
		search = fmt.Sprintf("%s s%2de%2d", tk.Subscribe, tk.Season, tk.Episode)
	}
	url := ""
	if tk != nil {
		url = tk.URL
	}
	subFiles := downloadSubs(name, url, search, m.quit)
	select {
	case <-m.quit:
		m.w.SendHideMessage()
		return
	default:
		if len(subFiles) == 0 {
			m.w.SendShowMessage("No subtitle", true)
			return
		}
		m.setupSubtitles(subFiles)
		break
	}
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

	m := movie{}
	app.m = &m

	m.quit = make(chan bool)
	m.p = CreateOrGetPlaying(name)

	m.open(app.w, filename, subs)

	if len(subs) == 0 {
		go SearchDownloadSubtitle(app.m)
	}

	go m.decode(name)
	go m.v.Play()

	return len(filename) > 0
}

func (app *appDelegate) WillTerminate() {
	if app.m != nil {
		app.m.p.LastPos = app.m.c.GetTime() - time.Second
		SavePlaying(app.m.p)
	}
}
func (app *appDelegate) SearchSubtitleMenuItemClick() {
	log.Print("SearchSubtitleMenuItemClick")

	go SearchDownloadSubtitle(app.m)
}
func (app *appDelegate) OnOpenOpenPanel() {
	if app.m != nil {
		app.m.c.Pause()
	}
}
func (app *appDelegate) OnCloseOpenPanel(filename string) {
	if app.m != nil {
		app.m.c.Resume()
	}

	if len(filename) > 0 {
		app.OpenFile(filename)
	}
}
func main() {
	// go website.Run()

	dbHelper.Init("sqlite3", path.Join(util.ReadConfig("dir"), "vger.db"))

	filelock.DefaultLock, _ = filelock.New("/tmp/vger.db.lock.txt")

	runtime.LockOSThread()

	util.SetCookie("gdriveid", util.ReadConfig("gdriveid"), "http://xunlei.com")

	// NetworkInit()
	app := &appDelegate{}
	gui.Initialize(app)
	app.w = gui.NewWindow("V'ger", 1024, 576)
	gui.PollEvents()
	return
}
