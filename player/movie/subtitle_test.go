package movie

import (
	"fmt"
	"testing"
	"time"
)

func TestLevenshtein(t *testing.T) {
	a := "Turn.S01E08.720p.HDTV.x264-IMMERSE.mkv"
	//b := "Turn.S01E07.720p.HDTV.x264-DIMEMSION.srt"
	b := "Turn.S01E08.720p.HDTV.X264-IMMERSE.srt"
	println(getNameDistance(a, b))
}

func TestDownloadSub(t *testing.T) {
	unarPath = "unar"
	movieName := "Under the Dome s03e03"
	url := ""
	quit := make(chan struct{})
	go func() {
		select {
		case <-quit:
			fmt.Println("search quit")
			break
		case <-time.After(SearchSubtitleTimeout):
			fmt.Println("search timeout")
			close(quit)
			break
		}
	}()
	downloadSubs(movieName, url, movieName, quit)
}
