package seeking

import (
	"log"
	"time"
)

type seekArg struct {
	t        time.Duration
	isEnd    bool
	isOffset bool
}

type Seeking struct {
	v VideoSeeker
	h SeekHandler

	chSeek chan *seekArg
	chQuit chan struct{}
}

type SeekHandler interface {
	OnSeekStarted() time.Duration
	OnSeek(time.Duration, []byte)
	OnSeekPaused(time.Duration)
	OnSeekEnded(time.Duration)
}

type VideoSeeker interface {
	Seek(time.Duration) (time.Duration, []byte, error)
}

func NewSeeking(v VideoSeeker, h SeekHandler, chQuit chan struct{}) *Seeking {
	s := &Seeking{v, h, make(chan *seekArg), chQuit}
	s.startSeekRoutine()
	return s
}

func (s *Seeking) seek(t time.Duration) time.Duration {
	var img []byte
	var err error
	t, img, err = s.v.Seek(t)
	if err != nil {
		log.Print(err)
		return t
	}

	s.h.OnSeek(t, img)
	return t
}
func (s *Seeking) handleSeek(t time.Duration, arg *seekArg) time.Duration {
	if t < 0 {
		t = s.h.OnSeekStarted()
		log.Print("seek start:", t.String())
	}
	if arg.isOffset {
		if arg.t != 0 {
			t = s.seek(t + arg.t)
		}
	} else {
		t = s.seek(arg.t)
	}

	if arg.isEnd {
		s.h.OnSeekEnded(t)
		t = -1
	}

	return t
}

func (s *Seeking) startSeekRoutine() {
	var chSeek chan *seekArg
	s.chSeek, chSeek = recentPipe(s.chQuit)

	go func(chSeek chan *seekArg) {
		var arg *seekArg
		t := time.Duration(-1)
		for {
			select {
			case <-s.chQuit:
				return
			case arg = <-chSeek:
				t = s.handleSeek(t, arg)
			case <-time.After(30 * time.Millisecond):
				if arg != nil && !arg.isEnd {
					s.h.OnSeekPaused(t)
				}
			}
		}
	}(chSeek)
}

func (s *Seeking) SendSeek(t time.Duration) {
	select {
	case s.chSeek <- &seekArg{t, false, false}:
	case <-s.chQuit:
	}
}

func (s *Seeking) SendSeekOffset(t time.Duration) {
	select {
	case s.chSeek <- &seekArg{t, false, true}:
	case <-s.chQuit:
	}
}

func (s *Seeking) SendEndSeekOffset() {
	select {
	case <-s.chQuit:
	case s.chSeek <- &seekArg{0, true, true}:
	}
}

func (s *Seeking) SendEndSeek(t time.Duration) {
	select {
	case <-s.chQuit:
	case s.chSeek <- &seekArg{t, true, false}:
	}
}
