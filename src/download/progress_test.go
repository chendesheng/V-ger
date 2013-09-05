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

	if sr.totalSize() != 10 {
		t.Errorf("total size expected 10, but %d", sr.totalSize())
		return
	}

	if sr.i != 1 {
		t.Errorf("total size expected 1, but %d", sr.i)
		return
	}
}
