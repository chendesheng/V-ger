package subtitle

type Subtitles []*Subtitle

func (s Subtitles) Len() int {
	return len([]*Subtitle(s))
}

func (s Subtitles) Less(i, j int) bool {
	a := s[i]
	b := s[j]

	compares := []int{
		compareLang(a.Lang1, a.Lang2, b.Lang1, b.Lang2),
		compareDistance(a.Distance, b.Distance),
		compareType(a.Type, b.Type),
	}

	for _, c := range compares {
		if c != 0 {
			return c > 0
		}
	}

	return false
}

func (s Subtitles) Swap(i, j int) {
	t := s[i]
	s[i] = s[j]
	s[j] = t
}

func (s Subtitles) Select() (a *Subtitle, b *Subtitle) {
	subs := ([]*Subtitle)(s)
	if len(subs) == 1 || subs[0].IsTwoLangs() {
		a = subs[0]
		b = nil
	} else {
		if subs[0].Lang1 == "en" {
			a = subs[0]
			b = nil
		} else {
			a = subs[0]
			for _, c := range subs {
				if c.Lang1 == "en" {
					b = c
					break
				}
			}
		}
	}

	return
}

func compareLang(a1, a2, b1, b2 string) int {
	if a1 == b1 && a2 == b2 {
		return 0
	}
	//multi lang > signle lang
	if len(a2) > 0 && len(b2) == 0 {
		return 1
	}
	if len(a2) == 0 && len(b2) > 0 {
		return -1
	}
	//cn > en
	if len(a2) == 0 && len(b2) == 0 {
		if a1 == "chs" {
			return 1
		}
		if b1 == "chs" {
			return -1
		}
		if a1 == "cht" {
			return 1
		}
		if b1 == "cht" {
			return -1
		}
		return 1
	}

	if a2 == "chs" {
		return 1
	}
	if b2 == "chs" {
		return -1
	}
	if a2 == "cht" {
		return 1
	}
	if b2 == "cht" {
		return -1
	}
	return 1
}

func compareType(a, b string) int {
	if a == b {
		return 0
	} else if a == "srt" {
		return 1
	} else {
		return -1
	}
}

func compareDistance(a, b int) int {
	switch {
	case a == b:
		return 0
	case a < b:
		return 1
	case a > b:
		return -1
	default:
		return 0
	}
}
