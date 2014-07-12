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
	"sort"
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

	// data = bytes.Replace(data, []byte{'\\', 'N'}, []byte{'\n'}, -1)
	// data = bytes.Replace(data, []byte{'\\', 'n'}, []byte{'\n'}, -1)

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

			if util.CheckExt(subname, ".rar", ".zip") {
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
func readSubtitlesFromDir(movieName, dir string) {
	util.WalkFiles(dir, func(filename string) {
		log.Print(filename)

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
			log.Printf("insert subtitle %s", name)
			InsertSubtitle(&Sub{movieName, name, 0, utf8Text, path.Ext(filename)[1:], "", ""})
		} else {
			log.Print(err)
		}
	}, ".srt", ".ass")
}
func downloadSubs(movieName string, url string, search string, quit chan struct{}) {
	chSubs := make(chan subtitles.Subtitle)
	thunder.Login()
	go subtitles.SearchSubtitlesMaxCount(search, url, chSubs, 2, quit)

	dir := path.Join(util.ReadConfig("dir"), "subs", movieName)
	util.MakeSurePathExists(dir)
	if receiveAndExtractSubtitles(chSubs, dir, quit) {
		readSubtitlesFromDir(movieName, dir)
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

	search, url := m.getSubtitleSearch()
	downloadSubs(m.p.Movie, url, search, m.quit)

	name, content := subtitles.Addic7edSubtitle(search)
	if len(name) > 0 && len(content) > 0 {
		sub := &Sub{
			Movie:   m.p.Movie,
			Name:    name,
			Content: content,
		}
		log.Print("insert subtitle:", sub.Name)
		InsertSubtitle(sub)
	}

	select {
	case <-m.quit:
		w.SendHideMessage()
	default:
		subs := GetSubtitlesMap(m.p.Movie)
		if len(subs) == 0 {
			w.SendShowMessage("No subtitle", true)
		} else {
			m.setupSubtitles(subs)
			w.SendHideMessage()
		}
	}
}
func (m *Movie) getSub(name string) *Subtitle {
	for _, s := range m.subs {
		if s.Name == name {
			return s
		}
	}
	return nil
}
func (m *Movie) setupDefaultSubtitles() {
	if len(m.p.Sub1) == 0 && len(m.p.Sub2) > 0 {
		m.p.Sub1 = m.p.Sub2
		m.p.Sub2 = ""
	}

	var s1, s2 *Subtitle
	if len(m.p.Sub1) > 0 {
		s1 = m.getSub(m.p.Sub1)
		if s1 != nil {
			m.s = s1
		} else {
			m.p.Sub1 = ""
		}
	}

	if len(m.p.Sub2) > 0 {
		s2 = m.getSub(m.p.Sub2)
		if s2 != nil {
			m.s2 = s2
		} else {
			m.p.Sub2 = ""
		}
	}

	if m.s == nil && m.s2 == nil {
		m.s, m.s2 = Subtitles(m.subs).Select()
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

	SavePlayingAsync(m.p)
}
func (m *Movie) setupSubtitlesMenu() {
	HideSubtitleMenu()

	tags := make([]int32, 0)
	names := make([]string, 0)

	selected1 := -1
	selected2 := -1

	for i, sub := range m.subs {
		tags = append(tags, int32(i))
		names = append(names, filepath.Base(sub.Name))

		if m.s != nil && sub.Name == m.s.Name {
			selected1 = i
		}

		if m.s2 != nil && sub.Name == m.s2.Name {
			selected2 = i
		}
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
		for _, sub := range m.subs {
			sub.Stop()
		}

		m.subs = nil
		width, height := m.v.Width, m.v.Height
		for _, sub := range subs {
			m.subs = append(m.subs, NewSubtitle(sub, m.w, m.c, float64(width), float64(height)))
		}

		sort.Sort(Subtitles(m.subs))

		m.setupDefaultSubtitles()
		m.setupSubtitlesMenu()
	}
}
