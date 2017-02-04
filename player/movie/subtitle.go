package movie

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"sort"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"
	"vger/download"
	"vger/player/shared"
	"vger/player/subtitle"
	"vger/subtitles"
	"vger/task"
	"vger/toutf8"
	"vger/util"
)

var unarPath string //for unit test
var SearchSubtitleTimeout = 40 * time.Second

func extract(subFile string) {
	dir := path.Dir(os.Args[0])
	var unar string
	if dir == "." {
		unar = "./unar"
	} else {
		if len(unarPath) > 0 {
			unar = unarPath
		} else {
			unar = path.Join(dir, "unar")
		}
	}

	log.Print("Extract file: ", subFile)
	util.Extract(unar, subFile)
}
func saveToDisk(subFile string, data []byte) {
	data = bytes.Replace(data, []byte{'+'}, []byte{' '}, -1)
	//replace chinese space to ascii space
	spaceBytes := make([]byte, 4)
	n := utf8.EncodeRune(spaceBytes, '　')
	spaceBytes = spaceBytes[:n]
	data = bytes.Replace(data, spaceBytes, []byte{' '}, -1)

	// data = bytes.Replace(data, []byte{'\\', 'N'}, []byte{'\n'}, -1)
	// data = bytes.Replace(data, []byte{'\\', 'n'}, []byte{'\n'}, -1)

	err := ioutil.WriteFile(subFile, data, 0666)
	if err != nil {
		log.Print(err)
	}
}
func receiveAndExtractSubtitles(chSubs chan subtitles.Subtitle, dir string, quit chan struct{}) (finded bool) {
	for {
		select {
		case s, ok := <-chSubs:
			if !ok {
				return
			}

			log.Printf("%v", s)
			// text, _ := json.Marshal(s)
			// io.WriteString(ws, string(text))
			_, subname, _, data, err := download.GetDownloadInfoN(s.URL, s.Context, 2, true, quit)
			if err != nil {
				log.Print(err)
				break
			}

			if s.Source == "Kankan" {
				if strings.HasSuffix(s.Description, ".srt") {
					subname = s.Description
				} else {
					subname = s.Description + ".srt" //always use srt
				}

				subname = subname[strings.LastIndex(subname, "\\")+1:]
			}

			subFile := path.Join(dir, subname)

			if util.CheckExt(subname, ".rar", ".zip") {
				err := ioutil.WriteFile(subFile, data, 0666)
				if err != nil {
					log.Print(err)
				}

				extract(subFile)

				finded = true
			} else {
				saveToDisk(subFile, data)

				finded = true
			}
		case <-quit:
			return
		}
	}

	return
}
func readSubtitlesFromDir(movieName, dir string) {
	util.WalkFiles(dir, func(filename string) error {
		f, err := os.OpenFile(filename, os.O_RDONLY, 0666)
		if err != nil {
			log.Print(err)
			return nil //still continue
		}

		utf8Text, _, err := toutf8.ConverToUTF8(f)
		if err == nil {
			err := ioutil.WriteFile(filename, []byte(utf8Text), 0666)
			if err != nil {
				log.Print(err)
				return nil
			}

			name := path.Base(filename)

			// lang1, lang2 := cld.DetectLanguage2(utf8Text)
			dis := getNameDistance(name, movieName)
			log.Printf("insert subtitle %s, dis:%d", name, dis)
			shared.InsertSubtitle(&shared.Sub{movieName, name, 0, utf8Text, path.Ext(filename)[1:], "", "", dis})
		} else {
			log.Print(err)
		}

		return nil
	}, ".srt", ".ass")
}

func downloadSubs(movieName string, url string, search string, quit chan struct{}) {
	chSubs := make(chan subtitles.Subtitle)

	dir := path.Join("/tmp/vgersubs", movieName)
	go subtitles.SearchSubtitles(search, url, chSubs, quit)

	err, _ := util.MakeSurePathExists("/tmp/vgersubs")
	if err != nil {
		log.Print(err)
		return
	}

	err, _ = util.MakeSurePathExists(dir)
	if err != nil {
		log.Print(err)
		return
	}

	if receiveAndExtractSubtitles(chSubs, dir, quit) {
		readSubtitlesFromDir(movieName, dir)
	}
}

func (m *Movie) ImportSubtitle(filename string) {
	f, err := os.OpenFile(filename, os.O_RDONLY, 0666)
	if err != nil {
		log.Print(err)
		return
	}

	utf8Text, _, err := toutf8.ConverToUTF8(f)
	// log.Print(utf8Text)
	movieName := m.p.Movie
	name := path.Base(filename)

	// lang1, lang2 := cld.DetectLanguage2(utf8Text)
	dis := getNameDistance(name, movieName)

	shared.InsertSubtitle(&shared.Sub{movieName, name, 0, utf8Text, path.Ext(filename)[1:], "", "", dis})
	go m.w.SendShowMessage(fmt.Sprintf(`"%s" Imported`, name), true)
	log.Printf("insert subtitle %s, dis:%d", name, dis)

	subs := shared.GetSubtitlesMap(m.p.Movie)
	if len(subs) > 0 {
		m.setupSubtitles(subs, name)
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
func (m *Movie) ToggleSearchSubtitle() {
	if m.subQuit != nil {
		q := m.subQuit
		m.subQuit = nil
		close(q)
	} else {
		m.searchDownloadSubtitle()
	}
}
func (m *Movie) IsSearchingSubtitle() bool {
	return m.subQuit != nil
}

func (m *Movie) searchDownloadSubtitle() {
	quit := make(chan struct{})
	m.subQuit = quit

	w := m.w
	w.SendShowMessage("Downloading subtitles...", false)

	search, url := m.getSubtitleSearch()

	go func() {
		select {
		case <-quit:
			break
		case <-time.After(SearchSubtitleTimeout):
			log.Printf("Search subtitle %s timeout", search)
			q := m.subQuit
			m.subQuit = nil
			close(q)
			break
		}
	}()

	downloadSubs(m.p.Movie, url, search, quit)

	subs := shared.GetSubtitlesMap(m.p.Movie)
	if len(subs) == 0 {
		w.SendShowMessage("No subtitle", true)
	} else {
		m.setupSubtitles(subs, "")
		w.SendHideMessage()
	}
}
func (m *Movie) getSub(name string) *subtitle.Subtitle {
	for _, s := range m.subs {
		if s.Name == name {
			return s
		}
	}
	return nil
}
func (m *Movie) setupDefaultSubtitles(selectedSub string) {
	m.stopPlayingSubs()

	var s1 *subtitle.Subtitle
	var s2 *subtitle.Subtitle

	if len(selectedSub) == 0 {
		s1, s2 = m.getSub(m.p.Sub1), m.getSub(m.p.Sub2)
		if s1 == nil && s2 == nil {
			s1, s2 = subtitle.Subtitles(m.subs).Select()
		}
	} else {
		s1 = m.getSub(selectedSub)
		s2 = nil
	}

	switch {
	case s1 != nil && s2 != nil:
		m.p.Sub1, m.p.Sub2 = s1.Name, s2.Name
		s1.IsMainSub, s2.IsMainSub = true, false

	case s1 != nil && s2 == nil:
		s1.IsMainSub = true
		m.p.Sub1 = s1.Name

	case s1 == nil && s2 != nil:
		s1 = s2
		s1.IsMainSub = true
		m.p.Sub1 = s1.Name

	case s1 == nil && s2 == nil:
	}

	m.setPlayingSubs(s1, s2)

	shared.SavePlayingAsync(m.p)
}

func (m *Movie) setupSubtitles(subs map[string]*shared.Sub, selectedSub string) {
	if len(subs) > 0 {
		m.subs = nil
		width, height := m.v.Width, m.v.Height
		for _, sub := range subs {
			m.subs = append(m.subs, subtitle.NewSubtitle(sub, m.w, float64(width), float64(height)))
		}

		sort.Sort(subtitle.Subtitles(m.subs))

		m.setupDefaultSubtitles(selectedSub)
	}
}

func getNameDistance(from, to string) int {
	if i := strings.LastIndex(from, "."); i >= 0 {
		from = from[:i]
	}

	if i := strings.LastIndex(to, "."); i >= 0 {
		to = to[:i]
	}

	return levenshtein(from, to)
}

//This version uses dynamic programming with time complexity of O(mn) where m and n are lengths of a and b,
//	and the space complexity is n+1 of integers plus some constant space(i.e. O(n)).
//copy from:
//http://en.wikibooks.org/wiki/Algorithm_Implementation/Strings/Levenshtein_distance#Go
func levenshtein(a, b string) int {
	f := make([]int, utf8.RuneCountInString(b)+1)

	for j := range f {
		f[j] = j
	}

	for _, ca := range a {
		j := 1
		fj1 := f[0] // fj1 is the value of f[j - 1] in last iteration
		f[0]++
		for _, cb := range b {
			mn := minint(f[j]+1, f[j-1]+1) // delete & insert
			if unicode.ToLower(cb) != unicode.ToLower(ca) {
				mn = minint(mn, fj1+1) // change
			} else {
				mn = minint(mn, fj1) // matched
			}

			fj1, f[j] = f[j], mn // save f[j] to fj1(j is about to increase), update f[j] to mn
			j++
		}
	}

	return f[len(f)-1]
}

func minint(a, b int) int {
	if a <= b {
		return a
	} else {
		return b
	}
}
