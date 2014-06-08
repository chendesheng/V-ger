package download

import (
	"io"
	"log"
	"runtime/debug"
	"sync"
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

type basicFilter struct {
	input  chan block
	output chan block
	quit   chan bool
}

func (f *basicFilter) connect(next *basicFilter) {
	next.input = f.output
}

func (f *basicFilter) writeOutput(b block) bool {
	if f.output != nil {
		select {
		case f.output <- b:
		case <-f.quit:
			return true
		}
	}

	return false
}

//wait until input is closed or quit
func (f *basicFilter) drainInput() {
	for {
		select {
		case _, ok := <-f.input:
			if !ok {
				return
			}
		case <-f.quit:
			return
		}
	}
}

func (f *basicFilter) closeOutput() {
	if f.output != nil {
		close(f.output)
	}
}

func (f *basicFilter) closeQuit() {
	log.Print(string(debug.Stack()))
	ensureQuit(f.quit)
}

func (f *basicFilter) wait(d time.Duration) {
	select {
	case <-f.quit:
		return
	case <-time.After(d):
		break
	}
}

var traceLock sync.Mutex

func trace(output string) {
	traceLock.Lock()
	defer traceLock.Unlock()

	print("[trace]")
	println(output)
}

func activeFilters(filters []filter) {
	log.Printf("filters: %v", filters)
	lastIndex := len(filters) - 1
	for _, f := range filters[:lastIndex] {
		go f.active()
	}

	filters[lastIndex].active()
}

func doDownload(t *task.Task, w io.WriterAt, from, to int64,
	maxSpeed int64, chMaxSpeed chan int64, restartTimeout time.Duration, quit chan bool) {
	url := t.URL

	maxConnections := util.ReadIntConfig("max-connection")

	gf := &generateFilter{
		basicFilter{nil, make(chan block, maxConnections*2), quit},
		from,
		to,
		maxSpeed,
		chMaxSpeed,
		maxConnections * 2,
	}

	df := &downloadFilter{
		basicFilter{nil, make(chan block), quit},
		url,
		false,
		maxConnections,
	}

	sf := &sortFilter{
		basicFilter{nil, make(chan block), quit},
		from,
	}

	wf := &writeFilter{
		basicFilter{nil, make(chan block), quit},
		t.Name,
		w,
		restartTimeout,
	}

	pf := &progressFilter{
		basicFilter{nil, make(chan block), quit},
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

func streaming(url string, w WriterAtQuit, from, to int64,
	sm SpeedMonitor, quit chan bool) {

	maxConnections := util.ReadIntConfig("max-connection")

	gf := &generateFilter{
		basicFilter{nil, make(chan block, 2*maxConnections), quit},
		from,
		to,
		0,
		nil,
		maxConnections * 2,
	}

	df := &downloadFilter{
		basicFilter{nil, make(chan block), quit},
		url,
		true,
		maxConnections,
	}

	sf := &sortFilter{
		basicFilter{nil, make(chan block), quit},
		from,
	}

	swf := &simpleWriteFilter{
		basicFilter{nil, make(chan block), quit},
		w,
	}

	spf := &speedFilter{
		basicFilter{nil, make(chan block), quit},
		sm,
	}

	gf.connect(&df.basicFilter)
	df.connect(&sf.basicFilter)
	sf.connect(&swf.basicFilter)
	swf.connect(&spf.basicFilter)
	spf.connect(&gf.basicFilter) //circle

	activeFilters([]filter{gf, df, sf, swf, spf})

	return
}
