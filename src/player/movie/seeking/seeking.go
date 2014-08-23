package seeking

import (
	"log"
	"time"
)

type seekArg struct {
	t     time.Duration
	isEnd bool
}

type Seeking struct {
	v VideoSeeker
	h SeekHandler

	chSeek chan *seekArg
	chQuit chan struct{}
}

type SeekHandler interface {
	OnSeekStarted()
	OnSeek(time.Duration, []byte)
	OnSeekPaused(time.Duration)
	OnSeekEnded(time.Duration)
}

type VideoSeeker interface {
	Seek(time.Duration) (time.Duration, []byte, error)
	SeekOffset(time.Duration) (time.Duration, []byte, error)
	FlushBuffer()
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

func (s *Seeking) startSeekRoutine() {
	var chSeek chan *seekArg
	s.chSeek, chSeek = recentPipe(s.chQuit)

	go func(chSeek chan *seekArg) {
		var t time.Duration
		var arg *seekArg
		started := false
		for {
			select {
			case <-s.chQuit:
				return
			case arg = <-chSeek:
				if !started {
					started = true
					s.h.OnSeekStarted()
				}

				t = s.seek(arg.t)

				if arg.isEnd {
					started = false
					s.h.OnSeekEnded(t)
				}
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
	case s.chSeek <- &seekArg{t, false}:
	case <-s.chQuit:
	}
}
func (s *Seeking) SendEndSeek(t time.Duration) {
	select {
	case <-s.chQuit:
	case s.chSeek <- &seekArg{t, true}:
	}
}
