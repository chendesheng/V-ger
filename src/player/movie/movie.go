package movie

import (
	"block"
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

	quit        chan struct{}
	finishClose chan bool

	subs []*Sub

	audioStreams []AVStream

	size int64

	httpBuffer *buffer

	// chSeekEnd      chan time.Duration
	chSeekProgress chan *seekArg
	chPause        chan chan time.Duration
	chProgress     chan time.Duration
	chSpeed        chan float64
}

type seekArg struct {
	t     time.Duration
	isEnd bool
}

func NewMovie() *Movie {
	m := &Movie{}
	m.quit = make(chan struct{})
	return m
}

func updateSubscribeDuration(movie string, duration time.Duration) {
	if t, _ := task.GetTask(movie); t != nil {
		println("get subscribe:", t.Subscribe)
		if subscr := subscribe.GetSubscribe(t.Subscribe); subscr != nil && subscr.Duration == 0 {
			subscribe.UpdateDuration(t.Subscribe, duration)
		}
	}
}

func (m *Movie) Open(w *Window, file string) {
	println("open ", file)

	m.w = w
	w.InitEvents()
	m.uievents()

	var ctx AVFormatContext
	var filename string

	if strings.HasPrefix(file, "http://") ||
		strings.HasPrefix(file, "https://") {
		ctx, filename = m.openHttp(file)
		if ctx.IsNil() {
			log.Fatal("open failed: ", file)
			return
		}
		ctx.FindStreamInfo()
	} else {
		ctx = NewAVFormatContext()
		ctx.OpenInput(file)
		if ctx.IsNil() {
			log.Fatal("open failed:", file)
			return
		}

		filename = filepath.Base(file)

		ctx.FindStreamInfo()
		ctx.DumpFormat()
	}

	m.ctx = ctx

	duration := ctx.Duration2()
	if m.httpBuffer != nil {
		var capacity int64
		if duration < 10*time.Minute {
			capacity = m.httpBuffer.size
		} else {
			//overflow, divide before cross
			capacity = int64(float64(m.httpBuffer.size) / float64(duration) * 10 * float64(time.Minute))
			println(capacity)
		}

		if capacity > 200*block.MB {
			capacity = 200 * block.MB
		}
		if capacity < 10*block.MB {
			capacity = 10 * block.MB
		}
		m.httpBuffer.SetCapacity(capacity)
	}

	m.c = NewClock(duration)

	m.setupVideo()

	m.p = CreateOrGetPlaying(filename)
	log.Print("video duration:", duration.String(), m.p.LastPos)
	var start time.Duration
	if m.p.LastPos > time.Second && m.p.LastPos < duration-50*time.Millisecond {
		var img []byte
		start, img, _ = m.v.Seek(m.p.LastPos)
		w.SendDrawImage(img)

		if m.httpBuffer != nil {
			m.httpBuffer.Wait(3 * 1024 * 1024)
		}
	}

	m.p.LastPos = start
	m.p.Duration = duration

	go updateSubscribeDuration(m.p.Movie, m.p.Duration)

	go func() {
		if strings.Contains(file, "googlevideo.com/videoplayback") {
			return
		}

		subs := GetSubtitlesMap(filename)
		log.Printf("%v", subs)
		if len(subs) == 0 {
			m.SearchDownloadSubtitle()
		} else {
			println("setupSubtitles")
			m.setupSubtitles(subs)

			if m.s != nil {
				m.s.Seek(m.c.GetTime())
			}
			if m.s2 != nil {
				m.s2.Seek(m.c.GetTime())
			}
		}
	}()

	w.SendSetTitle(filename)
	w.SendSetSize(m.v.Width, m.v.Height)
	m.v.SetRender(m.w)

	m.setupAudio()

	m.c.SetTime(start)

	m.showProgressInner(start)

	w.SendSetCursor(false)
	w.SendHideSpinning()
}

func (m *Movie) SavePlaying() {
	SavePlaying(m.p)
}

func (m *Movie) Close() {
	m.w.FlushImageBuffer()
	m.w.RefreshContent(nil)
	m.w.ShowStartupView()

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
	go m.showProgressPerSecond()
	go m.decode(m.p.Movie)
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

func (m *Movie) ResumeClock() {
	m.c.Resume()
}

func (m *Movie) PauseClock() {
	m.c.Pause()
}
