package gui

// #cgo CFLAGS: -x objective-c
// #cgo LDFLAGS: -framework Cocoa -framework OpenGL -framework QuartzCore
import "C"
import "vger/player/gui/cocoa"

func init() {
	cocoa.OnMouseMove = onMouseMove
	cocoa.OnDraw = onDraw
	cocoa.OnFullscreenChanged = onFullscreenChanged
	cocoa.OnKeyDown = onKeyDown
	cocoa.OnMenuClicked = onMenuClicked
	cocoa.OnMouseWheel = onMouseWheel
	cocoa.OnOpenFile = onOpenFile
	cocoa.OnOpenOpenPanel = onOpenOpenPanel
	cocoa.OnCloseOpenPanel = onCloseOpenPanel
	cocoa.OnProgressChanged = onProgressChanged
	cocoa.OnTimerTick = onTimerTick
}

type Window struct {
	cocoa.NativeWindow
	window
}

func newWindow(title string, width, height int) cocoa.NativeWindow {
	return cocoa.NewWindow(title, width, height)
}

func run() {
	cocoa.Run()
}

func getScreenSize() (int, int) {
	return cocoa.GetScreenSize()
}
