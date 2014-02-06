package download

import (
	// "errors"
	// "fmt"
	"io"
	"task"
	// "time"
	// "bytes"
	"log"
	// "bytes"
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

	doDownload(t, writerAtWrap{w}, from, to, 0, nil, play_quit)
}

type writerAtWrap struct {
	w io.Writer
}

func (w writerAtWrap) WriteAt(p []byte, off int64) (int, error) {
	return w.w.Write(p)
}
