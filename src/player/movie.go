package main

import (
	"fmt"
	"path/filepath"
	. "player/clock"
	. "player/libav"
	. "player/subtitle"
	"strings"
	"task"
	// . "player/shared"
	"log"
	. "player/gui"
	. "player/video"
	"time"
	// "util"
)

type movie struct {
	ctx AVFormatContext
	v   *Video
	a   *audio
	s   *Subtitle
	s2  *Subtitle
	c   *Clock
	w   *Window

	chSeek chan time.Duration
}

func (m *movie) setupAudio() {
	ctx := m.ctx

	audioStreams := ctx.AudioStream()

	audioStreamNames := make([]string, 0)
	audioStreamIndexes := make([]int32, 0)

	if len(audioStreams) > 0 {

		for _, stream := range audioStreams {
			dic := stream.MetaData()
			m := dic.Map()
			title := m["title"]                        //dic.AVDictGet("title", AVDictionaryEntry{}, 2).Value()
			language := strings.ToLower(m["language"]) //dic.AVDictGet("language", AVDictionaryEntry{}, 2).Value()

			// println(title, language)
			audioStreamNames = append(audioStreamNames, fmt.Sprintf("[%s] %s", language, title))
			audioStreamIndexes = append(audioStreamIndexes, int32(stream.Index()))
		}

		selected := audioStreams[0].Index()
		for _, stream := range audioStreams {
			dic := stream.MetaData()
			m := dic.Map()
			language := strings.ToLower(m["language"])
			if strings.Contains(language, "eng") {
				selected = stream.Index()
				break
			}
		}

		m.a = &audio{streams: audioStreams}
		m.a.setCurrentStream(selected)
		m.a.c = m.c

		if len(audioStreams) > 1 {
			m.w.InitAudioMenu(audioStreamNames, audioStreamIndexes, m.a.stream.Index())
		}
	}
}

func (m *movie) setupSubtitles(subFiles []string) {
	if len(subFiles) > 0 {
		tags := make([]int32, 0)
		names := make([]string, 0)
		for i, n := range subFiles {
			tags = append(tags, int32(i))
			names = append(names, filepath.Base(n))
		}
		m.w.InitSubtitleMenu(names, tags, 0)
		m.w.FuncSubtitleMenuClicked = append(m.w.FuncSubtitleMenuClicked, func(index int, showOrHide bool) {
			go func(m *movie, subFiles []string) {
				if showOrHide {
					// m.s.Stop()
					s := NewSubtitle(subFiles[index], m.w, m.c)
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
		m.s = NewSubtitle(subFiles[0], m.w, m.c)
		go m.s.Play(0)
	}
}

func (m *movie) setupVideo() {
	ctx := m.ctx
	videoStream := ctx.VideoStream()
	if !videoStream.IsNil() {
		var err error
		m.v, err = NewVideo(ctx, videoStream, m.c)
		if err != nil {
			log.Fatal(err)
			return
		}

	} else {
		log.Fatal("No video stream find.")
	}
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

	m.chSeek = make(chan time.Duration)

	m.ctx = ctx
	m.c = NewClock(time.Duration(float64(ctx.Duration()) / AV_TIME_BASE * float64(time.Second)))

	m.setupVideo()
	m.w = NewWindow(filepath.Base(file), m.v.Width, m.v.Height)
	m.v.SetRender(m.w)

	m.uievents()

	m.setupAudio()

	m.setupSubtitles(subFiles)

	start = m.v.Seek(start)

	m.c.Reset()
	m.c.SetTime(start)

	if m.s != nil {
		m.s.Seek(start)
	}

	go m.showProgress(filepath.Base(file))
}
func (m *movie) SeekTo(t time.Duration) time.Duration {
	// t = t / time.Second * time.Second

	if m.a != nil {
		m.a.flushBuffer()
	}

	if m.v != nil {
		m.v.FlushBuffer()
		t = m.v.Seek(t)
	}

	println("seek to", t.String())

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
func (m *movie) drawCurrentFrame() time.Duration {
	ctx := m.ctx
	v := m.v
	if v == nil {
		return time.Duration(0)
	}

	packet := AVPacket{}

	// frame := v.frame
	for ctx.ReadFrame(&packet) >= 0 {
		// if v.stream.Index() == packet.StreamIndex() {
		if frameFinished, t, bytes := v.DecodeAndScale(&packet); frameFinished {
			packet.Free()

			m.w.RefreshContent(bytes)

			return t
		} else {
			packet.Free()
		}
	}

	return time.Duration(0)
}

func (m *movie) SendPacket(index int, ch chan *AVPacket, packet AVPacket) bool {
	if index == packet.StreamIndex() {
		pkt := packet
		pkt.Dup()
		select {
		case ch <- &pkt:
			return true
		case t := <-m.chSeek:
			packet.Free()
			m.c.SetTime(m.SeekTo(t))
			return true
		}
	}
	return false
}
func (m *movie) showProgress(name string) {
	p := m.c.CalcPlayProgress(m.c.GetPercent())

	t, err := task.GetTask(name)
	if err == nil {
		p.Percent2 = float64(t.DownloadedSize) / float64(t.Size)
	} //else {
	// log.Print(err)
	//}

	m.w.SendShowProgress(p)
}
func (m *movie) decode(name string) {
	packet := AVPacket{}
	ctx := m.ctx

	ticker := time.NewTicker(time.Second)

	for {
		// m.c.WaitUtilRunning()
		select {
		case <-ticker.C:
			go m.showProgress(name)
		default:
		}

		if ctx.ReadFrame(&packet) >= 0 {
			// if m.SendPacket(m.v.StreamIndex, m.v.ChanPacket, packet) {
			// 	continue
			// }
			if m.v.StreamIndex == packet.StreamIndex() {
				if frameFinished, pts, img := m.v.DecodeAndScale(&packet); frameFinished {
					//make sure seek operations not happens before one frame finish decode
					//if not, segment fault & crash
					select {
					case m.v.ChanDecoded <- &VideoFrame{pts, img}:
						break
					case t := <-m.chSeek:
						m.c.SetTime(m.SeekTo(t))
						break
					}
				}
				packet.Free()
				continue
			}

			if m.a != nil {
				if m.SendPacket(m.a.stream.Index(), m.a.ch, packet) {
					continue
				}
			}

			packet.Free()
		} else {
			m.stop()
			return
		}
	}
}

func (m *movie) play() {
	go m.v.Play()

	if m.w != nil {
		PollEvents()
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
