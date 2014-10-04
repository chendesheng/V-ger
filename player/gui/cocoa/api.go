package cocoa

// #include "gui.h"
// #include <stdlib.h>
import "C"
import (
	"log"
	"unsafe"
)

func Run() {
	C.initialize()
	C.pollEvents()
}

func GetScreenSize() (int, int) {
	sz := C.getScreenSize()
	return int(sz.width), int(sz.height)
}

type NativeWindow uintptr

func (NativeWindow) HideSubtitleMenu() {
	C.hideSubtitleMenu()
}

func (NativeWindow) HideAudioMenu() {
	C.hideAudioMenu()
}

func (NativeWindow) SetSubtitleMenuItem(t1, t2 int) {
	C.setSubtitleMenuItem(C.int(t1), C.int(t2))
}

func (w NativeWindow) InitAudioMenu(names []string, tags []int32, selected int) {
	if len(names) == 0 {
		return
	}

	cnames := make([]*C.char, 0)
	for _, name := range names {
		cnames = append(cnames, C.CString(name))
	}

	C.initAudioMenu(unsafe.Pointer(w), (**C.char)(&cnames[0]), (*C.int32_t)(unsafe.Pointer(&tags[0])), C.int(len(cnames)), C.int(selected))

	for _, cname := range cnames {
		C.free(unsafe.Pointer(cname))
	}
}

func (w NativeWindow) InitSubtitleMenu(names []string, tags []int32, selected1 int, selected2 int) {
	if len(names) == 0 {
		return
	}

	cnames := make([]*C.char, 0)
	for _, name := range names {
		cnames = append(cnames, C.CString(name))
	}

	log.Printf("selected1:%d, selected2:%d", selected1, selected2)

	C.initSubtitleMenu(unsafe.Pointer(w), (**C.char)(&cnames[0]), (*C.int32_t)(unsafe.Pointer(&tags[0])), C.int(len(cnames)), C.int32_t(selected1), C.int32_t(selected2))

	for _, cname := range cnames {
		C.free(unsafe.Pointer(cname))
	}
}

func (w NativeWindow) RefreshContent() {
	C.refreshWindowContent(unsafe.Pointer(w))
}

func (w NativeWindow) GetSize() (int, int) {
	return int(C.getWindowWidth(unsafe.Pointer(w))), int(C.getWindowHeight(unsafe.Pointer(w)))
}

func (w NativeWindow) SetTitle(title string) {
	ctitle := C.CString(title)
	defer C.free(unsafe.Pointer(ctitle))

	C.setWindowTitle(unsafe.Pointer(w), ctitle)
}

func (w NativeWindow) SetSize(width, height int) {
	C.setWindowSize(unsafe.Pointer(w), C.int(width), C.int(height))
}

func NewWindow(title string, width, height int) NativeWindow {
	ctitle := C.CString(title)
	defer C.free(unsafe.Pointer(ctitle))

	return NativeWindow(unsafe.Pointer(C.newWindow(ctitle, C.int(width), C.int(height))))
}

func (w NativeWindow) Show() {
	C.showWindow(unsafe.Pointer(w))
}

func (w NativeWindow) MakeCurrentContext() {
	C.makeWindowCurrentContext(unsafe.Pointer(w))
}

func (w NativeWindow) ToggleFullScreen() {
	C.windowToggleFullScreen(unsafe.Pointer(w))
}

func (w NativeWindow) HideStartupView() {
	C.windowHideStartupView(unsafe.Pointer(w))
}
func (w NativeWindow) ShowStartupView() {
	C.windowShowStartupView(unsafe.Pointer(w))
}
func (w NativeWindow) ShowProgress(left, right string, percent float64) {
	cleft := C.CString(left)
	defer C.free(unsafe.Pointer(cleft))

	cright := C.CString(right)
	defer C.free(unsafe.Pointer(cright))

	C.showWindowProgress(unsafe.Pointer(w), cleft, cright, C.double(percent))
}
func (w NativeWindow) ShowBufferInfo(speed string, percent float64) {
	cspeed := C.CString(speed)
	defer C.free(unsafe.Pointer(cspeed))

	C.showWindowBufferInfo(unsafe.Pointer(w), cspeed, C.double(percent))
}

func (w NativeWindow) ShowText(items []struct {
	Content string
	Style   int //0 -normal, 1 -italic, 2 -bold, 3 italic and bold
	Color   uint
}, posType int, x, y float64) uintptr {

	if len(items) == 0 {
		return 0
	}

	ctexts := make([]C.AttributedString, 0)
	for _, str := range items {
		cstr := C.CString(str.Content)
		defer C.free(unsafe.Pointer(cstr))

		ctexts = append(ctexts, C.AttributedString{cstr, C.int(str.Style), C.uint(str.Color)})
	}

	citem := &C.SubItem{&ctexts[0], C.int(len(ctexts)), C.int(posType), C.double(x), C.double(y)}

	return uintptr(C.showText(unsafe.Pointer(w), citem))
}

func (w NativeWindow) HideText(ptr uintptr) {
	C.hideText(unsafe.Pointer(w), unsafe.Pointer(ptr))
}

func (w NativeWindow) HideCursor() {
	C.hideCursor(unsafe.Pointer(w))
}
func (w NativeWindow) ShowCursor() {
	C.showCursor(unsafe.Pointer(w))
}

func (w NativeWindow) ShowSpinning() {
	C.showSpinning(unsafe.Pointer(w))
}
func (w NativeWindow) HideSpinning() {
	C.hideSpinning(unsafe.Pointer(w))
}

func (w NativeWindow) SetVolume(volume byte) {
	C.setVolume(unsafe.Pointer(w), C.int(volume))
}

func (w NativeWindow) SetVolumeDisplay(b bool) {
	if b {
		C.setVolumeDisplay(unsafe.Pointer(w), 1)
	} else {
		C.setVolumeDisplay(unsafe.Pointer(w), 0)
	}
}
