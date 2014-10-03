package movie

import (
	"log"
	"os"
	"path"
	"time"
	. "vger/player/libav"
	"vger/player/shared"
	"vger/task"
	"vger/util"
)

func (m *Movie) Hold() time.Duration {
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

func (m *Movie) Unhold(t time.Duration) {
	m.v.Unhold()

	select {
	case m.chHold <- t:
	case <-m.quit:
	}
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

	if ok, file := getNextEpisode(m.filename); ok {
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

func (m *Movie) sendPacket(index int, ch chan *AVPacket, packet *AVPacket) bool {
	if index == packet.StreamIndex() {
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
	}
	return false
}

func (m *Movie) decode(name string) {
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
		m.w.SendDrawImage(img)

		if m.waitBuffer(3 * 1024 * 1024) {
			return
		}
	}

	m.w.SendSetSize(m.v.Width, m.v.Height)
	m.w.SendSetCursor(true)

	ctx := m.ctx

	m.c.SetTime(start)

	go m.v.Play()

	for {
		select {
		case m.chProgress <- m.c.GetTime():
		case <-m.quit:
			return
		case <-time.After(50 * time.Millisecond):
			log.Print("write m.chProgress timeout")
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

		packet := AVPacket{}
		resCode := ctx.ReadFrame(&packet)
		if resCode >= 0 {
			if m.sendPacket(m.v.StreamIndex, m.v.ChPackets, &packet) {
				m.seekPlayingSubs(m.c.GetTime(), false)
				continue
			}

			if m.a != nil {
				if m.sendPacket(m.a.StreamIndex(), m.a.ChPackets, &packet) {
					continue
				}
			}

			packet.Free()
		} else {
			log.Printf("read frame error: %x", resCode)
			if resCode == AVERROR_EOF && (m.c.TotalTime()-m.c.GetTime() < 2*time.Second) {
				m.c.SetTime(m.c.TotalTime())
				if m.playNextEpisode() {
					return
				}
			} else {
				if m.httpBuffer == nil && resCode == AVERROR_INVALIDDATA {
					t := m.c.GetTime()
					t = t / (500 * time.Millisecond) * 500 * time.Millisecond
					m.c.SetTime(t)
				} else {
					m.v.FlushBuffer()
					if m.a != nil {
						m.a.FlushBuffer()
					}

					t, _, err := m.v.Seek(m.c.GetTime())
					if err == nil {
						log.Print("seek success:", t.String())

						if m.waitBuffer(2 * 1024 * 1024) {
							return
						}

						m.c.SetTime(t)

						continue
					} else {
						log.Print("seek error:", err)
					}
				}
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
