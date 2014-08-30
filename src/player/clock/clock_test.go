package clock

import (
	"math"
	"math/rand"
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

func TestPauseResume(t *testing.T) {
	c := NewClock(time.Hour)
	if c == nil {
		t.Error("new clock fail")
	}
	<-time.After(time.Second)
	d := c.GetTime()
	c.Pause()
	<-time.After(time.Second)

	diff := d - c.GetTime()
	if absAndCheck(diff) {
		t.Error("paused clock should not accumulate time, diff: ", diff.String())
	}

	c.Resume()
	d = c.GetTime()
	time.Sleep(time.Second)
	diff = c.GetTime() - d - time.Second
	if absAndCheck(diff) {
		t.Error("resume not working")
	}
}

func TestCalcPos(t *testing.T) {
	c := NewClock(time.Hour)
	p := c.CalcPlayProgress(time.Minute)

	if p.Left != "00:01:00" {
		t.Error("p.Left should be 00:01:00 but %s", p.Left)
	}

	if p.Right != "-00:59:00" {
		t.Error("p.Right should be -00:59:00 but %s", p.Right)
	}

	if math.Abs(p.Percent-1.0/60.0) > 1e-5 {
		t.Errorf("p.Left should 0.016666667 but %f", p.Percent)
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
func TestAddTime(t *testing.T) {
	c := NewClock(time.Hour)
	c.SetTime(10 * time.Second)
	c.AddTime(time.Second)
	diff := c.GetTime() - 11*time.Second
	if absAndCheck(diff) {
		t.Error("wrong current time:", diff.String())
	}
}

func TestWaitUntilRunning(t *testing.T) {
	c := NewClock(time.Hour)
	go func() {
		if c.WaitUntilRunning(nil) {
			t.Error("no quit")
		}
	}()

	d := time.Duration(rand.Int63n(int64(time.Second))) + time.Second
	c.Pause()
	go func() {
		ti := time.Now()
		c.WaitUntilRunning(nil)
		diff := time.Now().Sub(ti) - d
		if absAndCheck(diff) {
			t.Errorf("wrong diff time: %s", diff.String())
		}
	}()
	<-time.After(d)
	c.Resume()

	<-time.After(time.Second)
}

func TestWaitUntil(t *testing.T) {
	c := NewClock(time.Hour)
	d := time.Duration(rand.Int63n(int64(time.Second))) + time.Second

	ti := time.Now()
	c.WaitUntilWithQuit(d, nil)
	diff := time.Now().Sub(ti) - d
	if absAndCheck(diff) {
		t.Errorf("wrong diff time: %s", diff.String())
	}
}
