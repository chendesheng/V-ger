package subtitles

import (
	"testing"
)

func TestAddic7edGet(t *testing.T) {
	name, content := Addic7edSubtitle("rake s01e11")
	println(name, content)
	if content == "" {
		t.Error("should not be empty")
	}
}

func TestAddic7edNotExists(t *testing.T) {
	_, content := Addic7edSubtitle("rake s01e1122")
	if content != "" {
		t.Error("should be empty")
	}
}
