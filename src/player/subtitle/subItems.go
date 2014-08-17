package subtitle

import (
	. "player/shared"
	"time"
)

type subItems struct {
	nooverlap []*SubItem
	others    []*SubItem
}

const GAP = 1000

func newSubItems(items []*SubItem) *subItems {
	if len(items) <= 1 {
		return &subItems{items, nil}
	}

	nooverlap := make([]*SubItem, 0, len(items))
	nooverlap = append(nooverlap, items[0])

	others := make([]*SubItem, 0)

	for i, item := range items[1:] {
		if item.IsInDefaultPosition() && items[i].To < item.From {
			nooverlap = append(nooverlap, item)
		} else {
			others = append(others, item)
		}
	}

	for i, item := range nooverlap {
		item.Id = i
	}

	l := len(nooverlap)
	for i, item := range others {
		item.Id = l + GAP + i
	}

	return &subItems{nooverlap, others}
}

func (si *subItems) get(t time.Duration) []*SubItem {
	ret := make([]*SubItem, 0)

	if i, ok := si.findPos(t); ok {
		ret = append(ret, si.nooverlap[i])
	}

	for _, item := range si.others {
		if item.Contains(t) {
			ret = append(ret, item)
		}
	}

	return ret
}

func (si *subItems) findPos(t time.Duration) (int, bool) {
	nooverlap := si.nooverlap
	p := 0
	q := len(nooverlap) - 1

	for p <= q {

		i := (p + q) / 2
		current := nooverlap[i]

		if current.Contains(t) {
			return i, true
		} else if t < current.From {
			q = i - 1
		} else {
			p = i + 1
		}
	}

	return q, false
}

func (si *subItems) getById(id int) *SubItem {
	l := len(si.nooverlap)
	l1 := len(si.others)

	switch {
	case id < 0:
		return nil
	case id < l:
		return si.nooverlap[id]
	case id < l+GAP:
		return nil
	case id < l+l1:
		return si.others[id-l1]
	default:
		return nil
	}
}

func (si *subItems) each(fn func(*SubItem)) {
	for _, item := range si.nooverlap {
		fn(item)
	}

	for _, item := range si.others {
		fn(item)
	}
}
