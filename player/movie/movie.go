package movie

import (
	"errors"
	"fmt"
	"log"
	"path/filepath"
	"strings"
	"sync"
	"time"
	"vger/block"
	"vger/download"
	"vger/player/clock"
	"vger/player/gui"
	"vger/player/libav"
	. "vger/player/movie/audio"
	. "vger/player/movie/seeking"
	. "vger/player/movie/video"
	"vger/player/shared"
	"vger/player/subtitle"
)

type Movie struct {
	ctx libav.AVFormatContext
	v   *Video
	a   *Audio
	c   *clock.Clock
	w   *gui.Window
	p   *shared.Playing
	movieSubs

	quit        chan struct{}
	subQuit     chan struct{}
	finishClose chan bool

	subs []*subtitle.Subtitle

	audioStreams []libav.AVStream

	httpBuffer *buffer
	streaming  *download.Streaming

	chSeek     chan *seekArg
	chHold     chan time.Duration
	chProgress chan time.Duration
	chSpeed    chan float64

	Filename string
	seeking  *Seeking
}

type movieSubs struct {
	sync.Mutex

	s1 *subtitle.Subtitle
	s2 *subtitle.Subtitle
}

func (m *movieSubs) getPlayingSubs() (*subtitle.Subtitle, *subtitle.Subtitle) {
	m.Lock()
	defer m.Unlock()

	return m.s1, m.s2
}
func (m *movieSubs) setPlayingSubs(s1, s2 *subtitle.Subtitle) {
	m.Lock()
	defer m.Unlock()

	m.s1, m.s2 = s1, s2
}
func (m *movieSubs) SeekPlayingSubs(t time.Duration) {
	s1, s2 := m.getPlayingSubs()
	if s1 != nil {
		s1.Seek(t)
	}
	if s2 != nil {
		s2.Seek(t)
	}
}
func (m *movieSubs) stopPlayingSubs() {
	s1, s2 := m.getPlayingSubs()
	m.setPlayingSubs(nil, nil)

	if s1 != nil {
		s1.Stop()
	}
	if s2 != nil {
		s2.Stop()
	}
}

type seekArg struct {
	t     time.Duration
	isEnd bool
}

func New() *Movie {
	log.Print("New movie")

	m := &Movie{}
	m.Reset()
	return m
}

func (m *Movie) Reset() {
	m.quit = make(chan struct{})
	m.chProgress = make(chan time.Duration)
	m.finishClose = make(chan bool)
	m.Filename = ""
	m.a = nil
	m.v = nil
	m.p = nil
	m.c = nil
	m.httpBuffer = nil
	m.s1 = nil
	m.s2 = nil
	m.subs = nil
	m.chSeek = nil
	m.chHold = nil
	m.chSpeed = nil
	m.seeking = nil
}

func checkDownloadSubtitle(m *Movie, file string, filename string) {
	if strings.Contains(file, "googlevideo.com/videoplayback") {
		return
	}

	subs := shared.GetSubtitlesMap(filename)
	log.Printf("%v", subs)
	if len(subs) == 0 {
		if m.p.FirstOpen {
			m.searchDownloadSubtitle()
		}
	} else {
		log.Print("setupSubtitles")
		m.setupSubtitles(subs)

		m.SeekPlayingSubs(m.c.GetTime())
	}
}

func setBufferCapacity(buf *buffer, duration time.Duration) {
	var capacity int64
	if duration < 10*time.Minute {
		capacity = buf.size
	} else {
		//overflow, divide before cross
		capacity = int64(float64(buf.size) / float64(duration) * 1 * float64(time.Hour))
		log.Print(capacity)
	}

	if capacity > 1000*block.MB {
		capacity = 1000 * block.MB
	}
	if capacity < 10*block.MB {
		capacity = 10 * block.MB
	}
	buf.SetCapacity(capacity)
}

func (m *Movie) setupContext(file string) (filename string, duration time.Duration, err error) {
	log.Print("setupContext")

	var ctx libav.AVFormatContext

	if strings.HasPrefix(file, "http://") ||
		strings.HasPrefix(file, "https://") {

		ctx, filename, err = m.openHttp(file)
		if err != nil {
			return
		}
	} else {
		filename = filepath.Base(file)

		ctx = libav.NewAVFormatContext()
		if err = ctx.OpenInput(file); err != nil {
			return
		}
	}

	if err = ctx.FindStreamInfo(); err != nil {
		return
	}
	ctx.DumpFormat()
	m.ctx = ctx

	duration = ctx.Duration()
	return
}

func (m *Movie) Open(w *gui.Window, file string) (err error) {
	w.SendShowSpinning()
	defer w.SendHideSpinning(false)

	defer func() {
		if err != nil {
			close(m.finishClose)
		}
	}()

	log.Print("open ", file)
	m.w = w
	m.uievents()

	filename, duration, err := m.setupContext(file)
	if err != nil {
		return
	}
	m.Filename = filename

	if m.httpBuffer != nil {
		setBufferCapacity(m.httpBuffer, duration)
	}

	m.c = clock.New(duration)

	m.p = shared.CreateOrGetPlaying(filename)
	m.p.Duration = duration

	err = m.setupVideo()
	if err != nil {
		return
	}

	err = m.setupAudio()
	if err != nil {
		return
	}

	m.seeking = NewSeeking(m.v, m, m.quit)
	m.uiProgressBarEvents()

	go checkDownloadSubtitle(m, file, filename)
	return
}

func (m *Movie) SavePlaying() {
	if m.p != nil {
		shared.SavePlaying(m.p)
	}
}

func (m *Movie) Close() {
	if m.subQuit != nil {
		close(m.subQuit)
	}

	close(m.quit)

	m.w.ClearEvents()

	m.stopPlayingSubs()

	if m.httpBuffer != nil {
		m.httpBuffer.Close()
		m.streaming.Stop()
	}

	<-m.finishClose

	m.w.SendDestoryRender()
}

func (m *Movie) PlayAsync() {
	log.Print("movie play async")

	go m.decode()
	go m.showProgressPerSecond()
}

func (m *Movie) setupVideo() error {
	log.Print("setup video")

	ctx := m.ctx
	videoStream := ctx.VideoStream()
	if !videoStream.IsNil() {
		var err error
		m.v, err = NewVideo(ctx, videoStream, m.c, m.w, m)
		if err != nil {
			return err
		}
	} else {
		return fmt.Errorf("No video stream find.")
	}

	return nil
}

func (m *Movie) IsPlaying() bool {
	if m.c == nil {
		return false
	} else {
		return m.c.IsRunning()
	}
}

func (m *Movie) GetSubtitleNames() (names []string) {
	if len(m.subs) > 0 {
		names = make([]string, len(m.subs))
		for i, sub := range m.subs {
			names[i] = sub.Name
		}
	}

	return
}
func (m *Movie) GetPlayingSubtitles() (firstSub int, secondSub int) {
	firstSub = -1
	secondSub = -1

	if len(m.subs) > 0 {
		s1, s2 := m.getPlayingSubs()
		for i, sub := range m.subs {
			if s1 == sub {
				firstSub = i
			} else if s2 == sub {
				secondSub = i
			}
		}
	}

	return
}

func (m *Movie) TogglePlay() {
	if m.c != nil {
		m.c.Toggle()
	}
}

func (m *Movie) GetAllAudioTracks() (names []string) {
	if len(m.audioStreams) > 0 {
		names = make([]string, len(m.audioStreams))
		for i, stream := range m.audioStreams {
			dic := stream.MetaData()
			mp := dic.Map()
			title := strings.TrimSpace(mp["title"])
			language := strings.TrimSpace(strings.ToLower(mp["language"]))

			if len(title) == 0 && len(language) == 0 {
				names[i] = fmt.Sprintf("Track %d", i+1)
			} else if len(language) > 0 {
				names[i] = fmt.Sprintf("Track %d [%s]", i+1, language)
			} else if len(title) > 0 {
				names[i] = fmt.Sprintf("Track %d - %s", i+1, title)
			} else {
				names[i] = fmt.Sprintf("Track %d [%s] - %s", i+1, title, language)
			}
		}
	}
	return
}

func (m *Movie) GetPlayingAudioTrack() int {
	for i, stream := range m.audioStreams {
		if m.a.StreamIndex() == stream.Index() {
			return i
		}
	}

	return -1
}

func (m *Movie) SeekBySubtitle(forward bool) {
	var offset time.Duration
	s1, _ := m.getPlayingSubs()
	var r int
	if forward {
		r = 1
	} else {
		r = -1
	}

	if s1 != nil {
		t := m.c.GetTime()
		subTime := s1.GetSubTime(t, r)
		log.Print("subtime:", subTime)

		if subTime == 0 {
			offset = time.Duration(r) * 10 * time.Second
		} else {
			offset = subTime - t
		}
	} else {
		offset = time.Duration(r) * 10 * time.Second
	}
	m.SeekOffset(offset)
}

// increase or decrease volume
// return -1 means audio not ready
func (m *Movie) AddVolume(d int) int {
	if m.a == nil {
		return -1
	}

	var volume int
	if d < 0 {
		volume = m.a.DecreaseVolume()
	} else {
		volume = m.a.IncreaseVolume()
	}

	m.p.Volume = volume
	shared.SavePlayingAsync(m.p)

	go func() {
		select {
		case chVolume <- struct{}{}:
		case <-m.quit:
		case <-time.After(100 * time.Millisecond):
		}
	}()

	return volume
}

func (m *Movie) SyncMainSubtitle(d time.Duration) (time.Duration, error) {
	s1, _ := m.getPlayingSubs()
	if s1 != nil {
		offset := s1.AddOffset(d)
		shared.UpdateSubtitleOffsetAsync(s1.Name, offset)

		return offset, nil
	} else {
		return 0, errors.New("main subtitle not exists")
	}
}

func (m *Movie) SyncSecondSubtitle(d time.Duration) (time.Duration, error) {
	_, s2 := m.getPlayingSubs()
	if s2 != nil {
		offset := s2.AddOffset(d)
		shared.UpdateSubtitleOffsetAsync(s2.Name, offset)

		return offset, nil
	} else {
		return 0, errors.New("second subtitle not exists")
	}
}

func (m *Movie) SyncAudio(d time.Duration) time.Duration {
	return m.a.AddOffset(d)
}

func (m *Movie) ToggleSubtitle(index int) {
	log.Print("toggle subtitle:", index)

	subs := m.subs
	clicked := subs[index]

	var s1, s2 *subtitle.Subtitle
	ps1, ps2 := m.getPlayingSubs()

	if ps1 == nil && ps2 == nil {
		//add playing s1
		s1 = clicked
		// go s1.Play()

		m.p.Sub1 = s1.Name
		m.p.Sub2 = ""

	} else if ps1 == clicked {
		//remove playing s1
		ps1.Stop()
		if ps2 != nil {
			s1 = ps2
			s1.IsMainSub = true

			m.p.Sub1 = s1.Name
			m.p.Sub2 = ""
		}
	} else if ps2 == clicked {
		//remove playing s2
		ps2.Stop()
		s1 = ps1

		m.p.Sub1 = s1.Name
		m.p.Sub2 = ""
	} else {
		//replace playing subtitle
		if clicked.IsTwoLangs() {
			s1 = clicked
			s2 = nil
		} else if ps1.IsTwoLangs() {
			s1 = clicked
			s2 = nil
		} else if isLangEqual(ps1.Lang1, clicked.Lang1) {
			s1 = clicked
			s2 = ps2
		} else if ps2 == nil {
			s1 = ps1
			s2 = clicked
		} else if isLangEqual(ps2.Lang1, clicked.Lang1) {
			s1 = ps1
			s2 = clicked
		} else { //third language which is impossible for now
			s1 = ps1
			s2 = clicked
		}

		if s1 != ps1 {
			ps1.Stop()

			s1.IsMainSub = true
			// go s1.Play()

			m.p.Sub1 = s1.Name
		}

		if s2 != nil {
			if s2 != ps2 {
				if ps2 != nil {
					ps2.Stop()
				}

				s2.IsMainSub = false
				// go s2.Play()

				m.p.Sub2 = s2.Name
			}
		} else {
			if ps2 != nil {
				ps2.Stop()
			}

			m.p.Sub2 = ""
		}
	}

	m.setPlayingSubs(s1, s2)
	shared.SavePlayingAsync(m.p)
}

func (m *Movie) SetAudioTrack(i int) {
	if m.audioStreams[i].Index() == m.a.StreamIndex() {
		return
	}

	m.a.Close()
	err := m.a.Open(m.audioStreams[i])
	if err != nil {
		log.Print(err)
	} else {
		m.p.SoundStream = m.a.StreamIndex()
		shared.SavePlayingAsync(m.p)
	}
}
