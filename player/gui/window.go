package gui

import (
	"log"
	"math"
	"runtime"
	"time"
	"vger/player/shared"

	"github.com/go-gl/gl"
)

var w *Window // current window

type window struct {
	FuncTimerTick         []func()
	FuncKeyDown           []func(int) bool
	FuncOnProgressChanged []func(int, float64)

	chFunc chan func()
	chDraw chan drawArg

	chShowSpinning chan bool

	img []byte

	originalWidth  int
	originalHeight int

	currentMessagePtr uintptr
	currentMessage    *shared.SubItem

	render imageRender

	forceRatio float64

	showMessageDeadline time.Time

	chDelayShowSpinning      chan bool
	chDelayForceShowSpinning chan struct{}

	displayingTexts map[int]uintptr
}
type drawArg struct {
	data []byte
	res  chan struct{}
}

type imageRender interface {
	draw(img []byte, width, height int)
	delete()
}

func (w *Window) SendSetControlsVisible(b bool, autoHide bool) {
	w.chFunc <- func() {
		w.SetControlsVisible(b, autoHide)
	}
}

func (w *Window) DestoryRender() {
	if w.render != nil {
		w.render.delete()
		w.render = nil
	}
}
func (w *Window) SendShowText(s shared.SubItem) {
	arg := &s
	w.chFunc <- func() {
		ptr := w.ShowSubtitle(arg)
		w.displayingTexts[arg.Handle] = ptr
	}
}

func (w *Window) SendDestoryRender() {
	log.Print("SendDestoryRender")

	w.chFunc <- func() {
		w.DestoryRender()
	}
}

func fequal(a, b float64) bool {
	return math.Abs(a-b) < 1e-5
}

func (w *Window) ToggleForceScreenRatio() {
	sw, sh := getScreenSize()
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
	w.drawClear()

	log.Printf("set window size:%d %d", width, height)

	w.render = NewYUVRender(width, height)

	w.originalWidth, w.originalHeight = width, height

	if !w.IsFullScreen() {
		sw, sh := getScreenSize()
		if width > int(0.9*float64(sw)) || height > int(0.9*float64(sh)) {
			ratio := float64(height) / float64(width)
			width = int(float64(sw) * 0.9)
			height = int(float64(sw) * 0.9 * ratio)

			w.NativeWindow.SetSize(width, height)
		} else {
			w.NativeWindow.SetSize(width, height)
		}
	}
}

func (w *Window) SetForceRatio(forceRatio float64) {
	width, height := w.originalWidth, w.originalHeight
	w.forceRatio = forceRatio

	if forceRatio > 0 {
		w.NativeWindow.SetSize(int(float64(height)*forceRatio+0.5), height)
	} else {
		sw, sh := getScreenSize()
		if width > int(0.8*float64(sw)) || height > int(0.8*float64(sh)) {
			w.NativeWindow.SetSize(int(0.8*float64(width)), int(0.8*float64(height)))
		} else {
			w.NativeWindow.SetSize(width, height)
		}
	}
}

func NewWindow(title string, width, height int) *Window {
	w = &Window{
		newWindow(title, width, height),
		window{
			chShowSpinning: make(chan bool),

			originalWidth:  width,
			originalHeight: height,

			displayingTexts: make(map[int]uintptr),
			chFunc:          make(chan func()),
			chDraw:          make(chan drawArg),
		},
	}

	log.Print("NewWindow:", w.NativeWindow)

	w.Show()
	w.MakeCurrentContext() //must make current context before do texture bind or we will get a all white window

	gl.Init()
	gl.ClearColor(0, 0, 0, 1)

	w.drawClear()

	go func() {
		runtime.LockOSThread()

		for input := range w.chDraw {
			w.MakeGLCurrentContext()
			w.draw(input.data, w.originalWidth, w.originalHeight)
			w.FlushBuffer()
			input.res <- struct{}{}
		}
	}()

	return w
}

func (w *Window) ClearEvents() {
	w.FuncOnProgressChanged = nil
	w.FuncKeyDown = nil
}

func (w *Window) fitToWindow(imgWidth, imgHeight int) (int, int, int, int) {
	width, height := w.GetSize()

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
	if len(img) == 0 {
		log.Print("draw: no image")
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
}

func (w *Window) SendShowProgress(left string, right string, percent float64) {
	info := &struct {
		left    string
		right   string
		percent float64
	}{left, right, percent}
	w.chFunc <- func() {
		w.UpdatePlaybackInfo(info.left, info.right, info.percent)
	}
}
func (w *Window) SendShowBufferInfo(speed string, percent float64) {
	w.chFunc <- func() {
		w.UpdateBufferInfo(speed, percent)
	}
}

func createMessageSubItem(msg string) shared.SubItem {
	s := shared.SubItem{}
	s.PositionType = 7
	s.X = 20
	s.Y = 20
	s.Content = make([]shared.AttributedString, 0)
	s.Content = append(s.Content, shared.AttributedString{msg, 3, 0xffffff})

	return s
}
func (w *Window) SendShowMessage(msg string, autoHide bool) {
	w.chFunc <- func() {
		w.ShowMessage(msg, autoHide)
	}
}

func (w *Window) ShowMessage(msg string, autoHide bool) {
	s := createMessageSubItem(msg)
	w.showMessage(&s, autoHide)
}

func (w *Window) showMessage(s *shared.SubItem, autoHide bool) {
	if w.currentMessagePtr != 0 {
		w.HideSubtitle(w.currentMessagePtr)
	}

	w.currentMessagePtr = w.ShowSubtitle(s)
	w.currentMessage = s

	if autoHide {
		w.showMessageDeadline = time.Now().Add(2 * time.Second)
	} else {
		w.showMessageDeadline = time.Now().Add(1000 * time.Hour)
	}
}

func (w *Window) HideMessage() {
	if w.currentMessagePtr != 0 {
		w.HideSubtitle(w.currentMessagePtr)
		w.currentMessage = nil
		w.currentMessagePtr = 0
	}
	w.currentMessage = nil
}

func (w *Window) SendHideMessage() {
	w.chFunc <- func() {
		w.HideMessage()
	}
}

func (w *Window) SendSetSize(width, height int) {
	w.chFunc <- func() {
		w.SetSize(width, height)
	}
}

func (w *Window) SendSetTitle(title string) {
	w.chFunc <- func() {
		w.SetTitle(title)
	}
}

func (w *Window) ShowSubtitle(s *shared.SubItem) uintptr {
	strs := s.Content
	items := make([]struct {
		Content string
		Style   int
		Color   uint
	}, len(strs))

	for i, str := range strs {
		items[i].Content = str.Content
		items[i].Style = str.Style
		items[i].Color = str.Color
	}

	return w.NativeWindow.ShowSubtitle(items, s.PositionType, s.X, s.Y)
}

func (w *Window) SendHideText(handle int) {
	w.chFunc <- func() {
		if ptr, ok := w.displayingTexts[handle]; ok {
			w.HideSubtitle(ptr)
			delete(w.displayingTexts, handle)
		}
	}
}

func (w *Window) SendShowSpinning() {
	// log.Print(string(debug.Stack()))

	if w.chDelayShowSpinning == nil {
		w.chDelayShowSpinning = make(chan bool)
		w.chDelayForceShowSpinning = make(chan struct{})
		go func() {
			b := false
			bsaved := false
			for {
				select {
				case <-time.After(500 * time.Millisecond):
					if b != bsaved {
						bsaved = b
						w.chShowSpinning <- b
					}
				case b = <-w.chDelayShowSpinning:
				case <-w.chDelayForceShowSpinning:
					b = false
					if b != bsaved {
						bsaved = b
						w.chShowSpinning <- b
					}
					b = <-w.chDelayShowSpinning
				}
			}
		}()
	}
	w.chDelayShowSpinning <- true
}
func (w *Window) SendHideSpinning(forceHide bool) {
	// log.Print(string(debug.Stack()))

	if forceHide {
		w.chDelayForceShowSpinning <- struct{}{}
	} else {
		w.chDelayShowSpinning <- false
	}
}

func (w *Window) SendSetVolume(volume int) {
	w.chFunc <- func() {
		w.SetVolume(volume)
	}
}

func (w *Window) SendSetVolumeVisible(b bool) {
	w.chFunc <- func() {
		w.SetVolumeVisible(b)
	}
}

func (w *Window) Draw(img []byte) {
	if len(img) > 0 {
		res := make(chan struct{})
		w.chDraw <- drawArg{img, res}
		<-res
	}
}
func (w *Window) drawClear() {
	w.MakeGLCurrentContext()
	gl.Clear(gl.COLOR_BUFFER_BIT)
	w.FlushBuffer()
}
func onTimerTick() {
	if w != nil {
		select {
		case b := <-w.chShowSpinning:
			w.SetSpinningVisible(b)
		case fn := <-w.chFunc:
			fn()
			break
		default:
			if w.currentMessagePtr != 0 && time.Now().After(w.showMessageDeadline) {
				w.HideMessage()
			}
		}
	}
}

func onKeyDown(keycode int) bool {
	ret := false
	if w != nil {
		for _, fn := range w.FuncKeyDown {
			b := fn(keycode)
			if b {
				ret = true
			}
		}
	}
	return ret
}

func onProgressChange(typ int, position float64) {
	if w != nil {
		for _, fn := range w.FuncOnProgressChanged {
			fn(typ, position)
		}
	}
}

func (w *Window) SendAlert(str string) {
	w.chFunc <- func() {
		w.Alert(str)
	}
}

func SendAddRecentOpenedFile(filename string) {
	w.chFunc <- func() {
		AddRecentOpenedFile(filename)
	}
}

func SendSetSubFontSize(sz float64) {
	w.chFunc <- func() {
		w.SetSubFontSize(sz)
	}
}
