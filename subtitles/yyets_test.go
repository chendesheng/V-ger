package subtitles

import (
	"log"
	"testing"
)

func TestYYetsSearch(t *testing.T) {
	yyets := &yyetsSearch{"The Vampire Diaries S06E10", 1, nil}
	result := make(chan Subtitle)
	go func() {
		for s := range result {
			log.Print(s)
		}
	}()
	err := yyets.search(result)
	if err != nil {
		log.Print(err)
	}
	close(result)
}
