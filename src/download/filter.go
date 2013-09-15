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

type aFilter struct {
	input  chan *block
	output chan *block
	quit   chan bool
}

func (f *aFilter) connect(next *aFilter) {
	next.input = f.output
}

type generateFilter struct {
	aFilter
	from        int64
	to          int64
	chBlockSize chan int64
}

func (gf *generateFilter) active() {
	generateBlock(gf.input, gf.output, gf.chBlockSize, gf.from, gf.to, gf.quit)
}

type downloadFilter struct {
	aFilter
	url string
}

func (df *downloadFilter) active() {
	concurrentDownload(df.url, df.input, df.output, df.quit)
}

type writeFilter struct {
	aFilter
	w io.Writer
}

func (wf *writeFilter) active() {
	writeOutput(wf.w, wf.input, wf.output, wf.quit)
}

type progressFilter struct {
	aFilter
	t *task.Task
}

func (pf *progressFilter) active() {
	handleProgress(pf.input, pf.output, pf.t, pf.quit)
}

type sortFilter struct {
	aFilter
	from int64
}

func (sf *sortFilter) active() {
	sortOutput(sf.input, sf.output, sf.quit, sf.from)
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
	maxSpeed int64, quit chan bool) {
	url := t.URL
	// for {
	// 	finalUrl, _, _, err := GetDownloadInfo(url)
	// 	if err == nil {
	// 		url = finalUrl
	// 		break
	// 	}

	// 	select {
	// 	case <-quit:
	// 		return
	// 	default:
	// 		time.Sleep(time.Second * 2)
	// 	}
	// }

	gf := &generateFilter{
		aFilter{nil, make(chan *block), quit},
		from,
		to,
		make(chan int64),
	}

	df := &downloadFilter{
		aFilter{nil, make(chan *block), quit},
		url,
	}

	sf := &sortFilter{
		aFilter{nil, make(chan *block), quit},
		from,
	}

	wf := &writeFilter{
		aFilter{nil, make(chan *block), quit},
		w,
	}

	pf := &progressFilter{
		aFilter{nil, make(chan *block), quit},
		t,
	}

	gf.connect(&df.aFilter)
	lf.connect(&gf.aFilter, &df.aFilter, quit, gf.chBlockSize)
	df.connect(&sf.aFilter)
	sf.connect(&wf.aFilter)
	wf.connect(&pf.aFilter)
	pf.connect(&gf.aFilter) //circle

	activeFilters([]filter{gf, lf, df, sf, wf, pf})

	return
}
