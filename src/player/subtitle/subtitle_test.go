package subtitle

import (
	"sort"
	"testing"
)

func TestCompareLang(t *testing.T) {
	if compareLang("chs", "", "cht", "") != 1 {
		t.Error("should >")
	}

	if compareLang("chs", "", "en", "") != 1 {
		t.Error("should >")
	}

	if compareLang("cht", "", "en", "") != 1 {
		t.Error("should >")
	}
}

func TestCompareLang2(t *testing.T) {
	if compareLang("en", "chs", "chs", "") != 1 {
		t.Error("should >")
	}
	if compareLang("en", "chs", "cht", "") != 1 {
		t.Error("should >")
	}
	if compareLang("en", "chs", "en", "") != 1 {
		t.Error("should >")
	}
}

func TestCompareLang3(t *testing.T) {
	if compareLang("en", "chs", "en", "cht") != 1 {
		t.Error("should >")
	}
}

func TestCompareLangEqual(t *testing.T) {
	if compareLang("en", "", "en", "") != 0 {
		t.Error("should =")
	}
	if compareLang("en", "chs", "en", "chs") != 0 {
		t.Error("should =")
	}
}

func TestCompareFormat(t *testing.T) {
	if compareFormat("ass", "srt") != -1 {
		t.Error("should <")
	}
	if compareFormat("ass", "ass") != 0 {
		t.Error("should =")
	}
	if compareFormat("srt", "srt") != 0 {
		t.Error("should =")
	}
}

func TestSortSubtitle(t *testing.T) {
	subs := make([]*Subtitle, 0, 10)
	subs = append(subs, &Subtitle{
		Name:   "1",
		Lang1:  "en",
		Lang2:  "",
		Format: "ass",
	})
	subs = append(subs, &Subtitle{
		Name:   "2",
		Lang1:  "chs",
		Lang2:  "",
		Format: "ass",
	})
	subs = append(subs, &Subtitle{
		Name:   "3",
		Lang1:  "cht",
		Lang2:  "",
		Format: "srt",
	})
	subs = append(subs, &Subtitle{
		Name:   "4",
		Lang1:  "en",
		Lang2:  "cht",
		Format: "ass",
	})
	subs = append(subs, &Subtitle{
		Name:   "5",
		Lang1:  "en",
		Lang2:  "cht",
		Format: "srt",
	})
	subs = append(subs, &Subtitle{
		Name:   "6",
		Lang1:  "en",
		Lang2:  "chs",
		Format: "ass",
	})

	sort.Sort(Subtitles(subs))
	var order string
	for _, s := range subs {
		order += s.Name
	}
	if order != "654231" {
		t.Errorf("Expect 654231 but %s", order)
	}
}

func TestSelectSorted(t *testing.T) {
	subs := make([]*Subtitle, 0, 10)
	subs = append(subs, &Subtitle{
		Name:   "4",
		Lang1:  "en",
		Lang2:  "cht",
		Format: "ass",
	})
	subs = append(subs, &Subtitle{
		Name:   "5",
		Lang1:  "en",
		Lang2:  "cht",
		Format: "srt",
	})
	subs = append(subs, &Subtitle{
		Name:   "6",
		Lang1:  "en",
		Lang2:  "chs",
		Format: "ass",
	})

	sort.Sort(Subtitles(subs))
	a, b := Subtitles(subs).Select()
	if a == nil || b != nil || a.Name != "6" {
		t.Errorf("Expect 6 but %s", a.Name)
	}
}
func TestSelectSorted2(t *testing.T) {
	subs := make([]*Subtitle, 0, 10)
	subs = append(subs, &Subtitle{
		Name:   "1",
		Lang1:  "en",
		Lang2:  "",
		Format: "ass",
	})
	subs = append(subs, &Subtitle{
		Name:   "2",
		Lang1:  "chs",
		Lang2:  "",
		Format: "ass",
	})
	subs = append(subs, &Subtitle{
		Name:   "3",
		Lang1:  "cht",
		Lang2:  "",
		Format: "srt",
	})

	sort.Sort(Subtitles(subs))
	a, b := Subtitles(subs).Select()
	if a == nil || b == nil || a.Name != "2" || b.Name != "1" {
		t.Errorf("Expect 2,3 but %s,%s", a.Name, b.Name)
	}
}
