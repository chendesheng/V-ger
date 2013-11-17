package gui

/*
#include "gui.h"
#include <stdlib.h>
*/
import "C"
import (
	"github.com/go-gl/gl"
	"unsafe"

	. "player/shared"
)

const (
	draw = iota
	KeyPress
	DrawSub
	DrawLeftTime
	TrackPositionChanged
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

	texture gl.Texture

	ChanDraw         chan []byte
	ChanShowText     chan *SubItem
	ChanShowProgress chan *PlayProgressInfo

	img []byte

	originalWidth  int
	originalHeight int
}

// func (w *Window) Show() {
// 	C.showWindow(w.ptr)
// }

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

		ChanDraw:         make(chan []byte),
		ChanShowProgress: make(chan *PlayProgressInfo),
		ChanShowText:     make(chan *SubItem),

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
	// gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	// gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)

	gl.TexParameterf(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	gl.TexParameterf(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
	gl.TexParameterf(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameterf(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGB, width, height, 0,
		gl.RGB, gl.UNSIGNED_BYTE, make([]byte, width*height*3)) //alloc memory

	gl.Enable(gl.TEXTURE_2D)
	gl.Disable(gl.DEPTH_TEST) //disable 3d
	gl.ShadeModel(gl.SMOOTH)

	w.texture = texture

	// w.RefreshContent()

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

	w.HideStartupView()
}

func (w *Window) HideStartupView() {
	C.windowHideStartupView(w.ptr)
}

func (w *Window) ShowProgress(left string, right string, percent float64) {
	cleft := C.CString(left)
	defer C.free(unsafe.Pointer(cleft))

	cright := C.CString(right)
	defer C.free(unsafe.Pointer(cright))

	C.showWindowProgress(w.ptr, cleft, cright, C.double(percent))
}
func (w *Window) ShowText(s *SubItem) {
	strs := s.Content

	items := make([]C.SubItem, 0)
	for _, str := range strs {
		cstr := C.CString(str.Content)
		defer C.free(unsafe.Pointer(cstr))

		items = append(items, C.SubItem{cstr, C.int(str.Style), C.uint(str.Color)})
	}

	var p *C.SubItem
	if len(strs) > 0 {
		p = (*C.SubItem)(unsafe.Pointer(&items[0]))
	}

	// var t = C.int(0)
	// if withPosition {
	// 	t = 1
	// }

	C.showText(w.ptr, p, C.int(len(items)), 0, 0)
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
	case img := <-w.ChanDraw:
		w.RefreshContent(img)
	default:
	}

	select {
	case p := <-w.ChanShowProgress:
		w.ShowProgress(p.Left, p.Right, p.Percent)
	default:
	}

	select {
	case s := <-w.ChanShowText:
		// width, height := w.GetWindowSize()
		w.ShowText(s)
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
