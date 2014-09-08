package gui

/*
#include "gui.h"
#include <stdlib.h>
#include <stdint.h>
*/
import "C"
import (
	"log"
	"unsafe"
)

func HideSubtitleMenu() {
	C.hideSubtitleMenu()
}

func HideAudioMenu() {
	C.hideAudioMenu()
}

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

func SetSubtitleMenuItem(t1, t2 int) {
	C.setSubtitleMenuItem(C.int(t1), C.int(t2))
}

func (w *Window) InitSubtitleMenu(names []string, tags []int32, selected1 int, selected2 int) {
	if len(names) == 0 {
		return
	}

	cnames := make([]*C.char, 0)
	for _, name := range names {
		cnames = append(cnames, C.CString(name))
	}

	log.Printf("selected1:%d, selected2:%d", selected1, selected2)

	C.initSubtitleMenu(w.ptr, (**C.char)(&cnames[0]), (*C.int32_t)(unsafe.Pointer(&tags[0])), C.int(len(cnames)), C.int32_t(selected1), C.int32_t(selected2))

	for _, cname := range cnames {
		C.free(unsafe.Pointer(cname))
	}
}
