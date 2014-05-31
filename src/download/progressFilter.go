package download

import (
	"fmt"
	"log"
	"task"
	"time"
)

type progressFilter struct {
	basicFilter
	t *task.Task
}

func (pf *progressFilter) active() {
	defer pf.closeOutput()

	t := pf.t

	size, downloaded := t.Size, t.DownloadedSize

	timer := time.NewTicker(time.Second)

	speed := float64(0)
	sr := newSegRing(40)

	for {
		select {
		case b, ok := <-pf.input:
			if !ok {
				fmt.Println("progress filter finish")
				saveProgress(t.Name, 0, 0, downloaded)
				return
			}

			length := b.to - b.from
			downloaded = b.to

			pf.writeOutput(b)

			sr.add(length)
			break
		case <-timer.C:
			sr.add(0)
			speed = sr.calcSpeed()

			est := calcEst(downloaded, size, speed)
			saveProgress(t.Name, speed, est, downloaded)
			break
		case <-pf.quit:
			fmt.Println("progress quit")
			return
		}
	}
}

func calcEst(downloaded, size int64, speed float64) (est time.Duration) {
	if speed == 0 {
		est = 0
	} else {
		est = time.Duration(float64((size-downloaded))/speed) * time.Millisecond
	}
	return
}

func saveProgress(name string, speed float64, est time.Duration, downloaded int64) {
	if t, err := task.GetTask(name); err == nil {
		t.Speed = speed
		t.Est = est
		t.DownloadedSize = downloaded

		if err := task.SaveTask(t); err != nil {
			log.Print(err)
			task.SaveTask(t)
		}
	}
}
