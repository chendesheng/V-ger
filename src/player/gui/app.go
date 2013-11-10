package gui

//#include "gui.h"
import "C"

func PollEvents() {
	C.pollEvents()
}
