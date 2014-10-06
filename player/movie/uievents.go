package movie

import (
	"log"
	"time"
	"vger/player/gui"
	. "vger/player/shared"
	// . "vger/player/subtitle"
)

var chVolume = make(chan struct{})

func (m *Movie) uievents() {
	log.Print("movie uievents")

	m.w.FuncAudioMenuClick = append(m.w.FuncAudioMenuClick, func(i int) {
		go func() {
			log.Printf("Audio menu click:%d", i)

			if m.audioStreams[i].Index() == m.a.StreamIndex() {
				return
			}

			m.a.Close()
			err := m.a.Open(m.audioStreams[i])
			if err != nil {
				log.Print(err)
			} else {
				m.p.SoundStream = m.a.StreamIndex()
				SavePlayingAsync(m.p)
			}
		}()
	})

	m.w.FuncKeyDown = append(m.w.FuncKeyDown, func(keycode int) bool {
		SavePlayingAsync(m.p)

		switch keycode {
		case gui.KEY_R:
			if !m.w.IsFullScreen() {
				m.w.ToggleForceScreenRatio()
			}
			break
		// case gui.KEY_MINUS:
		// 	log.Print("key minus pressed")
		// 	go func() {
		// 		s1, _ := m.getPlayingSubs()
		// 		if s1 != nil {
		// 			offset := s1.AddOffset(-200 * time.Millisecond)
		// 			m.w.SendShowMessage(fmt.Sprint("Subtitle offset ", offset.String()), true)

		// 			UpdateSubtitleOffsetAsync(s1.Name, offset)
		// 		}
		// 	}()
		// 	break
		// case gui.KEY_EQUAL:
		// 	log.Print("key equal pressed")
		// 	go func() {
		// 		s1, _ := m.getPlayingSubs()
		// 		if s1 != nil {
		// 			offset := s1.AddOffset(200 * time.Millisecond)
		// 			m.w.SendShowMessage(fmt.Sprint("Subtitle offset ", offset.String()), true)

		// 			UpdateSubtitleOffsetAsync(s1.Name, offset)
		// 		}
		// 	}()
		// 	break
		// case gui.KEY_LEFT_BRACKET:
		// 	log.Print("left bracket pressed")
		// 	go func() {
		// 		// if m.s != nil {
		// 		// 	offset := m.s.AddOffset(-1000 * time.Millisecond)
		// 		// 	m.w.SendShowMessage(fmt.Sprint("Subtitle offset ", offset.String()), true)
		// 		// }
		// 		_, s2 := m.getPlayingSubs()
		// 		if s2 != nil {
		// 			offset := s2.AddOffset(-200 * time.Millisecond)
		// 			m.w.SendShowMessage(fmt.Sprint("Subtitle 2 offset ", offset.String()), true)

		// 			UpdateSubtitleOffsetAsync(s2.Name, offset)
		// 		}
		// 	}()
		// 	break
		// case gui.KEY_RIGHT_BRACKET:
		// 	log.Print("right bracket pressed")
		// 	go func() {
		// 		// if m.s != nil {
		// 		// 	offset := m.s.AddOffset(1000 * time.Millisecond)
		// 		// 	m.w.SendShowMessage(fmt.Sprint("Subtitle offset ", offset.String()), true)
		// 		// }
		// 		_, s2 := m.getPlayingSubs()
		// 		if s2 != nil {
		// 			offset := s2.AddOffset(200 * time.Millisecond)
		// 			m.w.SendShowMessage(fmt.Sprint("Subtitle 2 offset ", offset.String()), true)

		// 			UpdateSubtitleOffsetAsync(s2.Name, offset)
		// 		}
		// 	}()
		// 	break
		// case gui.KEY_COMMA:
		// 	m.a.AddOffset(-100 * time.Millisecond)
		// 	break
		// case gui.KEY_PERIOD:
		// 	m.a.AddOffset(100 * time.Millisecond)
		// 	break
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

	// m.w.FuncSubtitleMenuClick = func(index int) {
	// 	log.Print("toggle subtitle:", index)

	// 	subs := m.subs
	// 	clicked := subs[index]

	// 	var s1, s2 *Subtitle
	// 	ps1, ps2 := m.getPlayingSubs()

	// 	if ps1 == nil && ps2 == nil {
	// 		//add playing s1
	// 		s1 = clicked
	// 		// go s1.Play()

	// 		m.p.Sub1 = s1.Name
	// 		m.p.Sub2 = ""

	// 	} else if ps1 == clicked {
	// 		//remove playing s1
	// 		ps1.Stop()
	// 		if ps2 != nil {
	// 			s1 = ps2
	// 			s1.IsMainSub = true

	// 			m.p.Sub1 = s1.Name
	// 			m.p.Sub2 = ""
	// 		}
	// 	} else if ps2 == clicked {
	// 		//remove playing s2
	// 		ps2.Stop()
	// 		s1 = ps1

	// 		m.p.Sub1 = s1.Name
	// 		m.p.Sub2 = ""
	// 	} else {
	// 		//replace playing subtitle
	// 		if clicked.IsTwoLangs() {
	// 			s1 = clicked
	// 			s2 = nil
	// 		} else if ps1.IsTwoLangs() {
	// 			s1 = clicked
	// 			s2 = nil
	// 		} else if isLangEqual(ps1.Lang1, clicked.Lang1) {
	// 			s1 = clicked
	// 			s2 = ps2
	// 		} else if ps2 == nil {
	// 			s1 = ps1
	// 			s2 = clicked
	// 		} else if isLangEqual(ps2.Lang1, clicked.Lang1) {
	// 			s1 = ps1
	// 			s2 = clicked
	// 		} else { //third language which is impossible for now
	// 			s1 = ps1
	// 			s2 = clicked
	// 		}

	// 		if s1 != ps1 {
	// 			ps1.Stop()

	// 			s1.IsMainSub = true
	// 			// go s1.Play()

	// 			m.p.Sub1 = s1.Name
	// 		}

	// 		if s2 != nil {
	// 			if s2 != ps2 {
	// 				if ps2 != nil {
	// 					ps2.Stop()
	// 				}

	// 				s2.IsMainSub = false
	// 				// go s2.Play()

	// 				m.p.Sub2 = s2.Name
	// 			}
	// 		} else {
	// 			if ps2 != nil {
	// 				ps2.Stop()
	// 			}

	// 			m.p.Sub2 = ""
	// 		}
	// 	}

	// 	m.setPlayingSubs(s1, s2)
	// 	SavePlayingAsync(m.p)
	// }

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
			m.w.UpdatePlaybackInfo(p.Left, p.Right, p.Percent)

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
