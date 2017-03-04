package subtitles

import (
	"log"
	"net/http"
	"sync"
	// "fmt"
	// "io"

	// "os"
)

type subSearcher interface {
	search(chan Subtitle) error
}

type Subtitle struct {
	URL         string
	Description string
	Source      string
	Context     http.Header
}

func (s *Subtitle) String() string {
	return s.Description
}

func SearchSubtitles(name string, url string, result chan Subtitle, quit chan struct{}) {
	searchers := []subSearcher{
		&yyetsSearch{name, 1, quit},
		//&kankanSearch{url, quit},
		&addic7ed{name, quit},
	}

	w := sync.WaitGroup{}
	w.Add(len(searchers))
	for _, s := range searchers {
		go func(s subSearcher) {
			if err := s.search(result); err != nil {
				log.Print(err)
			}
			w.Done()
		}(s)
	}
	w.Wait()
	close(result)
}
