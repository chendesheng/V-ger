package download

import (
	"fmt"
	"time"
)

type simpleWriteFilter struct {
	basicFilter
	w  WriterAtQuit
	sm SpeedMonitor
}

func (swf *simpleWriteFilter) active() {
	defer swf.closeOutput()

	timer := time.NewTicker(time.Second)
	sr := newSegRing(40)
	for {
		select {
		case b, ok := <-swf.input:
			if !ok {
				fmt.Println("close simple write output")
				return
			}

			// trace(fmt.Sprint("simple write filter input:", b.from, b.to))

			swf.w.WriteAtQuit(b, swf.quit)

			swf.writeOutput(b)
			// trace(fmt.Sprint("simple write filter output:", b.from, b.to))

			sr.add(int64(len(b.Data)))
			break
		case <-timer.C:
			sr.add(0)
			if swf.sm != nil {
				swf.sm.SetSpeed(sr.calcSpeed())
			} else {
				println("speed monitor is nil")
			}
		case <-swf.quit:
			fmt.Println("simple write output quit")
			return
		}
	}

	fmt.Println("simpleWriteOutput end")
}
