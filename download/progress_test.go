package download

import (
	// "fmt"
	"testing"
)

func TestSegRing(t *testing.T) {
	sr := newSegRing(2)

	if len(sr.segs) != 2 {
		t.Errorf("segs len expected 2, but %d", len(sr.segs))
		return
	}

	sr.add(10)
	// fmt.Printf("%v", sr)

	if _, l := sr.total(); l != 10 {
		t.Errorf("total length expected 10, but %d", l)
		return
	}
	if sr.i != 1 {
		t.Errorf("index expected 1, but %d", sr.i)
		return
	}

	sr.add(10)
	if _, l := sr.total(); l != 20 {
		t.Errorf("total length expected 20, but %d", l)
		return
	}
	if sr.i != 0 {
		t.Errorf("index expected 1, but %d", sr.i)
		return
	}

	sr.add(20)
	if _, l := sr.total(); l != 30 {
		t.Errorf("total length expected 30, but %d", l)
		return
	}
	if sr.i != 1 {
		t.Errorf("index expected 1, but %d", sr.i)
		return
	}
}
