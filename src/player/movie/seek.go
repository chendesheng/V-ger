package movie

import (
	"time"
)

func (m *Movie) seekOffset(offset time.Duration) {
	t := m.c.GetTime() + offset

	m.SeekBegin()

	var img []byte
	var err error
	t, img, err = m.v.SeekOffset(t)
	if err != nil {
		return
	}

	// m.w.RefreshContent(img)
	go m.w.SendDrawImage(img)

	m.c.SetTime(t)
	percent := m.c.GetPercent()
	m.w.ShowProgress(m.c.CalcPlayProgress(percent))

	if m.s != nil {
		m.s.Seek(t)
	}
	if m.s2 != nil {
		m.s2.Seek(t)
	}
	m.SeekEnd(t)
}

func (m *Movie) SeekBegin() {
	m.chSeekPause <- -1
	m.v.FlushBuffer()
	m.a.FlushBuffer()
}

func (m *Movie) Seek(t time.Duration) time.Duration {
	var img []byte
	var err error
	println("seek:", t.String())
	t, img, err = m.v.Seek(t)
	println("end seek:", t.String())
	if err != nil {
		return t
	}
	// println("seek refresh")
	if len(img) > 0 {
		m.w.RefreshContent(img)
	}

	m.c.SetTime(t)
	percent := m.c.GetPercent()
	m.w.ShowProgress(m.c.CalcPlayProgress(percent))

	if m.s != nil {
		m.s.Seek(t)
	}
	if m.s2 != nil {
		m.s2.Seek(t)
	}
	return t
}

func (m *Movie) SeekEnd(t time.Duration) {
	m.chSeekPause <- t
	println("seek end:", t.String())
}
