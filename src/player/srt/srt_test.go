package srt

import (
	// "fmt"
	"io/ioutil"
	"testing"
)

// func TestParse(t *testing.T) {
// 	data, err := ioutil.ReadFile("a.srt")
// 	if err != nil {
// 		t.Error(err)
// 		return
// 	}

// 	text := string(data)
// 	items := Parse(text)
// 	if len(items) == 0 {
// 		t.Errorf("parse error")
// 	}
// }
func TestParse2(t *testing.T) {
	data, err := ioutil.ReadFile("b.srt")
	if err != nil {
		t.Error(err)
		return
	}

	text := string(data)
	items := Parse(text)

	// fmt.Printf("%v", items)
	if len(items) == 0 {
		t.Errorf("parse error")
	}
}
func TestParseItalic(t *testing.T) {
	text := `1
00:00:00,000 --> 00:00:01,111
a<i>你好</i>c`
	items := Parse(text)
	if len(items) != 1 {
		t.Error("should be only 1 item")
		return
	}

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
	items := Parse(text)
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
	items := Parse(text)
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
	items := Parse(text)
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
	items := Parse(text)
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
	items := Parse(text)
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
	items := Parse(text)
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
	items := Parse(text)
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
	items := Parse(text)

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

func TestParsePosition(t *testing.T) {
	text := `1
00:00:00,001 --> 00:00:04,336
{\pos(1.12,20.4)}President reagan:`
	items := Parse(text)
	if items[0].Content[0].Content != "President reagan:" {
		t.Error("first content should be 'President reagan:' but", items[0].Content[0].Content)
	}
	if !(items[0].PositionType == 10 && items[0].X == 1 && items[0].Y == 20) {
		t.Errorf("parse position faild except (1,20) but (%d, %d)", items[0].X, items[0].Y)
	}
}

func TestParsePosition2(t *testing.T) {
	text := `1
00:00:00,001 --> 00:00:04,336
{\an8}President reagan:`
	items := Parse(text)
	if items[0].Content[0].Content != "President reagan:" {
		t.Error("first content should be 'President reagan:' but", items[0].Content[0].Content)
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
	items := Parse(text)
	if items[0].Content[0].Content != "9141" {
		t.Error("first content should be '9141' but", items[0].Content[0].Content)
	}
	if items[1].Content[0].Content != "11462份申请" {
		t.Error("first content should be '11462份申请' but", items[0].Content[0].Content)
	}
}

func TestParsePanic(t *testing.T) {
	bytes, _ := ioutil.ReadFile("c.srt")
	str := string(bytes)

	Parse(str)
}
