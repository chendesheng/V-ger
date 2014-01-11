package main

import (
	"io/ioutil"
	"log"
	// "os"
	"path"
	"runtime"
	"strings"
	"task"
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

	dir := util.ReadConfig("dir")
	subs := findSubs(path.Join(dir, "subs", name))

	m := movie{}

	t, err := task.GetTask(name)
	lastPlaying := time.Duration(0)
	if err == nil {
		lastPlaying = t.LastPlaying

		go func() {
			ticker := time.Tick(3 * time.Second)
			for _ = range ticker {
				if m.c == nil {
					continue
				}

				t, err := task.GetTask(name)
				if err != nil {
					log.Print(err)
					return
				}

				t.LastPlaying = m.c.GetTime()

				task.SaveTask(t)

				m.c.WaitUtilRunning()
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
