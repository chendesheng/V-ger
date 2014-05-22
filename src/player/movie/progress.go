package movie

import (
	"fmt"
	"log"
	"task"
	"time"
)

func (m *Movie) showProgress() {
	name := m.p.Movie
	m.p.LastPos = m.c.GetTime()

	p := m.c.CalcPlayProgress(m.c.GetPercent())
	if m.httpBuffer != nil {
		p.Speed = fmt.Sprintf("%d KB/s", int(m.p.Speed))
	}

	println("download speed:", p.Speed)

	done := make(chan struct{})
	go func() {
		t, err := task.GetTask(name)

		if err == nil {
			if t.Status == "Finished" {
				p.Percent2 = 1
			} else {
				if m.httpBuffer != nil {
					p.Percent2 = float64(m.httpBuffer.currentPos+m.httpBuffer.SizeAhead()) / float64(t.Size)
				} else {
					p.Percent2 = float64(t.BufferedPosition) / float64(t.Size)
				}
			}
		} else {
			log.Print(err)
		}
		close(done)
	}()

	select {
	case <-done:
		break
	case <-time.After(100 * time.Millisecond):
		break
	}

	m.w.SendShowProgress(p)
}

var chShowProgress chan struct{}

func (m *Movie) showProgressPerSecond() {
	chShowProgress = make(chan struct{})

	ticker := time.NewTicker(time.Second)
	for {
		if m.c.WaitUtilRunning(m.quit) {
			return
		}

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
