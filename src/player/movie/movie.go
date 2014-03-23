package movie

import (
	"log"
	"path/filepath"
	. "player/audio"
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
)

type Movie struct {
	ctx AVFormatContext
	v   *Video
	a   *Audio
	s   *Subtitle
	s2  *Subtitle
	c   *Clock
	w   *Window
	p   *Playing

	chSeekPause chan time.Duration

	quit        chan bool
	finishClose chan bool

	subFiles []string

	audioStreams []AVStream
}

func NewMovie() *Movie {
	m := &Movie{}
	m.quit = make(chan bool)
	return m
}

func (m *Movie) Open(w *Window, file string, subFiles []string) {
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

	m.p = CreateOrGetPlaying(filename)
	m.chSeekPause = make(chan time.Duration)

	m.ctx = ctx
	// dur := ctx.Duration()
	// if dur < 0 {
	// 	dur = -dur
	// }
	duration := time.Duration(float64(ctx.Duration()) / AV_TIME_BASE * float64(time.Second))
	m.c = NewClock(duration)
	m.c.Pause()
	func() {
		time.After(100 * time.Millisecond)
		m.c.Resume()
	}()

	m.setupVideo()
	m.w = w
	if len(subFiles) == 0 {
		go m.SearchDownloadSubtitle()
	}
	w.InitEvents()
	w.SetTitle(filename)
	w.SetSize(m.v.Width, m.v.Height)
	m.v.SetRender(m.w)

	println("audio")
	m.setupAudio()
	println("setupSubtitles")
	m.setupSubtitles(subFiles)

	m.uievents()

	start, _, _ := m.v.Seek(m.p.LastPos)
	// start := m.p.LastPos
	// start := time.Duration(0)
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
	println("open return")
}
func (m *Movie) Close() {
	m.w.FlushImageBuffer()
	m.w.RefreshContent(nil)
	m.w.ShowStartupView()

	SavePlaying(m.p)

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
func (m *Movie) PlayAsync() {
	go m.v.Play()
	go m.decode(m.p.Movie)
}
func (m *Movie) Resume() {
	m.c.Resume()
}
func (m *Movie) Pause() {
	m.c.Pause()
}

func tabs(t time.Duration) time.Duration {
	if t < 0 {
		t = -t
	}
	return t
}

func (m *Movie) setupVideo() {
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

func (m *Movie) SendPacket(index int, ch chan *AVPacket, packet AVPacket) bool {
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
func (m *Movie) showProgress(name string) {
	m.p.LastPos = m.c.GetTime()

	p := m.c.CalcPlayProgress(m.c.GetPercent())

	t, err := task.GetTask(name)
	if err == nil {
		if t.Status == "Finished" {
			p.Percent2 = 1
		} else {
			p.Percent2 = float64(t.BufferedPosition) / float64(t.Size)
		}
	} else {
		log.Print(err)
	}

	m.w.SendShowProgress(p)
}

func (m *Movie) decode(name string) {
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
				if m.SendPacket(m.a.StreamIndex(), m.a.PacketChan, packet) {
					continue
				}
			}

			packet.Free()
		} else {
			bufferring = true
			m.c.Pause()

			m.a.FlushBuffer()
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
