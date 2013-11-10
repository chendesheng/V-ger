package main

import (
	. "player/gui"
)

func main() {
	w := NewWindow("hello", 1280, 720)
	w.Show()

	PollEvents()
}