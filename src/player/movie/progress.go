package movie

import (
	"fmt"
	"time"
)

func (m *Movie) showProgress() {
	m.p.LastPos = m.c.GetTime()

	p := m.c.CalcPlayProgress(m.c.GetPercent())
	if m.httpBuffer != nil {
		p.Speed = fmt.Sprintf("%d KB/s", int(m.p.Speed))
		p.Percent2 = float64(m.httpBuffer.currentPos+m.httpBuffer.SizeAhead()) / float64(m.httpBuffer.size)
	}

	m.w.SendShowProgress(p)
}

var chShowProgress chan struct{}

func (m *Movie) showProgressPerSecond() {
	chShowProgress = make(chan struct{})

	ticker := time.NewTicker(time.Second)
	for {
		// if m.c.WaitUtilRunning(m.quit) {
		// 	return
		// }

		select {
		case <-chShowProgress:
			<-chShowProgress
			break
		case <-ticker.C:
			m.showProgress()
			break
		case <-m.quit:
			return
		}
	}
}

func (m *Movie) toggleShowProgress() {
	chShowProgress <- struct{}{}
}
