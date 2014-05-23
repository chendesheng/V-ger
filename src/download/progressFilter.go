package download

import (
	"fmt"
	"log"
	"math"
	"task"
	"time"
)

type progressFilter struct {
	basicFilter
	t *task.Task
}

func (pf *progressFilter) active() {
	handleProgress(pf.input, pf.output, pf.t, pf.quit)
}

func handleProgress(progress chan *block, output chan *block, t *task.Task, quit chan bool) {
	size, total, elapsedTime := t.Size, t.DownloadedSize, t.ElapsedTime
	if t.Status == "Playing" {
		total = t.BufferedPosition
	}

	timer := time.NewTicker(time.Second)

	speed := float64(0)
	sr := newSegRing(40)

	for {
		select {
		case b, ok := <-progress:
			if !ok {
				saveProgress(t.Name, 0, total, elapsedTime, 0)
				if output != nil {
					close(output)
				}
				return
			}

			length := b.to - b.from
			total += length

			if output != nil {
				select {
				case output <- b:
				case <-quit:
					return
				}
			}
			sr.add(length)
			break
		case <-timer.C:
			elapsedTime += time.Second

			sr.add(0)

			totalDurtion, totalLength := sr.total()
			// fmt.Printf("totalDurtion %s, totalLength %d\n", totalDurtion, totalLength)
			speed = math.Floor(float64(totalLength)*float64(time.Second)/float64(totalDurtion)/1024.0 + 0.5)
			_, est := calcProgress(total, size, speed)
			saveProgress(t.Name, speed, total, elapsedTime, est)

			if total == size {
				fmt.Println("progress return")
				return
			}
		case <-quit:
			fmt.Println("progress quit")
			return
		}
	}

}
func calcProgress(total, size int64, speed float64) (percentage float64, est time.Duration) {
	percentage = float64(total) / float64(size) * 100
	if speed == 0 {
		est = 0
	} else {
		est = time.Duration(float64((size-total))/speed) * time.Millisecond
	}
	return
}
func saveProgress(name string, speed float64, total int64, elapsedTime time.Duration, est time.Duration) {
	if t, err := task.GetTask(name); err == nil {
		if t.Status != "Playing" {
			t.DownloadedSize = total
			t.BufferedPosition = total
		} else {
			t.BufferedPosition = total
		}

		t.ElapsedTime = elapsedTime
		t.Speed = speed
		t.Est = est

		if err := task.SaveTask(t); err != nil {
			log.Print(err)
			task.SaveTask(t)
		}
	}
}
