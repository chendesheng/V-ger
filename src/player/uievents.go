package main

import (
	"fmt"
	. "player/shared"
	// "log"
	"player/gui"
	// . "player/video"
	// . "player/libav"
	"log"
	"time"
)

func (m *movie) uievents() {
	m.w.FuncAudioMenuClicked = append(m.w.FuncAudioMenuClicked, func(i int) {
		go func() {
			log.Printf("Audio menu click:%d", i)

			m.p.SoundStream = i
			SavePlayingAsync(m.p)

			m.a.Close()
			m.a.setCurrentStream(i)
		}()
	})

	// var chPausing chan seekArg
	// chPausing = nil
	// var pausingTime time.Duration
	m.w.FuncKeyDown = append(m.w.FuncKeyDown, func(keycode int) {
		SavePlayingAsync(m.p)

		switch keycode {
		case gui.KEY_SPACE:
			m.c.Toggle()
			break
		case gui.KEY_LEFT:
			var offset time.Duration
			if m.s != nil {
				t := m.c.GetTime()
				subTime := m.s.GetSubtime(t, -1)

				if subTime == 0 {
					offset = -10 * time.Second
				} else {
					offset = subTime - t
				}
			} else {
				offset = -10 * time.Second
			}
			m.seekOffset(offset)
			break
		case gui.KEY_RIGHT:
			var offset time.Duration
			if m.s != nil {
				t := m.c.GetTime()
				subTime := m.s.GetSubtime(t, 1)
				println("subtime:", subTime)

				if subTime == 0 {
					offset = 10 * time.Second
				} else {
					offset = subTime - t
				}
			} else {
				offset = 10 * time.Second
			}
			m.seekOffset(offset)
			break
		case gui.KEY_UP:
			m.seekOffset(-5 * time.Second)
			break
		case gui.KEY_DOWN:
			m.seekOffset(5 * time.Second)
			break
		case gui.KEY_MINUS:
			println("key minus pressed")
			go func() {
				if m.s != nil {
					offset := m.s.AddOffset(-200 * time.Millisecond)
					m.w.SendShowMessage(fmt.Sprint("Subtitle offset ", offset.String()), true)

					UpdateSubtitleOffsetAsync(m.s.Name, offset)
				}
			}()
			break
		case gui.KEY_EQUAL:
			println("key equal pressed")
			go func() {
				if m.s != nil {
					offset := m.s.AddOffset(200 * time.Millisecond)
					m.w.SendShowMessage(fmt.Sprint("Subtitle offset ", offset.String()), true)

					UpdateSubtitleOffsetAsync(m.s.Name, offset)
				}
			}()
			break
		case gui.KEY_LEFT_BRACKET:
			println("left bracket pressed")
			go func() {
				// if m.s != nil {
				// 	offset := m.s.AddOffset(-1000 * time.Millisecond)
				// 	m.w.SendShowMessage(fmt.Sprint("Subtitle offset ", offset.String()), true)
				// }
				if m.s2 != nil {
					offset := m.s2.AddOffset(-200 * time.Millisecond)
					m.w.SendShowMessage(fmt.Sprint("Subtitle 2 offset ", offset.String()), true)

					UpdateSubtitleOffsetAsync(m.s2.Name, offset)
				}
			}()
			break
		case gui.KEY_RIGHT_BRACKET:
			println("right bracket pressed")
			go func() {
				// if m.s != nil {
				// 	offset := m.s.AddOffset(1000 * time.Millisecond)
				// 	m.w.SendShowMessage(fmt.Sprint("Subtitle offset ", offset.String()), true)
				// }
				if m.s2 != nil {
					offset := m.s2.AddOffset(200 * time.Millisecond)
					m.w.SendShowMessage(fmt.Sprint("Subtitle 2 offset ", offset.String()), true)

					UpdateSubtitleOffsetAsync(m.s2.Name, offset)
				}
			}()
			break
		case gui.KEY_ESCAPE:
			m.w.ToggleFullScreen()
			break
		}
	})

	var lastSeekTime time.Duration
	// var lastText uintptr
	// var chPause chan seekArg

	m.w.FuncOnProgressChanged = append(m.w.FuncOnProgressChanged, func(typ int, percent float64) { //run in main thread, safe to operate ui elements
		switch typ {
		case 0:
			m.SeekBegin()
			t := m.Seek(m.c.CalcTime(percent))

			lastSeekTime = t
			break
		case 2:
			t := lastSeekTime
			m.SeekEnd(t)

			m.p.LastPos = m.c.GetTime()
			SavePlayingAsync(m.p)
			break
		case 1:
			t := m.c.CalcTime(percent)
			t = m.Seek(t)
			lastSeekTime = t

			break
		}

	})

	m.w.FuncOnFullscreenChanged = append(m.w.FuncOnFullscreenChanged, func(b bool) {
		if m.s != nil {
			t := m.c.GetTime()
			m.s.SeekRefresh(t)
		}

		if m.s2 != nil {
			t := m.c.GetTime()
			m.s2.SeekRefresh(t)
		}
	})
}
