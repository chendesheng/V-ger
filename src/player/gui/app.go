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

//export goOnWillTerminate
func goOnWillTerminate() {
	if appDelegate != nil {
		appDelegate.WillTerminate()
	}
}

//export goOnSearchSubtitleMenuItemClick
func goOnSearchSubtitleMenuItemClick() {
	if appDelegate != nil {
		appDelegate.SearchSubtitleMenuItemClick()
	}
}

//export goOnOpenOpenPanel
func goOnOpenOpenPanel() {
	if appDelegate != nil {
		appDelegate.OnOpenOpenPanel()
	}
}

//export goOnCloseOpenPanel
func goOnCloseOpenPanel(filename *_Ctype_char) {
	if appDelegate != nil {
		name := C.GoString(filename)
		appDelegate.OnCloseOpenPanel(name)
	}
}

type AppDelegate interface {
	OpenFile(string) bool
	OnOpenOpenPanel()
	OnCloseOpenPanel(filename string)
	WillTerminate()
	SearchSubtitleMenuItemClick()
}

var appDelegate AppDelegate

func Initialize(d AppDelegate) {
	appDelegate = d

	log.Println("before initialize")
	C.initialize()
}

func GetScreenSize() (int, int) {
	sz := C.getScreenSize()
	return int(sz.width), int(sz.height)
}
