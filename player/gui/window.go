package gui

import (
	"log"
	"math"
	"time"
	. "vger/player/shared"

	"github.com/go-gl/gl"
)

var w *Window // current window

type window struct {
	FuncTimerTick           []func()
	FuncKeyDown             []func(int) bool
	FuncOnProgressChanged   []func(int, float64)
	FuncAudioMenuClicked    []func(int)
	FuncSubtitleMenuClicked []func(int)
	FuncMouseWheelled       []func(float64)

	chAlert chan string

	ChanDraw     chan []byte
	ChanShowText chan SubItemArg
	ChanSetSize  chan argSize
	ChanSetTitle chan string

	ChanShowMessage chan SubItemArg
	ChanHideMessage chan uintptr

	ChanDestoryRender chan struct{}

	ChanShowProgress chan *struct {
		left    string
		right   string
		percent float64
	}
	ChanShowSpeed chan *BufferInfo

	ChanSetCursor        chan bool
	ChanSetVolume        chan byte
	ChanSetVolumeDisplay chan bool

	ChanShowSpinning chan bool

	img []byte

	originalWidth  int
	originalHeight int

	currentMessagePtr uintptr
	currentMessage    *SubItem

	render imageRender

	forceRatio float64

	showMessageDeadline time.Time

	chCursor         chan struct{}
	chCursorAutoHide chan struct{}

	chDelayShowSpinning chan int
}

type imageRender interface {
	draw(img []byte, width, height int)
	delete()
}
type argSize struct {
	width, height int
}

func (w *Window) SendDrawImage(img []byte) {
	w.ChanDraw <- img
}
func (w *Window) SendSetCursor(b bool) {
	go func() {
		w.ChanSetCursor <- b
	}()

	if b {
		w.chCursor <- struct{}{}
	} else {
		<-w.chCursor
	}
}

func (w *Window) DestoryRender() {
	if w.render != nil {
		w.render.delete()
		w.render = nil
	}
}

func (w *Window) SendDestoryRender() {
	log.Print("SendDestoryRender")

	w.ChanDestoryRender <- struct{}{}
}

func (w *Window) IsFullScreen() bool {
	width, height := w.GetSize()
	swidth, sheight := getScreenSize()

	return width == swidth && height == sheight
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
	w.ShowStartupView()

	log.Printf("set window size:%d %d", width, height)

	w.ChanDraw = make(chan []byte)

	if width%4 != 0 {
		gl.PixelStorei(gl.UNPACK_ALIGNMENT, 1)
	}

	w.render = NewYUVRender(width, height)

	w.originalWidth, w.originalHeight = width, height

	if w.IsFullScreen() {
		return
	}

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
			ChanDraw: make(chan []byte),
			ChanShowProgress: make(chan *struct {
				left    string
				right   string
				percent float64
			}),
			ChanShowSpeed: make(chan *BufferInfo),
			ChanShowText:  make(chan SubItemArg, 20), //the buffer is required because show&hide must handles in the same order
			ChanSetSize:   make(chan argSize),
			ChanSetTitle:  make(chan string),

			ChanShowMessage: make(chan SubItemArg),
			ChanHideMessage: make(chan uintptr),

			ChanSetCursor:       make(chan bool),
			ChanShowSpinning:    make(chan bool),
			chDelayShowSpinning: nil,

			ChanSetVolume:        make(chan byte),
			ChanSetVolumeDisplay: make(chan bool),

			ChanDestoryRender: make(chan struct{}),

			originalWidth:  width,
			originalHeight: height,

			chCursor:         make(chan struct{}),
			chCursorAutoHide: make(chan struct{}),
			chAlert:          make(chan string),
		},
	}

	log.Print("NewWindow:", w.NativeWindow)

	w.Show()
	w.MakeCurrentContext() //must make current context before do texture bind or we will get a all white window
	gl.Init()
	// gl.ClearColor(0, 0, 0, 1)

	w.initEvents()

	return w
}

func (w *Window) initEvents() {
	w.FuncOnProgressChanged = append(w.FuncOnProgressChanged, func(typ int, percent float64) { //run in main thread, safe to operate ui elements
		if typ == 0 || typ == 2 {
			w.chCursorAutoHide <- struct{}{}
		}
	})

	go func() {
		for {
			select {
			case <-time.After(2 * time.Second):
				w.SendSetCursor(false)
				break
			case <-w.chCursor:
				break
			case <-w.chCursorAutoHide:
				<-w.chCursorAutoHide
				break
			}
		}
	}()
}

func (w *Window) ClearEvents() {
	w.FuncOnProgressChanged = w.FuncOnProgressChanged[:1]
	w.FuncKeyDown = nil
	w.FuncAudioMenuClicked = nil
	w.FuncSubtitleMenuClicked = nil
	w.FuncMouseWheelled = nil
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

	// gl.Clear(gl.COLOR_BUFFER_BIT)
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

func (w *Window) SendShowProgress(left string, right string, percent float64) {
	w.ChanShowProgress <- &struct {
		left    string
		right   string
		percent float64
	}{left, right, percent}
}
func (w *Window) SendShowBufferInfo(info *BufferInfo) {
	w.ChanShowSpeed <- info
}

// func (w *Window) ShowProgress(p *PlayProgressInfo) {
// 	w.ShowWindowProgress(w.ptr, p.Right, p.Left, p.Percent)
// }

func (w *Window) SendShowText(s SubItemArg) {
	// res := make(chan SubItemExtra)
	w.ChanShowText <- s
	// return <-res
}
func createMessageSubItem(msg string) SubItem {
	s := SubItem{}
	s.PositionType = 7
	s.X = 20
	s.Y = 20
	s.Content = make([]AttributedString, 0)
	s.Content = append(s.Content, AttributedString{msg, 3, 0xffffff})

	return s
}
func (w *Window) SendShowMessage(msg string, autoHide bool) {
	s := createMessageSubItem(msg)

	w.ChanShowMessage <- SubItemArg{s, autoHide, nil}
}

func (w *Window) ShowMessage(msg string, autoHide bool) {
	s := createMessageSubItem(msg)
	w.showMessage(&s, autoHide)
}

func (w *Window) showMessage(s *SubItem, autoHide bool) {
	if w.currentMessagePtr != 0 {
		w.HideText(w.currentMessagePtr)
	}

	w.currentMessagePtr = w.ShowText(s)
	w.currentMessage = s

	if autoHide {
		w.showMessageDeadline = time.Now().Add(2 * time.Second)
	} else {
		w.showMessageDeadline = time.Now().Add(1000 * time.Hour)
	}
}

func (w *Window) HideMessage() {
	if w.currentMessagePtr != 0 {
		w.HideText(w.currentMessagePtr)
		w.currentMessage = nil
		w.currentMessagePtr = 0
	}
	w.currentMessage = nil
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

	return w.NativeWindow.ShowText(items, s.PositionType, s.X, s.Y)
}

func (w *Window) SendHideText(arg SubItemArg) {
	w.ChanShowText <- arg
}

func (w *Window) SendShowSpinning() {
	// log.Print(string(debug.Stack()))

	if w.chDelayShowSpinning == nil {
		w.chDelayShowSpinning = make(chan int)
		go func() {
			w.ChanShowSpinning <- true
			i := 1
			delta := <-w.chDelayShowSpinning
			if delta == 0 {
				i = 0
			} else {
				i += delta
			}

			for {
				// log.Print(i)
				select {
				case <-time.After(500 * time.Millisecond):
					w.ChanShowSpinning <- (i > 0)
					delta := <-w.chDelayShowSpinning
					if delta == 0 {
						i = 0
					} else {
						i += delta
					}
				case delta := <-w.chDelayShowSpinning:
					if delta == 0 {
						i = 0
					} else {
						i += delta
					}
				}
			}
		}()
	} else {
		w.chDelayShowSpinning <- 1
	}
}
func (w *Window) SendHideSpinning(forceHide bool) {
	// log.Print(string(debug.Stack()))

	if forceHide {
		w.ChanShowSpinning <- false
		if w.chDelayShowSpinning != nil {
			w.chDelayShowSpinning <- 0
		}
	} else {
		if w.chDelayShowSpinning != nil {
			w.chDelayShowSpinning <- -1
		}
	}
}

func (w *Window) SendSetVolume(volume byte) {
	w.ChanSetVolume <- volume
}

func (w *Window) SendSetVolumeDisplay(b bool) {
	w.ChanSetVolumeDisplay <- b
}

func (w *Window) Refresh(img []byte) {
	w.img = img
	w.RefreshContent()
}

func onDraw() {
	if w != nil {
		w.draw(w.img, w.originalWidth, w.originalHeight)
	}
}

func onTimerTick() {
	if w != nil {
		select {
		case img, ok := <-w.ChanDraw:
			if ok {
				w.Refresh(img)
			}
		case <-w.ChanDestoryRender:
			w.DestoryRender()
		case str := <-w.chAlert:
			w.Alert(str)
		default:
		}

		select {
		case b := <-w.ChanShowSpinning:
			if b {
				w.ShowSpinning()
			} else {
				w.HideSpinning()
			}
		case p := <-w.ChanShowProgress:
			w.ShowProgress(p.left, p.right, p.percent)
		case info := <-w.ChanShowSpeed:
			w.ShowBufferInfo(info.Speed, info.BufferPercent)
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
		case volume := <-w.ChanSetVolume:
			w.SetVolume(volume)
		case b := <-w.ChanSetVolumeDisplay:
			w.SetVolumeDisplay(b)
		case arg := <-w.ChanSetSize:
			w.SetSize(arg.width, arg.height)
			break
		case title := <-w.ChanSetTitle:
			w.SetTitle(title)
			break
		case arg := <-w.ChanShowText:
		skip:
			for {
				if arg.Handle == 0 || arg.Handle == 1 {
					arg.Result <- SubItemExtra{arg.Id, w.ShowText(&arg.SubItem)}
				} else {
					w.HideText(arg.Handle)
				}
				select {
				case arg = <-w.ChanShowText:
				default:
					break skip
				}
			}
			break
		case arg := <-w.ChanShowMessage:
			w.showMessage(&arg.SubItem, arg.AutoHide)
		case <-w.ChanHideMessage:
			w.HideMessage()
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

func onMenuClick(typ int, tag int) {
	if w != nil {
		switch typ {
		case 0:
			for _, fn := range w.FuncAudioMenuClicked {
				fn(tag)
			}
		case 1:
			for _, fn := range w.FuncSubtitleMenuClicked {
				fn(tag)
			}
		case 2:
			onSearchSubtitleMenuItemClick()
		}
	}
}

func onMouseWheel(deltaX float64, deltaY float64) {
	if w != nil {
		for _, fn := range w.FuncMouseWheelled {
			fn(deltaY)
		}
	}
}

func onMouseMove(x, y int) {
	if w != nil {
		w.SendSetCursor(true)
	}
}

func (w *Window) SendAlert(str string) {
	w.chCursorAutoHide <- struct{}{}
	w.chAlert <- str
}
