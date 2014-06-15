package download

import (
	"block"
	"fmt"
	"io"
	"log"
	"native"
	"os"
	"path/filepath"
	"time"
)

type writeFilter struct {
	basicFilter
	taskName       string
	w              io.WriterAt
	restartTimeout time.Duration
}

func (wf *writeFilter) active() {
	defer wf.closeOutput()

	restartTimeout := wf.restartTimeout

	if restartTimeout == 0 {
		restartTimeout = time.Duration(1<<63 - 1) //max duration
	}

	timerRestart := time.NewTimer(restartTimeout)
	for {
		select {
		case b, ok := <-wf.input:
			if !ok {
				log.Print("close write output")
				return
			}
			timerRestart.Reset(restartTimeout)

			err := wf.mustWrite(b)
			if err != nil {
				return
			}
			break
		case <-timerRestart.C:
			wf.closeQuit()
			return
		case <-wf.quit:
			fmt.Println("write output quit")
			return

		}
	}

	log.Print("writeOutput end")
}

func (wf *writeFilter) mustWrite(b block.Block) error {
	pathErrNotifyTimes := 0
	for {
		_, err := wf.w.WriteAt(b.Data, b.From)

		if err == nil {
			pathErrNotifyTimes = 0

			wf.writeOutput(b)
			break
		} else if perr, ok := err.(*os.PathError); ok {
			if pathErrNotifyTimes == 0 { //only report once
				native.SendNotification("Error write "+filepath.Base(perr.Path), perr.Err.Error())
			}
			if pathErrNotifyTimes%300 == 0 { //write error log every 10 mins
				log.Print(err)
			}
			pathErrNotifyTimes++

			wf.wait(2 * time.Second)
		} else {
			log.Print(err)

			wf.closeQuit()
			return err
		}
	}

	return nil
}
