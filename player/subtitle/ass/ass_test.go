package ass

import (
	"player/shared"
	// "io/ioutil"
	"log"
	"os"
	"strings"
	"testing"
)

func TestParse(t *testing.T) {
	f, _ := os.Open("a.ass")
	// bytes, _ := ioutil.ReadAll(f)
	subs, err := Parse(f, 100, 100)

	if err != nil {
		t.Error(err)
	}
	if len(subs) != 904 {
		t.Errorf("Expect 904 but %d", len(subs))
	}
}

func TestDropSvg(t *testing.T) {
	s := `
[Events]
Format: Layer, Start, End, Style, Name, MarginL, MarginR, MarginV, Effect, Text
Dialogue: 1,0:00:40.72,0:00:44.70,Default,NTP,0,0,0,,{\bord0\shad0\fscx11\fscy14\p1\c&HECB000&\pos(86,248)}m 325 0 l 355 0 l 332 47 l 354 47 l 399 136 l 369 136 l 328 53 l 288 136 l 257 136 l 325 0 m 467 0 l 497 0 l 474 47 l 496 47 l 541 136 l 511 136 l 470 53 l 430 136 l 399 136 l 467 0 m 545 1 l 583 1 l 583 14 l 568 14 l 568 19 l 583 19 l 583 30 l 568 30 l 568 34 l 599 34 l 599 30 l 583 30 l 583 19 l 599 19 l 599 14 l 583 14 l 583 1 l 611 1 b 616 1 622 6 622 10 l 622 36 l 652 0 l 678 0 l 644 41 l 622 41 l 622 47 l 596 47 l 597 54 l 625 54 l 625 68 l 541 68 l 541 54 l 572 54 l 571 47 l 545 47 l 545 1 m 583 72 l 583 85 l 569 85 l 569 90 l 598 90 l 598 85 l 583 85 l 583 72 l 611 72 b 615 72 621 78 621 82 l 653 44 l 678 44 l 644 86 l 621 86 l 621 103 l 597 103 l 597 136 l 570 136 l 564 126 l 562 136 l 542 136 l 548 107 l 568 107 l 565 121 l 571 121 l 571 103 l 547 103 l 547 72 l 583 72 m 600 107 l 620 107 l 624 124 l 653 89 l 679 89 l 642 136 l 615 136 l 618 132 l 606 132 l 600 107 m 689 0 l 716 0 l 721 15 l 732 15 l 732 30 l 718 56 l 731 56 l 735 100 l 721 100 l 717 59 l 714 64 l 714 136 l 693 136 l 693 79 l 676 79 l 707 30 l 679 30 l 679 15 l 694 15 l 689 0 m 738 0 l 804 0 b 807 0 813 6 813 9 l 813 87 l 794 87 l 794 14 l 756 14 l 756 87 l 763 77 l 763 21 l 787 21 l 787 91 l 798 91 l 798 120 l 820 120 l 812 136 l 778 136 l 778 90 l 748 136 l 723 136 l 756 87 l 738 87 l 738 0 m 257 151 l 275 151 l 297 182 l 319 151 l 337 151 l 304 197 l 304 227 l 290 227 l 290 197 l 257 151 m 337 151 l 355 151 l 377 182 l 399 151 l 417 151 l 384 197 l 384 227 l 370 227 l 370 197 l 337 151 m 425 192 l 447 192 b 445 181 427 181 425 192 l 410 192 b 414 162 458 162 462 192 l 462 206 l 425 206 b 425 210 429 212 433 213 l 462 213 l 462 227 l 433 227 b 412 227 408 203 410 192 m 462 151 l 532 151 l 532 165 l 504 165 l 504 227 l 489 227 l 489 165 l 462 165 l 462 151 m 580 172 l 580 186 l 549 186 b 543 186 543 192 549 192 l 565 192 b 589 192 589 227 565 227 l 532 227 l 532 213 l 565 213 b 570 213 570 206 565 206 l 549 206 b 524 206 524 172 549 172 l 580 172 m 592 213 l 606 213 l 606 227 l 592 227 l 592 213 m 639 172 l 665 172 l 665 186 l 639 186 b 623 186 623 213 639 213 l 665 213 l 665 227 l 639 227 b 603 227 603 172 639 172 m 700 184 b 722 184 722 215 700 215 l 700 229 b 740 229 740 170 700 170 b 660 170 660 229 700 229 l 700 215 b 680 215 680 184 700 184 m 737 172 l 782 172 b 803 172 813 177 813 198 l 813 228 l 799 228 l 799 198 b 799 186 793 187 782 187 l 782 228 l 768 228 l 768 187 l 752 187 l 752 228 l 737 228 l 737 172{\p0}`
	items, _ := Parse(strings.NewReader(s), 100, 100)
	println(items)
	if items[0].Content[0].Content != "" {
		t.Errorf("Expect empty but %s", items[0].Content[0].Content)
	}
}

func TestParseB(t *testing.T) {
	f, _ := os.Open("b.ass")
	// bytes, _ := ioutil.ReadAll(f)
	subs, err := Parse(f, 100, 100)

	if err != nil {
		t.Error(err)
	}
	// if len(subs) != 904 {
	// t.Errorf("Expect 904 but %d", len(subs))
	// }
	// printSubs(subs)
	println(subs)
}

func BenchmarkParseB(b *testing.B) {
	b.StopTimer()
	f, _ := os.Open("b.ass")
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		_, err := Parse(f, 100, 100)
		if err != nil {
			b.Error(err)
			break
		}
	}
}

func printSubs(subs []*shared.SubItem) {
	for _, s := range subs {
		log.Printf("%v", s)
	}
}
