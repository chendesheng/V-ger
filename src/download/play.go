package download

import (
	"errors"
	"fmt"
	"io"
	// "time"
	// "bytes"
	// "log"
	"bytes"
)

var play_quit chan bool

func Play(url string, w io.Writer, from, to int64) {
	fmt.Println("from ", from, " to ", to)
	if play_quit != nil {
		for i := 0; i < 100; i++ {
			play_quit <- true
		}
		fmt.Println("quit previous download")
	}
	control := make(chan int)
	play_quit = make(chan bool, 100)

	input := make(chan block)
	output := make(chan *dataBlock)

	go func(quit chan bool) {
		generateBlock(input, from, to, 0, control, quit)
	}(play_quit)

	go func(quit chan bool) {
		concurrentDownload(url, input, output, quit, from, to)
	}(play_quit)

	for {
		select {
		case db, ok := <-output:
			fmt.Printf("write to client %d-%d\n", db.from, db.to)
			if !ok {
				fmt.Println("play output finish")
				return
			}
			_, err := w.Write(db.data)
			if err != nil {
				fmt.Println(err)
				for i := 0; i < 100; i++ {
					play_quit <- true
				}
				play_quit = nil

				return
			}
			if db.to == to {
				fmt.Println("play finish")
				return
			}
		case <-play_quit:
			fmt.Println("play_quit quit")
			return
		}
	}
}

type OnlineVideo struct {
	pos         int64
	downloading bool

	buf  *bytes.Buffer
	quit chan bool

	size int64
	url  string
}

func NewOnlineVideo(url string) *OnlineVideo {
	url, _, size := GetDownloadInfo(url)

	return &OnlineVideo{0, false, nil, nil, size, url}
}
func (o *OnlineVideo) stopDownload() {
	o.downloading = false
	if o.quit != nil {
		for i := 0; i < 100; i++ {
			o.quit <- true
		}
	}
}
func (o *OnlineVideo) startDownload() {
	o.downloading = true
	o.buf = &bytes.Buffer{}

	o.quit = make(chan bool, 100)

	input := make(chan block)
	output := make(chan *dataBlock)

	go func(quit chan bool) {
		generateBlock(input, o.pos, o.size, 0, make(chan int), quit)
	}(o.quit)

	go func(quit chan bool) {
		concurrentDownload(o.url, input, output, quit, o.pos, o.size)
	}(o.quit)

	for {
		select {
		case db, ok := <-output:
			if !ok {
				return
			}
			o.buf.Write(db.data)
		case <-o.quit:
			return
		}
	}
}

func (o *OnlineVideo) Read(bytes []byte) (int, error) {
	if !o.downloading {
		o.startDownload()
	}

	n, err := o.buf.Read(bytes)
	o.pos += int64(n)
	return n, err
}

var errWhence = errors.New("Seek: invalid whence")
var errOffset = errors.New("Seek: invalid offset")

func (o *OnlineVideo) Seek(offset int64, whence int) (ret int64, err error) {
	o.stopDownload()

	switch whence {
	default:
		return 0, errWhence
	case 0:
		offset += 0
	case 1:
		offset += o.pos
	case 2:
		offset += o.size
	}
	if offset < 0 || offset > o.size {
		return 0, errOffset
	}
	o.pos = offset
	return offset - 0, nil
}
