package cocoa

// #cgo CFLAGS: -x objective-c
// #cgo LDFLAGS: -framework Cocoa -framework OpenGL -framework QuartzCore
import "C"
import "unsafe"

var (
	OnMouseMove         func(int, int)
	OnMenuClicked       func(int, int)
	OnMouseWheel        func(float64, float64)
	OnDrop              func(string)
	OnDraw              func()
	OnTimerTick         func()
	OnKeyDown           func(int) bool
	OnProgressChanged   func(int, float64)
	OnFullscreenChanged func(int)
	OnWillTerminate     func()
	OnOpenOpenPanel     func()
	OnCloseOpenPanel    func(string)
	OnOpenFile          func(string) bool
)

//export goOnMenuClicked
func goOnMenuClicked(typ int, tag int) {
	if OnMenuClicked != nil {
		OnMenuClicked(typ, tag)
	}
}

//export goOnMouseWheel
func goOnMouseWheel(deltaY float64) {
	if OnMouseWheel != nil {
		OnMouseWheel(0, deltaY)
	}
}

//export goOnMouseMove
func goOnMouseMove() {
	if OnMouseMove != nil {
		OnMouseMove(0, 0)
	}
}

//export goOnDraw
func goOnDraw() {
	if OnDraw != nil {
		OnDraw()
	}
}

//export goOnTimerTick
func goOnTimerTick() {
	if OnTimerTick != nil {
		OnTimerTick()
	}
}

//export goOnKeyDown
func goOnKeyDown(keycode int) C.int { //true if already handled
	if OnKeyDown != nil {
		ret := OnKeyDown(keycode)

		if ret {
			return 1
		} else {
			return 0
		}
	} else {
		return 0
	}
}

//export goOnProgressChanged
func goOnProgressChanged(typ int, position float64) {
	if OnProgressChanged != nil {
		OnProgressChanged(typ, position)
	}
}

//export goOnFullscreenChanged
func goOnFullscreenChanged(typ int) {
	if OnFullscreenChanged != nil {
		OnFullscreenChanged(typ)
	}
}

//export goOnWillTerminate
func goOnWillTerminate() {
	if OnWillTerminate != nil {
		OnWillTerminate()
	}
}

//export goOnOpenOpenPanel
func goOnOpenOpenPanel() {
	if OnOpenOpenPanel != nil {
		OnOpenOpenPanel()
	}
}

//export goOnCloseOpenPanel
func goOnCloseOpenPanel(filename *_Ctype_char) {
	if OnCloseOpenPanel != nil {
		name := C.GoString(filename)
		OnCloseOpenPanel(name)
	}
}

//export goOnOpenFile
func goOnOpenFile(cfilename unsafe.Pointer) C.int {
	filename := C.GoString((*C.char)(cfilename))

	if OnOpenFile != nil {
		if OnOpenFile(filename) {
			return 1
		} else {
			return 0
		}
	} else {
		return 0
	}
}
