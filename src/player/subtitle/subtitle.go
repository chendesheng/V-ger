package subtitle

import (
	"sync"
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
	sync.Mutex

	r     SubRender
	items *subItems

	c *Clock

	offset time.Duration
	quit   chan bool

	ChanSeek        chan time.Duration
	ChanSeekRefersh chan time.Duration

	ChanOffset chan durationArg

	chanGetSubTime chan subTimeArg

	chanStop chan bool

	Name string

	IsMainSub bool

	Lang1 string //one subtitle file may has double languages
	Lang2 string

	Format string

	chanRes chan SubItemExtra

	displaying []*SubItem
}

type displayingItem struct {
	pos    int
	handle uintptr
}

func (s *Subtitle) Seek(t time.Duration, refersh bool) {
	s.Lock()
	defer s.Unlock()

	t -= s.offset

	if s.chanRes == nil {
		s.chanRes = make(chan SubItemExtra, 20)
		go func() {
			for arg := range s.chanRes {
				s.Lock()

				item := s.items.getById(arg.Id)
				item.Handle = arg.Handle
				s.displaying = append(s.displaying, item)

				s.Unlock()
			}
		}()
	}

	if refersh {
		s.hideAll()
	} else {
		for i := len(s.displaying) - 1; i >= 0; i-- {
			item := s.displaying[i]

			if !item.Contains(t) {
				s.hideSubItem(item)

				s.displaying = sliceRemove(s.displaying, i)
			}
		}
	}

	for _, item := range s.items.get(t) {
		if item.Handle == 0 {
			s.showSubitem(*item)
			item.Handle = 1
		}
	}
}

//check https://code.google.com/p/go-wiki/wiki/SliceTricks
func sliceRemove(a []*SubItem, i int) []*SubItem {
	l := len(a)
	a[i], a[l-1], a = a[l-1], nil, a[:l-1]
	return a
}

func (s *Subtitle) Stop() {
	s.Lock()
	defer s.Unlock()

	s.hideAll()
}

func (s *Subtitle) hideAll() {
	for _, item := range s.displaying {
		s.hideSubItem(item)
	}
	s.displaying = nil
}

func (s *Subtitle) checkPos(pos int, t time.Duration) bool {
	item := s.items.getById(pos)
	if item == nil {
		return false
	} else {
		return item.Contains(t)
	}
}

func (s *Subtitle) showSubitem(item SubItem) {
	if !s.IsMainSub && item.PositionType != 10 {
		if item.IsInDefaultPosition() {
			item.PositionType = 10
		} else {
			return
		}
	}
	arg := SubItemArg{item, false, s.chanRes}
	go s.r.SendShowText(arg)
}
func (s *Subtitle) hideSubItem(item *SubItem) {
	go s.r.SendHideText(SubItemArg{*item, false, nil})
	item.Handle = 0
}

func detectLanguage(si *subItems) (string, string) {
	content := ""
	si.each(func(item *SubItem) {
		for _, attrStr := range item.Content {
			content += attrStr.Content
		}
	})

	return language.DetectLanguages(content)
}

func simplized(si *subItems) {
	si.each(func(item *SubItem) {
		for i, attrStr := range item.Content {
			item.Content[i].Content = language.Simplized(attrStr.Content)
		}
	})
}

func NewSubtitle(sub *Sub, r SubRender, c *Clock, width, height float64) *Subtitle {
	if sub == nil {
		return nil
	}
	s := &Subtitle{}

	err := s.parse(sub, width, height)
	if err != nil {
		log.Print(err)
		return nil
	}

	s.c = c
	s.offset = sub.Offset
	s.ChanSeek = make(chan time.Duration)
	s.ChanSeekRefersh = make(chan time.Duration)
	s.ChanOffset = make(chan durationArg)
	s.chanStop = make(chan bool)
	s.chanGetSubTime = make(chan subTimeArg)
	s.Name = sub.Name
	s.IsMainSub = true
	s.Lang1 = sub.Lang1
	s.Lang2 = sub.Lang2

	s.Format = sub.Type

	if len(sub.Lang1) == 0 && len(sub.Lang2) == 0 {
		s.Lang1, s.Lang2 = detectLanguage(s.items)
		UpdateSubtitleLanguage(s.Name, s.Lang1, s.Lang2)
	}

	simplized(s.items)

	s.quit = make(chan bool)
	s.r = r

	log.Printf("parse sub:%s; %d items", sub.Name, len(s.items.nooverlap))
	return s
}

func (s *Subtitle) parse(sub *Sub, width, height float64) error {
	var items []*SubItem
	var err error

	if s.Format == "ass" {
		items, err = ass.Parse(strings.NewReader(sub.Content), width, height)
		if err != nil {
			return err
		}
	} else {
		items, err = srt.Parse(strings.NewReader(sub.Content), width, height)
		if err != nil {
			return err
		}
	}

	s.items = newSubItems(items)
	return nil
}

func (s *Subtitle) AddOffset(d time.Duration) time.Duration {
	s.Lock()
	defer s.Unlock()

	s.offset += d
	return s.offset
}

func (s *Subtitle) GetSubtime(t time.Duration, offset int) time.Duration {
	s.Lock()
	defer s.Unlock()

	t -= s.offset

	pos, ok := s.items.findPos(t)
	if ok || offset > 0 {
		pos += offset
	}

	if item := s.items.getById(pos); item != nil {
		return item.From
	} else {
		return 0
	}
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
