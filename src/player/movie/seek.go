package movie

import (
	// . "player/libav"
	"log"
	. "player/shared"
	"time"
)

func (m *Movie) SeekOffset(offset time.Duration) {
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

	m.OnSeekStarted()

	t, img, err := m.v.SeekOffset(t)
	if err != nil {
		log.Print(err)
		return
	}

	m.showProgressInner(t)
	m.OnSeek(t, img)

	if len(m.w.FuncMouseMoved) > 0 {
		m.w.FuncMouseMoved[0]()
	}

	m.OnSeekEnded(t)
}

func (m *Movie) OnSeekStarted() {
	m.hold()
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
