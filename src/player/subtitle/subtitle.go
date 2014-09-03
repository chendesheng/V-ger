package subtitle

import (
	"sync"
	// "io/ioutil"
	// "cld"
	"log"
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

	r          SubRender
	items      *subItems
	offset     time.Duration
	Name       string
	IsMainSub  bool
	Lang1      string //one subtitle file may has two languages
	Lang2      string
	Type       string
	chanRes    chan SubItemExtra
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

func NewSubtitle(sub *Sub, r SubRender, width, height float64) *Subtitle {
	if sub == nil {
		return nil
	}
	s := &Subtitle{}

	err := s.parse(sub, width, height)
	if err != nil {
		log.Print(err)
		return nil
	}

	s.offset = sub.Offset
	s.Name = sub.Name
	s.IsMainSub = true
	s.Lang1 = sub.Lang1
	s.Lang2 = sub.Lang2

	s.Type = sub.Type

	if len(sub.Lang1) == 0 && len(sub.Lang2) == 0 {
		s.Lang1, s.Lang2 = detectLanguage(s.items)
		UpdateSubtitleLanguage(s.Name, s.Lang1, s.Lang2)
	}

	simplized(s.items)

	s.r = r

	log.Printf("parse sub:%s, %s, %s; %d items", sub.Name, sub.Lang1, sub.Lang2, len(s.items.nooverlap))
	return s
}

func (s *Subtitle) parse(sub *Sub, width, height float64) error {
	var items []*SubItem
	var err error

	if sub.Type == "ass" {
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

func (s *Subtitle) GetSubTime(t time.Duration, offset int) time.Duration {
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
