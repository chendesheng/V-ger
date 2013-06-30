package thunder

import (
	"regexp"
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
