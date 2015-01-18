package thunder

import (
	"fmt"
	"regexp"
	// "strconv"
	"strings"
	"testing"
)

func TestRegexp(t *testing.T) {
	regexUrlQuery := regexp.MustCompile(`queryUrl\((-?[0-9]*),'([^']*)','([^']*)','([^']*)','([^']*)',new Array\((.+)\),new Array\(([^)]*)\),new Array\(([^)]*)\),new Array\(([^)]*)\),new Array\(([^)]*)\),new Array\(([^)]*)\),'[^']*','[^']*'\)`)

	text := "queryUrl(1,'9087C4A034B7082D2A984748BA9991CBDAF9E9F8','13100042557','Carlito\\'s.Way.1993.720p.BluRay.DTS.x264-dxva.mkv','0',new Array('Carlito\\'s.Way.1993.720p.BluRay.DTS.x264-dxva.mkv'),new Array('12.2G'),new Array('13100042557'),new Array('1'),new Array('RMVB'),new Array('0'),'','0')"

	text = strings.Replace(text, "\\'", "", -1)

	matches := regexUrlQuery.FindStringSubmatch(text)

	if len(matches) == 0 {
		t.Error("not match")
	}
}

func TestDecodeFid(t *testing.T) {
	fid := "Zb0SrmlbQ2AQrCCwk52l570ZspG5oDukAQAAAEy6ExeTT35g5NKTyq6cayOh/x8+"
	cid, size, gcid := ParseFid(fid)

	if cid != "65BD12AE695B436010AC20B0939DA5E7BD19B291" {
		t.Errorf("Expect cid '65BD12AE695B436010AC20B0939DA5E7BD19B291', but '%s'", cid)
	}
	if size != 7050338489 {
		t.Errorf("Expect size '7050338489', but '%d'", size)
	}
	if gcid != "4CBA1317934F7E60E4D293CAAE9C6B23A1FF1F3E" {
		t.Errorf("Expect gcid '4CBA1317934F7E60E4D293CAAE9C6B23A1FF1F3E', but '%s'", "", gcid)
	}
}

func TestDecodeFid2(t *testing.T) {
	fid := "jS1NEVi6Jf+5ozPXoq0N9L4mNOILCQ0yAAAAAAX8xPViluQV9HVgg7q1YIMIOMa8"
	cid, size, gcid := ParseFid(fid)

	println(cid, size, gcid)
	// if cid != "65BD12AE695B436010AC20B0939DA5E7BD19B291" {
	// 	t.Errorf("Expect cid '65BD12AE695B436010AC20B0939DA5E7BD19B291', but '%s'", cid)
	// }
	// if size != 7050338489 {
	// 	t.Errorf("Expect size '7050338489', but '%d'", size)
	// }
	// if gcid != "4CBA1317934F7E60E4D293CAAE9C6B23A1FF1F3E" {
	// 	t.Errorf("Expect gcid '4CBA1317934F7E60E4D293CAAE9C6B23A1FF1F3E', but '%s'", "", gcid)
	// }
}

func TestMagnet(t *testing.T) {
	UserName = ""
	Password = ""

	tks, err := NewTask("magnet:?xt=urn:btih:66883EE64827A8D7C4DD04209568D55F5BB7572A", "")
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("tasks:%#v", tks)
}
