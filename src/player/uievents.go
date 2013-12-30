package main

import (
	"fmt"
	"player/gui"
	// . "player/libav"
	"time"
)

func (m *movie) uievents() {
	m.w.FuncAudioMenuClicked = append(m.w.FuncAudioMenuClicked, func(i int) {
		go func() {
			m.a.setCurrentStream(i)
		}()
	})
	m.w.FuncKeyDown = append(m.w.FuncKeyDown, func(keycode int) {
		switch keycode {
		case gui.KEY_SPACE:
			m.c.Toggle()
			break
		case gui.KEY_LEFT:
			println("key left pressed")
			// m.c.SetTime(m.SeekTo(m.c.GetSeekTime() - 10*time.Second))
			m.chSeek <- m.c.GetSeekTime() - 10*time.Second
			break
		case gui.KEY_RIGHT:
			// m.c.SetTime(m.SeekTo(m.c.GetSeekTime() + 10*time.Second))
			m.chSeek <- m.c.GetSeekTime() + 10*time.Second
			break
		case gui.KEY_UP:
			// m.c.SetTime(m.SeekTo(m.c.GetSeekTime() + time.Second))
			m.chSeek <- m.c.GetSeekTime() + time.Second
			break
		case gui.KEY_DOWN:
			// m.c.SetTime(m.SeekTo(m.c.GetSeekTime() - time.Second))
			m.chSeek <- m.c.GetSeekTime() - time.Second
			break
		case gui.KEY_MINUS:
			println("key minus pressed")
			go func() {
				if m.s != nil {
					offset := m.s.AddOffset(-200 * time.Millisecond)
					m.w.SendShowMessage(fmt.Sprint("Subtitle offset ", offset.String()))
				}
				if m.s2 != nil {
					m.s2.AddOffset(-200 * time.Millisecond)
				}
			}()
			break
		case gui.KEY_EQUAL:
			println("key equal pressed")
			go func() {
				if m.s != nil {
					offset := m.s.AddOffset(200 * time.Millisecond)
					m.w.SendShowMessage(fmt.Sprint("Subtitle offset ", offset.String()))
				}
				if m.s2 != nil {
					m.s2.AddOffset(200 * time.Millisecond)
				}
			}()
			break
		case gui.KEY_LEFT_BRACKET:
			println("left bracket pressed")
			go func() {
				if m.s != nil {
					offset := m.s.AddOffset(-1000 * time.Millisecond)
					m.w.SendShowMessage(fmt.Sprint("Subtitle offset ", offset.String()))
				}
				if m.s2 != nil {
					m.s2.AddOffset(-1000 * time.Millisecond)
				}
			}()
			break
		case gui.KEY_RIGHT_BRACKET:
			println("right bracket pressed")
			go func() {
				if m.s != nil {
					offset := m.s.AddOffset(1000 * time.Millisecond)
					m.w.SendShowMessage(fmt.Sprint("Subtitle offset ", offset.String()))
				}
				if m.s2 != nil {
					m.s2.AddOffset(1000 * time.Millisecond)
				}
			}()
			break
		}
	})

	var lastSeekTime time.Duration
	var lastText uintptr

	m.w.FuncOnProgressChanged = append(m.w.FuncOnProgressChanged, func(typ int, percent float64) { //run in main thread, safe to operate ui elements
		switch typ {
		case 0:
			// lastSeekTime = m.c.GetSeekTime()
			lastSeekTime = m.c.CalcTime(percent)

			m.c.Pause()
			m.a.Pause(true)
			break
		case 2:
			if lastText != 0 {
				m.w.HideText(lastText)
				lastText = 0
			}
			// t := m.c.CalcTime(percent)
			// m.a.Pause(false)
			// m.c.ResumeWithTime(lastSeekTime)
			m.c.Resume()
			m.chSeek <- lastSeekTime

			// m.c.ResumeWithTime(m.SeekTo(lastSeekTime))
			// time.Sleep(5 * time.Millisecond)
			break
		case 1:
			t := m.c.CalcTime(percent)
			// flags := AVSEEK_FLAG_FRAME
			// if t < lastSeekTime {
			// 	flags |= AVSEEK_FLAG_BACKWARD
			// }
			// m.ctx.SeekFrame(m.v.stream, t, flags)
			m.v.Seek(t)
			// m.ctx.SeekFile(t, flags)
			lastSeekTime = t

			m.drawCurrentFrame()

			if m.s != nil {
				if _, item := m.s.FindPos(t); item != nil {
					if lastText != 0 {
						m.w.HideText(lastText)
						lastText = 0
					}

					lastText = m.w.ShowText(item)
				} else {
					if lastText != 0 {
						m.w.HideText(lastText)
						lastText = 0
					}
				}
			}

			break
		}

		m.w.ShowProgress(m.c.CalcPlayProgress(percent))
	})
}
