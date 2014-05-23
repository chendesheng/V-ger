package download

import "time"

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
		n = 0
	}

	r := make([]*segment, 0, n)
	for i := 0; i < n; i++ {
		r = append(r, &segment{0, 0})
	}
	return segRing{r, 0, time.Now()}
}
func (sr *segRing) increaseAdd(d time.Duration, l int64) {
	sr.i++
	if sr.i == len(sr.segs) {
		sr.i = 0
	}

	s := sr.segs[sr.i]
	s.d = d
	s.l = l
}
func (sr *segRing) add(l int64) {
	now := time.Now()
	d := now.Sub(sr.lastTime)
	sr.lastTime = now

	s := sr.segs[sr.i]
	prel := s.l

	sr.increaseAdd(d, l)

	if prel == 0 && l == 0 { //accelerate speed down to zero
		sr.increaseAdd(0, 0)
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
