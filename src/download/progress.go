package download

import (
	"fmt"
	// "log"
	"math"
	"task"
	"time"
)

type segment struct {
	d time.Duration //takes 'd' time, download 'l' byte
	l int64
}

type segRing struct {
	segs     []*segment
	i        int
	lastTime time.Time
}

func newSegRing(n int) segRing {
	if n <= 0 {
		panic(fmt.Errorf("Init length must greater than zero."))
	}

	r := make([]*segment, 0)
	for i := 0; i < n; i++ {
		r = append(r, &segment{0, 0})
	}
	return segRing{r, 0, time.Now()}
}
func (sr *segRing) increase() {
	sr.i++
	if sr.i == len(sr.segs) {
		sr.i = 0
	}
}
func (sr *segRing) add(l int64) {
	d := time.Now().Sub(sr.lastTime)
	sr.lastTime = time.Now()

	s := sr.segs[sr.i]
	prel := s.l

	sr.increase()

	s = sr.segs[sr.i]
	s.d = d
	s.l = l

	if prel == 0 && l == 0 { //accelerate speed down to zero
		sr.increase()
		s = sr.segs[sr.i]
		s.d = 0
		s.l = 0
	}
}
func (sr *segRing) total() (time.Duration, int64) {
	var totalLength int64
	var totalDurtion time.Duration
	for _, s := range sr.segs {
		totalDurtion += s.d
		totalLength += s.l
	}
	return totalDurtion, totalLength
}

func handleProgress(progress chan *block, output chan *block, t *task.Task, quit <-chan bool) {
	size, total, elapsedTime := t.Size, t.DownloadedSize, t.ElapsedTime

	timer := time.NewTicker(time.Millisecond * 1000)

	speed := float64(0)
	part := int64(0)
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
			part += length

			if output != nil {
				go func() {
					select {
					case output <- b:
					case <-quit:
						return
					}
				}()
			}
			sr.add(length)

		case <-timer.C:
			elapsedTime += time.Millisecond * 1000

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
		}
		t.ElapsedTime = elapsedTime
		t.Speed = speed
		t.Est = est
		task.SaveTask(t)
	}
}
