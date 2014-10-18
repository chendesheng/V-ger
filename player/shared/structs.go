package shared

import (
	"sync/atomic"

	"time"
)

type AttributedString struct {
	Content string
	Style   int //0 -normal, 1 -italic, 2 -bold, 3 italic and bold
	Color   uint
}

type SubItem struct {
	From, To time.Duration
	Content  []AttributedString

	PositionType int
	Position

	SubItemExtra
}

func (s *SubItem) String() string {
	var res string
	for _, c := range s.Content {
		res += c.Content
	}
	return res
}

func (s *SubItem) IsEmpty() bool {
	return len(s.Content) == 0 //|| (s.To-s.From) < 100*time.Millisecond
}

func (s *SubItem) Contains(t time.Duration) bool {
	return s.From <= t && t < s.To
}

type Sub struct {
	Movie   string
	Name    string
	Offset  time.Duration
	Content string
	Type    string

	Lang1 string //one subtitle file may has double languages
	Lang2 string
}

type Playing struct {
	Movie       string
	LastPos     int64
	SoundStream int
	Sub1        string
	Sub2        string
	Duration    time.Duration
	Volume      int
	Speed       float64 //online video downlad speed
}

func (p *Playing) GetLastPos() time.Duration {
	return time.Duration(atomic.LoadInt64(&p.LastPos))
}

func (p *Playing) SetLastPos(t time.Duration) {
	atomic.StoreInt64(&p.LastPos, int64(t))
}

func (item *SubItem) IsInDefaultPosition() bool {
	return item.PositionType == 2 && item.X < 0 && item.Y < 0
}

type SubItemExtra struct {
	Id     int
	Handle int
}
type SubItemArg struct {
	SubItem
	AutoHide bool
}
type Position struct {
	X, Y float64
}

type PlayProgressInfo struct {
	Left    string
	Right   string
	Percent float64
}

type SubItems []*SubItem

func (s SubItems) Len() int {
	return len([]*SubItem(s))
}

func (s SubItems) Less(i, j int) bool {
	return s[i].From < s[j].From
}

func (s SubItems) Swap(i, j int) {
	t := s[i]
	s[i] = s[j]
	s[j] = t
}
