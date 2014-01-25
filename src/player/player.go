package main

import (
	"log"
	. "player/shared"
	// "os"
	"io/ioutil"
	"path"
	"runtime"
	"strings"
	// "task"
	"subtitles"
	"time"
	"toutf8"
	"util"

	// . "player/shared"
	// "website"
	. "logger"
	"player/gui"
	. "player/libav"
)

// var filename string
// var launchedFromGUI bool

func init() {
	InitLog(util.ReadConfig("playerlog"))

	log.Print("log initialized.")

	runtime.GOMAXPROCS(runtime.NumCPU() - 1)

	DbFile = util.ReadConfig("dir") + "/vger.db"

	// for i, s := range os.Args[1:] {
	// 	log.Printf("Args[%d]:%s", i+1, s)
	// 	//Mac OS X assigns a unique process serial number ("PSN") to all apps launched via GUI. Call flag.Prase will crash app if not remove it.
	// 	if strings.HasPrefix(s, "-psn") {
	// 		// os.Args[i] = ""
	// 		launchedFromGUI = true
	// 	} else {
	// 		filename = s
	// 		break
	// 	}
	// }

	// flag.Parse()
}
func findSubs(base string) []string {
	infoes, err := ioutil.ReadDir(base)
	if err == nil {
		res := make([]string, 0)
		for _, f := range infoes {
			filename := strings.ToLower(path.Join(base, f.Name()))
			log.Print(filename)

			if f.IsDir() {
				res = append(res, findSubs(filename)...)
			} else {
				if !util.CheckExt(filename, "srt", "ass") {
					continue
				}

				log.Print("try convert to utf8:", filename)

				utf8Text, err := toutf8.ConverToUTF8(filename)
				if err == nil {
					log.Print("convert to utf8 success")
					ioutil.WriteFile(filename, []byte(utf8Text), 0666)
					res = append(res, filename)
				} else {
					log.Print(err.Error())

					// lower := strings.ToLower(f.Name())
					// if strings.Contains(lower, "chs") || strings.Contains(lower, "gb") {
					// 	log.Print("guess encoding by file name:", lower)
					// 	text, err := toutf8.GB18030ToUTF8(filename)
					// 	if err == nil {
					// 		ioutil.WriteFile(filename, []byte(text), 0666)
					// 	} else {
					// 		log.Println(err.Error())
					// 	}
					// }
				}

				// res = append(res, filename)

			}
		}
		return res
	} else {
		return nil
	}
}

type appDelegate struct {
}

var mv *movie

func (a *appDelegate) OpenFile(filename string) bool {
	log.Println("open file:", filename)
	name := path.Base(filename)

	dir := util.ReadConfig("dir")
	subs := findSubs(path.Join(dir, "subs", name))
	for i, sub := range subs {
		sub = strings.ToLower(sub)
		bytes, err := ioutil.ReadFile(sub)
		if err == nil {
			InsertSubtitle(&Sub{name, path.Base(sub), 0, string(bytes), path.Ext(sub)[1:]})
		}

		subs[i] = path.Base(sub)
	}

	m := movie{}
	mv = &m
	m.p = CreateOrGetPlaying(name)

	// go func() {
	// 	ticker := time.Tick(3 * time.Second)
	// 	for _ = range ticker {
	// 		m.p.LastPos = m.c.GetTime()
	// 		SavePlaying(m.p)
	// 	}
	// }()

	// log.Print("sub: ", sub)
	m.open(filename, subs)

	go m.decode(name)

	go m.v.Play()

	return true
}

func (a *appDelegate) WillTerminate() {
	mv.p.LastPos = mv.c.GetTime() - time.Second
	SavePlaying(mv.p)
}
func (a *appDelegate) SearchSubtitleMenuItemClick() {
	log.Print("SearchSubtitleMenuItemClick")

	res := make(chan subtitles.Subtitle)
	subtitles.SearchSubtitles(mv.p.Movie, res)
	for sub := range res {
		// mv.w.ShowSubList()
		println(sub.URL)
	}
}
func main() {
	runtime.LockOSThread()

	// println(filename)
	// if len(filename) > 0 {
	NetworkInit()

	// 	app := &appDelegate{}
	// 	gui.Initialize(app)

	// 	app.OpenFile(filename)

	// 	gui.PollEvents()
	// } else {
	log.Println("open with file")

	app := &appDelegate{}
	gui.Initialize(app)
	gui.PollEvents()
	// }

	return
}
