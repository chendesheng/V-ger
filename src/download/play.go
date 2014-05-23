package download

import (
	"io"
	"log"
	"sync"
	"task"
)

var play_quit chan bool

func Play(t *task.Task, w io.Writer, from, to int64) {
	log.Print("playing download from ", from, " to ", to)

	if play_quit != nil {
		ensureQuit(play_quit)
	}

	t.Status = "Playing"
	task.SaveTask(t)

	play_quit = make(chan bool)

	doDownload(t, writerWrap{w}, from, to, 0, nil, 0, play_quit)
}

var downloadQuit chan bool
var lock sync.Mutex

func Streaming(t *task.Task, w WriterAtQuit, from int64, sm SpeedMonitor) {
	println("start download:", from)
	if downloadQuit != nil {
		ensureQuit(downloadQuit)
	}

	lock.Lock()
	defer lock.Unlock()

	downloadQuit = make(chan bool)

	t.BufferedPosition = from
	task.SaveTask(t)

	println("speed monitor:", sm)

	streaming(t, w, from, t.Size, sm, downloadQuit)

	println("stop download:", from)
}

type writerWrap struct {
	w io.Writer
}

func (w writerWrap) WriteAt(p []byte, off int64) (int, error) {
	return w.w.Write(p)
}
