package subtitle

import (
	"fmt"
	"math/rand"
	"os"
	"sort"
	"testing"
	"time"
	. "vger/player/shared"
	"vger/player/subtitle/srt"
)

type subItemsForTest struct {
	items []*SubItem
}

func newSubItemsForTest(items []*SubItem) *subItemsForTest {
	return &subItemsForTest{items}
}

func (si *subItemsForTest) get(t time.Duration) []*SubItem {
	ret := make([]*SubItem, 0)
	for _, item := range si.items {
		if item.Contains(t) {
			ret = append(ret, item)
		}
	}

	return ret
}
func parseTestFile(name string) []*SubItem {
	f, err := os.Open(fmt.Sprintf("srt/%s", name))
	if err != nil {
		panic(err)
	}

	items, err := srt.Parse(f, 384, 303)
	if err != nil {
		panic(err)
	}

	return items
}

func TestSubItemsNotFind(t *testing.T) {
	items := parseTestFile("b.srt")
	si := newSubItems(items)
	si2 := newSubItemsForTest(items)

	find := si.get(0)
	find2 := si2.get(0)
	if len(find) != len(find2) {
		t.Errorf("expect same len: %d <> %d", len(find), len(find2))
		return
	}
}

func TestSubItemsNotFindLarger(t *testing.T) {
	items := parseTestFile("b.srt")
	si := newSubItems(items)
	si2 := newSubItemsForTest(items)

	find := si.get(10 * time.Hour)
	find2 := si2.get(10 * time.Hour)
	if len(find) != len(find2) {
		t.Errorf("expect same len: %d <> %d", len(find), len(find2))
		return
	}
}

func TestSubItemsOverlap(t *testing.T) {
	items := parseTestFile("b.srt")
	si := newSubItems(items)
	si2 := newSubItemsForTest(items)

	find := si.get(17 * time.Second)
	find2 := si2.get(17 * time.Second)
	if len(find) != len(find2) {
		t.Errorf("expect same len: %d <> %d", len(find), len(find2))
		return
	}

	find = si.get(17 * time.Second)
	find2 = si2.get(17 * time.Second)
	if len(find) != len(find2) {
		t.Errorf("expect same len: %d <> %d", len(find), len(find2))
		return
	}
}

func TestSubItemsCompare(t *testing.T) {
	items := parseTestFile("b.srt")
	si := newSubItems(items)
	si2 := newSubItemsForTest(items)

	for i := 0; i < 10000; i++ {
		dur := time.Duration(rand.Int63n(int64(40 * time.Minute)))
		find := si.get(dur)
		find2 := si2.get(dur)

		if len(find) != len(find2) {
			t.Errorf("expect same len: %d <> %d", len(find), len(find2))
			return
		}

		// if len(find) > 0 {
		// 	println(dur.String())
		// }

		sort.Sort(SubItems(find))
		sort.Sort(SubItems(find2))
		for i, item := range find {
			item2 := find2[i]

			if item.String() != item2.String() {
				t.Errorf("expact same: %s <> %s at %s", item.String(), item2.String(), dur.String())
			}
		}
	}
}

func TestFindPosBeforeFirst(t *testing.T) {
	items := parseTestFile("b.srt")
	si := newSubItems(items)

	i, ok := si.findPos(0)
	if ok {
		t.Errorf("should not find: %d", i)
	}

	if i >= 0 {
		t.Errorf("should less than 0 but %d", i)
	}
}

func TestFindPosAfterLast(t *testing.T) {
	items := parseTestFile("b.srt")
	si := newSubItems(items)

	i, ok := si.findPos(10 * time.Hour)
	if ok {
		t.Errorf("should not find: %d", i)
	}

	if i != len(si.nooverlap)-1 {
		t.Errorf("expect %d but %d", len(si.nooverlap)-1, i)
	}
}
func TestFindPosNotFind(t *testing.T) {
	items := parseTestFile("b.srt")
	si := newSubItems(items)

	i, ok := si.findPos(50 * time.Second)
	if ok {
		t.Errorf("should not find: %d", i)
	}

	if i != 19 {
		t.Errorf("expect %d but %d", 19, i)
	}
}
func BenchmarkSubItems(b *testing.B) {
	b.StopTimer()

	items := parseTestFile("g.srt")
	si := newSubItems(items)

	rnd := make([]time.Duration, 0)
	for i := 0; i < b.N; i++ {
		rnd = append(rnd, time.Duration(rand.Int63n(int64(40*time.Minute))))
	}

	b.StartTimer()

	for i := 0; i < b.N; i++ {
		si.get(rnd[i])
	}
}

func BenchmarkSubItemsForTest(b *testing.B) {
	b.StopTimer()

	items := parseTestFile("g.srt")

	si := newSubItemsForTest(items)
	rnd := make([]time.Duration, 0)
	for i := 0; i < b.N; i++ {
		rnd = append(rnd, time.Duration(rand.Int63n(int64(40*time.Minute))))
	}

	b.StartTimer()

	for i := 0; i < b.N; i++ {
		si.get(rnd[i])
	}
}

func TestGetById(t *testing.T) {
	items := parseTestFile("g.srt")
	si := newSubItems(items)

	// println(si.nooverlap)
	// println(si.others)

	item := si.getById(1510)
	if item == nil {
		t.Error("should found")
	}
}
