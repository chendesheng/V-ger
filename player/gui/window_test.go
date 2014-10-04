package gui

import (
	"bytes"
	"io/ioutil"
	"runtime"
	"strconv"
	"strings"
	"testing"
	"time"
	. "vger/player/shared"
)

func TestShowWindow(t *testing.T) {
	w := NewWindow("hello", 1280, 720)

	w.FuncKeyDown = append(w.FuncKeyDown, func(key int) bool {
		w.SetSize(500, 500)
		w.SetTitle("test")

		return true
	})

	run()
}
func TestFitToWindowLarger(t *testing.T) {
	w := NewWindow("title", 1280, 720)
	x, y, width, height := w.fitToWindow(1280, 700)
	println(x, y, width, height)
	if !fequal(float64(width)/float64(height), 1280.0/700.0) {
		t.Error("Aspect radio should be equal")
	}
	if width < 1280 || height < 700 {
		t.Error("Must larger than before")
	}
}

func TestFitToWindowSmaller(t *testing.T) {
	w := NewWindow("title", 1280, 720)
	x, y, width, height := w.fitToWindow(1280, 300)
	println(x, y, width, height)
	if !fequal(float64(width)/float64(height), 1280.0/300.0) {
		t.Error("Aspect radio should be equal")
	}
	if width > 1280 || height > 300 {
		t.Error("Must smaller than before")
	}
}

func TestGetScreenSize(t *testing.T) {
	w, h := getScreenSize()
	if w != 1920 || h != 1080 {
		t.Errorf("Expect 1920x1080 but %dx%d", w, h)
	}
}

func TestIsFullScreen(t *testing.T) {
	w := NewWindow("hello", 1280, 720)

	if w.IsFullScreen() {
		t.Errorf("should not full screen")
	}

	// w.FuncOnFullscreenChanged = append(w.FuncOnFullscreenChanged, func(b int) {
	// 	if b != 0 {
	// 		if !w.IsFullScreen() {
	// 			t.Errorf("should full screen")
	// 		} else {
	// 			println("full screen")
	// 		}
	// 	}
	// })

	run()
}

func TestStartupView(t *testing.T) {
	w := NewWindow("hello", 1280, 720)

	w.FuncKeyDown = append(w.FuncKeyDown, func(key int) bool {
		if key == KEY_A {
			w.HideStartupView()
		} else if key == KEY_B {
			w.ShowStartupView()
		}
		return true
	})

	run()
}

func TestShowMessage(t *testing.T) {
	w := NewWindow("title", 1280, 720)
	w.SetSize(1280, 720)
	go w.SendShowMessage("Downloading...", false)
	// w.hideStartupView()
	go func() {
		time.Sleep(time.Second)
		black := make([]byte, 1280*720*3/2)
		// for i := 0; i < 1280*270; i++ {
		// 	black[i] = 16
		// }
		// for i := 1280 * 720; i < 1280*270*3/2; i++ {
		// 	black[i] = 128
		// }
		w.SendDrawImage(black)
	}()
	run()
}

func TestMenu(t *testing.T) {
	w := NewWindow("hello", 1280, 720)

	w.FuncKeyDown = append(w.FuncKeyDown, func(key int) bool {
		if key == KEY_A {
			w.HideSubtitleMenu()
			w.HideAudioMenu()
		} else if key == KEY_B {
			names := make([]string, 0)
			names = append(names, "sub1")

			tags := make([]int32, 0)
			tags = append(tags, 0)

			selected1 := -1
			selected2 := -1

			w.InitSubtitleMenu(names, tags, selected1, selected2)

			w.InitAudioMenu(names, tags, selected1)
		}

		return true
	})

	run()
}
func readPPMFile(file string) (int, int, []byte) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		println(err.Error())
	}
	println(string(data[:16]))

	l1 := bytes.IndexByte(data, '\n')
	println("l1", l1)
	l2 := bytes.IndexByte(data[l1+1:], '\n') + l1 + 1
	println("l2", l2)
	l3 := bytes.IndexByte(data[l2+1:], '\n') + l2 + 1
	println("l3", l3)

	str := string(data[l1+1 : l2])
	strs := strings.Split(str, " ")
	w, _ := strconv.Atoi(strs[0])
	h, _ := strconv.Atoi(strs[1])
	return w, h, data[l3+1:]
}

// func TestOpenGLShader(t *testing.T) {
// 	width, height, img := readPPMFile("window_test_image.ppm")

// 	w := NewWindow("GLSL", width, height)
// 	w.SetSize(width, height)
// 	w.RefreshContent(img)

// 	gl.Init()

// 	program := gl.CreateProgram()
// 	defer program.Delete()
// 	shaderAttachFromFile(program, gl.FRAGMENT_SHADER, `
// void main() {
// 	gl_FragColor=vec4(1.0,0.0,0.0,1.0);
// }
// 		`)
// 	program.Link()
// 	program.Use()

// 	run()
// }

func TestYUV2RGBShader(t *testing.T) {
	img, _ := ioutil.ReadFile("window_test_image.yuv")
	width, height := 1920, 1040

	w := NewWindow("GLSL", width, height)
	w.SetSize(width, height)

	w.Refresh(img)

	run()
}

func TestSpinningView(t *testing.T) {
	w := NewWindow("hello", 1280, 720)

	w.FuncKeyDown = append(w.FuncKeyDown, func(key int) bool {
		if key == KEY_A {
			w.ShowSpinning()
		} else if key == KEY_B {
			w.HideSpinning()
		}

		return true
	})

	run()
}

func TestVolumeView(t *testing.T) {
	w := NewWindow("hello", 1280, 720)

	w.FuncKeyDown = append(w.FuncKeyDown, func(key int) bool {
		if key == KEY_A {
			w.SetVolumeDisplay(true)
		} else if key == KEY_B {
			w.SetVolumeDisplay(false)
		}

		return true
	})

	run()
}

func TestTextView(t *testing.T) {
	runtime.LockOSThread()

	w := NewWindow("title", 1280, 720)
	for i := 1; i < 2; i++ {
		w.ShowText(&SubItem{
			Content:      []AttributedString{AttributedString{"Test subtitle", 0, 0}},
			Position:     Position{-1, -1},
			PositionType: i,
		})
	}

	w.HideStartupView()

	w.FuncKeyDown = append(w.FuncKeyDown, func(keycode int) bool {
		if keycode == KEY_ESCAPE {
			w.ToggleFullScreen()
		}
		return true
	})

	run()
}

func TestTextView2(t *testing.T) {
	runtime.LockOSThread()

	w := NewWindow("title", 1280, 720)
	for i := 1; i < 10; i++ {
		w.ShowText(&SubItem{
			Content:      []AttributedString{AttributedString{"Test subtitle", 0, 0}},
			Position:     Position{300, 300},
			PositionType: i,
		})
	}

	w.HideStartupView()

	w.FuncKeyDown = append(w.FuncKeyDown, func(keycode int) bool {
		if keycode == KEY_ESCAPE {
			w.ToggleFullScreen()
		}
		return true
	})

	run()
}

func TestAlert(t *testing.T) {
	runtime.LockOSThread()

	w := NewWindow("title", 390, 120)
	go w.SendAlert("test")
	run()
}
