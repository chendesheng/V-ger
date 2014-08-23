package movie

import (
	// . "player/libav"
	"log"
	. "player/shared"
	"time"
)

func (m *Movie) seekOffsetAsync(offset time.Duration) {
	go func() {
		if m.httpBuffer != nil {
			m.w.SendShowSpinning()
			defer m.w.SendHideSpinning()
		}

		m.seekOffset(offset)
	}()
}
func (m *Movie) seekOffset(offset time.Duration) {
	t := m.c.GetTime() + offset
	if t < 0 {
		t = 0
	}
	ch := m.Hold()

	t, img, err := m.v.SeekOffset(t)
	if err != nil {
		log.Print(err)
		return
	}

	m.showProgressInner(t)
	m.w.SendDrawImage(img)

	if len(m.w.FuncMouseMoved) > 0 {
		m.w.FuncMouseMoved[0]()
	}

	m.p.LastPos = t
	SavePlayingAsync(m.p)

	if m.httpBuffer != nil {
		defer m.w.SendHideMessage()
		m.httpBuffer.WaitQuit(1024*1024, m.quit)
	}

	select {
	case ch <- t:
	case <-m.quit:
	}
}

func (m *Movie) handleSeekProgress(ch chan time.Duration, arg *seekArg, chSeek chan *seekArg) (chan time.Duration, time.Duration) {
	if ch == nil {
		ch = m.Hold()
	}

	if m.httpBuffer != nil {
		m.w.SendShowBufferInfo(&BufferInfo{"-- KB/s", 0})
	}

	t := m.Seek(arg.t)

	if arg.isEnd {
		if m.httpBuffer != nil {
			m.w.SendShowSpinning()
			defer m.w.SendHideSpinning()
			m.httpBuffer.WaitQuit(1024*1024, m.quit)
			select {
			case arg := <-chSeek:
				return m.handleSeekProgress(ch, arg, chSeek)
			case <-m.quit:
				return nil, 0
			default:
			}
		}

		m.p.LastPos = t
		SavePlayingAsync(m.p)

		log.Print("seek end end time:", t.String())
		select {
		case ch <- t:
			ch = nil
		case <-m.quit:
			return nil, 0
		}
	}

	return ch, t
}
func (m *Movie) startSeekRoutine() {
	m.chSeek = make(chan *seekArg)
	chSeek := make(chan *seekArg)
	go recentPipe(m.chSeek, chSeek, m.quit)

	go func(chSeek chan *seekArg) {
		var ch chan time.Duration
		var t time.Duration
		for {
			select {
			case <-m.quit:
				return
			case arg := <-chSeek:
				ch, t = m.handleSeekProgress(ch, arg, chSeek)
			case <-time.After(30 * time.Millisecond):
				if ch != nil {
					m.showProgressInner(t)
				}
			}
		}
	}(chSeek)
}

func recentPipe(in chan *seekArg, out chan *seekArg, quit chan struct{}) {
	var recentValue *seekArg
	var sendout chan *seekArg
	for {
		select {
		case t, ok := <-in:
			if !ok {
				return
			}
			sendout = out
			recentValue = t
		case sendout <- recentValue:
			sendout = nil
		case <-quit:
			return
		}
	}
}

func (m *Movie) SeekAsync(t time.Duration) {
	//log.Print("seek async:", t.String())
	select {
	case m.chSeek <- &seekArg{t, false}:
	case <-m.quit:
	}
}

func (m *Movie) SeekAccurate(t time.Duration) time.Duration {
	log.Print("seek2:", t.String())

	var img []byte
	var err error

	t, img, err = m.v.SeekAccurate(t)

	if err != nil {
		return t
	}

	if len(img) > 0 {
		log.Print("send draw image:", t.String())
		m.w.SendDrawImage(img)
	}

	m.seekPlayingSubs(t, false)
	return t
}

func (m *Movie) Seek(t time.Duration) time.Duration {
	var img []byte
	var err error
	t, img, err = m.v.Seek(t)
	if err != nil {
		log.Print(err)
		return t
	}

	if len(img) > 0 {
		// log.Print("sendDrawImage")
		m.w.SendDrawImage(img)
	}

	m.seekPlayingSubs(t, false)
	return t
}

func (m *Movie) SeekEnd(t time.Duration) {
	select {
	case m.chSeek <- &seekArg{t, true}:
	case <-m.quit:
	}
}
