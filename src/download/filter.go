package download

import (
	"block"
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
	WriteAtQuit(bk block.Block, quit chan struct{})
}

type basicFilter struct {
	input  chan block.Block
	output chan block.Block
	quit   chan struct{}
}

func (f *basicFilter) connect(next *basicFilter) {
	next.input = f.output
}

func (f *basicFilter) writeOutput(b block.Block) bool {
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

	log.Print(output)
}

func activeFilters(filters []filter) {
	lastIndex := len(filters) - 1
	for _, f := range filters[:lastIndex] {
		go f.active()
	}

	filters[lastIndex].active()
}

func doDownload(t *task.Task, w io.WriterAt, from, to int64,
	maxSpeed int, chMaxSpeed chan int, restartTimeout time.Duration, quit chan struct{}) {
	url := t.URL

	maxConnections := util.ReadIntConfig("max-connection")

	gf := &generateFilter{
		basicFilter{nil, make(chan block.Block, maxConnections*2), quit},
		from,
		to,
		maxSpeed,
		chMaxSpeed,
		maxConnections * 2,
	}

	df := &downloadFilter{
		basicFilter{nil, make(chan block.Block), quit},
		url,
		false,
		maxConnections,
	}

	sf := &sortFilter{
		basicFilter{nil, make(chan block.Block), quit},
		from,
	}

	var tf *timeoutFilter
	if restartTimeout > 0 {
		tf = &timeoutFilter{
			basicFilter{nil, make(chan block.Block), quit},
			restartTimeout,
		}
	}

	wf := &writeFilter{
		basicFilter{nil, make(chan block.Block), quit},
		t.Name,
		w,
	}

	pf := &progressFilter{
		basicFilter{nil, make(chan block.Block), quit},
		t,
	}

	gf.connect(&df.basicFilter)
	// go lf.connect(&gf.basicFilter, &df.basicFilter) //will block unless lf is actived
	df.connect(&sf.basicFilter)
	sf.connect(&tf.basicFilter)
	tf.connect(&wf.basicFilter)
	wf.connect(&pf.basicFilter)
	pf.connect(&gf.basicFilter) //circle

	activeFilters([]filter{gf, df, sf, tf, wf, pf})

	return
}

func streaming(url string, w WriterAtQuit, from, to int64,
	sm SpeedMonitor, quit chan struct{}) {

	maxConnections := util.ReadIntConfig("max-connection")

	gf := &generateFilter{
		basicFilter{nil, make(chan block.Block, 2*maxConnections), quit},
		from,
		to,
		0,
		nil,
		maxConnections * 2,
	}

	df := &downloadFilter{
		basicFilter{nil, make(chan block.Block), quit},
		url,
		true,
		maxConnections,
	}

	sf := &sortFilter{
		basicFilter{nil, make(chan block.Block), quit},
		from,
	}

	swf := &simpleWriteFilter{
		basicFilter{nil, make(chan block.Block), quit},
		w,
	}
	spf := &speedFilter{
		basicFilter{nil, make(chan block.Block), quit},
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
