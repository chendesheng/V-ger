package movie

import (
	"fmt"
	"time"
)

func (m *Movie) showProgress() {
	m.showProgressInner(m.p.LastPos)
}

func (m *Movie) showProgressInner(t time.Duration) {
	p := m.c.CalcPlayProgress(t)
	if m.httpBuffer != nil {
		p.Percent2 = m.httpBuffer.BufferPercent()
	}

	println("showProgressInner", p.Left, p.Percent, p.Right)
	m.w.SendShowProgress(p)
}

//SpeedMonitor interface
func (m *Movie) SetSpeed(speed float64) {
	if m.chSpeed != nil {
		select {
		case m.chSpeed <- speed:
		case <-m.quit:
		case <-time.After(500 * time.Millisecond):
			println("write m.chSpeed timeout")
		}
	}
}

func (m *Movie) showProgressPerSecond() {
	m.chProgress = make(chan time.Duration)
	if m.httpBuffer != nil {
		m.chSpeed = make(chan float64)
		m.w.SendShowSpeed("0 KB/s")
	}

	var t time.Duration
	var lastTime time.Duration
	t = m.c.GetTime()
	lastTime = t

	var speed float64
	var lastSpeed float64

	for {
		select {
		case t = <-m.chProgress:
			if t/time.Second != lastTime/time.Second {
				lastTime = t
				m.p.LastPos = t
				m.showProgressInner(t)
			}
		case speed = <-m.chSpeed:
			if speed != lastSpeed {
				lastSpeed = speed
				m.p.Speed = speed
				m.w.SendShowSpeed(fmt.Sprintf("%.0f KB/s", speed))
			}
		case <-m.quit:
			return
		}
	}
}
