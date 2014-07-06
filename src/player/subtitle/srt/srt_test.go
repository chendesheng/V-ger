package srt

import (
	"player/shared"
	// "fmt"
	// "io/ioutil"
	"os"
	"strings"
	"testing"
)

// func TestParse(t *testing.T) {
// 	data, err := ioutil.ReadFile("a.srt")
// 	if err != nil {
// 		t.Error(err)
// 		return
// 	}

// 	text := string(data)
// 	items, _ := Parse(strings.NewReader(text), 384,303)
// 	if len(items) == 0 {
// 		t.Errorf("parse error")
// 	}
// }
func TestRemovePositionInfo(t *testing.T) {
	res := removePositionInfo("{\\pos(aaa)}{\\an1}bbb{\\test}")
	if res != "bbb" {
		t.Errorf("Expect 'bbb' but %s", res)
	}

}
func TestParse2(t *testing.T) {
	f, err := os.Open("b.srt")
	if err != nil {
		t.Error(err)
		return
	}

	items, _ := Parse(f, 384, 303)

	// fmt.Printf("%v", items)
	if len(items) != 723 {
		t.Errorf("Expect 723 items but %d items.", len(items))
	}
}
func TestParseItalic(t *testing.T) {
	text := `1
00:00:00,000 --> 00:00:01,111
a<i>你好</i>c`
	items, _ := Parse(strings.NewReader(text), 384, 303)
	if len(items) != 1 {
		t.Errorf("Expect 1 item but %d item(s).", len(items))
		return
	}
	// println(items[0].Content)

	if items[0].Content[0].Content != "a" {
		t.Errorf("first content should be 'a' but '%s'", items[0].Content[0].Content)
		return
	}
	if items[0].Content[1].Content != "你好" {
		t.Error("second content should be '你好' but ", items[0].Content[1].Content)
		return
	}
	if items[0].Content[1].Style != 1 {
		t.Error("second style should be '1' but ", items[0].Content[1].Style)
		return
	}
	if items[0].Content[2].Content != "c" {
		t.Error("third content should be 'c' but ", items[0].Content[2].Content)
		return
	}
}
func TestParseItalic2(t *testing.T) {
	text := `1
00:00:00,000 --> 00:00:01,111
cc
a<i>b</i>c`
	items, _ := Parse(strings.NewReader(text), 384, 303)
	if len(items) != 1 {
		t.Error("should be only 1 item")
		return
	}

	if items[0].Content[0].Content != `cc
a` {
		t.Errorf("first content should be 'cc\na' but '%s'", items[0].Content[0].Content)
		return
	}
	if items[0].Content[1].Content != "b" {
		t.Error("second content should be 'b' but ", items[0].Content[1].Content)
		return
	}
	if items[0].Content[1].Style != 1 {
		t.Error("second style should be '1' but ", items[0].Content[1].Style)
		return
	}
	if items[0].Content[2].Content != "c" {
		t.Error("third content should be 'c' but ", items[0].Content[2].Content)
		return
	}
}
func TestParseBold(t *testing.T) {
	text := `1
00:00:00,000 --> 00:00:01,111
cc
a<b>b</b>c`
	items, _ := Parse(strings.NewReader(text), 384, 303)
	if len(items) != 1 {
		t.Error("should be only 1 item")
		return
	}

	if items[0].Content[0].Content != `cc
a` {
		t.Errorf("first content should be 'cc\na' but '%s'", items[0].Content[0].Content)
		return
	}
	if items[0].Content[1].Content != "b" {
		t.Error("second content should be 'b' but ", items[0].Content[1].Content)
		return
	}
	if items[0].Content[1].Style != 2 {
		t.Error("second style should be '2' but ", items[0].Content[1].Style)
		return
	}
	if items[0].Content[2].Content != "c" {
		t.Error("third content should be 'c' but ", items[0].Content[2].Content)
		return
	}
}

func TestParseMixed(t *testing.T) {
	text := `1
00:00:00,000 --> 00:00:01,111
cc
a<b>b</b>c<i>d</i><b>e</b>`
	items, _ := Parse(strings.NewReader(text), 384, 303)
	if len(items) != 1 {
		t.Error("should be only 1 item")
		return
	}

	if items[0].Content[0].Content != `cc
a` {
		t.Errorf("first content should be 'cc\na' but '%s'", items[0].Content[0].Content)
		return
	}
	if items[0].Content[1].Content != "b" {
		t.Error("second content should be 'b' but ", items[0].Content[1].Content)
		return
	}
	if items[0].Content[1].Style != 2 {
		t.Error("second style should be '2' but ", items[0].Content[1].Style)
		return
	}
	if items[0].Content[2].Content != "c" {
		t.Error("third content should be 'c' but ", items[0].Content[2].Content)
		return
	}
	if items[0].Content[2].Style != 0 {
		t.Error("third style should be 0 but ", items[0].Content[2].Style)
		return
	}

	if items[0].Content[3].Content != "d" {
		t.Errorf("forth content should be 'd' but '%s'", items[0].Content[3].Content)
		return
	}
	if items[0].Content[3].Style != 1 {
		t.Errorf("forth style should be 1 but %d", items[0].Content[3].Style)
		return
	}

	if items[0].Content[4].Content != "e" {
		t.Errorf("forth content should be 'e' but '%s'", items[0].Content[4].Content)
		return
	}
	if items[0].Content[4].Style != 2 {
		t.Errorf("forth style should be 2 but %d", items[0].Content[4].Style)
		return
	}
}

func TestParseNested(t *testing.T) {
	text := `1
00:00:00,000 --> 00:00:01,111
<b>a<i>b</i></b>`
	items, _ := Parse(strings.NewReader(text), 384, 303)
	if len(items) != 1 {
		t.Error("should be only 1 item")
		return
	}
	if items[0].Content[0].Content != "a" {
		t.Error("first content should be 'a' but ", items[0].Content[0].Content)
		return
	}
	if items[0].Content[0].Style != 2 {
		t.Error("first style should be 2 but ", items[0].Content[0].Style)
		return
	}
	if items[0].Content[1].Content != "b" {
		t.Error("second content should be 'b' but ", items[0].Content[1].Content)
		return
	}
	if items[0].Content[1].Style != 3 {
		t.Error("second style should be 3 but ", items[0].Content[1].Style)
		return
	}
}

func TestParseFont(t *testing.T) {
	text := `1
00:00:00,000 --> 00:00:01,111
<font color="white">a</font>`
	items, _ := Parse(strings.NewReader(text), 384, 303)
	if len(items) != 1 {
		t.Error("should be only 1 item")
		return
	}
	if items[0].Content[0].Content != "a" {
		t.Error("first content should be 'a' but ", items[0].Content[0].Content)
		return
	}
	if items[0].Content[0].Color != 0xffffff {
		t.Errorf("first color should be 0xffffff but 0x%x", items[0].Content[0].Color)
		return
	}
}
func TestParseFontColor(t *testing.T) {
	text := `1
00:00:00,000 --> 00:00:01,111
<font color="#00ffee">a</font>`
	items, _ := Parse(strings.NewReader(text), 384, 303)
	if len(items) != 1 {
		t.Error("should be only 1 item")
		return
	}
	if items[0].Content[0].Color != 0x00ffee {
		t.Errorf("first color should be 0x00ffee but 0x%x", items[0].Content[0].Color)
		return
	}
}

func TestParseFull(t *testing.T) {
	text := `1
00:00:00,000 --> 00:00:01,111
<i><b><Font Color="white">a</font></b></i>b<b><font color="Black">c</font></b><i>d<i>`
	items, _ := Parse(strings.NewReader(text), 384, 303)
	if len(items) != 1 {
		t.Error("should be only 1 item")
		return
	}
	if items[0].Content[0].Content != "a" {
		t.Error("first content should be 'a' but ", items[0].Content[0].Content)
		return
	}
	if items[0].Content[0].Color != 0xffffff {
		t.Errorf("first color should be 0xffffff but 0x%x", items[0].Content[0].Color)
		return
	}

	if items[0].Content[1].Content != "b" {
		t.Error("second content should be 'b' but ", items[0].Content[1].Content)
		return
	}
	if items[0].Content[1].Color != 0xffffff {
		t.Errorf("second color should be 0xffffff but 0x%x", items[0].Content[1].Color)
		return
	}
	if items[0].Content[1].Style != 0 {
		t.Errorf("second style should be 0 but %d", items[0].Content[1].Style)
		return
	}

	if items[0].Content[2].Content != "c" {
		t.Error("third content should be 'c' but ", items[0].Content[2].Content)
		return
	}
	if items[0].Content[2].Color != 0 {
		t.Errorf("third color should be 0 but 0x%x", items[0].Content[2].Color)
		return
	}

	if items[0].Content[3].Content != "d" {
		t.Error("forth content should be 'd' but ", items[0].Content[3].Content)
		return
	}
	if items[0].Content[3].Style != 1 {
		t.Errorf("forth style should be 1 but %d", items[0].Content[3].Style)
		return
	}
}
func TextBreaklne(t *testing.T) {
	text := `1
00:00:00,001 --> 00:00:04,336
President reagan:
<i>Air and naval forces</i>
<i>of the United States</i>`
	items, _ := Parse(strings.NewReader(text), 384, 303)

	if items[0].Content[0].Content != `President reagan:
` {
		t.Error("first content should be 'resident reagan\n' but ", items[0].Content[0].Content)
		return
	}
	if items[0].Content[2].Content != "\n" {
		t.Error("third content should be '\\n' but ", items[0].Content[2].Content)
		return
	}
}

func floatEqual(a, b float64) bool {
	d := a - b
	return d < 1e5 && d > -1e5
}

func TestParsePosition(t *testing.T) {
	text := `1
00:00:00,001 --> 00:00:04,336
{\pos(1.12,20.4)}President reagan:`
	items, _ := Parse(strings.NewReader(text), 384, 303)
	if items[0].Content[0].Content != "President reagan:" {
		t.Error("first content should be 'President reagan:' but", items[0].Content[0].Content)
	}
	if !(items[0].PositionType == 2 && floatEqual(items[0].X, 1.0) && floatEqual(items[0].Y, 20.0)) {
		t.Errorf("parse position faild except (2, 1,20) but (%d, %f, %f)", items[0].PositionType, items[0].X, items[0].Y)
	}
}

func TestParsePosition2(t *testing.T) {
	text := `1
00:00:00,001 --> 00:00:04,336
{\an8}President reagan:`
	items, _ := Parse(strings.NewReader(text), 384, 303)
	if items[0].Content[0].Content != "President reagan:" {
		t.Error("first content should be 'President reagan:' but", items[0].Content[0].Content)
	}
	if !(items[0].PositionType == 8 && items[0].X == -1 && items[0].Y == -1) {
		t.Errorf("parse position faild except (8,-1,-1) but (%d, %f, %f)", items[0].PositionType, items[0].X, items[0].Y)
	}
}
func TestParsePosition3(t *testing.T) {
	text := `1
00:00:00,001 --> 00:00:04,336
{\an8}{\pos(1,2)}President reagan:`
	items, _ := Parse(strings.NewReader(text), 384, 303)
	if items[0].Content[0].Content != "President reagan:" {
		t.Error("first content should be 'President reagan:' but", items[0].Content[0].Content)
	}
	if !(items[0].PositionType == 8 && floatEqual(items[0].Position.X, 1) && floatEqual(items[0].Position.Y, 2)) {
		t.Errorf("type&position except (8, 1,2) but (%d, %f,%f)", items[0].PositionType, items[0].Position.X, items[0].Position.Y)
	}
	// if !(items[0].UsePosition && items[0].X == 1 && items[0].Y == 20) {
	// 	t.Errorf("parse position faild except (1,20) but (%d, %d)", items[0].X, items[0].Y)
	// }
}

func TestParseMulti(t *testing.T) {
	text := `2478
01:56:59,333 --> 01:57:01,963
9141

2479
01:57:03,464 --> 01:57:06,234
11462份申请`
	items, _ := Parse(strings.NewReader(text), 384, 303)
	if items[0].Content[0].Content != "9141" {
		t.Error("first content should be '9141' but", items[0].Content[0].Content)
	}
	if items[1].Content[0].Content != "11462份申请" {
		t.Error("first content should be '11462份申请' but", items[0].Content[0].Content)
	}
}

func TestParsePanic(t *testing.T) {
	// bytes, _ := ioutil.ReadFile("c.srt")
	// str := string(bytes)
	f, _ := os.Open("c.srt")
	Parse(f, 384, 303)
}

func TestDropSvg(t *testing.T) {
	text := `1
01:57:03,464 --> 01:57:06,234
m 325 0 l 355 0 l 332 47 l 354 47 l 399 136 l 369 136 l 328 53 l 288 136 l 257 136 l 325 0 m 467 0 l 497 0 l 474 47 l 496 47 l 541 136 l 511 136 l 470 53 l 430 136 l 399 136 l 467 0 m 545 1 l 583 1 l 583 14 l 568 14 l 568 19 l 583 19 l 583 30 l 568 30 l 568 34 l 599 34 l 599 30 l 583 30 l 583 19 l 599 19 l 599 14 l 583 14 l 583 1 l 611 1 b 616 1 622 6 622 10 l 622 36 l 652 0 l 678 0 l 644 41 l 622 41 l 622 47 l 596 47 l 597 54 l 625 54 l 625 68 l 541 68 l 541 54 l 572 54 l 571 47 l 545 47 l 545 1 m 583 72 l 583 85 l 569 85 l 569 90 l 598 90 l 598 85 l 583 85 l 583 72 l 611 72 b 615 72 621 78 621 82 l 653 44 l 678 44 l 644 86 l 621 86 l 621 103 l 597 103 l 597 136 l 570 136 l 564 126 l 562 136 l 542 136 l 548 107 l 568 107 l 565 121 l 571 121 l 571 103 l 547 103 l 547 72 l 583 72 m 600 107 l 620 107 l 624 124 l 653 89 l 679 89 l 642 136 l 615 136 l 618 132 l 606 132 l 600 107 m 689 0 l 716 0 l 721 15 l 732 15 l 732 30 l 718 56 l 731 56 l 735 100 l 721 100 l 717 59 l 714 64 l 714 136 l 693 136 l 693 79 l 676 79 l 707 30 l 679 30 l 679 15 l 694 15 l 689 0 m 738 0 l 804 0 b 807 0 813 6 813 9 l 813 87 l 794 87 l 794 14 l 756 14 l 756 87 l 763 77 l 763 21 l 787 21 l 787 91 l 798 91 l 798 120 l 820 120 l 812 136 l 778 136 l 778 90 l 748 136 l 723 136 l 756 87 l 738 87 l 738 0 m 257 151 l 275 151 l 297 182 l 319 151 l 337 151 l 304 197 l 304 227 l 290 227 l 290 197 l 257 151 m 337 151 l 355 151 l 377 182 l 399 151 l 417 151 l 384 197 l 384 227 l 370 227 l 370 197 l 337 151 m 425 192 l 447 192 b 445 181 427 181 425 192 l 410 192 b 414 162 458 162 462 192 l 462 206 l 425 206 b 425 210 429 212 433 213 l 462 213 l 462 227 l 433 227 b 412 227 408 203 410 192 m 462 151 l 532 151 l 532 165 l 504 165 l 504 227 l 489 227 l 489 165 l 462 165 l 462 151 m 580 172 l 580 186 l 549 186 b 543 186 543 192 549 192 l 565 192 b 589 192 589 227 565 227 l 532 227 l 532 213 l 565 213 b 570 213 570 206 565 206 l 549 206 b 524 206 524 172 549 172 l 580 172 m 592 213 l 606 213 l 606 227 l 592 227 l 592 213 m 639 172 l 665 172 l 665 186 l 639 186 b 623 186 623 213 639 213 l 665 213 l 665 227 l 639 227 b 603 227 603 172 639 172 m 700 184 b 722 184 722 215 700 215 l 700 229 b 740 229 740 170 700 170 b 660 170 660 229 700 229 l 700 215 b 680 215 680 184 700 184 m 737 172 l 782 172 b 803 172 813 177 813 198 l 813 228 l 799 228 l 799 198 b 799 186 793 187 782 187 l 782 228 l 768 228 l 768 187 l 752 187 l 752 228 l 737 228 l 737 172`
	items, _ := Parse(strings.NewReader(text), 300, 300)
	if items[0].Content[0].Content != "" {
		t.Errorf("Expect empty but %s", items[0].Content[0].Content)
	}
}

func TestParse_d(t *testing.T) {
	r, _ := os.Open("d.srt")
	items, _ := Parse(r, 300, 300)
	if len(items) != 709 {
		t.Errorf("Expect 709 items but %d items.", len(items))
	}
}
func TestParse_e(t *testing.T) {
	r, _ := os.Open("e.srt")
	items, _ := Parse(r, 300, 300)
	if len(items) != 701 {
		t.Errorf("Expect 709 items but %d items.", len(items))
	}
}

func TestParse_f(t *testing.T) {
	r, _ := os.Open("f.srt")
	items, err := Parse(r, 300, 300)
	if err != nil {
		t.Error(err)
	}

	printSubs(items)
}

func printSubs(subs []*shared.SubItem) {
	for _, s := range subs {
		for _, c := range s.Content {
			println(c.Content)
		}
	}
}
