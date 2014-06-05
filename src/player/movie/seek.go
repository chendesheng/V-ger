package movie

import (
	. "player/libav"
	. "player/shared"
	"runtime"
	"time"
)

func (m *Movie) seekOffsetAsync(offset time.Duration) {
	go func() {
		if m.httpBuffer != nil {
			m.w.SendShowMessage("Bufferring...", false)
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
	seekToBegin := t == 0

	m.SeekBegin()

	var img []byte
	var err error
	t, img, err = m.v.SeekOffset(t)
	if err != nil {
		return
	}

	go m.w.SendDrawImage(img)
	if seekToBegin {
		t = 0
		m.ctx.SeekFrame(m.ctx.VideoStream(), t, AVSEEK_FLAG_FRAME)
	}

	m.c.SetTime(t)
	// percent := m.c.GetPercent()
	// m.w.ShowProgress(m.c.CalcPlayProgress(percent))

	if m.s != nil {
		m.s.Seek(t)
	}
	if m.s2 != nil {
		m.s2.Seek(t)
	}
	m.SeekEnd(t)
}

func (m *Movie) SeekBegin() {
	println("seek begin")
	m.chSeekPause <- -1

	m.v.FlushBuffer()
	m.a.FlushBuffer()

	chanSeek = make(chan time.Duration, 50)
	go func() {
		var t time.Duration
		var ok bool
		var lastTime time.Duration
		for {
			select {
			case t, ok = <-chanSeek:
				if !ok {
					chanSeek = nil

					println("seek end:", lastTime.String())
					lastTime = m.Seek(lastTime)
					println("seek end2:", lastTime.String())

					if m.httpBuffer != nil {
						m.w.SendShowMessage("Bufferring...", false)
						defer m.w.SendHideMessage()

						m.httpBuffer.Wait(1024 * 1024)
					}

					m.chSeekPause <- lastTime
					// println("seek end2:", t.String())

					SavePlayingAsync(m.p)
					return
				} else {
					lastTime = t
				}
			default:
				if t >= 0 {
					if m.httpBuffer == nil {
						m.Seek(t)
					} else {
						m.Seek2(t)
					}
					t = -1
				}
				runtime.Gosched()
			}
		}
	}()
	println("seek begin")
}

var chanSeek chan time.Duration

func (m *Movie) SeekAsync(t time.Duration) {
	println("seek async:", t.String())
	if chanSeek != nil {
		chanSeek <- t
	}
}

func (m *Movie) Seek2(t time.Duration) time.Duration {
	var img []byte
	var err error

	t, img, err = m.v.Seek2(t)

	if err != nil {
		return t
	}

	if len(img) > 0 {
		m.w.SendDrawImage(img)
	}

	// m.c.SetTime(t)

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
	// println("seek:", t.String())
	t, img, err = m.v.Seek(t)
	// println("end seek:", t.String())
	if err != nil {
		return t
	}
	// println("seek refresh")
	if len(img) > 0 {
		// m.w.RefreshContent(img)
		go m.w.SendDrawImage(img)
	}

	// m.c.SetTime(t)

	if m.s != nil {
		m.s.Seek(t)
	}
	if m.s2 != nil {
		m.s2.Seek(t)
	}

	// println("end seek2:", t.String())

	return t
}

func (m *Movie) SeekEnd(t time.Duration) {
	println("seek end:", t.String())
	if chanSeek != nil {
		chanSeek <- t
		close(chanSeek)
	}
}
