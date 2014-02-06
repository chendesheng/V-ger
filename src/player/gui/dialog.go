package gui

//#include "gui.h"
//#include <stdlib.h>
import "C"
import (
	"unsafe"
)

func NewDialog(title string, width, height int) {
	ctitle := C.CString(title)
	defer C.free(unsafe.Pointer(ctitle))
	// ptr := unsafe.Pointer(C.newDialog(ctitle, C.int(width), C.int(height)))
	// println(ptr)

	C.newDialog(ctitle, C.int(width), C.int(height))
}
