package download

import (
	"log"
	"time"
	"vger/block"
)

type limitFilter struct {
	chMaxSpeed chan int

	isActive bool

	chFilter chan *basicFilter
}

type blockFilter struct {
	b block.Block
	f *basicFilter
}

func (lf *limitFilter) active() {
	if lf.isActive {
		return
	}
	lf.isActive = true

	ch := make(chan *blockFilter)

	maxSpeed := 0
	for {
		s := time.Now()
		select {
		case bf := <-ch:
			// println("limitfilter")
			b := bf.b
			f := bf.f
			if len(b.Data) == 0 {
				f.closeOutput()
				return
			}

			f.writeOutput(b)

			if maxSpeed > 0 {
				size := len(b.Data)
				d1 := time.Duration(float64(time.Second) * float64(size) / float64(maxSpeed*block.KB))
				d2 := time.Now().Sub(s)
				if d1 > d2 {
					time.Sleep(d1 - d2)
				}
			}
		case maxSpeed = <-lf.chMaxSpeed:
			log.Print("set max speed: ", maxSpeed)
		case f := <-lf.chFilter:
			go func(f *basicFilter) {
				for {
					select {
					case b, ok := <-f.input:
						select {
						case ch <- &blockFilter{b, f}:
						case <-f.quit:
							return
						}
						if !ok {
							return
						}
					case <-f.quit:
						return
					}
				}
			}(f)
		}
	}
}

func (lf *limitFilter) connect(f1 *basicFilter, f2 *basicFilter) {
	f2.input = make(chan block.Block)
	f := &basicFilter{
		f1.output, f2.input, f1.quit,
	}

	lf.chFilter <- f
}

var lf *limitFilter = &limitFilter{
	make(chan int), false,
	make(chan *basicFilter, 0),
}

func LimitSpeed(speed int) error {
	if !lf.isActive {
		return nil
	}

	go func() {
		lf.chMaxSpeed <- speed
		for _, tc := range taskControls {
			tc.chMaxSpeed <- speed
		}
	}()

	return nil
}
