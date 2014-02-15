package gui

import (
	"testing"
)

func TestShowWindow(t *testing.T) {
	NewWindow("hello", 1280, 720)

	PollEvents()
}
func fequal(a, b float64) bool {
	d := a - b
	return d < 1e5 && d > -1e5
}
func TestFitToWindowLarger(t *testing.T) {
	w := NewWindow("title", 1280, 720)
	x, y, width, height := w.fitToWindow(1280, 700)
	println(x, y, width, height)
	if !fequal(float64(width)/float64(height), 1280.0/700.0) {
		t.Error("Aspect radio should be equal")
	}
	if width < 1280 || height < 700 {
		t.Error("Must larger than before")
	}
}

func TestFitToWindowSmaller(t *testing.T) {
	w := NewWindow("title", 1280, 720)
	x, y, width, height := w.fitToWindow(1280, 300)
	println(x, y, width, height)
	if !fequal(float64(width)/float64(height), 1280.0/300.0) {
		t.Error("Aspect radio should be equal")
	}
	if width > 1280 || height > 300 {
		t.Error("Must smaller than before")
	}
}
