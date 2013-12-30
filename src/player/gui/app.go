package gui

//#include "gui.h"
import "C"
import (
	"log"
	"unsafe"
)

func PollEvents() {
	C.pollEvents()
}

//export goOnOpenFile
func goOnOpenFile(cfilename unsafe.Pointer) C.int {
	filename := C.GoString((*C.char)(cfilename))
	log.Println("goOnOpenFile:", filename)

	if appDelegate != nil {
		if appDelegate.OpenFile(filename) {
			return 1
		} else {
			return 0
		}
	} else {
		return 0
	}
}

type AppDelegate interface {
	OpenFile(string) bool
}

var appDelegate AppDelegate

func Initialize(d AppDelegate) {
	appDelegate = d

	log.Println("before initialize")
	C.initialize()
}
