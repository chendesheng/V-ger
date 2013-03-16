package download

import (
	// "errors"
	"fmt"
	"io"
	// "time"
	// "bytes"
	// "log"
	// "bytes"
)

var play_quit chan bool

func Play(url string, w io.Writer, from, to int64) {
	fmt.Println("from ", from, " to ", to)
	if play_quit != nil {
		go func(quit chan bool) {
			for i := 0; i < 50; i++ {
				quit <- true
			}
			fmt.Println("quit previous download")
		}(play_quit)
	}
	control := make(chan int)
	play_quit = make(chan bool, 50)

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
			// fmt.Printf("write to client %d-%d\n", db.from, db.to)
			if !ok {
				fmt.Println("play output finish")
				return
			}
			_, err := w.Write(db.data)
			if err != nil {
				fmt.Println(err)
				for i := 0; i < 50; i++ {
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
