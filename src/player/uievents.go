package main

import (
	"fmt"
	// "log"
	"player/gui"
	// . "player/video"
	// . "player/libav"
	"time"
)

func (m *movie) uievents() {
	m.w.FuncAudioMenuClicked = append(m.w.FuncAudioMenuClicked, func(i int) {
		go func() {
			m.a.setCurrentStream(i)
		}()
	})

	// var chPausing chan seekArg
	// chPausing = nil
	// var pausingTime time.Duration
	m.w.FuncKeyDown = append(m.w.FuncKeyDown, func(keycode int) {
		switch keycode {
		case gui.KEY_SPACE:
			m.c.Toggle()
			// if chPausing == nil {
			// 	pausingTime = m.c.GetTime()

			// 	res := make(chan interface{})
			// 	arg := ctrlArg{PAUSE, res}
			// 	m.chCtrl <- arg
			// 	m.v.FlushBuffer()
			// 	m.a.flushBuffer()
			// 	ch := <-arg.res
			// 	chPausing = ch.(chan seekArg)
			// } else {
			// 	m.v.SeekOffset(pausingTime) //this is because we have a video queue which always ahead than current time.
			// 	chPausing <- seekArg{pausingTime, nil}
			// 	chPausing = nil
			// }
			break
		case gui.KEY_LEFT:
			m.seekOffset(-10 * time.Second)
			break
		case gui.KEY_RIGHT:
			m.seekOffset(10 * time.Second)
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
	var chPause chan seekArg

	m.w.FuncOnProgressChanged = append(m.w.FuncOnProgressChanged, func(typ int, percent float64) { //run in main thread, safe to operate ui elements
		switch typ {
		case 0:
			// lastSeekTime = m.c.GetSeekTime()

			// time.Sleep(5 * time.Millisecond)

			// m.c.Pause()
			// m.c.SetTime(lastSeekTime)
			// m.a.Pause(true)

			res := make(chan interface{})
			arg := ctrlArg{PAUSE, res}
			m.chCtrl <- arg
			ch := <-arg.res
			m.v.FlushBuffer()
			m.a.flushBuffer()
			chPause = ch.(chan seekArg)

			t := m.c.CalcTime(percent)
			t, _ = m.v.Seek(t)
			// var img []byte
			// t, img = m.getCurrentFrame()
			// m.w.RefreshContent(img)
			lastSeekTime = t
			m.c.SetTime(t)
			percent = m.c.GetPercent()
			m.w.ShowProgress(m.c.CalcPlayProgress(percent))
			break
		case 2:
			if lastText != 0 {
				m.w.HideText(lastText)
				lastText = 0
			}
			chPause <- seekArg{lastSeekTime, nil}

			// t := m.c.CalcTime(percent)
			// m.a.Pause(false)
			// m.c.ResumeWithTime(lastSeekTime)
			// m.v.FlushBuffer()
			// res := make(chan bool)
			// m.chPause <- seekArg{lastSeekTime, res}
			// <-res
			// println("before resume")
			// m.c.Resume()

			// m.c.ResumeWithTime(m.SeekTo(lastSeekTime))
			// time.Sleep(5 * time.Millisecond)

			// close(chPause)
			break
		case 1:
			t := m.c.CalcTime(percent)
			var img []byte
			t, img = m.v.Seek(t)
			println("seeking pts:", t.String())
			// img1 := make([]byte, len(img))
			m.w.RefreshContent(img)
			lastSeekTime = t
			// chPause <- seekArg{t, nil}

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

			m.c.SetTime(t)
			percent = m.c.GetPercent()
			m.w.ShowProgress(m.c.CalcPlayProgress(percent))

			break
		}

	})
}
