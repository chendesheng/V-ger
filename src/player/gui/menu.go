package gui

/*
#include "gui.h"
#include <stdlib.h>
#include <stdint.h>
*/
import "C"
import (
	"unsafe"
)

func (w *Window) InitAudioMenu(names []string, tags []int32, selected int) {
	if len(names) == 0 {
		return
	}

	cnames := make([]*C.char, 0)
	for _, name := range names {
		cnames = append(cnames, C.CString(name))
	}

	C.initAudioMenu(w.ptr, (**C.char)(&cnames[0]), (*C.int32_t)(unsafe.Pointer(&tags[0])), C.int(len(cnames)), C.int(selected))

	for _, cname := range cnames {
		C.free(unsafe.Pointer(cname))
	}
}

func (w *Window) InitSubtitleMenu(names []string, tags []int32, selected int) {
	if len(names) == 0 {
		return
	}

	cnames := make([]*C.char, 0)
	for _, name := range names {
		cnames = append(cnames, C.CString(name))
	}

	C.initSubtitleMenu(w.ptr, (**C.char)(&cnames[0]), (*C.int32_t)(unsafe.Pointer(&tags[0])), C.int(len(cnames)), C.int(selected))

	for _, cname := range cnames {
		C.free(unsafe.Pointer(cname))
	}
}
