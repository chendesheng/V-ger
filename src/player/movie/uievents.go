package movie

import (
	"fmt"
	. "player/shared"
	. "player/subtitle"
	// "log"
	"player/gui"
	// . "player/video"
	// . "player/libav"
	"log"
	"time"
)

func (m *Movie) uievents() {
	m.w.FuncAudioMenuClicked = append(m.w.FuncAudioMenuClicked, func(i int) {
		go func() {
			log.Printf("Audio menu click:%d", i)

			m.p.SoundStream = i
			SavePlayingAsync(m.p)

			m.a.Close()
			m.a.Open(getStream(m.audioStreams, i))
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
		case gui.KEY_R:
			m.w.ToggleForceScreenRatio()
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
			m.seekOffsetAsync(offset)
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
			m.seekOffsetAsync(offset)
			break
		case gui.KEY_UP:
			m.seekOffsetAsync(5 * time.Second)
			break
		case gui.KEY_DOWN:
			m.seekOffsetAsync(-5 * time.Second)
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

	chCursor := make(chan struct{})
	chCursorAutoHide := make(chan struct{})

	var lastSeekTime time.Duration
	// var lastText uintptr
	// var chPause chan seekArg
	chProgress := make(chan time.Duration)
	chProgressRes := make(chan time.Duration)
	go func() {
		var t time.Duration
		for {
			select {
			case <-time.After(100 * time.Millisecond):
				if t > 0 {
					chProgressRes <- m.Seek(t)
					t = 0
				}
				break
			case t = <-chProgress:
				chProgressRes <- t
				break
			}
		}
	}()
	m.w.FuncOnProgressChanged = append(m.w.FuncOnProgressChanged, func(typ int, percent float64) { //run in main thread, safe to operate ui elements
		// select {
		// case chCursor <- struct{}{}:
		// 	break
		// case <-time.After(50 * time.Millisecond):
		// 	log.Print("stop hide cursor timeout2")
		// 	break
		// }

		switch typ {
		case 0:
			m.SeekBegin()
			t := m.c.CalcTime(percent)
			if m.httpBuffer == nil {
				t = m.Seek(t)
			}

			lastSeekTime = t

			chCursorAutoHide <- struct{}{}
			break
		case 2:
			go func() {
				if m.httpBuffer != nil {
					m.w.SendShowMessage("Bufferring...", false)
					defer m.w.SendHideMessage()

					lastSeekTime = m.Seek(lastSeekTime)
					m.httpBuffer.Wait(1024 * 1024)
				}
				m.SeekEnd(lastSeekTime)
				m.p.LastPos = m.c.GetTime()
				SavePlaying(m.p)
			}()

			chCursorAutoHide <- struct{}{}

			break
		case 1:
			t := m.c.CalcTime(percent)
			if m.httpBuffer == nil {
				t = m.Seek(t)
			}
			lastSeekTime = t

			m.c.SetTime(t)
			go m.showProgress()
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

	var chVolume chan byte
	m.w.FuncMouseWheelled = append(m.w.FuncMouseWheelled, func(deltaY float64) {
		if chVolume == nil {
			chVolume = make(chan byte, 100)
			go func() {
				for {
					select {
					case <-time.Tick(time.Second * 2):
						m.w.SendHideMessage()

						close(chVolume)
						chVolume = nil
						return
					case volume := <-chVolume:
						m.p.Volume = volume
						SavePlayingAsync(m.p)
						m.w.SendShowMessage(fmt.Sprintf("Volume: %d%%", volume), false)
						break
					}
				}
			}()
		}
		if deltaY == 0 {
			return
		}

		if deltaY > 0 {
			chVolume <- m.a.DecreaseVolume()
		} else {
			chVolume <- m.a.IncreaseVolume()
		}
	})

	m.w.FuncSubtitleMenuClicked = append(m.w.FuncSubtitleMenuClicked, func(index int, showOrHide bool) {
		go func() {
			subs := m.subs
			clicked := subs[index]
			log.Print("toggle subtitle:", clicked.Name, index)
			if showOrHide {
				// m.s.Stop()
				width, height := m.w.GetWindowSize()
				s := NewSubtitle(clicked, m.w, m.c, float64(width), float64(height))
				if s != nil {
					if m.s == nil {
						m.s = s
						s.IsMainOrSecondSub = true
					} else {
						m.s2 = s
						s.IsMainOrSecondSub = false
					}
					go s.Play()
				}
			} else {
				if (m.s != nil) && (m.s.Name == clicked.Name) {
					m.s.Stop()
					if m.s2 != nil {
						m.s = m.s2
						m.s.IsMainOrSecondSub = true
						m.s2 = nil
					} else {
						m.s = nil
					}
				} else if (m.s2 != nil) && (m.s2.Name == clicked.Name) {
					m.s2.Stop()
					m.s2 = nil
				}
			}

			if m.s != nil {
				m.p.Sub1 = m.s.Name
			} else {
				m.p.Sub1 = ""
			}

			if m.s2 != nil {
				m.p.Sub2 = m.s2.Name
			} else {
				m.p.Sub2 = ""
			}

			SavePlayingAsync(m.p)
		}()
	})

	go func() {
		for {
			select {
			case <-time.After(2 * time.Second):
				m.w.SendSetCursor(false)
				<-chCursor //prevent call SendSetCursor every 2 seconds
				break
			case <-chCursor:
				break
			case <-chCursorAutoHide:
				<-chCursorAutoHide
				break
			}
		}
	}()
	m.w.FuncMouseMoved = append(m.w.FuncMouseMoved, func() {
		select {
		case chCursor <- struct{}{}:
			break
		case <-time.After(50 * time.Millisecond):
			log.Print("stop hide cursor timeout")
			break
		}
	})
}
