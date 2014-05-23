package download

import (
	"log"
	"math"
	"time"
)

type SpeedMonitor interface {
	SetSpeed(float64)
}

type speedFilter struct {
	basicFilter
	sm SpeedMonitor
}

func (sf *speedFilter) active() {
	timer := time.NewTicker(time.Second)
	sr := newSegRing(40)
	for {
		select {
		case b, ok := <-sf.input:
			if !ok {
				if sf.output != nil {
					close(sf.output)
				}
				return
			}

			sf.writeOutput(b)

			sr.add(b.to - b.from)
			break
		case <-timer.C:
			sr.add(0)
			if sf.sm != nil {
				sf.sm.SetSpeed(calcSpeed(&sr))
			}
		case <-sf.quit:
			log.Print("speed quit")
			return
		}
	}
}

func calcSpeed(sr *segRing) float64 {
	dur, length := sr.total()
	return math.Floor(float64(length)*float64(time.Second)/float64(dur)/1024.0 + 0.5)
}
