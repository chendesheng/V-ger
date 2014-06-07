package movie

import (
	"container/list"
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
	data       *list.List
	size       int64
}

func (b *buffer) fromTo() (int64, int64) {
	front := b.data.Front()
	back := b.data.Back()
	if front == nil {
		return 0, 0
	} else {
		head := front.Value.(*block)
		tail := back.Value.(*block)
		return head.off, tail.off + int64(len(tail.p))
	}
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
	b.data = list.New()
	b.currentPos = 0
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

	for {
		e := b.data.Front()
		if e == nil {
			break
		}

		bk := e.Value.(*block)
		if b.currentPos >= bk.off+int64(len(bk.p)) {
			b.data.Remove(e)

			b.Put(bk)
		} else {
			break
		}
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

	for e := b.data.Front(); e != nil; e = e.Next() {
		bk := (e.Value).(*block)
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

	for b.sizeAhead() > 20*1024*1024 {
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

	b.data.PushBack(bk)

	return nil
}
func (b *buffer) SizeAhead() int64 {
	b.Lock()
	defer b.Unlock()

	return b.sizeAhead()
}

func (b *buffer) sizeAhead() int64 {
	pos := b.currentPos

	for e := b.data.Front(); e != nil; e = e.Next() {
		bk := (e.Value).(*block)
		if bk.inside(pos) {
			pos = bk.off + int64(len(bk.p))
		}
	}

	return pos - b.currentPos
}

func (b *buffer) IsFinish() bool {
	b.Lock()
	defer b.Unlock()

	e := b.data.Back()
	if e == nil {
		return false
	}

	bk := e.Value.(*block)
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
		b.data = list.New()
		return b.currentPos, b.currentPos
	}
}
