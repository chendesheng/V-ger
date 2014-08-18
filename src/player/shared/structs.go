package shared

import "time"

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
	LastPos     time.Duration
	SoundStream int
	Sub1        string
	Sub2        string
	Duration    time.Duration
	Volume      byte
	Speed       float64 //online video downlad speed
}

func (item *SubItem) IsInDefaultPosition() bool {
	return item.PositionType == 2 && item.X < 0 && item.Y < 0
}

type SubItemExtra struct {
	Id     int
	Handle uintptr
}
type SubItemArg struct {
	SubItem
	AutoHide bool
	Result   chan SubItemExtra
}
type Position struct {
	X, Y float64
}

type PlayProgressInfo struct {
	Left    string
	Right   string
	Percent float64
}
type BufferInfo struct {
	Speed         string
	BufferPercent float64
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
