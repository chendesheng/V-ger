package clock

import (
	"fmt"
	"log"
	"sync"
	"time"
	. "vger/player/shared"
)

type Clock struct {
	sync.Mutex

	base       time.Time
	pausedTime time.Duration
	running    bool

	wait chan struct{}

	totalTime time.Duration

	seeking bool
}

var chClosed chan struct{}

func init() {
	chClosed = make(chan struct{})
	close(chClosed)
}

func (c *Clock) CalcTime(percent float64) time.Duration {
	t := time.Duration(float64(c.totalTime) * percent)
	return t
}

func (c *Clock) CalcPlayProgress(t time.Duration) *PlayProgressInfo {
	percent := float64(t) / float64(c.totalTime)
	leftT := c.totalTime - t

	return &PlayProgressInfo{c.formatTime(t), "-" + c.formatTime(leftT), percent}
}

func (c *Clock) formatTime(t time.Duration) string {
	sign := ""
	if t < 0 {
		t = -t
		sign = "-"
	}

	h := int(t.Hours())
	m := int(t.Minutes()) % 60
	s := int(t.Seconds()) % 60

	if c.totalTime < time.Hour {
		return fmt.Sprintf("%s%02d:%02d", sign, m, s)
	} else {
		return fmt.Sprintf("%s%02d:%02d:%02d", sign, h, m, s)
	}
}

func (c *Clock) GetPercent() float64 {
	return float64(c.GetTime()) / float64(c.totalTime)
}

func (c *Clock) GetTime() time.Duration {
	c.Lock()
	defer c.Unlock()

	if !c.running {
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
	c.Lock()
	defer c.Unlock()

	if !c.running {
		c.pausedTime = t
	}

	c.base = time.Now().Add(-t)
}
func (c *Clock) AddTime(d time.Duration) {
	c.Lock()
	defer c.Unlock()

	if !c.running {
		c.pausedTime += d
	}

	base := c.base.Add(-d)
	if base.After(time.Now()) {
		base = time.Now()
	}

	c.base = base
}

func (c *Clock) Pause() {
	c.Lock()
	defer c.Unlock()

	c.pause()
}

func (c *Clock) IsRunning() bool {
	c.Lock()
	defer c.Unlock()

	return c.running
}

func (c *Clock) pause() {
	if c.running {
		c.running = false
		c.pausedTime = time.Since(c.base)
	}
	// c.wait = make(chan struct{})
}

func (c *Clock) Toggle() {
	c.Lock()
	defer c.Unlock()

	if c.running {
		c.pause()
	} else {
		c.resume()
	}
}

func (c *Clock) Resume() {
	c.Lock()
	defer c.Unlock()

	if !c.running {
		c.resume()
	}
}

func (c *Clock) resume() {
	c.base = time.Now().Add(-c.pausedTime)
	c.running = true

	close(c.wait)
	c.wait = make(chan struct{})
}

func (c *Clock) Reset() {
	c.Lock()
	defer c.Unlock()
	// log.Print(c)
	c.base = time.Now()
	c.pausedTime = 0
	c.running = false
}

func (c *Clock) WaitUntilRunning(quit chan struct{}) bool {
	var ch chan struct{}
	c.Lock()
	ch = c.wait
	c.Unlock()

	if !c.running {
		select {
		case <-ch:
			log.Print("after paused")
			return false
		case <-quit:
			return true
		}
	}

	return false
}

func (c *Clock) WaitUntilWithQuit(t time.Duration, quit chan struct{}) bool {
	now := c.GetTime()
	d := t - now
	if d > 0 {
		select {
		case <-time.After(d):
			return false
		case <-quit:
			return true
		}
	}

	return false
}

func (c *Clock) WaitUntil(t time.Duration) <-chan time.Time {
	return time.After(t - c.GetTime())
}
func (c *Clock) WaitRunning() chan struct{} {
	var ch chan struct{}
	c.Lock()
	ch = c.wait
	c.Unlock()

	if !c.running {
		return ch
	} else {
		return chClosed
	}
}

func (c *Clock) TotalTime() time.Duration {
	return c.totalTime
}

func New(totalTime time.Duration) *Clock {
	log.Print("NewClock:", totalTime.String())

	now := time.Now()
	c := &Clock{
		base:       now,
		pausedTime: 0,
		running:    true,
		totalTime:  totalTime,
		wait:       make(chan struct{}),
	}
	return c
}
