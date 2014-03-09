package main

import (
	// "bytes"
	"download"
	"fmt"
	"log"
	"os"
	"path/filepath"
	. "player/clock"
	. "player/gui"
	. "player/libav"
	. "player/shared"
	. "player/subtitle"
	. "player/video"
	"strings"
	"subscribe"
	"task"
	"time"
	"util"
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
	p   *Playing

	chCtrl      chan ctrlArg
	chSeekPause chan time.Duration

	quit        chan bool
	finishClose chan bool

	subFiles []string
}

func (m *movie) setupAudio() {
	ctx := m.ctx

	audioStreams := ctx.AudioStream()

	audioStreamNames := make([]string, 0)
	audioStreamIndexes := make([]int32, 0)

	log.Print("setupAudio:", len(audioStreams))

	if len(audioStreams) > 0 {

		for _, stream := range audioStreams {
			dic := stream.MetaData()
			mp := dic.Map()
			title := mp["title"]                        //dic.AVDictGet("title", AVDictionaryEntry{}, 2).Value()
			language := strings.ToLower(mp["language"]) //dic.AVDictGet("language", AVDictionaryEntry{}, 2).Value()

			// println(title, language)
			audioStreamNames = append(audioStreamNames, fmt.Sprintf("[%s] %s", language, title))
			audioStreamIndexes = append(audioStreamIndexes, int32(stream.Index()))
		}

		selected := audioStreams[0].Index()
		for i := len(audioStreams) - 1; i >= 0; i-- {
			stream := audioStreams[i]

			dic := stream.MetaData()
			mp := dic.Map()
			language := strings.ToLower(mp["language"])
			if strings.Contains(language, "eng") {
				selected = stream.Index()
			}

			if m.p.SoundStream == stream.Index() {
				selected = m.p.SoundStream
				break
			}
		}

		m.a = &audio{
			streams: audioStreams,
			volume:  m.p.Volume,
		}

		m.a.setCurrentStream(selected)
		m.a.c = m.c

		if len(audioStreams) > 1 {
			m.w.InitAudioMenu(audioStreamNames, audioStreamIndexes, m.a.stream.Index())
		} else {
			HideAudioMenu()
		}
	} else {
		HideAudioMenu()
	}
}

func (m *movie) setupSubtitles(subFiles []string) {
	if len(subFiles) > 0 {
		tags := make([]int32, 0)
		names := make([]string, 0)

		m.subFiles = subFiles

		m.w.FuncSubtitleMenuClicked = append(m.w.FuncSubtitleMenuClicked, func(index int, showOrHide bool) {
			go func(m *movie) {
				subFiles = m.subFiles
				if showOrHide {
					// m.s.Stop()
					width, height := m.w.GetWindowSize()
					s := NewSubtitle(subFiles[index], m.w, m.c, float64(width), float64(height))
					if s != nil {
						if m.s == nil {
							m.s = s
							s.IsMainOrSecondSub = true
						} else {
							m.s2 = s
							s.IsMainOrSecondSub = false
						}
						go s.Play()
					}
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

				if m.s != nil {
					m.p.Sub1 = m.s.Name
				} else {
					m.p.Sub1 = ""
				}

				if m.s2 != nil {
					m.p.Sub2 = m.s2.Name
				} else {
					m.p.Sub2 = ""
				}

				SavePlayingAsync(m.p)
			}(m)
		})

		println("play subtitle:", subFiles)
		width, height := m.w.GetWindowSize()

		if len(m.p.Sub1) == 0 && len(m.p.Sub2) > 0 {
			m.p.Sub1 = m.p.Sub2
			m.p.Sub2 = ""

			SavePlayingAsync(m.p)
		}

		var s1, s2 *Subtitle
		if len(m.p.Sub1) > 0 {
			s1 = NewSubtitle(m.p.Sub1, m.w, m.c, float64(width), float64(height))
			if s1 != nil {
				s1.IsMainOrSecondSub = true

				if s1 != nil {
					go s1.Play()
				}
			} else {
				m.p.Sub1 = ""
			}
		}

		if len(m.p.Sub2) > 0 {
			s2 = NewSubtitle(m.p.Sub1, m.w, m.c, float64(width), float64(height))
			if s2 != nil {
				s2.IsMainOrSecondSub = false

				if s2 != nil {
					go s2.Play()
				}
			} else {
				m.p.Sub2 = ""
			}
		}

		if s1 != nil {
			m.s = s1
			m.s2 = s2
		} else {
			m.s = s2
		}

		if m.s == nil && m.s2 == nil {
			println("auto select default subtitle")

			var en, cn, double *Subtitle
			for _, file := range subFiles {
				s := NewSubtitle(file, m.w, m.c, float64(width), float64(height))
				if s != nil {
					if en == nil && s.Lang1 == "en" && len(s.Lang2) == 0 {
						en = s
					}
					if cn == nil && s.Lang1 == "cn" && len(s.Lang2) == 0 {
						cn = s
					}

					if double == nil && len(s.Lang1) > 0 && len(s.Lang2) > 0 {
						double = s
					}
				}
			}

			if double != nil {
				m.s = double
			} else {
				if cn != nil {
					m.s = cn
					m.s2 = en
				} else {
					m.s = en
				}
			}

			if m.s != nil {
				m.s.IsMainOrSecondSub = true
				m.p.Sub1 = m.s.Name
				go m.s.Play()
			}

			if m.s2 != nil {
				m.s.IsMainOrSecondSub = false
				m.p.Sub2 = m.s2.Name
				go m.s2.Play()
			}

			SavePlayingAsync(m.p)
		}

		selected1 := -1
		selected2 := -1
		for i, n := range subFiles {
			tags = append(tags, int32(i))
			names = append(names, filepath.Base(n))

			if m.s != nil && n == m.s.Name {
				selected1 = i
			}

			if m.s2 != nil && n == m.s2.Name {
				selected2 = i
			}
		}

		if selected1 == -1 && selected2 == -1 {
			selected1 = 0
		}

		if len(names) > 0 {
			m.w.InitSubtitleMenu(names, tags, selected1, selected2)
		} else {
			HideSubtitleMenu()
		}
	} else {
		println("remove subtitle menu")
		HideSubtitleMenu()
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

func (m *movie) openHttp(file string) (AVFormatContext, string) {
	_, name, size, err := download.GetDownloadInfo(file)
	if err != nil {
		log.Fatal(err)
	}
	t, err := task.GetTask(name)
	if err != nil {
		// log.Fatal(err)
		t = &task.Task{}
		t.Name = name
		t.Size = size
		t.StartTime = time.Now().Unix()
		t.Status = "Stopped"
		t.URL = file
		task.SaveTask(t)
	}

	mbuf := &util.Buffer{}

	// currentPos := int64(0)
	buf := AVObject{}
	buf.Malloc(1024 * 32)
	ioctx := NewAVIOContext(buf, func(buf AVObject) int {
		if buf.Size() == 0 {
			return 0
		}

		size := buf.Size()
		// currentPos := mbuf.CurrentPos

		// println("read: currentPos11:", mbuf.CurrentPos)
		// for size > 0 && currentPos+int64(size) < t.Size {
		for {

			if bytes, err := mbuf.Read(size); err == nil {
				buf.Write(bytes)
				// size -= len(bytes)
				currentPos := mbuf.GetCurrentPos()
				currentPos += int64(len(bytes))
				mbuf.SetCurrentPos(currentPos)
				println("readfunc:", currentPos, len(bytes))
				return len(bytes)
			}

			time.Sleep(50 * time.Millisecond)
		}

		// println("read: currentPos:", mbuf.CurrentPos)
		// return buf.Size() - size
	}, func(pos int64, whence int) int64 {
		println("seekfunc:", pos, whence)

		// download.Play(t, w, from, to)
		// t, _ = task.GetTask(t.Name)
		currentPos := mbuf.GetCurrentPos()
		switch whence {
		case os.SEEK_SET:
			currentPos = pos
			break
		case os.SEEK_CUR:
			currentPos += pos
			break
		case os.SEEK_END:
			currentPos = t.Size + pos
			break
		// case AVSEEK_SIZE:
		default:
			return t.Size
		}

		if currentPos > t.Size {
			currentPos = t.Size
			return currentPos
		}

		if currentPos < 0 {
			return -1
		}
		mbuf.SetCurrentPos(currentPos)
		mbuf.ClearData()
		go download.Play(t, mbuf, currentPos, t.Size)
		return currentPos
	})

	ctx := NewAVFormatContext()
	ctx.SetPb(ioctx)

	download.Play(t, mbuf, 0, 1024*32)

	if bytes, err := mbuf.Read(1024 * 32); err == nil {
		println("bytes:", len(bytes))
		pd := NewAVProbeData()

		obj := AVObject{}
		obj.Malloc(len(bytes))
		obj.Write(bytes)

		pd.SetBuffer(obj)
		pd.SetFileName("")

		ctx.SetInputFormat(pd.InputFormat())

		mbuf.ClearData()

		go download.Play(t, mbuf, 0, t.Size)

		ctx.OpenInput("")
	}

	return ctx, name
}

func (m *movie) open(w *Window, file string, subFiles []string) {
	println("open ", file)

	var ctx AVFormatContext
	var filename string

	if strings.HasPrefix(file, "http://") {
		ctx, filename = m.openHttp(file)
	} else {
		ctx = AVFormatContext{}
		ctx.OpenInput(file)
		if ctx.IsNil() {
			log.Fatal("open failed:", file)
			return
		}

		filename = filepath.Base(file)

		ctx.FindStreamInfo()
		ctx.DumpFormat()

	}

	m.chCtrl = make(chan ctrlArg)
	m.chSeekPause = make(chan time.Duration)

	m.ctx = ctx
	// dur := ctx.Duration()
	// if dur < 0 {
	// 	dur = -dur
	// }
	duration := time.Duration(float64(ctx.Duration()) / AV_TIME_BASE * float64(time.Second))
	m.c = NewClock(duration)

	m.setupVideo()
	m.w = w
	// m.w = NewWindow(filename, m.v.Width, m.v.Height)
	w.InitEvents()
	w.SetTitle(filename)
	w.SetSize(m.v.Width, m.v.Height)
	m.v.SetRender(m.w)

	m.uievents()

	m.setupAudio()

	m.setupSubtitles(subFiles)

	start, _, _ := m.v.Seek(m.p.LastPos)
	// start := m.p.LastPos
	m.p.LastPos = start
	m.p.Duration = duration

	if t, _ := task.GetTask(m.p.Movie); t != nil {
		println("get subscribe:", t.Subscribe)
		if subscr := subscribe.GetSubscribe(t.Subscribe); subscr != nil && subscr.Duration == 0 {
			subscribe.UpdateDuration(t.Subscribe, duration)
		}
	}

	SavePlayingAsync(m.p)

	m.c.Reset()
	m.c.SetTime(start)

	if m.s != nil {
		m.s.Seek(start)
	}

	go m.showProgress(filename)
}

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
		case <-m.quit:
			return false
		}
	}
	return false
}
func (m *movie) showProgress(name string) {
	m.p.LastPos = m.c.GetTime()

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
	defer func() {
		if m.a != nil {
			m.a.Close()
		}
		if m.v != nil {
			m.v.Close()
		}
		m.c.Reset()
		m.ctx.CloseInput()

		if m.finishClose != nil {
			close(m.finishClose)
		}
	}()

	packet := AVPacket{}
	ctx := m.ctx
	go func() {
		ticker := time.NewTicker(time.Second)
		for {
			if m.c.WaitUtilRunning(m.quit) {
				return
			}

			select {
			case <-ticker.C:
				m.showProgress(name)
			case <-m.quit:
				return
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
			case <-m.quit:
				return
			}

		}
		// println(bufferring)
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

	// m.w.RefreshContent(img)
	go m.w.SendDrawImage(img)

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
	println("seek:", t.String())
	t, img, err = m.v.Seek(t)
	println("end seek:", t.String())
	if err != nil {
		return t
	}
	// println("seek refresh")
	if len(img) > 0 {
		m.w.RefreshContent(img)
	}

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
	println("seek end:", t.String())
}

func (m *movie) Close() {
	m.w.FlushImageBuffer()
	m.w.RefreshContent(nil)
	m.w.ShowStartupView()

	SavePlayingAsync(m.p)

	m.finishClose = make(chan bool)
	close(m.quit)
	// time.Sleep(100 * time.Millisecond)

	m.w.ClearEvents()

	if m.s != nil {
		m.s.Stop()
		m.s = nil
	}

	if m.s2 != nil {
		m.s2.Stop()
		m.s2 = nil
	}

	<-m.finishClose
}
