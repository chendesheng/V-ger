package clock

import (
	"testing"
	"time"
)

func TestClock(t *testing.T) {
	c := NewClock(time.Hour)
	if c == nil {
		t.Error("new clock fail")
	}
	<-time.After(time.Second)
	d := c.GetTime()
	diff := d - time.Second
	if absAndCheck(diff) { //less than 0.2 percent
		t.Error("accurecy fail, diff: ", diff.String())
	}
}

func absAndCheck(d time.Duration) bool {
	return d > 10*time.Millisecond || d < -10*time.Millisecond
}

func TestClockPause(t *testing.T) {
	c := NewClock(time.Hour)
	if c == nil {
		t.Error("new clock fail")
	}
	<-time.After(time.Second)
	d := c.GetTime()
	c.Pause()
	<-time.After(time.Second)

	diff := d - time.Second
	if absAndCheck(diff) {
		t.Error("paused clock should not accumulate time, diff: ", diff.String())
	}
	go func() {
		<-time.After(time.Second)
		c.Resume()
	}()

	c.GetTime() //blocked

	<-time.After(time.Second)
	diff = c.GetTime() - 2*time.Second
	if absAndCheck(diff) {
		t.Error("diff should be 2s but ", diff.String())
	}

	b := time.Now()
	c.GetTime()
	if time.Since(b) > time.Millisecond {
		t.Error("should not longer than one millisecond.")
	}
}

func TestSetTime(t *testing.T) {
	c := NewClock(time.Hour)
	c.SetTime(time.Second)
	diff := c.GetTime() - time.Second
	if absAndCheck(diff) {
		t.Error("wrong current time:", diff.String())
	}
}
