package movie

import (
	// . "player/libav"
	"log"
	. "player/shared"
	"runtime"
	"time"
)

func (m *Movie) seekOffsetAsync(offset time.Duration) {
	go func() {
		if m.httpBuffer != nil {
			m.w.SendShowMessage("Buffering...", false)
			defer m.w.SendHideMessage()
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
	m.w.SendSetCursor(true)
	m.w.FuncMouseMoved[1]() //TODO.....

	m.p.LastPos = t
	SavePlayingAsync(m.p)
	ch <- t
}

func (m *Movie) SeekBegin() {
	println("seek begin")

	if m.chSeekQuit != nil {
		close(m.chSeekQuit)
	}
	m.chSeekQuit = make(chan struct{})

	ch := m.Pause(true)

	m.chSeekProgress = make(chan time.Duration, 500)
	go func(ch chan time.Duration) {
		var ok bool
		var t time.Duration
		t = -1
		var lastTime time.Duration

		for {
			select {
			case <-m.quit:
				return
			case <-m.chSeekQuit:
				return
			case t, ok = <-m.chSeekProgress:
				if !ok {
					m.chSeekProgress = nil

					lastTime = m.SeekAccurate(lastTime)

					if m.httpBuffer != nil {
						waitSize := int64(1024 * 1024)
						if m.httpBuffer.BufferFinish(waitSize) {
							m.w.SendShowMessage("Buffering...", false)
							defer m.w.SendHideMessage()

							m.httpBuffer.WaitQuit(waitSize, m.chSeekQuit)
						}
					}

					m.p.LastPos = lastTime
					SavePlayingAsync(m.p)

					println("seek end send time:", lastTime.String())
					select {
					case ch <- lastTime:
					case <-m.chSeekQuit:
					case <-m.quit:
					}
					return
				} else {
					lastTime = t
				}
			default:
				if t >= 0 {
					m.Seek(t)
					t = -1
				}
				runtime.Gosched()
			}
		}
	}(ch)
}

func (m *Movie) SeekAsync(t time.Duration) {
	println("seek async:", t.String())
	if m.chSeekProgress != nil {
		select {
		case m.chSeekProgress <- t:
			SavePlayingAsync(m.p)
		case <-time.After(20 * time.Millisecond):
		}
	}
}

func (m *Movie) SeekAccurate(t time.Duration) time.Duration {
	println("seek2:", t.String())

	var img []byte
	var err error

	t, img, err = m.v.SeekAccurate(t)

	if err != nil {
		return t
	}

	if len(img) > 0 {
		println("send draw image:", t.String())
		m.w.SendDrawImage(img)
	}

	if m.s != nil {
		m.s.Seek(t)
	}
	if m.s2 != nil {
		m.s2.Seek(t)
	}

	return t
}

func (m *Movie) Seek(t time.Duration) time.Duration {
	var img []byte
	var err error
	t, img, err = m.v.Seek(t)
	if err != nil {
		return t
	}

	if len(img) > 0 {
		println("sendDrawImage")
		m.w.SendDrawImage(img)
	}

	if m.s != nil {
		m.s.Seek(t)
	}
	if m.s2 != nil {
		m.s2.Seek(t)
	}

	return t
}

func (m *Movie) SeekEnd(t time.Duration) {
	println("begin SeekEnd:", t.String())
	if m.chSeekProgress != nil {
		select {
		case m.chSeekProgress <- t:
			close(m.chSeekProgress)
		case <-m.chSeekQuit:
		}
	}
	println("end SeekEnd:", t.String())
}
