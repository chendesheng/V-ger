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

func (m *Movie) Pause(flushBuffer bool) (ch chan time.Duration) {
	if m.chPause == nil {
		return
	}

	println("send pause movie")

	ch = make(chan time.Duration)
	select {
	case m.chPause <- ch:
		if flushBuffer {
			m.a.FlushBuffer()
			m.v.FlushBuffer()
		}
	case <-m.quit:
	}

	return
}

func (m *Movie) decodeVideo(packet *AVPacket) {
	if frameFinished, pts, img := m.v.DecodeAndScale(packet); frameFinished {
		//make sure seek operations not happens before one frame finish decode
		//if not, segment fault & crash
		select {
		case m.v.ChanDecoded <- &VideoFrame{pts, img}:
			break
		case ch := <-m.chPause:
			println("pause movie")
			select {
			case t := <-ch:
				println("resume movie", t.String())
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
	m.chPause = make(chan chan time.Duration)

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

	go m.seekRoutine()

	packet := AVPacket{}
	ctx := m.ctx

	var lastTime time.Duration
	for {
		// println("send chProgress:", m.c.GetTime().String())
		select {
		case m.chProgress <- m.c.GetTime():
		case <-m.quit:
		case <-time.After(50 * time.Millisecond):
			println("write m.chProgress timeout")
		}
		// println(buffering)

		resCode := ctx.ReadFrame(&packet)
		if resCode >= 0 {
			if lastTime > 0 {
				m.c.SetTime(lastTime)
				lastTime = 0
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
			if resCode == AVERROR_EOF && (m.c.TotalTime()-m.c.GetTime() < time.Second) {
				m.c.SetTime(m.c.TotalTime())
			} else {
				lastTime = m.c.GetTime()

				m.a.FlushBuffer()
				m.v.FlushBuffer()

				t, _, err := m.v.Seek(lastTime)
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
			case ch := <-m.chPause:
				println("pause movie")
				select {
				case t := <-ch:
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
