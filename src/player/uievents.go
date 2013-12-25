package main

import (
	"fmt"
	"player/gui"
	. "player/libav"
	"time"
)

func (m *movie) uievents() {
	m.v.window.FuncAudioMenuClicked = append(m.v.window.FuncAudioMenuClicked, func(i int) {
		go func() {
			m.a.setCurrentStream(i)
		}()
	})
	m.v.window.FuncKeyDown = append(m.v.window.FuncKeyDown, func(keycode int) {
		switch keycode {
		case gui.KEY_SPACE:
			m.c.Toggle()
			break
		case gui.KEY_LEFT:
			println("key left pressed")
			m.c.SetTime(m.SeekTo(m.c.GetSeekTime() - time.Second))
			break
		case gui.KEY_RIGHT:
			m.c.SetTime(m.SeekTo(m.c.GetSeekTime() + time.Second))
			break
		case gui.KEY_UP:
			m.c.SetTime(m.SeekTo(m.c.GetSeekTime() + 10*time.Second))
			break
		case gui.KEY_DOWN:
			m.c.SetTime(m.SeekTo(m.c.GetSeekTime() - 10*time.Second))
			break
		case gui.KEY_MINUS:
			println("key minus pressed")
			go func() {
				if m.s != nil {
					offset := m.s.AddOffset(200 * time.Millisecond)
					m.v.window.SendShowMessage(fmt.Sprint("Subtitle offset ", offset.String()))
				}
				if m.s2 != nil {
					m.s2.AddOffset(200 * time.Millisecond)
				}
			}()
			break
		case gui.KEY_EQUAL:
			println("key equal pressed")
			go func() {
				if m.s != nil {
					offset := m.s.AddOffset(-200 * time.Millisecond)
					m.v.window.SendShowMessage(fmt.Sprint("Subtitle offset ", offset.String()))
				}
				if m.s2 != nil {
					m.s2.AddOffset(-200 * time.Millisecond)
				}
			}()
			break
		case gui.KEY_LEFT_BRACKET:
			println("left bracket pressed")
			go func() {
				if m.s != nil {
					offset := m.s.AddOffset(1000 * time.Millisecond)
					m.v.window.SendShowMessage(fmt.Sprint("Subtitle offset ", offset.String()))
				}
				if m.s2 != nil {
					m.s2.AddOffset(1000 * time.Millisecond)
				}
			}()
			break
		case gui.KEY_RIGHT_BRACKET:
			println("right bracket pressed")
			go func() {
				if m.s != nil {
					offset := m.s.AddOffset(-1000 * time.Millisecond)
					m.v.window.SendShowMessage(fmt.Sprint("Subtitle offset ", offset.String()))
				}
				if m.s2 != nil {
					m.s2.AddOffset(-1000 * time.Millisecond)
				}
			}()
			break
		}
	})

	var lastSeekTime time.Duration
	var lastText uintptr

	m.v.window.FuncOnProgressChanged = append(m.v.window.FuncOnProgressChanged, func(typ int, percent float64) { //run in main thread, safe to operate ui elements
		switch typ {
		case 0:
			lastSeekTime = m.c.GetSeekTime()

			m.c.Pause()
			time.Sleep(5 * time.Millisecond)
			break
		case 2:
			if lastText != 0 {
				m.v.window.HideText(lastText)
				lastText = 0
			}
			t := m.c.CalcTime(percent)
			m.c.ResumeWithTime(m.SeekTo(t))
			// time.Sleep(5 * time.Millisecond)
			break
		case 1:
			t := m.c.CalcTime(percent)
			flags := AVSEEK_FLAG_FRAME
			if t < lastSeekTime {
				flags |= AVSEEK_FLAG_BACKWARD
			}
			// m.ctx.SeekFrame(m.v.stream, t, flags)
			m.ctx.SeekFile(t, flags)
			lastSeekTime = t

			codec := m.v.stream.Codec()
			codec.FlushBuffer()
			m.drawCurrentFrame()

			if m.s != nil {
				if _, item := m.s.FindPos(t); item != nil {
					if lastText != 0 {
						m.v.window.HideText(lastText)
						lastText = 0
					}

					lastText = m.v.window.ShowText(item)
				} else {
					if lastText != 0 {
						m.v.window.HideText(lastText)
						lastText = 0
					}
				}
			}

			break
		}

		m.v.window.ShowProgress(m.c.CalcPlayProgress(percent))
	})
}
