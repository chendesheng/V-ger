package main

import (
	"github.com/go-gl/gl"
	"log"
	// . "player/clock"
	"player/glfw"
	"player/srt"
	"time"
)

type TrackStatus struct {
	time    string
	left    string
	percent float64
}

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
)

type Window struct {
	*glfw.Window
	cursorX, cursorY float64 //for drag moving window

	chEvents chan Event
	quit     chan bool

	EventHandlers []EventHandlerFunc

	Texture gl.Texture
}

func NewWindow(width, height int, title string) *Window {
	if !glfw.Init() {
		log.Fatal("init glfw failed")
	}
	glfw.SwapInterval(1)

	println("create window")

	w := &Window{}
	w.chEvents = make(chan Event)
	w.quit = make(chan bool)
	w.EventHandlers = make([]EventHandlerFunc, 0)

	// glfw.WindowHint(glfw.Decorated, 0)
	var err error
	w.Window, err = glfw.CreateWindow(width, height, title)
	if err != nil {
		log.Fatal(err)
	}
	w.Show()
	w.MakeContextCurrent()
	// w.Hide()

	texture := gl.GenTexture()
	texture.Bind(gl.TEXTURE_2D)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGB, width, height, 0,
		gl.RGB, gl.UNSIGNED_BYTE, make([]byte, width*height*3)) //alloc memory

	gl.Enable(gl.TEXTURE_2D)
	gl.Disable(gl.DEPTH_TEST) //disable 3d
	gl.ShadeModel(gl.SMOOTH)

	w.Texture = texture

	w.SetKeyCallback(func(win *glfw.Window, k glfw.Key, scanCode int, action glfw.Action, m glfw.ModifierKey) {
		w.onKeyAction(k, scanCode, action, m)
	})
	w.SetMouseButtonCallback(func(win *glfw.Window, button glfw.MouseButton, action glfw.Action, mod glfw.ModifierKey) {
		w.onMouseButtonAction(button, action, mod)
	})
	w.SetCursorPositionCallback(func(win *glfw.Window, x float64, y float64) {
		w.onCursorPositionChange(x, y)
	})
	w.SetSizeCallback(func(win *glfw.Window, width int, height int) {
		w.onSizeChange(width, height)
	})
	w.SetTimerCallback(func(win *glfw.Window) {
		w.timerTick()
	})
	w.SetDrawCallback(func(win *glfw.Window) {
		w.fireEvent(Event{Draw, nil})
	})

	w.SetNeedsDisplay(true)
	w.ShowOrHideStartupView(false)
	return w
}

func (w *Window) AddEventHandler(fn EventHandlerFunc) {
	w.EventHandlers = append(w.EventHandlers, fn)
}

func (w *Window) fireEvent(e Event) {
	for _, fn := range w.EventHandlers {
		fn(e)
	}
}

func (w *Window) onKeyAction(k glfw.Key, scanCode int, action glfw.Action, m glfw.ModifierKey) {

	if action != glfw.Press {
		return
	}
	key := glfw.Key(k)
	// switch key {
	// // case glfw.KeyEscape:
	// // 	w.SetShouldClose(true)
	// // 	break
	// case glfw.KeySpace:
	// 	w.fireEvent(Event{KeyPress, key})
	// 	break
	// }
	w.fireEvent(Event{KeyPress, key})
}

//main thread
func (w *Window) onMouseButtonAction(button glfw.MouseButton, action glfw.Action, mod glfw.ModifierKey) {
	if button == glfw.MouseButtonLeft {
		switch action {
		case glfw.Press:
			w.cursorX, w.cursorY = w.GetCursorPosition()
			w.SetInputMode(glfw.Cursor, glfw.CursorHidden)
			break
			// case glfw.Release:
			// 	w.SetInputMode(glfw.Cursor, glfw.CursorNormal)
			// 	break
		}
	}
}

//main thread
func (w *Window) onCursorPositionChange(x, y float64) {
	action := w.GetMouseButton(glfw.MouseButtonLeft)
	if action == glfw.Press {
		wx, wy := w.GetPosition()
		w.SetPosition(int(float64(wx)+x-w.cursorX), int(float64(wy)+y-w.cursorY))
	} else if action == glfw.Release {
		w.SetInputMode(glfw.Cursor, glfw.CursorNormal)
	}
}

func (w *Window) timerTick() {
	select {
	case e := <-w.chEvents:
		if e.Kind == Draw {
			w.SetNeedsDisplay(true)
		} else if e.Kind == DrawSub {
			w.ShowText(e.Data.([]srt.AttributedString))
		} else if e.Kind == DrawLeftTime {
			arg := e.Data.(TrackStatus)
			w.ShowLeftTime(arg.time, arg.left, arg.percent)
		}
		break
	default:
		break
	}
}

//main thread
func (w *Window) onSizeChange(width int, height int) {
	if (width == 1920) && (height == 1080) {
		w.SetInputMode(glfw.Cursor, glfw.CursorHidden)
	}

	w.fireEvent(Event{Draw, nil})
}

func (w *Window) SetCursorAutoHide() {
	go func() {
		for {
			x, y := w.GetCursorPosition()

			<-time.After(time.Second)
			x1, y1 := w.GetCursorPosition()
			if x == x1 && y == y1 &&
				w.GetInputMode(glfw.Cursor) == glfw.CursorNormal {
				w.SetInputMode(glfw.Cursor, glfw.CursorHidden) //?? thread safe
			}
		}
	}()
}

func (w *Window) DrawClear(imgWidth, imgHeight int) {
	w.Draw(make([]byte, imgWidth*imgHeight*3), imgWidth, imgHeight)
	// gl.ClearColor(0, 0, 0, 0)
	// gl.Clear(gl.COLOR_BUFFER_BIT)

	// w.SwapBuffers()
}

func (w *Window) Draw(img []byte, imgWidth, imgHeight int) {
	// println("draw:", len(img))

	w.Texture.Bind(gl.TEXTURE_2D)

	gl.TexSubImage2D(gl.TEXTURE_2D, 0, 0, 0, int(imgWidth), int(imgHeight), gl.RGB,
		gl.UNSIGNED_BYTE, img)

	// width, height := w.GetFramebufferSize()
	// gl.Viewport(0, 0, width, height)
	width, height := w.GetFramebufferSize()
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

	w.SwapBuffers()
}

//run immediately
//must not call in main thread (deadlock)
func (w *Window) PostEvent(e Event) {
	w.chEvents <- e
}

func (w *Window) EventLoop() {
	println("begin event loop")
	for !w.ShouldClose() {
		glfw.WaitEvents()
	}

	println("end event loop")
}

func (w *Window) Destory() {
	w.Texture.Delete()
	w.Window.Destroy()
}
