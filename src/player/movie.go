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

type seekArg struct {
	t   time.Duration
	res chan bool
}
type ctrlArg struct {
	c   int
	res chan interface{}
}
type movie struct {
	ctx AVFormatContext
	v   *Video
	a   *audio
	s   *Subtitle
	s2  *Subtitle
	c   *Clock
	w   *Window

	chCtrl      chan ctrlArg
	chSeekPause chan time.Duration
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

					s.Play()
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
		go m.s.Play()
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

	m.chCtrl = make(chan ctrlArg)
	m.chSeekPause = make(chan time.Duration)

	m.ctx = ctx
	m.c = NewClock(time.Duration(float64(ctx.Duration()) / AV_TIME_BASE * float64(time.Second)))

	m.setupVideo()
	m.w = NewWindow(filepath.Base(file), m.v.Width, m.v.Height)
	m.v.SetRender(m.w)

	m.uievents()

	m.setupAudio()

	m.setupSubtitles(subFiles)

	start, _, _ = m.v.Seek(start)

	m.c.Reset()
	m.c.SetTime(start)

	if m.s != nil {
		m.s.Seek(start)
	}

	go m.showProgress(filepath.Base(file))
}

// func (m *movie) SeekTo(t time.Duration) (time.Duration, []byte) {
// 	// t = t / time.Second * time.Second
// 	var img []byte

// 	if m.a != nil {
// 		m.a.flushBuffer()
// 	}

// 	if m.v != nil {
// 		// m.v.FlushBuffer()

// 		var err error
// 		t, img, err = m.v.Seek(t)
// 		if err != nil {
// 			// log.
// 		}
// 	}

// 	log.Print("seek to:", t.String())

// 	if m.s != nil {
// 		m.s.Seek(t)
// 	}
// 	if m.s2 != nil {
// 		m.s2.Seek(t)
// 	}

// 	return t, img
// }

func tabs(t time.Duration) time.Duration {
	if t < 0 {
		t = -t
	}
	return t
}

func (m *movie) getCurrentFrame() (time.Duration, []byte) {
	ctx := m.ctx
	v := m.v
	if v == nil {
		return time.Duration(0), nil
	}

	packet := AVPacket{}

	// frame := v.frame
	for ctx.ReadFrame(&packet) >= 0 {
		// if v.stream.Index() == packet.StreamIndex() {
		if frameFinished, t, bytes := v.DecodeAndScale(&packet); frameFinished {
			packet.Free()

			// m.w.RefreshContent(bytes)

			return t, bytes
		} else {
			packet.Free()
		}
	}

	return time.Duration(0), nil
}

func (m *movie) SendPacket(index int, ch chan *AVPacket, packet AVPacket) bool {
	if index == packet.StreamIndex() {
		pkt := packet
		pkt.Dup()
		select {
		case ch <- &pkt:
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

const (
	PAUSE = iota
	RESUME
)

func (m *movie) decode(name string) {
	packet := AVPacket{}
	ctx := m.ctx
	go func() {
		ticker := time.NewTicker(time.Second)
		for {
			select {
			case <-ticker.C:
				m.showProgress(name)
			}
		}
	}()

	bufferring := false
	for {
		resCode := ctx.ReadFrame(&packet)
		if resCode >= 0 {
			if bufferring {
				bufferring = false
				m.c.Resume()
			}
			if m.v.StreamIndex == packet.StreamIndex() {
				if frameFinished, pts, img := m.v.DecodeAndScale(&packet); frameFinished {
					//make sure seek operations not happens before one frame finish decode
					//if not, segment fault & crash
					select {
					case m.v.ChanDecoded <- &VideoFrame{pts, img}:
						break
					case t := <-m.chSeekPause:
						if t != -1 {
							break
						}
						for {
							t := <-m.chSeekPause
							if t >= 0 {
								m.c.SetTime(t)
								break
							}
						}
						break
					}

					t := m.c.GetTime()
					if m.s != nil {
						m.s.Seek(t)
					}
					if m.s2 != nil {
						m.s2.Seek(t)
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
			bufferring = true
			m.c.Pause()

			m.a.flushBuffer()
			m.v.FlushBuffer()

			t, _, err := m.v.Seek(m.c.GetTime())
			if err == nil {
				println("seek success:", t.String())
				m.c.SetTime(t)
				continue
			} else {
				log.Print("seek error:", err)
			}

			// println("seek to unfinished:", m.c.GetTime().String())
			log.Print("get frame error:", resCode)

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
			}

		}
		println(bufferring)
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

func (m *movie) seekOffset(offset time.Duration) {
	t := m.c.GetTime() + offset

	m.SeekBegin()

	var img []byte
	var err error
	t, img, err = m.v.SeekOffset(t)
	if err != nil {
		return
	}

	m.w.RefreshContent(img)

	m.c.SetTime(t)
	percent := m.c.GetPercent()
	m.w.ShowProgress(m.c.CalcPlayProgress(percent))

	if m.s != nil {
		m.s.Seek(t)
	}
	if m.s2 != nil {
		m.s2.Seek(t)
	}
	m.SeekEnd(t)
}

func (m *movie) SeekBegin() {
	m.chSeekPause <- -1
	m.v.FlushBuffer()
	m.a.flushBuffer()
}

func (m *movie) Seek(t time.Duration) time.Duration {
	var img []byte
	var err error
	t, img, err = m.v.Seek(t)
	if err != nil {
		return t
	}
	// println("seek refresh")

	m.w.RefreshContent(img)

	m.c.SetTime(t)
	percent := m.c.GetPercent()
	m.w.ShowProgress(m.c.CalcPlayProgress(percent))

	if m.s != nil {
		m.s.Seek(t)
	}
	if m.s2 != nil {
		m.s2.Seek(t)
	}
	return t
}

func (m *movie) SeekEnd(t time.Duration) {
	m.chSeekPause <- t
}
