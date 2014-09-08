package util

import (
	"testing"
)

type task struct {
	Name string
	URL  string
}

func TestRW(t *testing.T) {
	err := WriteJson("a.json", task{"name", "url"})
	if err != nil {
		t.Error(err)
	}

	tk := task{}
	err = ReadJson("a.json", &tk)
	if err != nil {
		t.Error(err)
	}

	if tk.Name != "name" || tk.URL != "url" {
		t.Error("Write and read not match")
	}

	err = WriteJson("a.json", task{"name1", "url1"})
	if err != nil {
		t.Error(err)
	}

	tk = task{}
	err = ReadJson("a.json", &tk)
	if err != nil {
		t.Error(err)
	}

	if tk.Name != "name1" || tk.URL != "url1" {
		t.Error("Write and read not match")
	}
}

func TestRError(t *testing.T) {
	tk := task{}
	err := ReadJson("x", tk)
	if err == nil {
		t.Error("should has error")
	}
}

func TestWError(t *testing.T) {
	tk := task{}
	err := ReadJson("/x/a.a", tk)
	if err == nil {
		t.Error("should has error")
	}
}
