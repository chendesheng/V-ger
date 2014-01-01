package main

import (
	"io/ioutil"
	"path"
	"strings"
	"toutf8"
	// "flag"
	"flag"
	"log"
	"os"
	"runtime"
	"task"
	"time"
	"util"

	// . "player/shared"
	// "website"
	. "logger"
	"player/gui"
)

// var filename = flag.String("file", "", "file name")
// var filename = flag.String("file", "", "file name")
// var taskName = flag.String("task", "The.Walking.Dead.4x01.30.Days.Without.An.Accident.720p.HDTV.x264-IMMERSE.[tvu.org.ru].mkv", "vger-task file name")

// var taskName = flag.String("task", "The.Rainmaker.1997.720p.WEB-DL.DD5.1.H.264-ViGi.mkv", "vger-task file name")

// var taskName = flag.String("task", "the.walking.dead.s04e07.proper.720p.hdtv.x264-killers.mkv", "vger-task file name")
// var taskName = flag.String("task", "Nikita.S04E03.720p.HDTV.X264-DIMENSION.mkv", "vger-task file name")
var taskName = flag.String("task", "", "vger-task file name")

// var taskName = flag.String("task", "LS and TSB_Rip1080_HDR.mkv", "vger-task file name")

// var taskName = flag.String("task", "The.Vampire.Diaries.S05E09.720p.HDTV.X264-DIMENSION.mkv", "vger-task file name")

// var taskName = flag.String("task", "Google IO 2013 - Advanced Go Concurrency Patterns [720p].mp4", "vger-task file name")

// var taskName = flag.String("task", "The.Mentalist.S06E05.720p.HDTV.X264-DIMENSION.mkv", "vger-task file name")
var launchedFromGUI bool

func init() {
	// if logPath := util.ReadConfig("playerlog"); logPath != "" {
	// 	f, err := os.OpenFile(logPath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0666)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	log.SetOutput(f)
	// 	os.Stderr = f
	// }

	// log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	InitLog(util.ReadConfig("playerlog"))

	log.Print("log initialized.")

	runtime.GOMAXPROCS(runtime.NumCPU() - 1)

	for i, s := range os.Args {
		log.Printf("Args[%d]:%s", i, s)
		//Mac OS X assigns a unique process serial number ("PSN") to all apps launched via GUI. Call flag.Prase will crash app if not remove it.
		if strings.HasPrefix(s, "-psn") {
			os.Args[i] = ""
			launchedFromGUI = true
		}
	}

	flag.Parse()
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
				if !util.CheckExt(filename, "srt") {
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

func (a *appDelegate) OpenFile(filename string) bool {
	log.Println("open file:", filename)
	name := path.Base(filename)

	base := util.ReadConfig("dir")
	subs := findSubs(path.Join(base, "subs", name))

	m := movie{}

	t, err := task.GetTask(name)
	lastPlaying := time.Duration(0)
	if err == nil {
		lastPlaying = t.LastPlaying

		go func() {
			ticker := time.Tick(3 * time.Second)
			for _ = range ticker {
				t, err := task.GetTask(name)
				if err != nil {
					log.Print(err)
					return
				}

				t.LastPlaying = m.c.GetSeekTime()

				task.SaveTask(t)
			}
		}()
	}

	// log.Print("sub: ", sub)
	m.open(filename, subs, lastPlaying)

	go m.decode(name)

	go m.v.Play()

	return true
}

func main() {
	runtime.LockOSThread()

	if *taskName != "" {

		t, err := task.GetTask(*taskName)
		if err != nil {
			log.Fatal(err)
		}

		base := util.ReadConfig("dir")

		subs := findSubs(path.Join(base, "subs", t.Name))

		m := movie{}
		// log.Print("sub: ", sub)
		m.open(path.Join(base, t.Name), subs, t.LastPlaying)

		go m.decode(*taskName)

		go func() {
			ticker := time.Tick(3 * time.Second)
			for _ = range ticker {
				t, err := task.GetTask(*taskName)
				if err != nil {
					log.Fatal(err)
				}

				t.LastPlaying = m.c.GetSeekTime()

				task.SaveTask(t)
			}
		}()

		m.play()
		m.w.Destory()
	} else {
		log.Println("open with file")

		app := &appDelegate{}
		gui.Initialize(app)
		gui.PollEvents()
	}

	return
}
