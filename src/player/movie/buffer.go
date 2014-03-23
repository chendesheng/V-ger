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
		} else {
			break
		}
	}
}

func (b *buffer) Read(w io.Writer, require int64) int {
	b.Lock()
	defer b.Unlock()

	if b.currentPos+require > b.size {
		require = b.size - b.currentPos
	}
	ret := require

	nextPosition := b.currentPos + require
	for {
		b.read(w, require)
		if b.currentPos < nextPosition {
			require = nextPosition - b.currentPos
			b.Unlock()
			time.Sleep(100 * time.Millisecond)
			b.Lock()
		} else {
			break
		}
	}

	return int(ret)
}

func (b *buffer) read(w io.Writer, require int64) {
	nextPosition := b.currentPos + require
	for e := b.data.Front(); e != nil; e = e.Next() {
		bk := (e.Value).(*block)
		if bk.inside(b.currentPos) {
			from := b.currentPos - bk.off
			to := min(int64(len(bk.p)), nextPosition-bk.off)
			if w != nil {
				w.Write(bk.p[from:to])
			}

			b.currentPos = bk.off + to

			if b.currentPos >= nextPosition {
				break
			}
		}
	}
}

func (b *buffer) WriteAtQuit(p []byte, off int64, quit chan bool) error {
	// println("WriteAt:", off, len(p))

	b.Lock()
	defer b.Unlock()

	for off > b.currentPos && off-b.currentPos > 10*1024*1024 {
		//pause downloading if it is 10M ahead,
		b.Unlock()
		select {
		case <-time.After(100 * time.Millisecond):
			break
		case <-quit:
			return nil
		}
		b.Lock()
	}

	b.data.PushBack(&block{off, p})

	return nil
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
