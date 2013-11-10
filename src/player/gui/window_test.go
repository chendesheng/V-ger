package gui

import (
	"testing"
)

func TestShowWindow(t *testing.T) {
	w := NewWindow("hello", 1280, 720)
	w.Show()

	PollEvents()
}
