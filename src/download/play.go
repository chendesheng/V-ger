package download

import (
	"io"
	"log"
	"sync"
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

var downloadQuit chan struct{}
var lock sync.Mutex

func Streaming(url string, size int64, w WriterAtQuit, from int64, sm SpeedMonitor) {
	println("start download:", from)
	if downloadQuit != nil {
		ensureQuit(downloadQuit)
	}

	lock.Lock()
	defer lock.Unlock()

	downloadQuit = make(chan struct{})

	println("speed monitor:", sm)

	streaming(url, w, from, size, sm, downloadQuit)

	println("stop download:", from)
}

type writerWrap struct {
	w io.Writer
}

func (w writerWrap) WriteAt(p []byte, off int64) (int, error) {
	return w.w.Write(p)
}
