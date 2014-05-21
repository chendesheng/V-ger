package download

import (
	"io"
	"log"
	"task"
	"time"
	"util"
)

type filter interface {
	active()
}

type WriterAtQuit interface {
	//this method should return (nil) asap after close(quit)
	WriteAtQuit(p []byte, off int64, quit chan bool) error
}

type writeAtWrap struct {
	iow io.WriterAt
}

func (w writeAtWrap) WriteAtQuit(p []byte, off int64, quit chan bool) error {
	_, err := w.iow.WriteAt(p, off)
	return err
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

func doDownload(t *task.Task, w WriterAtQuit, from, to int64,
	maxSpeed int64, chMaxSpeed chan int64, restartTimeout time.Duration, m ProgressMonitor, quit chan bool) {
	url := t.URL

	maxConnections := util.ReadIntConfig("max-connection")

	gf := &generateFilter{
		basicFilter{nil, make(chan *block, maxConnections), quit},
		from,
		to,
		maxSpeed,
		chMaxSpeed,
		maxConnections * 2,
	}

	df := &downloadFilter{
		basicFilter{nil, make(chan *block), quit},
		url,
		maxConnections,
	}

	sf := &sortFilter{
		basicFilter{nil, make(chan *block), quit},
		from,
	}

	wf := &writeFilter{
		basicFilter{nil, make(chan *block), quit},
		w,
		restartTimeout,
	}

	pf := &progressFilter{
		basicFilter{nil, make(chan *block), quit},
		t,
		m,
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
