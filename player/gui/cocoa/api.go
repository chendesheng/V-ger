package cocoa

// #include "gui.h"
// #include <stdlib.h>
import "C"
import (
	"log"
	"reflect"
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

func AddRecentOpenedFile(str string) {
	cstr := C.CString(str)
	defer C.free(unsafe.Pointer(cstr))

	C.addRecentOpenedFile(cstr)
}

func (w NativeWindow) Alert(str string) {
	cstr := C.CString(str)
	defer C.free(unsafe.Pointer(cstr))

	C.alert(w.ptr, cstr)
}

type NativeWindow struct {
	ptr unsafe.Pointer
}

func (w NativeWindow) RefreshContent() {
	C.refreshWindowContent(w.ptr)
}

func (w NativeWindow) GetSize() (int, int) {
	sz := C.getWindowSize(w.ptr)
	return int(sz.width), int(sz.height)
}

func (w NativeWindow) SetTitleWithRepresentFilename(title string) {
	ctitle := C.CString(title)
	defer C.free(unsafe.Pointer(ctitle))

	C.setWindowTitleWithRepresentedFilename(w.ptr, ctitle)
}

func (w NativeWindow) SetTitle(title string) {
	ctitle := C.CString(title)
	defer C.free(unsafe.Pointer(ctitle))

	C.setWindowTitle(w.ptr, ctitle)
}

func (w NativeWindow) SetSize(width, height int) {
	C.setWindowSize(w.ptr, C.int(width), C.int(height))
}

func NewWindow(title string, width, height int) NativeWindow {
	ctitle := C.CString(title)
	defer C.free(unsafe.Pointer(ctitle))

	return NativeWindow{unsafe.Pointer(C.newWindow(ctitle, C.int(width), C.int(height)))}
}

func (w NativeWindow) Show() {
	C.showWindow(w.ptr)
}

func (w NativeWindow) MakeCurrentContext() {
	C.initWindowCurrentContext(w.ptr)
}

func (w NativeWindow) ToggleFullScreen() {
	C.toggleFullScreen(w.ptr)
}

func (w NativeWindow) IsFullScreen() bool {
	return C.isFullScreen(w.ptr) != 0
}

func (w NativeWindow) UpdatePlaybackInfo(left, right string, percent float64) {
	cleft := C.CString(left)
	defer C.free(unsafe.Pointer(cleft))

	cright := C.CString(right)
	defer C.free(unsafe.Pointer(cright))

	C.updatePlaybackInfo(w.ptr, cleft, cright, C.double(percent))
}
func (w NativeWindow) UpdateBufferInfo(speed string, percent float64) {
	cspeed := C.CString(speed)
	defer C.free(unsafe.Pointer(cspeed))

	C.updateBufferInfo(w.ptr, cspeed, C.double(percent))
}

func (w NativeWindow) ShowSubtitle(items []struct {
	Content string
	Style   int //0 -normal, 1 -italic, 2 -bold, 3 italic and bold
	Color   uint
}, posType int, x, y float64) uintptr {

	if len(items) == 0 {
		return 0
	}

	p := (*C.struct_AttributedString)(C.malloc(C.size_t(C.sizeof_struct_AttributedString * len(items))))

	var ctexts []C.AttributedString
	ctextsHeader := (*reflect.SliceHeader)(unsafe.Pointer(&ctexts))
	ctextsHeader.Cap = len(items)
	ctextsHeader.Len = len(items)
	ctextsHeader.Data = uintptr(unsafe.Pointer(p))
	defer C.free(unsafe.Pointer(p))

	for i, str := range items {
		// log.Println(str.Content)
		cstr := C.CString(str.Content)
		defer C.free(unsafe.Pointer(cstr))

		ctext := &ctexts[i]
		ctext.str = cstr
		ctext.style = C.int(str.Style)
		ctext.color = C.uint(str.Color)
	}

	csubitem := (*C.struct_SubItem)(C.malloc(C.sizeof_struct_SubItem))
	defer C.free(unsafe.Pointer(csubitem))

	csubitem.texts = (*C.struct_AttributedString)(p)
	csubitem.length = C.int(len(ctexts))
	csubitem.align = C.int(posType)
	csubitem.x = C.double(x)
	csubitem.y = C.double(y)

	return uintptr(C.showSubtitle(w.ptr, csubitem))
}

func (w NativeWindow) HideSubtitle(ptr uintptr) {
	C.hideSubtitle(w.ptr, C.long(ptr))
}

func (w NativeWindow) SetControlsVisible(b bool, autoHide bool) {
	C.setControlsVisible(w.ptr, b2i(b), b2i(autoHide))
}

func (w NativeWindow) SetSpinningVisible(b bool) {
	C.setSpinningVisible(w.ptr, b2i(b))
}

func (w NativeWindow) SetVolume(volume int) {
	C.setVolume(w.ptr, C.int(volume))
}

func (w NativeWindow) SetVolumeVisible(b bool) {
	C.setVolumeVisible(w.ptr, b2i(b))
}

func (w NativeWindow) Close() {
	C.closeWindow(w.ptr)
}

func (w NativeWindow) FlushBuffer() {
	C.flushBuffer(w.ptr)
}

func (w NativeWindow) MakeGLCurrentContext() {
	C.makeCurrentContext(w.ptr)
}

func (w NativeWindow) SetSubFontSize(sz float64) {
	log.Print("SetSubFontSize:", sz)

	C.setSubFontSize(w.ptr, C.double(sz))
}
