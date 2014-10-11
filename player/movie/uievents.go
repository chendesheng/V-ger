package movie

import (
	"log"
	"time"
	"vger/player/gui"
	. "vger/player/shared"
)

var chVolume = make(chan struct{})

func (m *Movie) uievents() {
	log.Print("movie uievents")

	m.w.FuncKeyDown = append(m.w.FuncKeyDown, func(keycode int) bool {
		SavePlayingAsync(m.p)

		switch keycode {
		case gui.KEY_R:
			if !m.w.IsFullScreen() {
				m.w.ToggleForceScreenRatio()
			}
			break
		default:
			return false
		}

		return true
	})

	m.w.FuncMouseWheelled = append(m.w.FuncMouseWheelled, func(deltaY float64) {
		if deltaY == 0 {
			return
		}
		m.AddVolume(int(deltaY * -10))

	})

	go func() {
		for {
			select {
			case <-time.After(time.Second):
				m.w.SendSetVolumeVisible(false)
				<-chVolume //prevent call SendSetVolumeDisplay every second
			case <-chVolume:
			case <-m.quit:
				return
			}
		}
	}()
}

func (m *Movie) uiProgressBarEvents() {

	m.w.FuncOnProgressChanged = append(m.w.FuncOnProgressChanged, func(typ int, percent float64) { //run in main thread, safe to operate ui elements
		if m.c == nil {
			return
		}
		switch typ {
		case 0:
			m.w.SetControlsVisible(true, false)
			fallthrough
		case 1:
			t := m.c.CalcTime(percent)
			t = t / time.Second * time.Second
			p := m.c.CalcPlayProgress(t)
			m.w.UpdatePlaybackInfo(p.Left, p.Right, p.Percent)

			m.seeking.SendSeek(t)
		case 2:
			m.w.SetControlsVisible(true, true)

			t := m.c.CalcTime(percent)
			t = t / time.Second * time.Second

			log.Print("release dragging:", t.String())

			m.seeking.SendEndSeek(t)
		}
	})
}

func isLangEqual(l1, l2 string) bool {
	if l1 == l2 {
		return true
	} else if l1 == "en" || l2 == "en" {
		return false
	} else {
		return true
	}
}
