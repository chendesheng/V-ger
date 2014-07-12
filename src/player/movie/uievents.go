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
			if !m.w.IsFullScreen() {
				m.w.ToggleForceScreenRatio()
			}
			break
		case gui.KEY_LEFT:
			var offset time.Duration
			s1, _ := m.getPlayingSubs()
			if s1 != nil {
				t := m.c.GetTime()
				subTime := s1.GetSubtime(t, -1)

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
			s1, _ := m.getPlayingSubs()
			if s1 != nil {
				t := m.c.GetTime()
				subTime := s1.GetSubtime(t, 1)
				log.Print("subtime:", subTime)

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
			log.Print("key minus pressed")
			go func() {
				s1, _ := m.getPlayingSubs()
				if s1 != nil {
					offset := s1.AddOffset(-200 * time.Millisecond)
					m.w.SendShowMessage(fmt.Sprint("Subtitle offset ", offset.String()), true)

					UpdateSubtitleOffsetAsync(s1.Name, offset)
				}
			}()
			break
		case gui.KEY_EQUAL:
			log.Print("key equal pressed")
			go func() {
				s1, _ := m.getPlayingSubs()
				if s1 != nil {
					offset := s1.AddOffset(200 * time.Millisecond)
					m.w.SendShowMessage(fmt.Sprint("Subtitle offset ", offset.String()), true)

					UpdateSubtitleOffsetAsync(s1.Name, offset)
				}
			}()
			break
		case gui.KEY_LEFT_BRACKET:
			log.Print("left bracket pressed")
			go func() {
				// if m.s != nil {
				// 	offset := m.s.AddOffset(-1000 * time.Millisecond)
				// 	m.w.SendShowMessage(fmt.Sprint("Subtitle offset ", offset.String()), true)
				// }
				_, s2 := m.getPlayingSubs()
				if s2 != nil {
					offset := s2.AddOffset(-200 * time.Millisecond)
					m.w.SendShowMessage(fmt.Sprint("Subtitle 2 offset ", offset.String()), true)

					UpdateSubtitleOffsetAsync(s2.Name, offset)
				}
			}()
			break
		case gui.KEY_RIGHT_BRACKET:
			log.Print("right bracket pressed")
			go func() {
				// if m.s != nil {
				// 	offset := m.s.AddOffset(1000 * time.Millisecond)
				// 	m.w.SendShowMessage(fmt.Sprint("Subtitle offset ", offset.String()), true)
				// }
				_, s2 := m.getPlayingSubs()
				if s2 != nil {
					offset := s2.AddOffset(200 * time.Millisecond)
					m.w.SendShowMessage(fmt.Sprint("Subtitle 2 offset ", offset.String()), true)

					UpdateSubtitleOffsetAsync(s2.Name, offset)
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

	m.w.FuncOnProgressChanged = append(m.w.FuncOnProgressChanged, func(typ int, percent float64) { //run in main thread, safe to operate ui elements
		switch typ {
		case 0:
			chCursorAutoHide <- struct{}{}
			fallthrough
		case 1:
			t := m.c.CalcTime(percent)
			p := m.c.CalcPlayProgress(t)
			m.w.ShowProgress(p)

			m.SeekAsync(t)
		case 2:
			t := m.c.CalcTime(percent)
			p := m.c.CalcPlayProgress(t)
			m.w.ShowProgress(p)

			log.Print("release dragging:", t.String())

			m.SeekEnd(t)
			chCursorAutoHide <- struct{}{}
		}
	})

	m.w.FuncOnFullscreenChanged = append(m.w.FuncOnFullscreenChanged, func(b bool) {
		m.seekPlayingSubs(m.c.GetTime(), true)
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
			chVolume <- byte(m.a.DecreaseVolume() * 100)
		} else {
			chVolume <- byte(m.a.IncreaseVolume() * 100)
		}
	})

	m.w.FuncSubtitleMenuClicked = append(m.w.FuncSubtitleMenuClicked, func(index int) {
		go func() {
			subs := m.subs
			clicked := subs[index]
			log.Print("toggle subtitle:", clicked.Name, index)

			var s1, s2 *Subtitle
			ps1, ps2 := m.getPlayingSubs()

			if ps1 == nil && ps2 == nil {
				//add playing s1
				s1 = clicked
				go s1.Play()

				m.p.Sub1 = s1.Name
				m.p.Sub2 = ""

			} else if ps1 == clicked {
				//remove playing s1
				ps1.Stop()
				if ps2 != nil {
					s1 = ps2
					s1.IsMainOrSecondSub = true

					m.p.Sub1 = s1.Name
					m.p.Sub2 = ""
				}
			} else if ps2 == clicked {
				//remove playing s2
				ps2.Stop()
				s1 = ps1

				m.p.Sub1 = s1.Name
				m.p.Sub2 = ""
			} else {
				//replace playing subtitle
				if clicked.IsTwoLangs() {
					s1 = clicked
					s2 = nil
				} else if ps1.IsTwoLangs() {
					s1 = clicked
					s2 = nil
				} else if isLangEqual(ps1.Lang1, clicked.Lang1) {
					s1 = clicked
					s2 = ps2
				} else if ps2 == nil {
					s1 = ps1
					s2 = clicked
				} else if isLangEqual(ps2.Lang1, clicked.Lang1) {
					s1 = ps1
					s2 = clicked
				} else { //third language which is impossible for now
					s1 = ps1
					s2 = clicked
				}

				if s1 != ps1 {
					ps1.Stop()

					s1.IsMainOrSecondSub = true
					go s1.Play()

					m.p.Sub1 = s1.Name
				}

				if s2 != nil {
					if s2 != ps2 {
						if ps2 != nil {
							ps2.Stop()
						}

						s2.IsMainOrSecondSub = false
						go s2.Play()

						m.p.Sub2 = s2.Name
					}
				} else {
					if ps2 != nil {
						ps2.Stop()
					}

					m.p.Sub2 = ""
				}
			}

			m.setPlayingSubs(s1, s2)

			t1, t2 := -1, -1
			for i, s := range m.subs {
				if s1 == s {
					t1 = i
				}
				if s2 == s {
					t2 = i
				}
			}
			gui.SetSubtitleMenuItem(t1, t2)

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

func isLangEqual(l1, l2 string) bool {
	if l1 == l2 {
		return true
	} else if l1 == "en" || l2 == "en" {
		return false
	} else {
		return true
	}
}
