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
	ch := m.Pause(true)

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

func (m *Movie) handleSeekProgress(ch chan time.Duration, arg *seekArg, chSeekProgress chan *seekArg) chan time.Duration {
	if ch == nil {
		ch = m.Pause(true)
	}

	if m.httpBuffer != nil {
		m.w.SendShowBufferInfo(&BufferInfo{"-- KB/s", 0})
	}

	t := arg.t

	log.Print("seekProgress:", arg.t.String())
	t = m.Seek(t)

	if arg.isEnd {
		if m.httpBuffer != nil {
			m.w.SendShowSpinning()
			defer m.w.SendHideSpinning()
			m.httpBuffer.Wait(1024 * 1024)
			select {
			case arg := <-chSeekProgress:
				return m.handleSeekProgress(ch, arg, chSeekProgress)
			case <-m.quit:
				return nil
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
			return nil
		}
	}

	return ch
}
func (m *Movie) startSeekRoutine() {
	m.chSeekProgress = make(chan *seekArg)
	chSeekProgress := make(chan *seekArg)
	go recentPipe(m.chSeekProgress, chSeekProgress, m.quit)

	go func(chSeekProgress chan *seekArg) {
		var ch chan time.Duration
		for {
			select {
			case <-m.quit:
				return
			case arg := <-chSeekProgress:
				ch = m.handleSeekProgress(ch, arg, chSeekProgress)
			}
		}
	}(chSeekProgress)
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
	case m.chSeekProgress <- &seekArg{t, false}:
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
		log.Print("sendDrawImage")
		m.w.SendDrawImage(img)
	}

	m.seekPlayingSubs(t, false)
	return t
}

func (m *Movie) SeekEnd(t time.Duration) {
	log.Print("begin SeekEnd:", t.String())
	select {
	case <-m.quit:
	case m.chSeekProgress <- &seekArg{t, true}:
	}
	log.Print("end SeekEnd:", t.String())
}
