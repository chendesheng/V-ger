package movie

import (
	"log"
	. "player/libav"
	. "player/video"
	"time"
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
func (m *Movie) decodeVideo(packet *AVPacket) {
	if frameFinished, pts, img := m.v.DecodeAndScale(packet); frameFinished {
		//make sure seek operations not happens before one frame finish decode
		//if not, segment fault & crash
		select {
		case m.v.ChanDecoded <- &VideoFrame{pts, img}:
			break
		case t := <-m.chSeekPause:
			if t != -1 {
				break
			}
			m.toggleShowProgress()
			for {
				t := <-m.chSeekPause
				if t >= 0 {
					m.toggleShowProgress()

					m.c.SetTime(t)
					break
				}
			}
			break
		case <-m.quit:
			packet.Free()
			return
		}

		t := m.c.GetTime()
		if m.s != nil {
			m.s.Seek(t)
		}
		if m.s2 != nil {
			m.s2.Seek(t)
		}
	}
}

func (m *Movie) decode(name string) {
	defer func() {
		if m.a != nil {
			m.a.Close()
		}
		if m.v != nil {
			m.v.Close()
		}
		m.ctx.CloseInput()

		if m.finishClose != nil {
			close(m.finishClose)
		}
	}()

	packet := AVPacket{}
	ctx := m.ctx

	bufferring := false
	for {
		resCode := ctx.ReadFrame(&packet)
		if resCode >= 0 {
			if bufferring {
				bufferring = false
				m.c.Resume()
			}
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
			if resCode == AVERROR_EOF && (m.c.GetTime()-m.c.TotalTime()<time.Second) {
				m.c.SetTime(m.c.TotalTime())
			} else {
				bufferring = true
				m.c.Pause()

				m.a.FlushBuffer()
				m.v.FlushBuffer()

				t, _, err := m.v.Seek(m.c.GetTime())
				if err == nil {
					println("seek success:", t.String())
					if m.httpBuffer != nil {
						m.httpBuffer.Wait(2 * 1024 * 1024)
					}
					m.c.SetTime(t)
					continue
				} else {
					log.Print("seek error:", err)
				}

				// println("seek to unfinished:", m.c.GetTime().String())
				log.Print("get frame error:", resCode)
			}
			select {
			case t := <-m.chSeekPause:
				println("seek to unfinished2")
				if t != -1 {
					continue
				}
				for {
					println("seek to unfinished3")
					t := <-m.chSeekPause
					println("seek to unfinished4")
					if t >= 0 {
						m.c.SetTime(t)
						break
					}
				}
			case <-time.After(100 * time.Millisecond):
				break
			case <-m.quit:
				return
			}

		}
		// println(bufferring)
	}
}
