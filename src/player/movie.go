package main

import (
	// "fmt"
	. "player/clock"
	. "player/libav"
	// . "player/shared"
	// "log"
	"player/gui"
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

	chSeek chan time.Duration
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

	m.chSeek = make(chan time.Duration)

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

	start = SeekFrame(ctx, videoStream, audioStream, m.s, start, false)

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
		m.v.window.FuncKeyDown = append(m.v.window.FuncKeyDown, func(keycode int) {
			switch keycode {
			case gui.KEY_SPACE:
				m.c.Toggle()
				break
			case gui.KEY_LEFT:
				println("key left pressed")

				m.c.Pause()
				m.c.ResumeWithTime(m.SeekTo(m.c.GetSeekTime() - 10*time.Second))
				break
			case gui.KEY_RIGHT:
				println("key right pressed")

				m.c.Pause()
				m.c.ResumeWithTime(m.SeekTo(m.c.GetSeekTime() + 10*time.Second))
				break
			case gui.KEY_UP:
				m.c.Pause()
				m.c.ResumeWithTime(m.SeekTo(m.c.GetSeekTime() + time.Minute))
				break
			case gui.KEY_DOWN:
				m.c.Pause()
				m.c.ResumeWithTime(m.SeekTo(m.c.GetSeekTime() - time.Minute))
				break
			}
		})
		m.v.window.FuncOnProgressChanged = append(m.v.window.FuncOnProgressChanged, func(typ int, percent float64) { //run in main thread, safe to operate ui elements
			lastSeekTime := time.Duration(0)

			switch typ {
			case 0:
				lastSeekTime = m.c.GetSeekTime()

				m.c.Pause()
				break
			case 2:
				println("mouse up:", percent)
				t := m.c.CalcTime(percent)
				m.c.ResumeWithTime(m.SeekTo(t))
				break
			case 1:
				t := m.c.CalcTime(percent)
				// m.ctx.SeekFile(t, 0)
				flags := AVSEEK_FLAG_FRAME
				if t < lastSeekTime {
					flags |= AVSEEK_FLAG_BACKWARD
				}
				m.ctx.SeekFrame(m.v.stream, t, flags)

				lastSeekTime = t

				if m.v != nil {
					codec := m.v.stream.Codec()
					codec.FlushBuffer()
					m.drawCurrentFrame()
				}

				if m.s != nil {
					m.v.window.ShowText(m.s.seek(t))
				}

				m.v.window.ShowProgress(m.c.CalcPlayProgress(percent))

				break
			}
		})
	}
}
func (m *movie) SeekTo(t time.Duration) time.Duration {
	backward := false
	if t < m.c.GetSeekTime() {
		backward = true
	}

	timeAfterSeek := SeekFrame(m.ctx, m.v.stream, m.a.stream, m.s, t, backward)

	println("seek to", timeAfterSeek.String())

	if m.a != nil {
		m.a.flushBuffer()
	}

	if m.s != nil {
		go m.s.play()
	}

	return timeAfterSeek
}

func SeekFrame(ctx AVFormatContext, videoStream AVStream, audioStream AVStream, s *subtitle, t time.Duration, backward bool) time.Duration {
	//seek audio is very very slow, it takes 30 seconds to seek to about 28m in movie (720p)
	// b := time.Now()
	// ctx.SeekFrame(audioStream, t, AVSEEK_FLAG_FRAME)
	// println(time.Since(b).String())
	flags := AVSEEK_FLAG_FRAME
	if backward {
		flags |= AVSEEK_FLAG_BACKWARD
	}
	ctx.SeekFrame(videoStream, t, flags)

	frame := AllocFrame()
	ret, _ := readOneFrame(ctx, videoStream, frame)

	if s != nil {
		s.seek(ret)
	}

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

//only call from UI thread
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

				v.swsCtx.Scale(frame, v.pictureRGB)

				v.window.RefreshContent(v.pictureRGB.RGBBytes(v.width, v.height))
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
		m.c.WaitUtilRunning()

		streamIndex := packet.StreamIndex()
		if m.v != nil {
			if m.v.stream.Index() == streamIndex {
				// println("decode video")
				m.v.decode(&packet)
				packet.Free()
			}
		}

		if m.a != nil {
			if m.a.stream.Index() == streamIndex {
				// println("decode audio")
				pkt := packet
				pkt.Dup()
				m.a.ch <- &pkt
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
	if m.a != nil {
		close(m.a.ch)
	}
}
