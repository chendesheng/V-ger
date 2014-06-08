package movie

import (
	"io"
	"os"
	"sync"
	"time"
)

type block struct {
	off int64
	p   []byte
}

func (bk *block) inside(position int64) bool {
	return bk.off <= position && position < bk.off+int64(len(bk.p))
}

type buffer struct {
	sync.Mutex
	sync.Pool

	currentPos int64
	data       []*block
	size       int64
	capacity   int64
}

func (b *buffer) fromTo() (int64, int64) {
	if len(b.data) == 0 {
		return 0, 0
	}

	back := b.data[len(b.data)-1]
	return b.data[0].off, back.off + int64(len(back.p))
}

func min(a, b int64) int64 {
	if a < b {
		return a
	} else {
		return b
	}
}
func NewBuffer(size int64) *buffer {
	b := &buffer{}
	b.size = size
	b.data = make([]*block, 0, 50)
	b.currentPos = 0
	if b.size < 4*1024*1024*1024 {
		b.capacity = 20 * 1024 * 1024
	} else {
		b.capacity = 200 * 1024 * 1024
	}
	b.New = func() interface{} {
		return &block{0, make([]byte, 0)}
	}

	go func() {
		for _ = range time.Tick(30 * time.Second) {
			b.GC()
		}
	}()
	return b
}

func (b *buffer) GC() {
	b.Lock()
	defer b.Unlock()

	if len(b.data) == 0 {
		return
	}

	var bk *block
	var i int
	for i, bk = range b.data {
		if b.currentPos < bk.off+int64(len(bk.p)) {
			break
		} else {
			b.Put(bk)
		}
	}

	if i > 0 {
		copy(b.data, b.data[i:])
		b.data = b.data[:len(b.data)-i]
	}
}

func (b *buffer) Read(w io.Writer, require int64) int64 {
	if w == nil {
		return 0
	}

	b.Lock()
	defer b.Unlock()

	lastPos := b.currentPos

	nextPosition := b.currentPos + require
	if nextPosition > b.size {
		require = b.size - b.currentPos
		nextPosition = b.size
	}

	for _, bk := range b.data {
		if bk.inside(b.currentPos) {
			from := b.currentPos - bk.off
			to := min(int64(len(bk.p)), nextPosition-bk.off)

			w.Write(bk.p[from:to])
			b.currentPos += to - from

			if b.currentPos >= nextPosition {
				break
			}
		}
	}

	return b.currentPos - lastPos
}

func (b *buffer) WriteAtQuit(p []byte, off int64, quit chan bool) error {
	// println("WriteAt:", off, len(p))

	b.Lock()
	defer b.Unlock()

	for b.sizeAhead() > b.capacity {
		//pause downloading if it is 20M ahead,
		b.Unlock()
		select {
		case <-time.After(100 * time.Millisecond):
			b.Lock()
			break
		case <-quit:
			b.Lock()
			return nil
		}
	}

	bk := b.Get().(*block)
	if len(bk.p) < len(p) {
		bk.p = make([]byte, len(p))
	} else {
		bk.p = bk.p[:len(p)]
	}
	copy(bk.p, p)
	bk.off = off

	b.data = append(b.data, bk)

	return nil
}
func (b *buffer) SizeAhead() int64 {
	b.Lock()
	defer b.Unlock()

	return b.sizeAhead()
}

func (b *buffer) sizeAhead() int64 {
	pos := b.currentPos

	for _, bk := range b.data {
		if bk.inside(pos) {
			pos = bk.off + int64(len(bk.p))
		}
	}
	return pos - b.currentPos
}

func (b *buffer) IsFinish() bool {
	b.Lock()
	defer b.Unlock()

	if len(b.data) == 0 {
		return false
	}

	bk := b.data[len(b.data)-1]
	return b.size <= bk.off+int64(len(bk.p))
}
func (b *buffer) Wait(size int64) {
	println("Wait:", b.SizeAhead(), b.IsFinish())
	for !(b.SizeAhead() >= size || b.IsFinish()) {
		time.Sleep(100 * time.Millisecond)
	}
}

func (b *buffer) Seek(offset int64, whence int) (int64, int64) {

	b.Lock()
	defer b.Unlock()

	switch whence {
	case os.SEEK_SET:
		b.currentPos = offset
		break
	case os.SEEK_CUR:
		b.currentPos += offset
		break
	case os.SEEK_END:
		b.currentPos = b.size + offset
		break
	}

	if b.currentPos > b.size {
		b.currentPos = b.size
	}

	from, to := b.fromTo()
	if b.currentPos >= from && b.currentPos < to {
		return b.currentPos, -1
	} else {
		for _, bk := range b.data {
			b.Put(bk)
		}
		b.data = b.data[0:0]
		return b.currentPos, b.currentPos
	}
}
