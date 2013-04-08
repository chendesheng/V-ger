package bencoding

import (
	"io/ioutil"
	"log"
	"os"
	"testing"
)

func init() {
	f, _ := os.OpenFile("test.log", os.O_CREATE|os.O_APPEND, 0666)
	log.SetOutput(f)
}
func TestParseString(t *testing.T) {
	pos, str := parseString([]byte("40:spamspamspamspamspamspamspamspamspamspamspamspam"), 0)

	if pos != 43 {
		t.Error("pos is not 43")
	}
	if str != "spamspamspamspamspamspamspamspamspamspam" {
		t.Error("str is not spamspamspamspamspamspamspamspamspamspam")
	}
}

func TestParseInt(t *testing.T) {
	pos, i := parseInteger([]byte("i1234567890e"), 0)
	if pos != 12 {
		t.Error("pos should be 12")
	}

	if i != 1234567890 {
		t.Error("i should be 1234567890")
	}
}

func TestParseObjectString(t *testing.T) {
	pos, obj := parseObject([]byte("4:spam"), 0)

	if pos != 6 {
		t.Error("pos should be 6")
	}

	if obj.(string) != "spam" {
		t.Error("obj should be spam")
	}
}
func TestParseObjectInt(t *testing.T) {
	pos, obj := parseObject([]byte("i1234e"), 0)

	if pos != 6 {
		t.Error("pos should be 6")
	}

	if obj.(int64) != 1234 {
		t.Error("obj should be 1234")
	}
}

func TestParseList(t *testing.T) {
	pos, l := parseList([]byte("li123e4:123e4:eeeee"), 0)
	if pos != 19 {
		t.Error("pos should be 13")
	}
	if l[0].(int64) != 123 {
		t.Error("l[0] should be 123")
	}
	if l[1].(string) != "123e" {
		t.Error("l[1] should be 123e")
	}
	if l[2].(string) != "eeee" {
		t.Error("l[2] should be eeee")
	}
}
func TestParseDict(t *testing.T) {
	_, d := parseDict([]byte("d4:spami123e5:spam15:valuee"), 0)
	if d["spam"].(int64) != 123 {
		t.Error("d[\"spam\"] should be 123")
	}
	if d["spam1"].(string) != "value" {
		t.Error("d[\"spam1\"] should be value")
	}
}
func TestParse(t *testing.T) {
	res, _ := Parse([]byte("d4:spami123e5:spam15:value5:spam2li222e4:spamee"))
	obj := res.(map[string]interface{})
	if obj["spam"].(int64) != 123 {
		t.Error("d[\"spam\"] should be 123")
	}
	if obj["spam1"].(string) != "value" {
		t.Error("d[\"spam1\"] should be value")
	}

	l := obj["spam2"].([]interface{})
	if l[0].(int64) != 222 {
		t.Error("d[\"spam2\"][0] should be 222")
	}
	if l[1].(string) != "spam" {
		t.Error("d[\"spam2\"][1] should be spam")
	}
}
func TestDumpTorrent(t *testing.T) {
	f, _ := os.OpenFile("test.torrent", os.O_RDONLY, 0666)
	bytes, _ := ioutil.ReadAll(f)
	res, _ := Parse(bytes)
	pieces := []byte(res.(map[string]interface{})["info"].(map[string]interface{})["pieces"].(string))
	res.(map[string]interface{})["info"].(map[string]interface{})["pieces"] = nil
	t.Errorf("%v\n", res)
	t.Error(len(pieces))
	t.Error(pieces)
	// t.Error(res.(map[string]interface{})["announce"])
	// t.Error(res.(map[string]interface{})["announce-list"])
	// t.Error(res.(map[string]interface{})["creation date"])
	// t.Error(res.(map[string]interface{})["comment"])
	// t.Error(res.(map[string]interface{})["created by"])
	// t.Error(res.(map[string]interface{})["encoding"])
	// t.Error(res.(map[string]interface{})["info"].(map[string]interface{})["piece length"])
	// // t.Error(res.(map[string]interface{})["info"].(map[string]interface{})["pieces"])
	// t.Error(res.(map[string]interface{})["info"].(map[string]interface{})["private"])
	// t.Error(res.(map[string]interface{})["info"].(map[string]interface{})["name"])
	// t.Error(res.(map[string]interface{})["info"].(map[string]interface{})["files"])

	// t.Errorf("%v", res)
}
