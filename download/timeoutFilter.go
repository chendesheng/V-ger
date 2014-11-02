package download

import (
	"fmt"
	"log"
	"time"
)

type timeoutFilter struct {
	basicFilter
	timeout time.Duration
}

func (tf *timeoutFilter) active() {
	defer tf.closeOutput()

	timeout := tf.timeout

	timerQuit := time.NewTimer(timeout)
	for {
		select {
		case b, ok := <-tf.input:
			if !ok {
				log.Print("timeout filter output")
				return
			}
			tf.writeOutput(b)
			timerQuit.Reset(timeout)
			break
		case <-timerQuit.C:
			tf.closeQuit()
			return
		case <-tf.quit:
			fmt.Println("timeout filter quit")
			return

		}
	}
}
