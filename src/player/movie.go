package main

import (
	// "fmt"
	. "libav"
	// "player/glfw"
	// "log"
	. "player/clock"
	"time"
	// "util"
)

type movie struct {
	ctx AVFormatContext
	v   *video
	a   *audio
	s   *subtitle
	c   *Clock

	width, height float64
}

func (m *movie) open(file string, subFile string, start time.Duration) {
	println("open ", file)

	ctx := AVFormatContext{}
	ctx.OpenInput(file)
	if ctx.IsNil() {
		println("open failed.")
		return
	}

	ctx.FindStreamInfo()
	ctx.DumpFormat()

	m.ctx = ctx

	videoStream := ctx.VideoStream()

	if !videoStream.IsNil() {
		m.v = &video{}
		m.v.setup(ctx, videoStream, file, start)
	}

	audioStream := ctx.AudioStream()

	if !audioStream.IsNil() {
		m.a = &audio{}
		m.a.setup(ctx, audioStream)
	}
	// avformat_seek_file
	// av_rescale
	if m.v != nil && len(subFile) > 0 {
		println("play subtitle:", subFile)
		m.s = NewSubtitle(subFile, m.v.window)
		// m.v.a = m.a
	}

	SeekFrame(ctx, videoStream, audioStream, m.s, start)

	m.c = NewClock(time.Duration(float64(ctx.Duration()) / AV_TIME_BASE * float64(time.Second)))

	if m.v != nil {
		m.v.c = m.c
	}

	if m.a != nil {
		m.a.c = m.c
	}

	if m.s != nil {
		m.s.c = m.c
	}

	m.c.Reset()
	m.c.SetTime(start)

	if m.v != nil {
		m.v.window.FuncOnProgressChanged = append(m.v.window.FuncOnProgressChanged, func(typ int, percent float64) { //run in main thread, safe to operate ui elements
			switch typ {
			case 0:
				m.c.GotoPercent(percent)
				m.c.Pause()
				break
			case 2:
				println("mouse up:", percent)
				m.c.Resume()
				m.c.StartSeek(percent)

				if m.a != nil {
					codec := m.a.stream.Codec()
					codec.FlushBuffer()
				}
				break
			case 1:
				m.c.GotoPercent(percent)
				t := m.c.GetSeekTime()
				// m.ctx.SeekFile(t, 0)
				m.ctx.SeekFrame(m.v.stream, t, AVSEEK_FLAG_FRAME)

				if m.v != nil {
					codec := m.v.stream.Codec()
					codec.FlushBuffer()
					m.drawCurrentFrame()
				}

				break
			}
		})
	}
}
func SeekFrame(ctx AVFormatContext, videoStream AVStream, audioStream AVStream, s *subtitle, t time.Duration) time.Duration {
	if s != nil {
		s.seek(t)
	}

	ctx.SeekFrame(audioStream, t, AVSEEK_FLAG_FRAME)
	ctx.SeekFrame(videoStream, t, AVSEEK_FLAG_FRAME)

	frame := AllocFrame()
	ret, _ := readOneFrame(ctx, videoStream, frame)
	return ret
	// return dropVideoFrames(ctx, videoStream, t, frame)
	// return dropFrames(ctx, audioStream, t, frame)
}
func readOneFrame(ctx AVFormatContext, stream AVStream, frame AVFrame) (time.Duration, bool) {
	packet := AVPacket{}
	codecCtx := stream.Codec()

	for ctx.ReadFrame(&packet) >= 0 {
		if packet.StreamIndex() == stream.Index() {
			if codecCtx.DecodeVideo(frame, &packet) {
				tmp := packet.Pts()
				if tmp == AV_NOPTS_VALUE {
					tmp = 0
				}

				pts := time.Duration(float64(tmp) * stream.Timebase().Q2D() * float64(time.Second))
				println("pts:", pts.String())
				packet.Free()

				return pts, true
			} else {
				packet.Free()
			}
		}
	}

	return 0, false
}
func dropVideoFrames(ctx AVFormatContext, videoStream AVStream, t time.Duration, frame AVFrame) time.Duration {
	packet := AVPacket{}
	codecCtx := videoStream.Codec()

	for ctx.ReadFrame(&packet) >= 0 {
		if packet.StreamIndex() == videoStream.Index() {
			if codecCtx.DecodeVideo(frame, &packet) {

				tmp := packet.Pts()
				if tmp == AV_NOPTS_VALUE {
					tmp = 0
				}

				pts := time.Duration(float64(tmp) * videoStream.Timebase().Q2D() * float64(time.Second))
				println("pts:", pts.String())
				packet.Free()

				if t-pts < 10*time.Millisecond {
					return pts
				}
			} else {
				packet.Free()
			}
		}
	}

	return 0
}

func tabs(t time.Duration) time.Duration {
	if t < 0 {
		t = -t
	}
	return t
}
func (m *movie) drawCurrentFrame() {
	ctx := m.ctx
	v := m.v
	if v == nil {
		return
	}

	packet := AVPacket{}

	frame := v.frame
	for ctx.ReadFrame(&packet) >= 0 {
		if v.stream.Index() == packet.StreamIndex() {
			codecCtx := v.codecCtx

			frameFinished := codecCtx.DecodeVideo(frame, &packet)
			packet.Free()

			if frameFinished {
				frame.Flip(v.height)
				swsCtx := SwsGetCachedContext(v.width, v.height, codecCtx.PixelFormat(),
					v.width, v.height, AV_PIX_FMT_RGB24, SWS_BICUBIC)

				swsCtx.Scale(frame, v.pictureRGB)
				obj := v.pictureRGB.Layout(AV_PIX_FMT_RGB24, v.width, v.height)
				v.setPic(picture{obj, 0})
				v.window.RefreshContent()
				break
			}
		} else {
			packet.Free()
		}
	}
}

func (m *movie) decode() {
	packet := AVPacket{}
	ctx := m.ctx

	for ctx.ReadFrame(&packet) >= 0 {
		streamIndex := packet.StreamIndex()
		if m.v != nil {
			if m.v.stream.Index() == streamIndex {
				m.v.decode(&packet)
				packet.Free()
			}
		}

		if m.a != nil {
			if m.a.stream.Index() == streamIndex {
				pkt := packet
				pkt.Dup()
				m.a.ch <- &pkt
			}
		}

		if m.c.IsSeeking() {
			now := m.c.GetTime()

			timeAfterSeek := SeekFrame(ctx, m.v.stream, m.a.stream, m.s, now)
			m.c.EndSeek(timeAfterSeek)

			m.a.flushBuffer()
			println("after seek back")
			if m.s != nil {
				go m.s.play()
			}
		}
	}

	m.stop()
}

func (m *movie) play() {
	if m.s != nil {
		go m.s.play()
	}
	if m.v != nil {
		m.v.play()
	} else {
		for {
			<-time.After(time.Millisecond)
		}
	}
}

func (m *movie) stop() {
	if m.v != nil {
		close(m.v.ch)
	}
	if m.a != nil {
		close(m.a.ch)
	}
}
