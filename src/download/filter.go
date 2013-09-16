package download

import (
	"io"
	"log"
	"task"
	// "time"
)

type filter interface {
	active()
}

type basicFilter struct {
	input  chan *block
	output chan *block
	quit   chan bool
}

func (f *basicFilter) connect(next *basicFilter) {
	next.input = f.output
}

func activeFilters(filters []filter) {
	log.Printf("filters: %v", filters)
	lastIndex := len(filters) - 1
	for _, f := range filters[:lastIndex] {
		go f.active()
	}

	filters[lastIndex].active()
}

func doDownload(t *task.Task, w io.Writer, from, to int64,
	maxSpeed int64, chMaxSpeed chan int64, quit chan bool) {
	url := t.URL

	gf := &generateFilter{
		basicFilter{nil, make(chan *block), quit},
		from,
		to,
		maxSpeed,
		chMaxSpeed,
	}

	df := &downloadFilter{
		basicFilter{nil, make(chan *block), quit},
		url,
		5,
	}

	sf := &sortFilter{
		basicFilter{nil, make(chan *block), quit},
		from,
	}

	wf := &writeFilter{
		basicFilter{nil, make(chan *block), quit},
		w,
	}

	pf := &progressFilter{
		basicFilter{nil, make(chan *block), quit},
		t,
	}

	// gf.connect(&df.basicFilter)
	go lf.connect(&gf.basicFilter, &df.basicFilter) //will block unless lf is actived
	df.connect(&sf.basicFilter)
	sf.connect(&wf.basicFilter)
	wf.connect(&pf.basicFilter)
	pf.connect(&gf.basicFilter) //circle

	activeFilters([]filter{gf, lf, df, sf, wf, pf})

	return
}
