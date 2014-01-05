package clock

import (
	"fmt"
	"log"
	. "player/shared"
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

	seeking bool
}

//will be blocked if clock is paused
func (c *Clock) GetTime() time.Duration {
	// c.waitUntilRunning()

	return c.getTime()
}

func (c *Clock) CalcTime(percent float64) time.Duration {
	t := time.Duration(float64(c.totalTime) * percent)
	return t
}

func (c *Clock) CalcPlayProgress(percent float64) *PlayProgressInfo {
	t := c.CalcTime(percent)
	leftT := c.totalTime - t

	return &PlayProgressInfo{formatTime(t), formatTime(leftT), percent, 0}
}

func (c *Clock) GetSeekTime() time.Duration {
	return c.getTime()
}

func formatTime(t time.Duration) string {
	sign := ""
	if t < 0 {
		t = -t
		sign = "-"
	}

	h := int(t.Hours())
	m := int(t.Minutes()) % 60
	s := int(t.Seconds()) % 60

	if h == 0 {
		return fmt.Sprintf("%s%02d:%02d", sign, m, s)
	}

	return fmt.Sprintf("%s%02d:%02d:%02d", sign, h, m, s)
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
	// log.Println("clock set time:", t.String())
	c.Lock()
	defer c.Unlock()

	if c.status == "paused" {
		c.pausedTime = t
	}

	c.base = time.Now().Add(-t)
}

func (c *Clock) AddTime(d time.Duration) {
	c.Lock()
	defer c.Unlock()
	println("clock base:", c.base.String())
	c.base = c.base.Add(-d)
}

func (c *Clock) Pause() {
	c.Lock()
	defer c.Unlock()

	c.pause()
}

func (c *Clock) pause() {
	c.status = "paused"
	c.pausedTime = time.Since(c.base)
	// c.wait = make(chan bool)
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
	log.Print("clock resume:", c.getTime().String())

	c.Lock()
	defer c.Unlock()

	if c.status == "paused" {
		c.resume()
	}
}

func (c *Clock) ResumeWithTime(t time.Duration) {
	c.Lock()
	defer c.Unlock()

	if c.status == "paused" {
		c.status = "running"
		println("clock running")
		c.base = time.Now().Add(-t)

		close(c.wait)
		c.wait = make(chan bool)
		println("close wait")
	}
}

func (c *Clock) resume() {
	c.base = time.Now().Add(-c.pausedTime)
	c.status = "running"

	close(c.wait)
	c.wait = make(chan bool)
}

func (c *Clock) Reset() {
	c.Lock()
	defer c.Unlock()

	c.base = time.Now()
	c.pausedTime = 0
	c.status = "running"
}

func (c *Clock) waitUntilRunning() bool {
	var ch chan bool
	var status string
	c.Lock()
	ch = c.wait
	status = c.status
	c.Unlock()

	if status == "paused" {
		<-ch
		println("after paused")

		return true
	}

	return false
}

func (c *Clock) After(d time.Duration) {
	// b := c.GetTime()
	// c.waitUntilRunning()

	if d > time.Second {
		log.Print("clock wait long time:", d.String())
	}
	<-time.After(d) //time.After is not very accuracy which about one millisecond delay while wait one second

	// c.waitUntilRunning()
	// c.SetTime(b + d)
}

func (c *Clock) WaitUtilRunning() bool {
	return c.waitUntilRunning()
}

type beforeWait func()

func (c *Clock) WaitUtilRunning2(fn beforeWait) bool {
	var ch chan bool
	var status string
	c.Lock()
	ch = c.wait
	status = c.status
	c.Unlock()

	if status == "paused" {
		fn()
		<-ch
		println("after paused")

		return true
	}

	return false
}

func (c *Clock) WaitUtil(t time.Duration) {
	now := c.getTime()

	if t > now {
		c.After(t - now)
	}
}

func (c *Clock) AfterWithQuit(d time.Duration, quit chan bool) bool {
	// c.waitUntilRunning()

	select {
	case <-time.After(d):
		break
	case <-quit:
		return true
	}

	// c.waitUntilRunning()
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

func (c *Clock) TotalTime() time.Duration {
	return c.totalTime
}

func NewClock(totalTime time.Duration) *Clock {
	now := time.Now()
	c := &Clock{
		base:       now,
		pausedTime: 0,
		status:     "running",
		totalTime:  totalTime,
		wait:       make(chan bool),
	}
	c.Mutex = sync.Mutex{}
	return c
}
