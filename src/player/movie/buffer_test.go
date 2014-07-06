package movie

import (
	"block"
	"testing"
)

func TestBufferWrite(t *testing.T) {
	b := NewBuffer(1000)
	b.WriteAtQuit(block.Block{100, make([]byte, 100)}, nil)
	if len(b.data) != 1 {
		t.Errorf("list length should be 1 but %d", len(b.data))
	}

	b.WriteAtQuit(block.Block{200, make([]byte, 100)}, nil)
	if len(b.data) != 2 {
		t.Errorf("list length should be 1 but %d", len(b.data))
	}
}

func TestBufferFromTo(t *testing.T) {
	b := NewBuffer(1000)
	b.WriteAtQuit(block.Block{100, make([]byte, 100)}, nil)
	if len(b.data) != 1 {
		t.Errorf("list length should be 1 but %d", len(b.data))
	}

	b.WriteAtQuit(block.Block{200, make([]byte, 100)}, nil)
	if len(b.data) != 2 {
		t.Errorf("list length should be 2 but %d", len(b.data))
	}

	from, to := b.fromTo()
	if from != 100 {
		t.Errorf("from should be 100 but %d", from)
	}
	if to != 300 {
		t.Errorf("to should be 300 but %d", to)
	}
}

type testWriter struct {
	result []int
}

func (w *testWriter) Write(p []byte) (int, error) {
	println("write length:", len(p))
	w.result = append(w.result, len(p))
	return len(p), nil
}
func TestBufferRead(t *testing.T) {
	b := NewBuffer(1000)
	from := int64(0)
	for from < b.size {
		b.WriteAtQuit(block.Block{from, make([]byte, 100)}, make(chan struct{}))
		from += 100
	}

	b.Seek(50, 0)
	println(b.currentPos, b.size, len(b.data))
	w := &testWriter{}
	b.Read(w, 200)
	if len(w.result) != 3 {
		t.Errorf("write count should be 3 but %d", len(w.result))
	}
	if w.result[0] != 50 {
		t.Errorf("result[0] should be 50 but %d", w.result[0])
	}
	if w.result[1] != 100 {
		t.Errorf("result[1] should be 100 but %d", w.result[1])
	}
	if w.result[2] != 50 {
		t.Errorf("result[2] should be 50 but %d", w.result[2])
	}

	b.Seek(0, 0)
	w.result = nil
	b.Read(w, 200)
	if len(w.result) != 2 {
		t.Errorf("write count should be 3 but %d", len(w.result))
	}

	if w.result[0] != 100 {
		t.Errorf("result[0] should be 100 but %d", w.result[0])
	}
	if w.result[1] != 100 {
		t.Errorf("result[1] should be 100 but %d", w.result[1])
	}
}

func TestBufferReadBorder(t *testing.T) {
	b := NewBuffer(1000)
	from := int64(0)
	for from < b.size {
		b.WriteAtQuit(block.Block{from, make([]byte, 100)}, make(chan struct{}))
		from += 100
	}

	b.Seek(50, 0)

	w := &testWriter{}
	got := b.Read(w, 1000)
	if got != 950 {
		t.Errorf("expect got 950 but %d", got)
	}

	if b.currentPos != 1000 {
		t.Error("expect currentPost 1000 but %d", b.currentPos)
	}
}

func TestGC(t *testing.T) {
	b := NewBuffer(1000)
	from := int64(0)
	i := int64(1)
	for from < b.size {
		b.WriteAtQuit(block.Block{from, make([]byte, 50*i)}, make(chan struct{}))
		from += 50 * i //50+100+150+200+250+300(250)
		i++
	}
	b.currentPos = 456
	b.GC()

	if len(b.data) != 3 {
		t.Errorf("expect len(b.data)=3 but %d", len(b.data))
	}

	if len(b.data[0].Data) != 200 {
		t.Errorf("expect len(b.data[0].p)=200 but %d", len(b.data[0].Data))
	}
}

func TestLastPos(t *testing.T) {
	b := NewBuffer(1000)
	if b.LastPos() != 0 {
		t.Errorf("%d != 0", b.LastPos())
	}

	b.WriteAtQuit(block.Block{10, make([]byte, 10)}, nil)
	if b.LastPos() != 20 {
		t.Errorf("%d != 20", b.LastPos())
	}
}
