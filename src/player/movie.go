package main

import (
	// "fmt"
	. "libav"
	"player/glfw"
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
		m.v.setup(ctx, videoStream, file)
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

	Seek(ctx, audioStream, videoStream, m.s, start)

	// println(time.Since(b).String())
	// ctx.SeekFrame(audioStream, start, 0)
	// ctx.SeekFrame(videoStream, start, 0)

	// seekVideo(ctx, videoStream, audioStream, start)

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
		m.v.window.AddEventHandler(func(e Event) { //run in main thread, safe to operate ui elements
			switch e.Kind {
			case KeyPress:
				switch e.Data.(glfw.Key) {
				case glfw.KeyLeft:
					println("key left pressed")
					m.c.AddTime(-10 * time.Second)
					break
				case glfw.KeyRight:
					println("key right pressed")
					m.c.AddTime(10 * time.Second)
					break
				case glfw.KeyUp:
					m.c.AddTime(time.Minute)
					break
				case glfw.KeyDown:
					m.c.AddTime(-time.Minute)
					break
				}
				break
			case TrackPositionChanged:
				m.c.GotoPercent(e.Data.(float64))
				break
			}
		})
	}
}
func Seek(ctx AVFormatContext, audioStream, videoStream AVStream, s *subtitle, start time.Duration) {
	ctx.SeekFile(audioStream, start, AVSEEK_FLAG_FRAME)
	ctx.SeekFile(videoStream, start, AVSEEK_FLAG_FRAME)

	packet := AVPacket{}
	frame := AllocFrame()
	codecCtx := videoStream.Codec()
	for ctx.ReadFrame(&packet) >= 0 {
		if packet.StreamIndex() == videoStream.Index() {
			pts := time.Duration(float64(packet.Pts()) * videoStream.Timebase().Q2D() * (float64(time.Second)))

			if codecCtx.DecodeVideo(frame, &packet) {
				packet.Free()
				if start-pts < 10*time.Millisecond {
					break
				}
			} else {
				packet.Free()
			}
		}
	}

	if s != nil {
		s.seek(start)
	}
}
func (m *movie) decode() {
	packet := AVPacket{}
	ctx := m.ctx

	for ctx.ReadFrame(&packet) >= 0 {

		if m.v != nil {
			if m.v.stream.Index() == packet.StreamIndex() {
				m.v.decode(&packet)
			}
		}

		if m.a != nil {
			if m.a.stream.Index() == packet.StreamIndex() {
				pt := &packet
				pt.Dup()
				p := *pt
				m.a.ch <- &p
			}
		}

		now := m.c.GetTime()
		vc := time.Duration(m.v.videoClock * (float64(time.Second)))
		// if now > vc+time.Second {
		// 	Seek(ctx, m.a.stream, m.v.stream, now)
		// 	m.v.frame = AllocFrame()
		// 	m.a.flushBuffer()
		// 	println("after seek")
		// } else
		if now < vc-time.Second || now > vc+time.Second {
			// ctx.SeekFile(m.v.stream, now, AVSEEK_FLAG_FRAME|AVSEEK_FLAG_BACKWARD)
			// ctx.SeekFile(m.a.stream, now, AVSEEK_FLAG_FRAME|AVSEEK_FLAG_BACKWARD)

			Seek(ctx, m.a.stream, m.v.stream, m.s, now)
			m.v.frame = AllocFrame()
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
