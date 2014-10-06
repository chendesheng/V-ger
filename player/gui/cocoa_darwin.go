package gui

// #cgo CFLAGS: -x objective-c
// #cgo LDFLAGS: -framework Cocoa -framework OpenGL -framework QuartzCore
import "C"
import "vger/player/gui/cocoa"

func init() {
	cocoa.OnMouseMove = onMouseMove
	cocoa.OnDraw = onDraw
	cocoa.OnKeyDown = onKeyDown
	cocoa.OnMenuClick = onMenuClick
	cocoa.OnMouseWheel = onMouseWheel
	cocoa.OnOpenFile = onOpenFile
	cocoa.OnOpenOpenPanel = onOpenOpenPanel
	cocoa.OnCloseOpenPanel = onCloseOpenPanel
	cocoa.OnProgressChange = onProgressChange
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

func AddRecentOpenedFile(filename string) {
	cocoa.AddRecentOpenedFile(filename)
}
