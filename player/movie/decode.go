package movie

import (
	"log"
	"os"
	"path"
	"time"
	"vger/player/libav"
	"vger/player/shared"
	"vger/task"
	"vger/util"
)

func (m *Movie) ClockHold() {
	if m.c != nil {
		m.c.Hold()
	}
}

func (m *Movie) ClockUnhold() {
	if m.c != nil {
		m.c.Unhold()
	}
}

func (m *Movie) SeekHold() time.Duration {
	if m.chHold == nil {
		if m.c != nil {
			return m.c.GetTime()
		} else {
			return -1
		}
	}

	m.v.Hold()

	log.Print("send pause movie")

	select {
	case m.chHold <- 0:
		if m.a != nil {
			m.a.FlushBuffer()
		}
	case <-m.quit:
	}

	if m.c != nil {
		return m.c.GetTime()
	} else {
		return -1
	}
}

func (m *Movie) SeekUnhold(t time.Duration) {
	select {
	case m.chHold <- t:
	case <-m.quit:
	}

	m.v.Unhold()
}

func getNextEpisode(filename string) (bool, string) {
	t, err := task.GetTask(filename)
	if err != nil {
		log.Print(err)
		return false, ""
	}

	if len(t.Subscribe) == 0 {
		return false, ""
	}

	season := t.Season
	episode := t.Episode + 1
	url := ""
	name := ""
	status := ""
	if name, status, url, err = task.GetEpisodeTask(t.Subscribe, season, episode); err != nil {
		season = t.Season + 1
		episode = 1

		if name, status, url, err = task.GetEpisodeTask(t.Subscribe, season, episode); err != nil {
			return false, ""
		}
	}

	if status == "Finished" {
		file := path.Join(util.ReadConfig("dir"), t.Subscribe, name)
		_, err := os.Stat(file)
		if err != nil {
			log.Print(err)
			return false, ""
		} else {
			return true, file
		}
	} else if len(url) > 0 {
		return true, url
	}

	return false, ""
}

func (m *Movie) playNextEpisode() bool {

	if ok, file := getNextEpisode(m.Filename); ok {
		log.Print("playNextEpisode:", file)

		go func() {
			m.w.SendShowProgress("", "", 0)

			m.SavePlaying()
			m.Close()
			m.Reset()
			err := m.Open(m.w, file)
			if err == nil {
				m.PlayAsync()
			} else {
				log.Print(err)
			}
		}()

		return true
	}

	return false
}

func (m *Movie) sendPacket(ch chan libav.AVPacket, packet libav.AVPacket) bool {
	select {
	case ch <- packet:
		return true
	case <-m.chHold:
		log.Print("hold movie2")
		select {
		case t := <-m.chHold:
			log.Print("unhold movie2:", t.String())
			m.c.SetTime(t)
		case <-m.quit:
			return false
		}
	case <-m.quit:
		return false
	}
	return false
}

func (m *Movie) decode() {
	m.chHold = make(chan time.Duration)

	defer func() {
		if m.a != nil {
			m.a.Close()
		}
		if m.v != nil {
			m.v.Close()
		}
		m.ctx.Close()

		if m.finishClose != nil {
			close(m.finishClose)
		}
	}()

	m.w.SendSetSize(m.v.Width, m.v.Height)

	var start time.Duration
	lastPos := m.p.GetLastPos()
	if lastPos > time.Second && lastPos < m.p.Duration-50*time.Millisecond {
		var img []byte
		var err error

		start, img, err = m.v.Seek(lastPos)

		if err != nil {
			log.Print(err)

			start, img, err = m.v.Seek(0)

			m.p.SetLastPos(0)
			shared.SavePlayingAsync(m.p)

			if err != nil {
				log.Print(err)
				return
			}
		}

		m.showProgressInner(start)
		m.w.Draw(img)

		if m.waitBuffer(3 * 1024 * 1024) {
			return
		}
	}

	m.w.SendSetControlsVisible(true, true)
	m.w.SendSetTitle(m.Filename)

	ctx := m.ctx

	m.c.SetTime(start)
	go m.v.Play()

	for {
		select {
		case m.chProgress <- m.c.GetTime():
		case <-m.quit:
			return
		case <-time.After(50 * time.Millisecond):
			// log.Print("write m.chProgress timeout")
		case <-m.chHold:
			log.Print("hold movie")
			select {
			case t := <-m.chHold:
				log.Print("unhold movie:", t.String())
				m.c.SetTime(t)
			case <-m.quit:
				return
			}
		}

		packet := libav.NewAVPacket()
		resCode := ctx.ReadFrame(packet)
		if resCode >= 0 {
			if m.v.StreamIndex == packet.StreamIndex() {
				if m.sendPacket(m.v.ChPackets, packet) {
					continue
				}
			}

			if m.a != nil {
				if m.a.StreamIndex() == packet.StreamIndex() {
					if m.sendPacket(m.a.ChPackets, packet) {
						continue
					}
				}
			}
			packet.FreePacket()
			packet.Free()
		} else {
			select {
			case m.v.ChPackets <- libav.AVPacket{}:
			case <-m.quit:
				return
			}

			select {
			case <-time.After(100 * time.Millisecond):
				break
			case <-m.quit:
				return
			}
		}
	}
}
