package download

import (
	"io"
	"log"
	"task"
)

var play_quit chan struct{}

func Play(t *task.Task, w io.Writer, from, to int64) {
	log.Print("playing download from ", from, " to ", to)

	if play_quit != nil {
		ensureQuit(play_quit)
	}

	t.Status = "Playing"
	task.SaveTask(t)

	play_quit = make(chan struct{})

	doDownload(t, writerWrap{w}, from, to, 0, nil, 0, play_quit)
}

//guarantee only one streaming, and could restart any moment
type Streaming struct {
	url   string
	size  int64
	w     WriterAtQuit
	sm    SpeedMonitor
	quit  chan struct{}
	chArg chan int64
}

func (s *Streaming) run() {
	for {
		from := <-s.chArg
		s.quit = make(chan struct{})

		streaming(s.url, s.w, from, s.size, s.sm, s.quit)
	}
}
func (s *Streaming) Restart(pos int64) {
	log.Print("Streaming restart:", pos, s)

	if s.quit != nil {
		close(s.quit)
	}
	s.chArg <- pos
}

func StartStreaming(url string, size int64, w WriterAtQuit, sm SpeedMonitor) *Streaming {
	s := &Streaming{url, size, w, sm, nil, make(chan int64)}
	go s.run()

	return s
}

type writerWrap struct {
	w io.Writer
}

func (w writerWrap) WriteAt(p []byte, off int64) (int, error) {
	return w.w.Write(p)
}
