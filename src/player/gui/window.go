package gui

/*
#include "gui.h"
#include <stdlib.h>
*/
import "C"
import (
	"log"
	"math"
	. "player/shared"
	"time"
	"unsafe"
	"github.com/go-gl/gl"
)

var windows map[unsafe.Pointer]*Window

func init() {
	windows = make(map[unsafe.Pointer]*Window)
}

type imageRender interface {
	draw(img []byte, width, height int)
	delete()
}
type argSize struct {
	width, height int
}

type Window struct {
	ptr unsafe.Pointer

	FuncTimerTick           []func()
	FuncKeyDown             []func(int)
	FuncOnFullscreenChanged []func(bool)
	FuncOnProgressChanged   []func(int, float64)
	FuncAudioMenuClicked    []func(int)
	FuncSubtitleMenuClicked []func(int, bool)
	FuncMouseWheelled       []func(float64)
	FuncMouseMoved          []func()

	ChanDraw     chan []byte
	ChanShowText chan SubItemArg
	ChanHideText chan uintptr
	ChanSetSize  chan argSize
	ChanSetTitle chan string

	ChanShowMessage chan SubItemArg
	ChanHideMessage chan uintptr

	ChanShowProgress chan *PlayProgressInfo

	ChanSetCursor chan bool

	img []byte

	originalWidth  int
	originalHeight int

	currentMessagePtr uintptr
	currentMessage    *SubItem

	render imageRender

	forceRatio float64
}

// func (w *Window) Show() {
// 	C.showWindow(w.ptr)
// }
func (w *Window) SendDrawImage(img []byte) {
	w.ChanDraw <- img
}
func (w *Window) SendSetCursor(b bool) {
	w.ChanSetCursor <- b
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
	w.render.delete()
}

func (w *Window) GetWindowSize() (int, int) {
	return int(C.getWindowWidth(w.ptr)), int(C.getWindowHeight(w.ptr))
}

func (w *Window) IsFullScreen() bool {
	width, height := w.GetWindowSize()
	swidth, sheight := GetScreenSize()

	return width == swidth && height == sheight
}

func (w *Window) SetTitle(title string) {
	ctitle := C.CString(title)
	defer C.free(unsafe.Pointer(ctitle))

	C.setWindowTitle(w.ptr, ctitle)
}

func fequal(a, b float64) bool {
	return math.Abs(a-b) < 1e-5
}
func (w *Window) ToggleForceScreenRatio() {
	sw, sh := GetScreenSize()
	if fequal(float64(w.originalWidth)/float64(w.originalHeight), float64(sw)/float64(sh)) {
		return
	}

	if w.forceRatio != 0 {
		w.SetForceRatio(0)
	} else {
		w.SetForceRatio(float64(sw) / float64(sh))
	}
}
func (w *Window) SetSize(width, height int) {
	w.ShowStartupView()

	println("set size")

	if width%4 != 0 {
		gl.PixelStorei(gl.UNPACK_ALIGNMENT, 1)
	}

	if w.render != nil {
		w.render.delete()
	}

	println("NewYUVRender")
	w.render = NewYUVRender(w.img, width, height)

	w.originalWidth, w.originalHeight = width, height

	if w.IsFullScreen() {
		w.ToggleFullScreen()
	}

	sw, sh := GetScreenSize()
	if width > int(0.8*float64(sw)) || height > int(0.8*float64(sh)) {
		C.setWindowSize(w.ptr, C.int(0.8*float64(width)), C.int(0.8*float64(height)))
	} else {
		C.setWindowSize(w.ptr, C.int(width), C.int(height))
	}
}

func (w *Window) SetForceRatio(forceRatio float64) {
	width, height := w.originalWidth, w.originalHeight
	w.forceRatio = forceRatio

	if forceRatio > 0 {
		C.setWindowSize(w.ptr, C.int(float64(height)*forceRatio+0.5), C.int(height))
	} else {
		sw, sh := GetScreenSize()
		if width > int(0.8*float64(sw)) || height > int(0.8*float64(sh)) {
			C.setWindowSize(w.ptr, C.int(0.8*float64(width)), C.int(0.8*float64(height)))
		} else {
			C.setWindowSize(w.ptr, C.int(width), C.int(height))
		}
	}
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
		ChanSetSize:      make(chan argSize),
		ChanSetTitle:     make(chan string),

		ChanShowMessage: make(chan SubItemArg),
		ChanHideMessage: make(chan uintptr),

		ChanSetCursor: make(chan bool),

		originalWidth:  width,
		originalHeight: height,
	}
	w.FuncMouseMoved = append(w.FuncMouseMoved, func() {
		w.ShowCursor()
	})

	println("window:", ptr)

	windows[ptr] = w

	C.showWindow(ptr)
	C.makeWindowCurrentContext(ptr) //must make current context before do texture bind or we will get a all white window
	gl.Init()
	gl.ClearColor(0, 0, 0, 1)
	return w
}

func (w *Window) InitEvents() {
	w.FuncOnFullscreenChanged = append(w.FuncOnFullscreenChanged, func(b bool) {
		if w.currentMessagePtr != 0 {
			w.HideText(w.currentMessagePtr)
			w.currentMessagePtr = w.ShowText(w.currentMessage)
		}
	})
}

func (w *Window) ClearEvents() {
	w.FuncOnFullscreenChanged = nil
	w.FuncOnProgressChanged = nil
	w.FuncKeyDown = nil
	w.FuncAudioMenuClicked = nil
	w.FuncSubtitleMenuClicked = nil
	w.FuncMouseWheelled = nil
}

func (w *Window) ToggleFullScreen() {
	C.windowToggleFullScreen(w.ptr)
}

func (w *Window) fitToWindow(imgWidth, imgHeight int) (int, int, int, int) {
	width, height := w.GetWindowSize()

	if w.forceRatio > 0 {
		return 0, 0, width, height
	}

	fwidth, fheight := float64(width), float64(height)

	ratio := float64(imgWidth) / float64(imgHeight)
	windowRatio := fwidth / fheight

	if ratio < windowRatio*1.15 && ratio > windowRatio*0.85 { //aspect radio is close enough
		if fwidth < ratio*fheight { //always larger
			fwidth = ratio * fheight
		} else {
			fheight = fwidth / ratio
		}
	} else {
		if fwidth > ratio*fheight { //always smaller
			fwidth = ratio * fheight
		} else {
			fheight = fwidth / ratio
		}
	}

	vwidth, vheight := int(fwidth+0.5), int(fheight+0.5)
	x, y := (width-vwidth)/2, (height-vheight)/2

	return x, y, vwidth, vheight
}

func (w *Window) draw(img []byte, imgWidth, imgHeight int) {
	println("draw:", imgWidth, imgHeight, img, w.render)

	if len(img) == 0 {
		return
	}

	if w.render == nil {
		return
	}

	w.render.draw(img, imgWidth, imgHeight)

	x, y, width, height := w.fitToWindow(imgWidth, imgHeight)
	gl.Viewport(x, y, width, height)

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
func (w *Window) ShowStartupView() {
	C.windowShowStartupView(w.ptr)
}
func (w *Window) SendShowProgress(p *PlayProgressInfo) {
	w.ChanShowProgress <- p
}
func (w *Window) ShowProgress(p *PlayProgressInfo) {
	cleft := C.CString(p.Left)
	defer C.free(unsafe.Pointer(cleft))

	cright := C.CString(p.Right)
	defer C.free(unsafe.Pointer(cright))

	cspeed := C.CString(p.Speed)
	defer C.free(unsafe.Pointer(cspeed))

	C.showWindowProgress(w.ptr, cleft, cright, C.double(p.Percent), C.double(p.Percent2), cspeed)
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

func (w *Window) SendHideMessage() {
	w.ChanHideMessage <- 0
}

func (w *Window) SendSetSize(width, height int) {
	w.ChanSetSize <- argSize{width, height}
}

func (w *Window) SendSetTitle(title string) {
	w.ChanSetTitle <- title
}

func (w *Window) ShowText(s *SubItem) uintptr {
	strs := s.Content
	items := make([]C.SubItem, 0)
	for _, str := range strs {
		cstr := C.CString(str.Content)
		defer C.free(unsafe.Pointer(cstr))

		// println("content:", str.Content)
		// println("color:", str.Color)
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

func (w *Window) HideCursor() {
	C.hideCursor(w.ptr)
}
func (w *Window) ShowCursor() {
	C.showCursor(w.ptr)
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
	case b := <-w.ChanSetCursor:
		if b {
			w.ShowCursor()
		} else {
			w.HideCursor()
		}
		break
	case arg := <-w.ChanSetSize:
		w.SetSize(arg.width, arg.height)
		break
	case title := <-w.ChanSetTitle:
		w.SetTitle(title)
		break
	case arg := <-w.ChanShowText:
		var arg1 SubItemArg
	begin:
		item := arg.SubItem
		if item.Handle == 0 || item.Handle == 1 {
		skip:
			for {
				select {
				case arg1 = <-w.ChanShowText:
					if arg1.SubItem.Handle != item.Handle || arg1.From != item.From {
						arg.Result <- SubItemExtra{item.Id, w.ShowText(&item)}
						arg = arg1

						goto begin
					}
					break
				default:
					break skip
				}
			}
			// println("show text:", arg.Result)
			arg.Result <- SubItemExtra{item.Id, w.ShowText(&item)}
			// println("show text2")
		} else {
			// println("hide text:", item.Handle)
			w.HideText(item.Handle)
			// println("hide text2:", item.Handle)
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
		w.currentMessage = &item

		arg.Result <- SubItemExtra{0, w.currentMessagePtr}
	case <-w.ChanHideMessage:
		// if (ptr != 0) && (w.currentMessagePtr == ptr) {
		w.HideText(w.currentMessagePtr)
		w.currentMessage = nil
		w.currentMessagePtr = 0
		// }
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

//export goOnMouseWheel
func goOnMouseWheel(ptr unsafe.Pointer, deltaY float64) {
	w := windows[ptr]

	for _, fn := range w.FuncMouseWheelled {
		fn(deltaY)
	}
}

//export goOnMouseMove
func goOnMouseMove(ptr unsafe.Pointer) {
	w := windows[ptr]

	for _, fn := range w.FuncMouseMoved {
		fn()
	}
}
