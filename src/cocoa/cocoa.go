package cocoa

// #cgo darwin LDFLAGS: -framework gococoa
// #include "gococoa.h"
// #include "stdlib.h"
import "C"
import (
// "runtime"
// "unsafe"
)

func NSAppRun() {
	C.NSAppRun()
}

func NSAppStop() {
	C.NSAppStop()
}

// func SetTitle(title string) {
// 	ctitle := C.CString(title)
// 	defer C.free(unsafe.Pointer(ctitle))
// 	C.SetTitle(ctitle)
// }

// func SendNotification() {
// 	// ctitle := C.CString(title)
// 	// defer C.free(unsafe.Pointer(ctitle))

// 	// cmessage := C.CString(message)
// 	// defer C.free(unsafe.Pointer(cmessage))

// 	C.SendNotification()
// }
