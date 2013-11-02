package main

import (
	// "fmt"
	. "libav"
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

	// b := time.Now()

	ctx.SeekFile(audioStream, start, AVSEEK_FLAG_FRAME)
	ctx.SeekFile(videoStream, start, AVSEEK_FLAG_FRAME)

	packet := AVPacket{}
	for ctx.ReadFrame(&packet) >= 0 {
		if packet.StreamIndex() == videoStream.Index() {
			pts := time.Duration(float64(packet.Pts()) * videoStream.Timebase().Q2D() * (float64(time.Second)))

			if m.v.codecCtx.DecodeVideo(m.v.frame, &packet) {
				packet.Free()
				if start-pts < 10*time.Millisecond {
					break
				}
			} else {
				packet.Free()
			}
		}
	}

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

	if m.s != nil {
		m.s.seek(start)
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
