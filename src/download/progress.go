package download

import (
	"math"
	"fmt"
	// "log"
	"task"
	"time"
)

type segment struct {
	t    time.Time
	size int64
}
type segRing struct {
	segs  []*segment
	i     int
	start time.Time
}

func newSegRing(n int) segRing {
	if n <= 0 {
		panic(fmt.Errorf("Init size must greater than zero."))
	}

	r := make([]*segment, 0)
	for i := 0; i < n; i++ {
		r = append(r, &segment{time.Now(), 0})
	}
	return segRing{r, 0, time.Now()}
}

func (sr *segRing) add(size int64) {
	sr.i++
	if sr.i == len(sr.segs) {
		sr.i = 0
	}

	s := sr.segs[sr.i]
	sr.start = s.t
	s.t = time.Now()
	s.size = size
}
func (sr *segRing) totalSize() int64 {
	total := int64(0)
	for _, s := range sr.segs {
		total += s.size
	}
	return total
}
func (sr *segRing) currentSegStart() time.Time {
	return sr.segs[sr.i].t
}

func handleProgress(progress chan int64, t *task.Task, quit <-chan bool) {
	size, total, elapsedTime := t.Size, t.DownloadedSize, t.ElapsedTime

	timer := time.NewTicker(time.Millisecond * 1500)

	speed := float64(0)
	part := int64(0)
	sr := newSegRing(12)

	for {
		select {
		case length, ok := <-progress:
			if !ok {
				saveProgress(t.Name, speed, total, elapsedTime, 0)
				return
			}
			total += length
			part += length
		case <-timer.C:
			elapsedTime += time.Second * 2

			sr.add(part)
			if part == 0 {
				sr.add(0) //accelerate speed down to zero
			} else {
				part = 0
			}

			sum := sr.totalSize()
			speed = math.Floor(float64(sum) * float64(time.Second) / float64(time.Since(sr.start)) / 1024 + 0.5)
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
