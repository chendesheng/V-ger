package audio

import (
	"testing"
	"time"
)

func TestAppend(t *testing.T) {
	sb := &sampleBuffer{}
	sb.bytesPerSec = 1

	sb.append(&samples{make([]byte, 100), 0})
	if len(sb.buf) != 1 {
		t.Errorf("buf len expect 1 but %d", len(sb.buf))
	}

	sb.append(&samples{make([]byte, 200), 100})
	if len(sb.buf) != 2 {
		t.Errorf("buf len expect 2 but %d", len(sb.buf))
	}
}

func TestCut(t *testing.T) {
	sb := &sampleBuffer{}
	sb.bytesPerSec = 1

	sb.append(&samples{make([]byte, 100), 0})
	sb.append(&samples{make([]byte, 200), 100})

	if len(sb.buf) != 2 {
		t.Errorf("buf len expect 2 but %d", len(sb.buf))
	}

	pts := sb.pts()
	bytes := sb.cut(10)
	if len(bytes) != 10 {
		t.Errorf("len expect 10 but %d", len(bytes))
	}
	if pts != 0 {
		t.Errorf("pts expect 0 but %d", pts)
	}
}

func TestCut2(t *testing.T) {
	sb := &sampleBuffer{}
	sb.bytesPerSec = 1

	sb.append(&samples{make([]byte, 100), 0})
	sb.append(&samples{make([]byte, 200), 100 * time.Second})

	pts := sb.pts()
	bytes := sb.cut(104)
	if len(bytes) != 104 {
		t.Errorf("len expect 104 but %d", len(bytes))
	}
	if pts != 0 {
		t.Errorf("pts expect 0 but %d", pts)
	}

	if len(sb.buf) != 1 {
		t.Errorf("buf len expect 1 but %d", len(sb.buf))
	}

	if len(sb.buf[0].data) != 196 {
		t.Errorf("samples len expect 196 but %d", len(sb.buf[0].data))
	}
	if sb.buf[0].pts != 104*time.Second {
		t.Errorf("samples len expect 104s but %s", sb.buf[0].pts.String())
	}
}
func TestCutAll(t *testing.T) {
	sb := &sampleBuffer{}
	sb.bytesPerSec = 1

	sb.append(&samples{make([]byte, 100), 0})
	sb.append(&samples{make([]byte, 200), 100 * time.Second})

	pts := sb.pts()
	bytes := sb.cut(500)
	if len(bytes) != 300 {
		t.Errorf("len expect 300 but %d", len(bytes))
	}
	if pts != 0 {
		t.Errorf("pts expect 0 but %d", pts)
	}

	if !sb.empty() {
		t.Errorf("expect empty")
	}
}

func TestCutByTime(t *testing.T) {
	sb := &sampleBuffer{}
	sb.bytesPerSec = 1

	sb.append(&samples{make([]byte, 100), 0})
	sb.append(&samples{make([]byte, 200), 100 * time.Second})

	sb.cutByTime(299 * time.Second)

	if len(sb.buf[0].data) != 4 {
		t.Errorf("len expect 4 but %d", len(sb.buf[0].data))
	}
}
