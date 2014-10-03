package movie

import (
	"fmt"
	"log"
	"time"
	"vger/player/gui"
	. "vger/player/shared"
	. "vger/player/subtitle"
)

func (m *Movie) uievents() {
	log.Print("movie uievents")

	m.w.InitEvents()

	m.w.FuncAudioMenuClicked = append(m.w.FuncAudioMenuClicked, func(i int) {
		go func() {
			log.Printf("Audio menu click:%d", i)

			m.p.SoundStream = i
			SavePlayingAsync(m.p)

			m.a.Close()
			err := m.a.Open(getStream(m.audioStreams, i))
			if err != nil {
				log.Print(err)
			}
		}()
	})

	// var chPausing chan seekArg
	// chPausing = nil
	// var pausingTime time.Duration
	m.w.FuncKeyDown = append(m.w.FuncKeyDown, func(keycode int) bool {
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
				subTime := s1.GetSubTime(t, -1)

				if subTime == 0 {
					offset = -10 * time.Second
				} else {
					offset = subTime - t
				}
			} else {
				offset = -10 * time.Second
			}
			m.SeekOffset(offset)
			break
		case gui.KEY_RIGHT:
			var offset time.Duration
			s1, _ := m.getPlayingSubs()
			if s1 != nil {
				t := m.c.GetTime()
				subTime := s1.GetSubTime(t, 1)
				log.Print("subtime:", subTime)

				if subTime == 0 {
					offset = 10 * time.Second
				} else {
					offset = subTime - t
				}
			} else {
				offset = 10 * time.Second
			}
			m.SeekOffset(offset)
			break
		case gui.KEY_UP:
			m.SeekOffset(5 * time.Second)
			break
		case gui.KEY_DOWN:
			m.SeekOffset(-5 * time.Second)
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
		case gui.KEY_COMMA:
			m.a.AddOffset(-100 * time.Millisecond)
			break
		case gui.KEY_PERIOD:
			m.a.AddOffset(100 * time.Millisecond)
			break
		default:
			return false
		}

		return true
	})

	chVolume := make(chan struct{})

	m.w.FuncOnFullscreenChanged = append(m.w.FuncOnFullscreenChanged, func(b bool) {
		if m.c != nil {
			m.seekPlayingSubs(m.c.GetTime(), true)
		}
	})

	m.w.FuncMouseWheelled = append(m.w.FuncMouseWheelled, func(deltaY float64) {
		if deltaY == 0 {
			return
		}

		if m.a == nil {
			return
		}

		var volume byte
		if deltaY > 0 {
			volume = byte(m.a.DecreaseVolume() * 100)
		} else {
			volume = byte(m.a.IncreaseVolume() * 100)
		}

		m.p.Volume = volume
		SavePlayingAsync(m.p)
		// m.w.ShowMessage(fmt.Sprintf("Volume: %d%%", volume), true)
		m.w.SetVolume(volume)
		m.w.SetVolumeDisplay(true)

		select {
		case chVolume <- struct{}{}:
		case <-m.quit:
			return
		case <-time.After(100 * time.Millisecond):
		}

	})

	m.w.FuncSubtitleMenuClicked = append(m.w.FuncSubtitleMenuClicked, func(index int) {
		go func() {
			log.Print("toggle subtitle:", index)

			subs := m.subs
			clicked := subs[index]

			var s1, s2 *Subtitle
			ps1, ps2 := m.getPlayingSubs()

			if ps1 == nil && ps2 == nil {
				//add playing s1
				s1 = clicked
				// go s1.Play()

				m.p.Sub1 = s1.Name
				m.p.Sub2 = ""

			} else if ps1 == clicked {
				//remove playing s1
				ps1.Stop()
				if ps2 != nil {
					s1 = ps2
					s1.IsMainSub = true

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

					s1.IsMainSub = true
					// go s1.Play()

					m.p.Sub1 = s1.Name
				}

				if s2 != nil {
					if s2 != ps2 {
						if ps2 != nil {
							ps2.Stop()
						}

						s2.IsMainSub = false
						// go s2.Play()

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
			m.w.SetSubtitleMenuItem(t1, t2)

			SavePlayingAsync(m.p)
		}()
	})

	go func() {
		for {
			select {
			case <-time.After(time.Second):
				m.w.SendSetVolumeDisplay(false)
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
			fallthrough
		case 1:
			t := m.c.CalcTime(percent)
			t = t / time.Second * time.Second
			p := m.c.CalcPlayProgress(t)
			m.w.ShowProgress(p.Left, p.Right, p.Percent)

			m.seeking.SendSeek(t)
		case 2:
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
