package download

import (
	"fmt"
	"log"
	"time"
)

type limitFilter struct {
	f1 *aFilter
	f2 *aFilter

	chMaxSpeed chan int64

	isActive bool

	bs1 chan int64
	bs2 chan int64
}

func (lf *limitFilter) active() {
	if lf.isActive {
		return
	}
	lf.isActive = true
	defer func() {
		lf.isActive = false
	}()

	maxSpeed := int64(0)

	f1 := lf.f1
	f2 := lf.f2
	for {
		s := time.Now()
		select {
		case maxSpeed = <-lf.chMaxSpeed:
			if lf.bs1 != nil {
				lf.bs1 <- maxSpeed
			}
			if lf.bs2 != nil {
				lf.bs2 <- maxSpeed
			}

			log.Print("set max speed ", maxSpeed)
		case b, ok := <-f1.input:
			handleInput(b, ok, f1, s, maxSpeed)
		case b, ok := <-f2.input:
			handleInput(b, ok, f2, s, maxSpeed)
		case <-f1.quit:
			setFilterNil(f1)
			lf.bs1 = nil
		case <-f2.quit:
			setFilterNil(f2)
			lf.bs2 = nil
		}
		if f1.quit == nil && f2.quit == nil {
			return
		}
	}
}

func setFilterNil(f *aFilter) {
	if f.output != nil {
		close(f.output)
	}
	f.output = nil
	f.input = nil
	f.quit = nil
}
func handleInput(b *block, ok bool, f *aFilter, s time.Time, maxSpeed int64) {
	if !ok {
		setFilterNil(f)
		if f == lf.f1 {
			lf.bs1 = nil
		} else {
			lf.bs2 = nil
		}
		return
	}
	select {
	case f.output <- b:
		// log.Printf("maxSpeed %d", maxSpeed)

		if maxSpeed > 0 {
			size := b.to - b.from
			d1 := time.Duration(float64(time.Second) * float64(size) / float64(maxSpeed*1024))
			d2 := time.Now().Sub(s)
			if d1 > d2 {
				// log.Print("sleep ", d1-d2)
				time.Sleep(d1 - d2)
			}
		}
	case <-f.quit:
		setFilterNil(f)
		if f == lf.f1 {
			lf.bs1 = nil
		} else {
			lf.bs2 = nil
		}
	}
}

func initConnect(f, f1, f2 *aFilter, quit chan bool) bool {
	if f.input == nil {
		if f.output == nil {
			f.output = make(chan *block)
		}
		f.quit = quit

		f1.connect(f)
		f.connect(f2)
		return true
	}
	return false
}
func (lf *limitFilter) connect(f1, f2 *aFilter, quit chan bool, bs chan int64) bool {
	if !initConnect(lf.f1, f1, f2, quit) {
		if initConnect(lf.f2, f1, f2, quit) {
			lf.bs2 = bs
			return true
		}
		return false
	} else {
		lf.bs1 = bs
	}

	return true
}

var lf *limitFilter = &limitFilter{
	&aFilter{nil, make(chan *block), nil},
	&aFilter{nil, make(chan *block), nil},
	make(chan int64), false, nil, nil,
}

func LimitSpeed(speed int64) error {
	select {
	case lf.chMaxSpeed <- speed:
		break
	case <-time.After(time.Second * 5):
		return fmt.Errorf("Limit speed operation timeout")
	}

	return nil
}
