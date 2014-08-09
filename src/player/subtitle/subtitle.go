package subtitle

import (
	// "io/ioutil"
	// "cld"
	"log"
	. "player/clock"
	"player/language"
	. "player/shared"
	"player/subtitle/ass"
	"player/subtitle/srt"
	"strings"
	// "runtime"
	"time"
)

type durationArg struct {
	d   time.Duration
	res chan time.Duration
}
type subTimeArg struct {
	t      time.Duration
	offset int
	res    chan time.Duration
}

type Subtitle struct {
	r     SubRender
	items []*SubItem

	c *Clock

	offset time.Duration
	quit   chan bool

	ChanSeek        chan time.Duration
	ChanSeekRefersh chan time.Duration

	ChanOffset chan durationArg

	chanGetSubTime chan subTimeArg

	chanStop chan bool

	Name string

	IsMainOrSecondSub bool

	Lang1 string //one subtitle file may has double languages
	Lang2 string

	Format string
}

type displayingItem struct {
	pos    int
	handle uintptr
}

func (s *Subtitle) printSub(pos int) {
	if pos >= len(s.items) {
		return
	}

	item := s.items[pos]
	content := ""
	if len(item.Content) > 0 {
		content = item.Content[0].Content
	}

	log.Printf("%d:%s", pos, content)
}

func (s *Subtitle) Play() {
	for i, _ := range s.items {
		s.items[i].Id = i
	}

	chRes := make(chan SubItemExtra, 20)
	for {
		select {
		case arg := <-chRes:
			// log.Print("res from show:", arg.Id, arg.Handle)
			s.items[arg.Id].Handle = arg.Handle
			break
		case arg := <-s.ChanOffset:
			s.offset += arg.d
			arg.res <- s.offset
			break
		case arg := <-s.chanGetSubTime:
			pos, _ := s.findPos(arg.t)
			if pos >= len(s.items) {
				arg.res <- 0
				break
			}

			if arg.offset < 0 && s.calcFrom(pos) > arg.t {
				pos--
			}

			// pos += arg.offset
			// log.Print("pos:", pos)
			for s.checkPos(pos, arg.t) {
				pos += arg.offset
			}

			for {
				if pos < 0 {
					pos = 0
					break
				}
				if pos >= len(s.items) {
					pos = len(s.items) - 1
					break
				}

				item := s.items[pos]
				if !item.IsInDefaultPosition() {
					pos += arg.offset

				} else {
					break
				}
			}

			arg.res <- s.calcFrom(pos)
			break
		case t := <-s.ChanSeek:
			s.render(t, chRes, false)
			break
		case t := <-s.ChanSeekRefersh:
			s.render(t, chRes, true)
			break
		case <-s.chanStop:
			for _, item := range s.items {
				if item.Handle != 0 {
					s.hideSubItem(*item)
					item.Handle = 0
				}
			}
			return
		}
	}
}

func (s *Subtitle) render(t time.Duration, chRes chan SubItemExtra, refersh bool) {
	for i, item := range s.items {
		if refersh || !s.checkPos(i, t) {
			if item.Handle != 0 && item.Handle != 1 {
				// log.Print(t.String(), "hide sub: ", item.Id, item.Content[0].Content)
				s.hideSubItem(*item)
				item.Handle = 0
			}
		}
	}

	for i, item := range s.items {
		if s.checkPos(i, t) {
			if item.Handle == 0 {
				// log.Print(t.String(), "show sub: ", item.Id, item.Content[0].Content)
				s.showSubitem(*item, chRes)
				item.Handle = 1
			}
		}
	}
}

func (s *Subtitle) Stop() {
	s.chanStop <- true
}

func (s *Subtitle) calcFromTo(pos int) (time.Duration, time.Duration) {
	item := s.items[pos]
	return item.From + s.offset, item.To + s.offset
}

func (s *Subtitle) calcFrom(pos int) time.Duration {
	item := s.items[pos]
	return item.From + s.offset
}

func (s *Subtitle) checkPos(pos int, t time.Duration) bool {
	if pos >= len(s.items) || pos < 0 {
		return false
	}

	from, to := s.calcFromTo(pos)
	// log.Print("check pos:", pos, t.String(), from.String(), to.String())
	return t >= from && t < to
}

func (s *Subtitle) showSubitem(item SubItem, chRes chan SubItemExtra) {
	if !s.IsMainOrSecondSub && item.PositionType != 10 {
		if (item.PositionType != 2) || (item.X >= 0) || (item.Y >= 0) {
			return
		} else {
			item.PositionType = 10
		}
	}
	arg := SubItemArg{item, false, chRes}
	go s.r.SendShowText(arg)
}
func (s *Subtitle) hideSubItem(item SubItem) {
	if item.Handle != 0 && item.Handle != 1 {
		go s.r.SendHideText(SubItemArg{item, false, nil})
	}
}

func (s *Subtitle) findPos(t time.Duration) (int, *SubItem) {
	for i := 0; i < len(s.items); i++ {
		from, to := s.calcFromTo(i)
		if t < to {
			if t >= from {
				return i, s.items[i]
			} else {
				return i, nil
			}
		}
	}

	return 1 << 31, nil
}

func detectLanguage(items []*SubItem) (string, string) {
	content := ""
	for _, item := range items {
		for _, attrStr := range item.Content {
			content += attrStr.Content
		}
	}

	return language.DetectLanguages(content)
}

func simplized(items []*SubItem) {
	for _, item := range items {
		for i, attrStr := range item.Content {
			item.Content[i].Content = language.Simplized(attrStr.Content)
		}
	}
}

func NewSubtitle(sub *Sub, r SubRender, c *Clock, width, height float64) *Subtitle {
	var err error
	// sub := GetSubtitle(file)
	if sub == nil {
		return nil
	}
	s := &Subtitle{}
	s.c = c
	s.offset = sub.Offset
	s.ChanSeek = make(chan time.Duration)
	s.ChanSeekRefersh = make(chan time.Duration)
	s.ChanOffset = make(chan durationArg)
	s.chanStop = make(chan bool)
	s.chanGetSubTime = make(chan subTimeArg)
	s.Name = sub.Name
	s.IsMainOrSecondSub = true
	s.Lang1 = sub.Lang1
	s.Lang2 = sub.Lang2

	s.Format = sub.Type
	if s.Format == "ass" {
		s.items, err = ass.Parse(strings.NewReader(sub.Content), width, height)
	} else {
		s.items, err = srt.Parse(strings.NewReader(sub.Content), width, height)
	}

	if len(sub.Lang1) == 0 && len(sub.Lang2) == 0 {
		s.Lang1, s.Lang2 = detectLanguage(s.items)
		UpdateSubtitleLanguage(s.Name, s.Lang1, s.Lang2)
	}

	simplized(s.items)

	if err != nil {
		log.Print(err)
		return nil
	}

	s.quit = make(chan bool)
	s.r = r

	log.Printf("parse sub:%s; %d items", sub.Name, len(s.items))
	return s
}

func (s *Subtitle) Seek(t time.Duration, refresh bool) {
	var input chan<- time.Duration
	if refresh {
		input = s.ChanSeekRefersh
	} else {
		input = s.ChanSeek
	}

	select {
	case input <- t:
	case <-time.After(200 * time.Millisecond):
		log.Print("Seek sub timeout1")
	}
}

func (s *Subtitle) AddOffset(d time.Duration) time.Duration {
	res := make(chan time.Duration)
	s.ChanOffset <- durationArg{d, res}
	return <-res
}

func (s *Subtitle) GetSubtime(t time.Duration, offset int) time.Duration {
	res := make(chan time.Duration)
	arg := subTimeArg{t, offset, res}
	s.chanGetSubTime <- arg
	return <-arg.res
}

func (s *Subtitle) IsTwoLangs() bool {
	return len(s.Lang1) > 0 && len(s.Lang2) > 0
}

func compareLang(a1, a2, b1, b2 string) int {
	if a1 == b1 && a2 == b2 {
		return 0
	}
	//multi lang > signle lang
	if len(a2) > 0 && len(b2) == 0 {
		return 1
	}
	if len(a2) == 0 && len(b2) > 0 {
		return -1
	}
	//cn > en
	if len(a2) == 0 && len(b2) == 0 {
		if a1 == "chs" {
			return 1
		}
		if b1 == "chs" {
			return -1
		}
		if a1 == "cht" {
			return 1
		}
		if b1 == "cht" {
			return -1
		}
		return 1
	}

	if a2 == "chs" {
		return 1
	}
	if b2 == "chs" {
		return -1
	}
	if a2 == "cht" {
		return 1
	}
	if b2 == "cht" {
		return -1
	}
	return 1
}
func compareFormat(a, b string) int {
	if a == b {
		return 0
	} else if a == "srt" {
		return 1
	} else {
		return -1
	}

	return 1
}

type Subtitles []*Subtitle

func (s Subtitles) Len() int {
	return len([]*Subtitle(s))
}

func (s Subtitles) Less(i, j int) bool {
	a := s[i]
	b := s[j]

	c := compareLang(a.Lang1, a.Lang2, b.Lang1, b.Lang2)
	if c != 0 {
		return c > 0
	} else {
		c = compareFormat(a.Format, b.Format)
		if c != 0 {
			return c > 0
		}
	}

	return false
}

func (s Subtitles) Swap(i, j int) {
	t := s[i]
	s[i] = s[j]
	s[j] = t
}

func (s Subtitles) Select() (a *Subtitle, b *Subtitle) {
	subs := ([]*Subtitle)(s)
	if len(subs) == 1 || subs[0].IsTwoLangs() {
		a = subs[0]
		b = nil
	} else {
		if subs[0].Lang1 == "en" {
			a = subs[0]
			b = nil
		} else {
			a = subs[0]
			for _, c := range subs {
				if c.Lang1 == "en" {
					b = c
					break
				}
			}
		}
	}

	return
}
