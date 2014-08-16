package subtitles

import (
	"testing"
	"time"
)

func TestAddic7edGet(t *testing.T) {
	a := &addic7ed{"rake s01e11", nil}
	result := make(chan Subtitle)
	go a.search(result)

	select {
	case s := <-result:
		if s.Context == nil && len(s.Context.Get("Referrer")) > 0 {
			t.Error("s.Content should not be nil")
		}
		println(s.Context.Get("Referrer"))
		println(s.URL)
		println(s.Description)
	case <-time.After(10 * time.Second):
		t.Error("timeout")
	}
}

func TestAddic7edNoResult(t *testing.T) {
	a := &addic7ed{"rake s01e22", nil}
	result := make(chan Subtitle)
	go func() {
		a.search(result)
		close(result)
	}()

	select {
	case _, ok := <-result:
		if ok {
			t.Error("should not have result")
		}
	case <-time.After(10 * time.Second):
		t.Error("timeout")
	}
}
