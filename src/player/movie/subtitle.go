package movie

import (
	"bytes"
	"download"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	. "player/gui"
	. "player/shared"
	. "player/subtitle"
	"subtitles"
	"task"
	"thunder"
	"time"
	"toutf8"
	"unicode/utf8"
	"util"
)

func extract(subFile string) {
	dir := path.Dir(os.Args[0])
	var unar string
	if dir == "." {
		unar = "./unar"
	} else {
		unar = path.Join(dir, "unar")
	}
	log.Print(path.Dir(os.Args[0]))
	log.Print(unar)
	log.Print(subFile)
	util.Extract(unar, subFile)
}
func saveThunderSubtitle(subFile string, data []byte) {
	data = bytes.Replace(data, []byte{'+'}, []byte{' '}, -1)
	//replace chinese space to ascii space
	spaceBytes := make([]byte, 4)
	n := utf8.EncodeRune(spaceBytes, '　')
	spaceBytes = spaceBytes[:n]
	data = bytes.Replace(data, spaceBytes, []byte{' '}, -1)

	data = bytes.Replace(data, []byte{'\\', 'N'}, []byte{'\n'}, -1)
	data = bytes.Replace(data, []byte{'\\', 'n'}, []byte{'\n'}, -1)

	ioutil.WriteFile(subFile, data, 0666)
}
func receiveAndExtractSubtitles(chSubs chan subtitles.Subtitle, dir string, quit chan struct{}) bool {
	deadline := time.Now().Add(time.Minute)
	for {
		select {
		case s, ok := <-chSubs:
			if !ok {
				return true
			}

			log.Printf("%v", s)
			// text, _ := json.Marshal(s)
			// io.WriteString(ws, string(text))
			url, subname, _, err := download.GetDownloadInfo(s.URL)
			if err != nil {
				log.Print(err)
				break
			}

			if subname == "content" {
				subname = s.Description + ".srt" //always use srt
			}

			subFile := path.Join(dir, subname)
			data, err := subtitles.QuickDownload(url)
			if err != nil {
				log.Print(err)
				break
			}

			if util.CheckExt(subname, "rar", "zip") {
				ioutil.WriteFile(subFile, data, 0666)
				extract(subFile)
			} else {
				saveThunderSubtitle(subFile, data)
			}
		case <-quit:
			return false
		}

		if time.Now().After(deadline) {
			break
		}
	}

	return true
}
func readSubtitlesFromDir(movieName, dir string) []string {
	log.Print(dir)
	subs := make([]string, 0)
	util.EmulateFiles(dir, func(filename string) {
		f, err := os.OpenFile(filename, os.O_RDONLY, 0666)
		if err != nil {
			log.Print(err)
			return
		}

		utf8Text, _, err := toutf8.ConverToUTF8(f)
		if err == nil {
			ioutil.WriteFile(filename, []byte(utf8Text), 0666)
			name := path.Base(filename)

			// lang1, lang2 := cld.DetectLanguage2(utf8Text)
			// log.Printf("subtitle %s language:%s, %s", name, lang1, lang2)
			InsertSubtitle(&Sub{movieName, name, 0, utf8Text, path.Ext(filename)[1:], "", ""})
			subs = append(subs, name)
		} else {
			log.Print(err)
		}
	}, "srt", "ass")
	return subs
}
func downloadSubs(movieName string, url string, search string, quit chan struct{}) []string {
	chSubs := make(chan subtitles.Subtitle)
	thunder.Login()
	go subtitles.SearchSubtitlesMaxCount(search, url, chSubs, 2, quit)

	dir := path.Join(util.ReadConfig("dir"), "subs", movieName)
	util.MakeSurePathExists(dir)
	if receiveAndExtractSubtitles(chSubs, dir, quit) {
		return readSubtitlesFromDir(movieName, dir)
	} else {
		return nil
	}
}

func (m *Movie) getSubtitleSearch() (string, string) {
	name := m.p.Movie
	t, _ := task.GetTask(name)
	var search = util.CleanMovieName(name)
	if t != nil && len(t.Subscribe) != 0 && t.Season > 0 {
		search = fmt.Sprintf("%s s%02de%02d", t.Subscribe, t.Season, t.Episode)
	}
	url := ""
	if t != nil {
		url = t.URL
	}
	return search, url
}

func (m *Movie) SearchDownloadSubtitle() {
	w := m.w
	w.SendShowMessage("Downloading subtitles...", false)
	defer w.SendHideMessage()

	search, url := m.getSubtitleSearch()
	subFiles := downloadSubs(m.p.Movie, url, search, m.quit)
	// if len(subFiles) < 5 {
	name, content := subtitles.Addic7edSubtitle(search)
	if len(name) > 0 && len(content) > 0 {
		sub := &Sub{
			Movie:   m.p.Movie,
			Name:    name,
			Content: content,
		}
		InsertSubtitle(sub)

		subFiles = append(subFiles, name)
	}
	// }

	select {
	case <-m.quit:
		w.SendHideMessage()
		return
	default:
		subs := GetSubtitlesMap(m.p.Movie)
		if len(subs) == 0 {
			w.SendShowMessage("No subtitle", true)
			return
		}
		m.setupSubtitles(subs)
		break
	}
}
func (m *Movie) setupDefaultSubtitles(subs map[string]*Sub, width, height int) {
	var en, cn, double *Subtitle
	for _, sub := range subs {
		s := NewSubtitle(sub, m.w, m.c, float64(width), float64(height))
		if s != nil {
			if en == nil && s.Lang1 == "en" && len(s.Lang2) == 0 {
				en = s
			}
			if cn == nil && s.Lang1 == "zh" && len(s.Lang2) == 0 {
				cn = s
			}

			if double == nil && len(s.Lang1) > 0 && len(s.Lang2) > 0 {
				double = s
			}
		}
	}

	if double != nil {
		m.s = double
	} else {
		if cn != nil {
			m.s = cn
			m.s2 = en
		} else {
			m.s = en
		}
	}

	if m.s != nil {
		m.s.IsMainOrSecondSub = true
		m.p.Sub1 = m.s.Name
		go m.s.Play()
	}

	if m.s2 != nil {
		m.s.IsMainOrSecondSub = false
		m.p.Sub2 = m.s2.Name
		go m.s2.Play()
	}
}
func (m *Movie) setupSubtitlesMenu(subs []*Sub) {
	HideSubtitleMenu()

	tags := make([]int32, 0)
	names := make([]string, 0)

	selected1 := -1
	selected2 := -1
	i := 0
	for _, sub := range subs {
		tags = append(tags, int32(i))
		names = append(names, filepath.Base(sub.Name))

		if m.s != nil && sub.Name == m.s.Name {
			selected1 = i
		}

		if m.s2 != nil && sub.Name == m.s2.Name {
			selected2 = i
		}

		i++
	}

	if selected1 == -1 && selected2 == -1 {
		selected1 = 0
	}

	if len(names) > 0 {
		m.w.InitSubtitleMenu(names, tags, selected1, selected2)
	}
}
func getSubValues(subs map[string]*Sub) []*Sub {
	var values []*Sub
	for _, s := range subs {
		values = append(values, s)
	}

	return values
}

func (m *Movie) setupSubtitles(subs map[string]*Sub) {
	if len(subs) > 0 {
		m.subs = getSubValues(subs)

		println("play subtitle:", subs)
		width, height := m.v.Width, m.v.Height

		if len(m.p.Sub1) == 0 && len(m.p.Sub2) > 0 {
			m.p.Sub1 = m.p.Sub2
			m.p.Sub2 = ""
		}

		var s1, s2 *Subtitle
		if len(m.p.Sub1) > 0 {
			s1 = NewSubtitle(subs[m.p.Sub1], m.w, m.c, float64(width), float64(height))
			if s1 != nil {
				s1.IsMainOrSecondSub = true

				if s1 != nil {
					go s1.Play()
				}
				m.s = s1
			} else {
				m.p.Sub1 = ""
			}
		}

		if len(m.p.Sub2) > 0 {
			s2 = NewSubtitle(subs[m.p.Sub2], m.w, m.c, float64(width), float64(height))
			if s2 != nil {
				s2.IsMainOrSecondSub = false

				if s2 != nil {
					go s2.Play()
				}
				m.s2 = s2
			} else {
				m.p.Sub2 = ""
			}
		}

		if m.s == nil && m.s2 == nil {
			println("auto select default subtitle")
			m.setupDefaultSubtitles(subs, width, height)
		}

		SavePlayingAsync(m.p)
	}

	m.setupSubtitlesMenu(m.subs)
}
