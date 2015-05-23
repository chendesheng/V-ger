package movie

import (
	"testing"
)

func TestLevenshtein(t *testing.T) {
	a := "Turn.S01E08.720p.HDTV.x264-IMMERSE.mkv"
	//b := "Turn.S01E07.720p.HDTV.x264-DIMEMSION.srt"
	b := "Turn.S01E08.720p.HDTV.X264-IMMERSE.srt"
	println(getNameDistance(a, b))
}
