package block

import (
	"testing"
)

func TestBlock(t *testing.T) {
	b := &Block{100, make([]byte, 1000)}
	if !b.Inside(100) {
		t.Error("should inside")
	}
	if !b.Inside(101) {
		t.Error("should inside")
	}
	if !b.Inside(1099) {
		t.Error("should inside")
	}
	if b.Inside(1100) {
		t.Error("should not inside")
	}
	if b.Inside(1101) {
		t.Error("should not inside")
	}
	if b.Inside(99) {
		t.Error("shoud not inside")
	}
}
