package movie

import "testing"

func TestBlock(t *testing.T) {
	b := &block{100, make([]byte, 1000)}
	if !b.inside(100) {
		t.Error("should inside")
	}
	if b.inside(1100) {
		t.Error("should not inside")
	}
	if b.inside(1101) {
		t.Error("should not inside")
	}
	if b.inside(99) {
		t.Error("shoud not inside")
	}
}

func TestBufferWrite(t *testing.T) {
	b := NewBuffer(1000)
	b.WriteAtQuit(make([]byte, 100), 100, nil)
	if b.data.Len() != 1 {
		t.Errorf("list length should be 1 but %d", b.data.Len())
	}

	b.WriteAtQuit(make([]byte, 100), 200, nil)
	if b.data.Len() != 2 {
		t.Errorf("list length should be 1 but %d", b.data.Len())
	}
}

func TestBufferFromTo(t *testing.T) {
	b := NewBuffer(1000)
	b.WriteAtQuit(make([]byte, 100), 100, nil)
	if b.data.Len() != 1 {
		t.Errorf("list length should be 1 but %d", b.data.Len())
	}

	b.WriteAtQuit(make([]byte, 100), 200, nil)
	if b.data.Len() != 2 {
		t.Errorf("list length should be 1 but %d", b.data.Len())
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
		b.WriteAtQuit(make([]byte, 100), from, make(chan bool))
		from += 100
	}

	b.Seek(50, 0)
	println(b.currentPos, b.size)
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
		b.WriteAtQuit(make([]byte, 100), from, make(chan bool))
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
