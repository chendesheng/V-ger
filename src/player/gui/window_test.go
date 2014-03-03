package gui

import (
	"testing"
)

func TestShowWindow(t *testing.T) {
	w := NewWindow("hello", 1280, 720)

	w.FuncKeyDown = append(w.FuncKeyDown, func(key int) {
		w.SetSize(500, 500)
		w.SetTitle("test")
	})

	PollEvents()
}
func fequal(a, b float64) bool {
	d := a - b
	return d < 1e5 && d > -1e5
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
	w, h := GetScreenSize()
	if w != 1920 || h != 1080 {
		t.Errorf("Expect 1920x1080 but %dx%d", w, h)
	}
}

func TestIsFullScreen(t *testing.T) {
	w := NewWindow("hello", 1280, 720)

	if w.IsFullScreen() {
		t.Errorf("should not full screen")
	}

	w.FuncOnFullscreenChanged = append(w.FuncOnFullscreenChanged, func(b bool) {
		if b {
			if !w.IsFullScreen() {
				t.Errorf("should full screen")
			} else {
				println("full screen")
			}
		}
	})

	PollEvents()
}

func TestStartupView(t *testing.T) {
	w := NewWindow("hello", 1280, 720)

	w.FuncKeyDown = append(w.FuncKeyDown, func(key int) {
		if key == KEY_A {
			w.hideStartupView()
		} else if key == KEY_B {
			w.ShowStartupView()
		}
	})

	PollEvents()
}

func TestMenu(t *testing.T) {
	w := NewWindow("hello", 1280, 720)

	w.FuncKeyDown = append(w.FuncKeyDown, func(key int) {
		if key == KEY_A {
			HideSubtitleMenu()
			HideAudioMenu()
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
	})

	PollEvents()
}
