package movie

import (
	"log"
	"os"
	"path"
	. "player/libav"
	. "player/movie/video"
	"player/shared"
	"task"
	"time"
	"util"
)

func (m *Movie) sendPacket(index int, ch chan *AVPacket, packet AVPacket) bool {
	if index == packet.StreamIndex() {
		pkt := packet
		pkt.Dup()

		select {
		case ch <- &pkt:
			return true
		case <-m.quit:
			return false
		}
	}
	return false
}

func (m *Movie) hold() {
	if m.chHold == nil {
		return
	}

	log.Print("send pause movie")

	select {
	case m.chHold <- 0:
		m.v.FlushBuffer()
		m.a.FlushBuffer()
	case <-m.quit:
	}

	return
}

func (m *Movie) unHold(t time.Duration) {
	select {
	case m.chHold <- t:
	case <-m.quit:
	}
}

func (m *Movie) decodeVideo(packet *AVPacket) {
	if frameFinished, pts, img := m.v.DecodeAndScale(packet); frameFinished {
		//make sure seek operations not happens before one frame finish decode
		//if not, segment fault & crash
		select {
		case m.v.ChanDecoded <- &VideoFrame{pts, img}:
			break
		case <-m.chHold:
			log.Print("pause movie")
			select {
			case t := <-m.chHold:
				log.Print("resume movie:", t.String())
				m.c.SetTime(t)
			case <-m.quit:
				packet.Free()
				return
			}
			break
		case <-m.quit:
			packet.Free()
			return
		}

		m.seekPlayingSubs(m.c.GetTime(), false)
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
			m.w.SendShowProgress(&shared.PlayProgressInfo{})

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
	if m.p.LastPos > time.Second && m.p.LastPos < m.p.Duration-50*time.Millisecond {
		var img []byte
		start, img, _ = m.v.Seek(m.p.LastPos)

		m.showProgressInner(start)
		m.w.SendDrawImage(img)

		if m.httpBuffer != nil && m.httpBuffer.WaitQuit(3*1024*1024, m.quit) {
			return
		}
	}

	m.w.SendSetSize(m.v.Width, m.v.Height)

	packet := AVPacket{}
	ctx := m.ctx

	m.w.SendHideSpinning()
	m.c.SetTime(start)
	for {
		select {
		case m.chProgress <- m.c.GetTime():
		case <-m.quit:
			return
		case <-time.After(50 * time.Millisecond):
			log.Print("write m.chProgress timeout")
		}

		resCode := ctx.ReadFrame(&packet)
		if resCode >= 0 {
			if m.v.StreamIndex == packet.StreamIndex() {
				m.decodeVideo(&packet)
				packet.Free()
				continue
			}

			if m.a != nil {
				if m.sendPacket(m.a.StreamIndex(), m.a.PacketChan, packet) {
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
					m.a.FlushBuffer()

					t, _, err := m.v.Seek(m.c.GetTime())
					m.w.SendHideSpinning()
					if err == nil {
						log.Print("seek success:", t.String())

						if m.httpBuffer != nil {
							m.w.SendShowSpinning()
							if m.httpBuffer.WaitQuit(2*1024*1024, m.quit) {
								m.w.SendHideSpinning()
								return
							}
							m.w.SendHideSpinning()
						}

						m.c.SetTime(t)

						continue
					} else {
						log.Print("seek error:", err)
					}
				}

				// log.Print("seek to unfinished:", m.c.GetTime().String())
				// log.Print("get frame error:", resCode)
				// }
			}
			select {
			case <-m.chHold:
				log.Print("pause movie")
				select {
				case t := <-m.chHold:
					m.c.SetTime(t)
				case <-m.quit:
					return
				}
				break
			case <-time.After(100 * time.Millisecond):
				break
			case <-m.quit:
				return
			}

		}
	}
}
