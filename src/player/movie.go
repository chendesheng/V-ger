package main

import (
	// "fmt"
	. "player/clock"
	. "player/libav"
	. "player/subtitle"
	// . "player/shared"
	"log"
	"time"
	// "util"
)

type movie struct {
	ctx AVFormatContext
	v   *video
	a   *audio
	s   *Subtitle
	c   *Clock
}

func (m *movie) open(file string, subFile string, start time.Duration) {
	println("open ", file)

	ctx := AVFormatContext{}
	ctx.OpenInput(file)
	if ctx.IsNil() {
		log.Fatal("open failed:", file)
		return
	}

	ctx.FindStreamInfo()
	ctx.DumpFormat()

	m.ctx = ctx
	m.c = NewClock(time.Duration(float64(ctx.Duration()) / AV_TIME_BASE * float64(time.Second)))

	audioStream := ctx.AudioStream()
	if !audioStream.IsNil() {
		m.a = &audio{}
		m.a.setup(ctx, audioStream)
		m.a.c = m.c
	}

	videoStream := ctx.VideoStream()
	if !videoStream.IsNil() {
		m.v = &video{}
		m.v.setup(ctx, videoStream, file, start)
		m.v.c = m.c

		if len(subFile) > 0 {
			println("play subtitle:", subFile)
			m.s = NewSubtitle(subFile, m.v.window, m.c)
			go m.s.Play()
		}
		m.uievents()
		start = m.v.seek(start)

		m.c.Reset()
		m.c.SetTime(start)

		if m.s != nil {
			m.s.Seek(start)
		}
	} else {
		log.Fatal("No video stream find.")
	}
}
func (m *movie) SeekTo(t time.Duration) time.Duration {
	if m.v != nil {
		t = m.v.seek(t)
	}

	println("seek to", t.String())

	if m.a != nil {
		m.a.flushBuffer()
	}

	if m.s != nil {
		m.s.Seek(t)
	}

	return t
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
