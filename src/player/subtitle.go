package main

import (
	"fmt"
	"io/ioutil"
	"log"
	. "player/clock"
	"player/glfw"
	"player/srt"
	// "strings"
	"sync"
	"time"
)

type subtitle struct {
	sync.Locker

	w *Window

	items []*srt.SubItem
	pos   int

	c *Clock

	offset time.Duration
	quit   chan bool
}

func NewSubtitle(file string, w *Window) *subtitle {
	var err error
	bytes, err := ioutil.ReadFile(file)
	if err != nil {
		log.Print(err)
		return nil
	}
	s := &subtitle{}
	s.Locker = &sync.Mutex{}

	s.items = srt.Parse(string(bytes))
	if err != nil {
		log.Print(err)
		return nil
	}

	s.quit = make(chan bool)
	s.w = w

	log.Print("sub items:", len(s.items))
	w.AddEventHandler(func(e Event) { //run in main thread, safe to operate ui elements
		switch e.Kind {
		case KeyPress:
			switch e.Data.(glfw.Key) {
			case glfw.KeyMinus:
				println("key minus pressed")
				s.addOffset(-1000 * time.Millisecond)
				break
			case glfw.KeyEqual:
				println("key equal pressed")
				s.addOffset(1000 * time.Millisecond)
				break
			case glfw.KeyLeftBracket:
				println("left bracket pressed")
				s.addOffset(-200 * time.Millisecond)
				break
			case glfw.KeyRightBracket:
				println("right bracket pressed")
				s.addOffset(200 * time.Millisecond)
				break
			}
			break
		}
	})
	return s
}

func (s *subtitle) setPosition(pos int) {
	// atomic.StoreInt32(&s.pos, int32(pos))
	s.Lock()
	defer s.Unlock()

	s.pos = pos
}
func (s *subtitle) position() int {
	s.Lock()
	defer s.Unlock()

	return s.pos
}
func (s *subtitle) increasePosition() {
	s.Lock()
	defer s.Unlock()

	s.pos += 1
}
func (s *subtitle) seek(t time.Duration) {
	for i, item := range s.items {
		to := item.To + time.Duration(s.offset)*time.Second
		if to > t {
			log.Print("seek to ", to.String(), " i: ", i, " Content:", item.Content)
			s.setPosition(i)
			return
		}
	}
}

func (s *subtitle) addOffset(delta time.Duration) {
	s.Lock()
	s.offset += delta
	close(s.quit)
	s.pos = 0
	s.Unlock()

	go s.play()
}

func (s *subtitle) getOffset() time.Duration {
	s.Lock()
	s.Unlock()

	return s.offset
}
func (s *subtitle) playWithQuit(quit chan bool) {
	for s.position() < len(s.items) {
		item := s.items[s.position()]
		s.increasePosition()

		offset := s.getOffset()
		from := item.From + offset
		to := item.To + offset
		if to < s.c.GetTime() {
			continue
		}
		if s.c.WaitUtilWithQuit(from, quit) {
			return
		}

		fmt.Printf("play subtitle %v, from: %s, to: %s\n", item.Content, from.String(), to.String())

		s.w.PostEvent(Event{DrawSub, item.Content})

		nextFrom := to
		nextPos := s.position()
		if nextPos < len(s.items) {
			nextFrom = s.items[nextPos].From
		}

		go func(to, nextFrom time.Duration) { //overlap time, it's really nice with goroutine.
			if to > nextFrom {
				return
			}

			if s.c.WaitUtilWithQuit(to-50*time.Millisecond, quit) {
				return
			}

			s.w.PostEvent(Event{DrawSub, make([]srt.AttributedString, 0)})
		}(to, nextFrom)
	}
}
func (s *subtitle) play() {
	s.quit = make(chan bool)
	s.playWithQuit(s.quit)
}
