package cocoa

// #cgo CFLAGS: -x objective-c
// #cgo LDFLAGS: -framework Cocoa -framework OpenGL -framework QuartzCore
import "C"
import "unsafe"

var (
	OnMouseMove      func(int, int)
	OnMenuClick      func(int, int) int
	OnMouseWheel     func(float64, float64)
	OnDrop           func(string)
	OnDraw           func()
	OnTimerTick      func()
	OnKeyDown        func(int) bool
	OnProgressChange func(int, float64)
	OnWillTerminate  func()
	OnOpenOpenPanel  func()
	OnCloseOpenPanel func(string)
	OnOpenFile       func(string) bool

	P Player
)

type Player interface {
	IsPlaying() bool
}

//export goOnMenuClick
func goOnMenuClick(typ int, tag int) C.int {
	if OnMenuClick != nil {
		return C.int(OnMenuClick(typ, tag))
	} else {
		return 0
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
		return b2i(OnKeyDown(keycode))
	} else {
		return 0
	}
}

//export goOnPlaybackChange
func goOnPlaybackChange(typ int, position float64) {
	if OnProgressChange != nil {
		OnProgressChange(typ, position)
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
		return b2i(OnOpenFile(filename))
	} else {
		return 0
	}
}

//export goIsPlaying
func goIsPlaying() C.int {
	if P != nil {
		return b2i(P.IsPlaying())
	} else {
		return 0
	}
}

func b2i(b bool) C.int {
	if b {
		return 1
	} else {
		return 0
	}
}
