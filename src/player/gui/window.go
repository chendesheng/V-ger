package gui

/*
#include "gui.h"
#include <stdlib.h>
*/
import "C"
import (
	"github.com/go-gl/gl"
	"player/srt"
	"unsafe"
)

type Event struct {
	Kind int
	Data interface{}
}
type EventHandlerFunc func(Event)

const (
	Draw = iota
	KeyPress
	DrawSub
	DrawLeftTime
	TrackPositionChanged
)

type PlayProgressInfo struct {
	Left    string
	Right   string
	Percent float64
}

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

	chEvents chan Event

	texture gl.Texture
}

// func (w *Window) Show() {
// 	C.showWindow(w.ptr)
// }

func (w *Window) RefreshContent() {
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
		ptr:      ptr,
		chEvents: make(chan Event),
	}

	println("window:", ptr)

	windows[ptr] = w

	C.showWindow(ptr)
	C.makeWindowCurrentContext(ptr) //must make current context before do texture bind or we will get a all white window

	texture := gl.GenTexture()
	texture.Bind(gl.TEXTURE_2D)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
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

func (w *Window) Draw(img []byte, imgWidth, imgHeight int) {
	// for _, b := range img[:1000] {
	// 	println(b)
	// }

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

	gl.ClearColor(0, 255, 0, 1)
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

	w.HideStartupView()
}

func (w *Window) HideStartupView() {
	C.windowHideStartupView(w.ptr)
}

//run immediately
//must not call in main thread (deadlock)
func (w *Window) PostEvent(e Event) {
	w.chEvents <- e
}

func (w *Window) ShowProgress(left string, right string, percent float64) {
	cleft := C.CString(left)
	defer C.free(unsafe.Pointer(cleft))

	cright := C.CString(right)
	defer C.free(unsafe.Pointer(cright))

	C.showWindowProgress(w.ptr, cleft, cright, C.double(percent))
}
func (w *Window) ShowText(strs []srt.AttributedString, withPosition bool, x, y float64) {
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

	C.showText(w.ptr, p, C.int(len(items)), C.double(x), C.double(y))
}

//export goOnDraw
func goOnDraw(ptr unsafe.Pointer) {
	w := windows[ptr]

	for _, fn := range w.FuncDraw {
		fn()
	}
}

//export goOnTimerTick
func goOnTimerTick(ptr unsafe.Pointer) {
	w := windows[ptr]

	select {
	case e := <-w.chEvents:
		if e.Kind == Draw {
			w.RefreshContent()
		} else if e.Kind == DrawSub {
			s := e.Data.(*srt.SubItem)
			width, height := w.GetWindowSize()

			// fmt.Print("show sub:", s.UsePosition, s.X, s.Y, "\n")

			w.ShowText(s.Content, false, s.X/float64(width), 1-s.Y/float64(height))
		} else if e.Kind == DrawLeftTime {
			arg := e.Data.(PlayProgressInfo)
			w.ShowProgress(arg.Left, arg.Right, arg.Percent)
		}
		break
	default:
		break
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
