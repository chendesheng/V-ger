package movie

import (
	"fmt"
	"log"
	"time"
)

func (m *Movie) showProgressInner(t time.Duration) {
	p := m.c.CalcPlayProgress(t)

	// log.Print("showProgressInner", p.Left, p.Percent, p.Right)
	m.w.SendShowProgress(p.Left, p.Right, p.Percent)
}

//SpeedMonitor interface
func (m *Movie) SetSpeed(speed float64) {
	// log.Print("set speed:", speed)

	if m.chSpeed != nil {
		select {
		case m.chSpeed <- speed:
		case <-m.quit:
		case <-time.After(500 * time.Millisecond):
			log.Print("write m.chSpeed timeout")
		}
	}
}

func (m *Movie) showProgressPerSecond() {
	if m.httpBuffer != nil {
		m.w.SendShowBufferInfo("0 KB/s", 0)
	} else {
		m.w.SendShowBufferInfo("", 0)
	}

	var t time.Duration
	var lastTime time.Duration
	t = m.c.GetTime()
	lastTime = t

	var lastSpeed float64

	for {
		select {
		case t = <-m.chProgress:
			if t/time.Second != lastTime/time.Second {
				lastTime = t
				m.p.SetLastPos(t)
				m.showProgressInner(t)
			}
		case speed := <-m.chSpeed:
			if speed != lastSpeed {
				percent := m.httpBuffer.BufferPercent()
				// log.Print("send show speed:", speed)

				lastSpeed = speed
				m.p.Speed = speed
				m.w.SendShowBufferInfo(fmt.Sprintf("%.0f KB/s", speed), percent)
			}
		case <-m.quit:
			log.Print("quit show progress")
			return
		}
	}
}
