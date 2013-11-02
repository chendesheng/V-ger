package clock

import (
	"fmt"
	"sync"
	"time"
)

type Clock struct {
	sync.Mutex

	base       time.Time
	pausedTime time.Duration
	status     string

	wait chan bool

	totalTime time.Duration
}

//will be blocked if clock is paused
func (c *Clock) GetTime() time.Duration {
	c.waitUntilRunning()

	return c.getTime()
}

func addZero(i int) string {
	if i < 10 {
		return fmt.Sprint("0", i)
	} else {
		return fmt.Sprint(i)
	}
}

func (c *Clock) GetLeftTimeString() string {
	left := c.totalTime - c.GetTime()

	h := addZero(int(left.Hours()))
	m := addZero(int(left.Minutes()) % 60)
	s := addZero(int(left.Seconds()) % 60)

	if h == "00" {
		return fmt.Sprintf("-%s:%s", m, s)
	}

	return fmt.Sprintf("-%s:%s:%s", h, m, s)
}

func (c *Clock) GetTimeString() string {
	time := c.GetTime()

	h := addZero(int(time.Hours()))
	m := addZero(int(time.Minutes()) % 60)
	s := addZero(int(time.Seconds()) % 60)

	if h == "00" {
		return fmt.Sprintf("%s:%s", m, s)
	}

	return fmt.Sprintf("%s:%s:%s", h, m, s)
}

func (c *Clock) GetPercent() float64 {
	return float64(c.GetTime()) / float64(c.totalTime)
}

func (c *Clock) getTime() time.Duration {
	c.Lock()
	defer c.Unlock()

	if c.status == "paused" {
		return c.pausedTime
	} else {
		t := time.Since(c.base)
		if t > c.totalTime {
			t = c.totalTime
		}
		return t
	}
}

func (c *Clock) SetTime(t time.Duration) {
	// println("clock SetTime ", t.String())

	c.Lock()
	defer c.Unlock()

	c.base = time.Now().Add(-t)
}

func (c *Clock) Pause() {
	c.Lock()
	defer c.Unlock()

	c.pause()
}

func (c *Clock) pause() {
	c.status = "paused"
	c.pausedTime = time.Since(c.base)
	c.wait = make(chan bool)
}

func (c *Clock) Toggle() {
	c.Lock()
	defer c.Unlock()

	if c.status == "running" {
		c.pause()
	} else {
		c.resume()
	}
}

func (c *Clock) Resume() {
	c.Lock()
	defer c.Unlock()

	if c.status == "paused" {
		c.resume()
	}
}

func (c *Clock) resume() {
	c.base = time.Now().Add(-c.pausedTime)
	c.status = "running"
	close(c.wait)
}

func (c *Clock) Reset() {
	c.Lock()
	defer c.Unlock()

	c.base = time.Now()
	c.pausedTime = 0
	c.status = "running"
}

func (c *Clock) waitUntilRunning() {
	var ch chan bool
	var status string
	c.Lock()
	ch = c.wait
	status = c.status
	c.Unlock()

	if status == "paused" {
		<-ch
	}
}

func (c *Clock) After(d time.Duration) {
	// b := c.GetTime()
	c.waitUntilRunning()

	<-time.After(d) //time.After is not very accuracy which about one millisecond delay while wait one second

	c.waitUntilRunning()
	// c.SetTime(b + d)
}

func (c *Clock) WaitUtilRunning() {
	c.waitUntilRunning()
}

func (c *Clock) WaitUtil(t time.Duration) {
	now := c.getTime()

	if t > now {
		c.After(t - now)
	}
}

func (c *Clock) AfterWithQuit(d time.Duration, quit chan bool) bool {
	c.waitUntilRunning()

	select {
	case <-time.After(d):
		break
	case <-quit:
		return true
	}

	c.waitUntilRunning()
	return false
}

func (c *Clock) WaitUtilWithQuit(t time.Duration, quit chan bool) bool {
	// println("wait until", t.String())
	now := c.getTime()

	if t > now {
		return c.AfterWithQuit(t-now, quit)
	}

	return false
}

func NewClock(totalTime time.Duration) *Clock {
	now := time.Now()
	c := &Clock{
		base:       now,
		pausedTime: 0,
		status:     "running",
		totalTime:  totalTime,
	}
	c.Mutex = sync.Mutex{}
	return c
}
