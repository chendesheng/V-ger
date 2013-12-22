package main

import (
	"io/ioutil"
	"path"
	// "strings"
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
)

// var filename = flag.String("file", "", "file name")
// var filename = flag.String("file", "", "file name")
// var taskName = flag.String("task", "The.Walking.Dead.4x01.30.Days.Without.An.Accident.720p.HDTV.x264-IMMERSE.[tvu.org.ru].mkv", "vger-task file name")

// var taskName = flag.String("task", "The.Rainmaker.1997.720p.WEB-DL.DD5.1.H.264-ViGi.mkv", "vger-task file name")

// var taskName = flag.String("task", "the.walking.dead.s04e07.proper.720p.hdtv.x264-killers.mkv", "vger-task file name")
// var taskName = flag.String("task", "Nikita.S04E03.720p.HDTV.X264-DIMENSION.mkv", "vger-task file name")
var taskName = flag.String("task", "LS and TSB_Rip1080_HDR.mkv", "vger-task file name")

// var taskName = flag.String("task", "Google IO 2013 - Advanced Go Concurrency Patterns [720p].mp4", "vger-task file name")

// var taskName = flag.String("task", "The.Mentalist.S06E05.720p.HDTV.X264-DIMENSION.mkv", "vger-task file name")

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU() - 1)
	flag.Parse()

	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	// if logPath := util.ReadConfig("log"); logPath != "" {
	f, err := os.OpenFile("player.log", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		log.Fatal(err)
	}
	log.SetOutput(f)
	// }

	log.Print("log initialized.")
}
func findSubs(base string) []string {
	infoes, err := ioutil.ReadDir(base)
	if err == nil {
		res := make([]string, 0)
		for _, f := range infoes {
			filename := path.Join(base, f.Name())
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
func main() {
	runtime.LockOSThread()

	if taskName == nil {
		return
	}

	t, err := task.GetTask(*taskName)
	if err != nil {
		log.Fatal(err)
	}

	base := util.ReadConfig("dir")

	sub := ""
	if len(t.Subs) > 0 {
		sub = t.Subs[0]
	} else {
		if subs := findSubs(path.Join(base, "subs", t.Name)); len(subs) > 0 {
			sub = subs[0]
		}
	}

	m := movie{}
	log.Print("sub: ", sub)
	m.open(path.Join(base, t.Name), sub, t.LastPlaying)

	go m.decode()

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

	go func() {
		if m.v == nil {
			return
		}

		for {
			c := m.v.c
			m.v.window.ChanShowProgress <- c.CalcPlayProgress(c.GetPercent())

			c.After(time.Second)
		}
	}()

	// go website.Run()

	// m.v.play()
	m.play()

	if m.v != nil {
		m.v.window.Destory()
	}

	return
}
