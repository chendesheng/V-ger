package movie

import (
	// . "player/libav"
	"log"
	. "player/shared"
	"time"
)

var timerEndSeek *time.Timer

func (m *Movie) SeekOffset(offset time.Duration) {
	if len(m.w.FuncMouseMoved) > 0 {
		m.w.FuncMouseMoved[0]()
	}

	m.seeking.SendSeekOffset(offset)

	if timerEndSeek == nil || !timerEndSeek.Reset(200*time.Millisecond) {
		if timerEndSeek == nil {
			timerEndSeek = time.NewTimer(200 * time.Millisecond)
		}
		go func() {
			select {
			case <-timerEndSeek.C:
				m.seeking.SendEndSeekOffset()
			case <-m.quit:
				timerEndSeek.Stop()
			}
		}()
	}
}
func (m *Movie) seekOffset(offset time.Duration) {
	t := m.c.GetTime() + offset
	if t < 0 {
		t = 0
	}

	m.OnSeekStarted()

	t, img, err := m.v.Seek(t)
	if err != nil {
		log.Print(err)
		return
	}

	m.OnSeek(t, img)

	if len(m.w.FuncMouseMoved) > 0 {
		m.w.FuncMouseMoved[0]()
	}

	m.OnSeekEnded(t)
	m.showProgressInner(t)
}

func (m *Movie) OnSeekStarted() time.Duration {
	t := m.c.GetTime()
	m.hold()

	return t / time.Second * time.Second
}

func (m *Movie) OnSeek(t time.Duration, img []byte) {
	if len(img) > 0 {
		m.w.SendDrawImage(img)
	}

	m.seekPlayingSubs(t, false)

}
func (m *Movie) OnSeekPaused(t time.Duration) {
	m.showProgressInner(t)
}
func (m *Movie) OnSeekEnded(t time.Duration) {
	if m.httpBuffer != nil {
		m.httpBuffer.WaitQuit(1024*1024, m.quit)
		m.w.SendHideSpinning()
	}

	m.unHold(t)

	m.p.LastPos = t
	SavePlayingAsync(m.p)

}
