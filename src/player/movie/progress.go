package movie

import (
	"log"
	"task"
	"time"
)

func (m *Movie) showProgress(name string) {
	m.p.LastPos = m.c.GetTime()

	p := m.c.CalcPlayProgress(m.c.GetPercent())

	done := make(chan struct{})
	go func() {
		t, err := task.GetTask(name)

		if err == nil {
			if t.Status == "Finished" {
				p.Percent2 = 1
			} else {
				p.Percent2 = float64(t.BufferedPosition) / float64(t.Size)
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

func (m *Movie) showProgressPerSecond(name string) {
	ticker := time.NewTicker(time.Second)
	for {
		if m.c.WaitUtilRunning(m.quit) {
			return
		}

		select {
		case <-ticker.C:
			m.showProgress(name)
		case <-m.quit:
			return
		}
	}
}
