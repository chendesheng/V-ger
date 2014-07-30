package movie

import (
	"block"
	"io"
	"log"
	"os"
	"runtime/debug"
	"sync"
	"time"
)

type buffer struct {
	sync.Mutex

	currentPos int64
	data       []*block.Block
	size       int64
	capacity   int64
}

func (b *buffer) fromTo() (int64, int64) {
	if len(b.data) == 0 {
		return 0, 0
	}

	bk := b.data[len(b.data)-1]
	return b.data[0].From, bk.To()
}

func min(a, b int64) int64 {
	if a < b {
		return a
	} else {
		return b
	}
}
func NewBuffer(size int64) *buffer {
	log.Print("NewBuffer:", size)

	b := &buffer{}
	b.size = size
	b.data = make([]*block.Block, 0, 50)
	b.currentPos = 0
	b.capacity = 20 * block.MB
	go func() {
		for _ = range time.Tick(30 * time.Second) {
			b.GC()
		}
	}()
	return b
}

func (b *buffer) SetCapacity(capacity int64) {
	log.Print("SetCapacity:", capacity)

	b.Lock()
	defer b.Unlock()

	b.capacity = capacity
}

func (b *buffer) GC() {
	b.Lock()
	defer b.Unlock()

	if len(b.data) == 0 {
		return
	}

	var bk *block.Block
	var i int
	for i, bk = range b.data {
		if b.currentPos < bk.To() {
			break
		} else {
			block.DefaultBlockPool.Put(bk)
		}
	}

	if i > 0 {
		l := copy(b.data, b.data[i:])
		b.data = b.data[:l]
	}
}

func (b *buffer) Read(w io.Writer, require int64) int64 {
	b.Lock()
	defer b.Unlock()

	lastPos := b.currentPos

	nextPosition := b.currentPos + require
	if nextPosition > b.size {
		require = b.size - b.currentPos
		nextPosition = b.size
	}

	for _, bk := range b.data {
		if bk.Inside(b.currentPos) {
			from := b.currentPos - bk.From
			to := min(int64(len(bk.Data)), nextPosition-bk.From)

			// log.Printf("http buffer read: %d-%d", from, to)
			w.Write(bk.Data[from:to])
			b.currentPos += to - from

			if b.currentPos >= nextPosition {
				break
			}
		}
	}

	return b.currentPos - lastPos
}

func (b *buffer) WriteAtQuit(bk block.Block, quit chan struct{}) {
	// log.Print("WriteAt:", bk.From, bk.To())

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
		}
	}

	// bk1 := block.DefaultBlockPool.Get(bk.From, len(bk.Data))
	// copy(bk1.Data, bk.Data)

	b.data = append(b.data, &bk)
}

func (b *buffer) CurrentPos() int64 {
	b.Lock()
	defer b.Unlock()

	return b.currentPos
}
func (b *buffer) SizeAhead() int64 {
	b.Lock()
	defer b.Unlock()

	return b.sizeAhead()
}

func (b *buffer) sizeAhead() int64 {
	pos := b.currentPos

	for _, bk := range b.data {
		if bk.Inside(pos) {
			pos = bk.To()
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
	return b.size <= bk.To()
}
func (b *buffer) Wait(size int64) {
	log.Print("Wait:", b.SizeAhead(), b.IsFinish())
	for b.BufferFinish(size) {
		time.Sleep(100 * time.Millisecond)
	}
}
func (b *buffer) BufferFinish(size int64) bool {
	return b.SizeAhead() < size && !b.IsFinish()
}
func (b *buffer) WaitQuit(size int64, quit chan struct{}) {
	// log.Print("WaitQuit:", b.SizeAhead(), b.IsFinish())
	for b.BufferFinish(size) {
		select {
		case <-time.After(100 * time.Millisecond):
		case <-quit:
			return
		}
	}
}

func (b *buffer) Seek(offset int64, whence int) (int64, int64) {
	b.Lock()
	defer b.Unlock()

	log.Print("buffer seek:", offset, whence)
	log.Print(string(debug.Stack()))

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

	from, to := b.fromTo()
	if b.currentPos >= from && b.currentPos < to {
		return b.currentPos, -1
	} else {
		b.Clear()
		return b.currentPos, b.currentPos
	}
}

func (b *buffer) BufferPercent() float64 {
	b.Lock()
	defer b.Unlock()

	res := float64(b.currentPos+b.sizeAhead()) / float64(b.size)
	if res >= 1.0 {
		log.Print("bufferPercent:", b.currentPos, b.sizeAhead(), b.size)
	}
	return res
}

func (b *buffer) Clear() {
	for _, bk := range b.data {
		block.DefaultBlockPool.Put(bk)
	}
	b.data = b.data[0:0]
}
