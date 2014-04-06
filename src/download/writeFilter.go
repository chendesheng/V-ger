package download

import (
	"fmt"
	"log"
	"native"
	"os"
	"path/filepath"
	"time"
)

type writeFilter struct {
	basicFilter
	w              WriterAtQuit
	restartTimeout time.Duration
}

func (wf *writeFilter) active() {
	writeOutput(wf.w, wf.input, wf.output, wf.restartTimeout, wf.quit)
}

func writeOutput(w WriterAtQuit, input <-chan *block, output chan *block, restartTimeout time.Duration, quit chan bool) {
	pathErrNotifyTimes := 0

	if restartTimeout == 0 {
		restartTimeout = time.Duration(1<<63 - 1) //max duration
	}

	timerRestart := time.NewTimer(restartTimeout)

	for {
		select {
		case b, ok := <-input:
			if !ok {
				fmt.Println("close write output")
				close(output)
				return
			}
			for {
				// println("writeAt:", b.from, len(b.data))
				if (b.to - b.from) != int64(len(b.data)) {
					log.Printf("wrong block:%d,%d,%d", b.from, b.to, len(b.data))
				}

				err := w.WriteAtQuit(b.data, b.from, quit)

				if err == nil {
					pathErrNotifyTimes = 0

					select {
					case output <- b:
						break
					case <-quit:
						return
					}
					break
				} else if perr, ok := err.(*os.PathError); ok {
					if pathErrNotifyTimes == 0 { //only report once
						native.SendNotification("Error write "+filepath.Base(perr.Path), perr.Err.Error())
					}
					if pathErrNotifyTimes%300 == 0 { //write error log every 10 mins
						log.Print(err)
					}
					pathErrNotifyTimes++

					select {
					case <-quit:
						return
					case <-time.After(time.Second * 2):
						break
					}
				} else {
					log.Print(err)
					ensureQuit(quit)
					return
				}
			}
			timerRestart.Reset(restartTimeout)
			break
		case <-timerRestart.C:
			log.Print("close quit after timeout")
			close(quit)
			return
		case <-quit:
			fmt.Println("write output quit")
			return

		}
	}

	fmt.Println("writeOutput end")
}
