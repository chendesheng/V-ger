package gui

/*
#include "gui.h"
#include <stdlib.h>
*/
import "C"
import (
	"github.com/go-gl/gl"
	. "player/shared"
	"time"
	"unsafe"
)

var windows map[unsafe.Pointer]*Window

func init() {
	windows = make(map[unsafe.Pointer]*Window)
}

type Window struct {
	ptr unsafe.Pointer

	FuncDraw                []func()
	FuncTimerTick           []func()
	FuncKeyDown             []func(int)
	FuncOnFullscreenChanged []func(bool)
	FuncOnProgressChanged   []func(int, float64)
	FuncAudioMenuClicked    []func(int)
	FuncSubtitleMenuClicked []func(int, bool)

	texture gl.Texture

	ChanDraw     chan []byte
	ChanShowText chan SubItemArg
	ChanHideText chan uintptr

	ChanShowMessage chan SubItemArg
	ChanHideMessage chan uintptr

	ChanShowProgress chan *PlayProgressInfo

	img []byte

	originalWidth  int
	originalHeight int

	currentMessagePtr uintptr
}

// func (w *Window) Show() {
// 	C.showWindow(w.ptr)
// }
func (w *Window) SendDrawImage(img []byte) {
	w.ChanDraw <- img
}
func (w *Window) FlushImageBuffer() {
	for {
		select {
		case <-w.ChanDraw:
			println("window drop image")
			break
		default:
			println("window flush image buffer return")
			return
		}
	}
}
func (w *Window) RefreshContent(img []byte) {
	w.img = img

	C.refreshWindowContent(w.ptr)
}

func (w *Window) Destory() {
	w.texture.Delete()
}

func (w *Window) DrawSubtitle() {

}

func (w *Window) GetWindowSize() (int, int) {
	return int(C.getWindowWidth(w.ptr)), int(C.getWindowHeight(w.ptr))
}

func NewWindow(title string, width, height int) *Window {
	ctitle := C.CString(title)
	defer C.free(unsafe.Pointer(ctitle))
	ptr := unsafe.Pointer(C.newWindow(ctitle, C.int(width), C.int(height)))

	w := &Window{
		ptr: ptr,

		ChanDraw:         make(chan []byte, 1),
		ChanShowProgress: make(chan *PlayProgressInfo),
		ChanShowText:     make(chan SubItemArg, 20), //the buffer is required because show&hide must handles in the same order
		ChanHideText:     make(chan uintptr),

		ChanShowMessage: make(chan SubItemArg),
		ChanHideMessage: make(chan uintptr),

		originalWidth:  width,
		originalHeight: height,
	}

	println("window:", ptr)

	windows[ptr] = w

	C.showWindow(ptr)
	C.makeWindowCurrentContext(ptr) //must make current context before do texture bind or we will get a all white window

	if width%4 != 0 {
		gl.PixelStorei(gl.UNPACK_ALIGNMENT, 1)
	}

	texture := gl.GenTexture()
	texture.Bind(gl.TEXTURE_2D)

	gl.TexParameterf(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameterf(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexParameterf(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameterf(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGB, width, height, 0,
		gl.RGB, gl.UNSIGNED_BYTE, make([]byte, width*height*3)) //alloc memory

	gl.Enable(gl.TEXTURE_2D)
	gl.Disable(gl.DEPTH_TEST) //disable 3d
	gl.ShadeModel(gl.SMOOTH)
	gl.ClearColor(0, 0, 0, 1)

	w.texture = texture

	w.FuncKeyDown = append(w.FuncKeyDown, func(keycode int) {
		if keycode == KEY_ESCAPE {
			C.windowToggleFullScreen(w.ptr)
		}
	})

	return w
}

func (w *Window) draw(img []byte, imgWidth, imgHeight int) {
	if len(img) == 0 {
		return
	}

	// println("width:", imgWidth, "height:", imgHeight)

	w.texture.Bind(gl.TEXTURE_2D)

	gl.TexSubImage2D(gl.TEXTURE_2D, 0, 0, 0, int(imgWidth), int(imgHeight), gl.RGB,
		gl.UNSIGNED_BYTE, img)

	// width, height := w.GetWindowSize()
	// gl.Viewport(0, 0, 1280, 720)
	width, height := w.GetWindowSize()
	fwidth, fheight := float64(width), float64(height)

	ratio := float64(imgWidth) / float64(imgHeight)
	if fwidth > ratio*fheight { //always smaller
		fwidth = ratio * fheight
	} else {
		fheight = fwidth / ratio
	}

	vwidth, vheight := int(fwidth+0.5), int(fheight+0.5)
	x, y := (width-vwidth)/2, (height-vheight)/2

	gl.Viewport(x, y, vwidth, vheight)

	gl.Clear(gl.COLOR_BUFFER_BIT)
	gl.MatrixMode(gl.PROJECTION)
	gl.LoadIdentity()

	gl.Begin(gl.QUADS)
	gl.TexCoord2d(0, 0)
	gl.Vertex2d(-1, -1)

	gl.TexCoord2f(1, 0)
	gl.Vertex2d(1, -1)

	gl.TexCoord2d(1, 1)
	gl.Vertex2d(1, 1)

	gl.TexCoord2d(0, 1)
	gl.Vertex2d(-1, 1)
	gl.End()

	w.hideStartupView()
}

func (w *Window) hideStartupView() {
	C.windowHideStartupView(w.ptr)
}
func (w *Window) SendShowProgress(p *PlayProgressInfo) {
	w.ChanShowProgress <- p
}
func (w *Window) ShowProgress(p *PlayProgressInfo) {
	cleft := C.CString(p.Left)
	defer C.free(unsafe.Pointer(cleft))

	cright := C.CString(p.Right)
	defer C.free(unsafe.Pointer(cright))

	C.showWindowProgress(w.ptr, cleft, cright, C.double(p.Percent), C.double(p.Percent2))
}
func (w *Window) SendShowText(s SubItemArg) {
	// res := make(chan SubItemExtra)
	w.ChanShowText <- s
	// return <-res
}
func (w *Window) SendShowMessage(msg string, autoHide bool) uintptr {
	s := SubItem{}
	s.PositionType = 7
	s.X = 20
	s.Y = 20
	s.Content = make([]AttributedString, 0)
	s.Content = append(s.Content, AttributedString{msg, 3, 0xffffff})

	res := make(chan SubItemExtra)
	w.ChanShowMessage <- SubItemArg{s, res}
	ptr := <-res

	if autoHide {
		time.Sleep(2 * time.Second)
		w.ChanHideMessage <- ptr.Handle

		return 0
	} else {
		return ptr.Handle
	}
}

func (w *Window) SendHideMessage(h uintptr) {
	w.ChanHideMessage <- h
}

func (w *Window) ShowText(s *SubItem) uintptr {
	strs := s.Content

	items := make([]C.SubItem, 0)
	for _, str := range strs {
		cstr := C.CString(str.Content)
		defer C.free(unsafe.Pointer(cstr))

		println("color:", str.Content)
		println("color:", str.Color)
		items = append(items, C.SubItem{cstr, C.int(str.Style), C.uint(str.Color)})
	}

	var p *C.SubItem
	if len(strs) > 0 {
		p = (*C.SubItem)(unsafe.Pointer(&items[0]))
	}

	return uintptr(C.showText(w.ptr, p, C.int(len(items)), C.int(s.PositionType), C.double(s.X), C.double(s.Y)))
}
func (w *Window) SendHideText(arg SubItemArg) {
	// w.ChanHideText <- ptr
	w.ChanShowText <- arg
}
func (w *Window) HideText(ptr uintptr) {
	C.hideText(w.ptr, unsafe.Pointer(ptr))
}
func (w *Window) ShowSubList(sub Sub) {
	// C.showSubList()
}

//export goOnDraw
func goOnDraw(ptr unsafe.Pointer) {
	w := windows[ptr]
	w.draw(w.img, w.originalWidth, w.originalHeight)
}

//export goOnTimerTick
func goOnTimerTick(ptr unsafe.Pointer) {
	w := windows[ptr]

	select {
	case img, ok := <-w.ChanDraw:
		if ok {
			w.RefreshContent(img)
		}
	default:
	}

	select {
	case p := <-w.ChanShowProgress:
		w.ShowProgress(p)
	default:
	}

	select {
	case arg := <-w.ChanShowText:
		item := arg.SubItem
		if item.Handle == 0 || item.Handle == 1 {
			println("show text:", arg.Result)
			arg.Result <- SubItemExtra{item.Id, w.ShowText(&item)}
			println("show text2")
		} else {
			println("hide text:", item.Handle)
			w.HideText(item.Handle)
			println("hide text2:", item.Handle)
		}
		break
	// case ptr := <-w.ChanHideText:
	// 	w.HideText(ptr)
	case arg := <-w.ChanShowMessage:
		if w.currentMessagePtr != 0 {
			w.HideText(w.currentMessagePtr)
		}

		item := arg.SubItem
		w.currentMessagePtr = w.ShowText(&item)
		arg.Result <- SubItemExtra{0, w.currentMessagePtr}
	case ptr := <-w.ChanHideMessage:
		if (ptr != 0) && (w.currentMessagePtr == ptr) {
			w.HideText(w.currentMessagePtr)
			w.currentMessagePtr = 0
		}
	default:
	}
}

//export goOnKeyDown
func goOnKeyDown(ptr unsafe.Pointer, keycode int) {
	w := windows[ptr]

	for _, fn := range w.FuncKeyDown {
		fn(keycode)
	}
}

//export goOnProgressChanged
func goOnProgressChanged(ptr unsafe.Pointer, typ int, position float64) {
	w := windows[ptr]

	for _, fn := range w.FuncOnProgressChanged {
		fn(typ, position)
	}
}

//export goOnFullscreenChanged
func goOnFullscreenChanged(ptr unsafe.Pointer, b int) {
	w := windows[ptr]

	for _, fn := range w.FuncOnFullscreenChanged {
		fn(b != 0)
	}
}

// //export goOnAudioStreamChanged
// func goOnAudioStreamChanged(cname *C.char) {
// 	name := C.GoString(cname)
// 	println(name)
// }

// //export goOnSubtitleChanged
// func goOnSubtitleChanged(name1 *C.char, name2 *C.char) {

// }

//export goOnAudioMenuClicked
func goOnAudioMenuClicked(ptr unsafe.Pointer, tag int) {
	w := windows[ptr]

	for _, fn := range w.FuncAudioMenuClicked {
		fn(tag)
	}
}

//export goOnSubtitleMenuClicked
func goOnSubtitleMenuClicked(ptr unsafe.Pointer, tag int, showOrHide int) {
	w := windows[ptr]

	for _, fn := range w.FuncSubtitleMenuClicked {
		fn(tag, showOrHide != 0)
	}
}
