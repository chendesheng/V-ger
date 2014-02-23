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
	"os/exec"
	"path"
	"player/gui"
	. "player/shared"
	"runtime"
	"strings"
	"subtitles"
	"task"
	"thunder"
	"time"
	"toutf8"
	"util"
)

func init() {
	InitLog(util.ReadConfig("playerlog"))

	log.Print("log initialized.")

	runtime.GOMAXPROCS(runtime.NumCPU() - 1)
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

				utf8Text, encoding, err := toutf8.ConverToUTF8(filename)
				if err == nil {
					log.Print("convert to utf8 success")
					ioutil.WriteFile(filename, []byte(utf8Text), 0666)

					res = append(res, filename)
					if encoding == "gb18030" || encoding == "utf-8" {
						tmp := res[0]
						res[0] = res[len(res)-1]
						res[len(res)-1] = tmp
					}
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

func downloadSubs(movieName, url string, search string) []string {
	chSubs := make(chan subtitles.Subtitle)
	thunder.Login()
	go subtitles.SearchSubtitles(search, url, chSubs)

	for s := range chSubs {
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

		subFileDir := path.Join(util.ReadConfig("dir"), "subs", movieName)
		util.MakeSurePathExists(subFileDir)
		subFile := path.Join(subFileDir, subname)

		println("subfile:", subFile)
		data, err := subtitles.QuickDownload(url)
		if err != nil {
			log.Print(err)
		} else {
			data = bytes.Replace(data, []byte{'+'}, []byte{' '}, -1)
			data = bytes.Replace(data, []byte{'\\', 'N'}, []byte{'\n'}, -1)

			ioutil.WriteFile(subFile, data, 0666)

			if util.CheckExt(subname, "rar", "zip") {
				log.Print(path.Join(path.Dir(os.Args[0]), "unar"))
				cmd := exec.Command(path.Join(path.Dir(os.Args[0]), "unar"), subFile, "-f", "-o", subFileDir)

				if err := cmd.Run(); err != nil {
					log.Print(err)
				} else {
					os.Remove(subFile)
				}
			}
		}
	}

	dir := util.ReadConfig("dir")
	subs := findSubs(path.Join(dir, "subs", movieName))
	for i, sub := range subs {
		sub = strings.ToLower(sub)
		bytes, err := ioutil.ReadFile(sub)
		if err == nil {
			InsertSubtitle(&Sub{movieName, path.Base(sub), 0, string(bytes), path.Ext(sub)[1:], "", ""})
		}
		subs[i] = path.Base(sub)
	}
	return subs
}

type appDelegate struct {
}

var mv *movie

func (a *appDelegate) OpenFile(filename string) bool {
	log.Println("open file:", filename)
	name := path.Base(filename)

	// dir := util.ReadConfig("dir")
	// subs := findSubs(path.Join(dir, "subs", name))
	// for i, sub := range subs {
	// 	sub = strings.ToLower(sub)
	// 	bytes, err := ioutil.ReadFile(sub)
	// 	if err == nil {
	// 		InsertSubtitle(&Sub{name, path.Base(sub), 0, string(bytes), path.Ext(sub)[1:]})
	// 	}

	// 	subs[i] = path.Base(sub)
	// }
	subs := make([]string, 0)
	local := GetSubtitles(name)
	if len(local) > 0 {
		for _, s := range local {
			subs = append(subs, s.Name)
		}
	}
	log.Printf("%v", subs)

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

	if len(subs) == 0 {
		go func() {
			m.w.SendShowMessage("Downloading subtitles...", false)
			defer m.w.SendHideMessage()
			tk, _ := task.GetTask(name)
			var search = util.CleanMovieName(name)
			if tk != nil && len(tk.Subscribe) != 0 && tk.Season > 0 {
				search = fmt.Sprintf("%s s%2de%2d", tk.Subscribe, tk.Season, tk.Episode)
			}
			subFiles := downloadSubs(name, tk.URL, search)
			if len(subFiles) == 0 {
				m.w.SendShowMessage("No subtitle", true)
			}
			m.setupSubtitles(subFiles)
		}()
	}

	go m.decode(name)

	go m.v.Play()

	return len(filename) > 0
}

func (a *appDelegate) WillTerminate() {
	mv.p.LastPos = mv.c.GetTime() - time.Second
	SavePlaying(mv.p)
}
func (a *appDelegate) SearchSubtitleMenuItemClick() {
	log.Print("SearchSubtitleMenuItemClick")

	res := make(chan subtitles.Subtitle)
	tk, _ := task.GetTask(mv.p.Movie)
	url := ""
	if tk != nil {
		url = tk.URL
	}
	subtitles.SearchSubtitles(mv.p.Movie, url, res)
	for sub := range res {
		// mv.w.ShowSubList()
		println(sub.URL)
	}
}

func main() {
	dbHelper.Init("sqlite3", path.Join(util.ReadConfig("dir"), "vger.db"))

	filelock.DefaultLock, _ = filelock.New("/tmp/vger.db.lock.txt")

	runtime.LockOSThread()

	util.SetCookie("gdriveid", util.ReadConfig("gdriveid"), "http://xunlei.com")

	// NetworkInit()
	app := &appDelegate{}
	gui.Initialize(app)
	gui.PollEvents()
	return
}
