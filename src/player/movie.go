package main

import (
	"fmt"
	"path/filepath"
	. "player/clock"
	. "player/libav"
	. "player/subtitle"
	"strings"
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
	s2  *Subtitle
	c   *Clock
}

func (m *movie) open(file string, subFiles []string, start time.Duration) {
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

	audioStreams := ctx.AudioStream()

	audioStreamNames := make([]string, 0)
	audioStreamIndexes := make([]int32, 0)

	if len(audioStreams) > 0 {

		selected := audioStreams[0].Index()
		for _, stream := range audioStreams {
			dic := stream.MetaData()
			m := dic.Map()
			title := m["title"]                        //dic.AVDictGet("title", AVDictionaryEntry{}, 2).Value()
			language := strings.ToLower(m["language"]) //dic.AVDictGet("language", AVDictionaryEntry{}, 2).Value()

			// println(title, language)
			audioStreamNames = append(audioStreamNames, fmt.Sprintf("[%s] %s", language, title))
			audioStreamIndexes = append(audioStreamIndexes, int32(stream.Index()))
			if strings.Contains(language, "eng") {
				selected = stream.Index()
			}
		}

		m.a = &audio{streams: audioStreams}
		m.a.setCurrentStream(selected)
		m.a.c = m.c
	}

	videoStream := ctx.VideoStream()
	if !videoStream.IsNil() {
		m.v = &video{}
		m.v.setup(ctx, videoStream, file, start)
		m.v.c = m.c

		if len(subFiles) > 0 {
			tags := make([]int32, 0)
			names := make([]string, 0)
			for i, n := range subFiles {
				tags = append(tags, int32(i))
				names = append(names, filepath.Base(n))
			}
			m.v.window.InitSubtitleMenu(names, tags, 0)
			m.v.window.FuncSubtitleMenuClicked = append(m.v.window.FuncSubtitleMenuClicked, func(index int, showOrHide bool) {
				go func(m *movie, subFiles []string) {
					if showOrHide {
						// m.s.Stop()
						s := NewSubtitle(subFiles[index], m.v.window, m.c)
						if m.s == nil {
							m.s = s
							s.IsMainOrSecondSub = true
						} else {
							m.s2 = s
							s.IsMainOrSecondSub = false
						}

						pos, _ := s.FindPos(m.c.GetSeekTime())
						s.Play(pos)
					} else {
						if (m.s != nil) && (m.s.Name == subFiles[index]) {
							m.s.Stop()
							if m.s2 != nil {
								m.s = m.s2
								m.s.IsMainOrSecondSub = true
								m.s2 = nil
							} else {
								m.s = nil
							}
						} else if (m.s2 != nil) && (m.s2.Name == subFiles[index]) {
							m.s2.Stop()
							m.s2 = nil
						}
					}
				}(m, subFiles)
			})

			println("play subtitle:", subFiles)
			m.s = NewSubtitle(subFiles[0], m.v.window, m.c)
			go m.s.Play(0)
		}
		m.uievents()
		start = m.v.seek(start)

		m.c.Reset()
		m.c.SetTime(start)

		if m.s != nil {
			m.s.Seek(start)
		}

		// for _, as := range audioStreams {
		// 	as.
		// }
		if len(audioStreams) > 1 {
			m.v.window.InitAudioMenu(audioStreamNames, audioStreamIndexes, m.a.stream.Index())
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
	if m.s2 != nil {
		m.s2.Seek(t)
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
