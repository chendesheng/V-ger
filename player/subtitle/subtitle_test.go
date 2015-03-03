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
	if compareType("ass", "srt") != -1 {
		t.Error("should <")
	}
	if compareType("ass", "ass") != 0 {
		t.Error("should =")
	}
	if compareType("srt", "srt") != 0 {
		t.Error("should =")
	}
}

func TestSortSubtitle(t *testing.T) {
	subs := make([]*Subtitle, 0, 10)
	subs = append(subs, &Subtitle{
		Name:  "1",
		Lang1: "en",
		Lang2: "",
		Type:  "ass",
	})
	subs = append(subs, &Subtitle{
		Name:  "2",
		Lang1: "chs",
		Lang2: "",
		Type:  "ass",
	})
	subs = append(subs, &Subtitle{
		Name:  "3",
		Lang1: "cht",
		Lang2: "",
		Type:  "srt",
	})
	subs = append(subs, &Subtitle{
		Name:  "4",
		Lang1: "en",
		Lang2: "cht",
		Type:  "ass",
	})
	subs = append(subs, &Subtitle{
		Name:  "5",
		Lang1: "en",
		Lang2: "cht",
		Type:  "srt",
	})
	subs = append(subs, &Subtitle{
		Name:  "6",
		Lang1: "en",
		Lang2: "chs",
		Type:  "ass",
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
		Name:  "4",
		Lang1: "en",
		Lang2: "cht",
		Type:  "ass",
	})
	subs = append(subs, &Subtitle{
		Name:  "5",
		Lang1: "en",
		Lang2: "cht",
		Type:  "srt",
	})
	subs = append(subs, &Subtitle{
		Name:  "6",
		Lang1: "en",
		Lang2: "chs",
		Type:  "ass",
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
		Name:  "1",
		Lang1: "en",
		Lang2: "",
		Type:  "ass",
	})
	subs = append(subs, &Subtitle{
		Name:  "2",
		Lang1: "chs",
		Lang2: "",
		Type:  "ass",
	})
	subs = append(subs, &Subtitle{
		Name:  "3",
		Lang1: "cht",
		Lang2: "",
		Type:  "srt",
	})

	sort.Sort(Subtitles(subs))
	a, b := Subtitles(subs).Select()
	if a == nil || b == nil || a.Name != "2" || b.Name != "1" {
		t.Errorf("Expect 2,3 but %s,%s", a.Name, b.Name)
	}
}

func TestSort(t *testing.T) {
	subs := make([]*Subtitle, 0, 10)

	subs = append(subs, &Subtitle{
		Name:     "house.of.cards.2013.s03e10.720p.webrip.x264",
		Lang1:    "cn",
		Lang2:    "en",
		Type:     "ass",
		Distance: 11,
	})

	subs = append(subs, &Subtitle{
		Name:     "House.of.Cards.S03E10.720p.WEBRip.DD5.1.x264",
		Lang1:    "cn",
		Lang2:    "en",
		Type:     "srt",
		Distance: 0,
	})

	sort.Sort(Subtitles(subs))

	if subs[0].Distance != 0 {
		t.Errorf("first sub's distance: expect %d but %d", 0, subs[0].Distance)
	}
}
