package srt

import (
	"encoding/hex"
	// "fmt"
	"log"
	"sort"
	// "io/ioutil"
	"bytes"
	"github.com/peterbourgon/html"
	"regexp"
	"strconv"
	"strings"
	"time"

	. "player/shared"
)

type SubItems []*SubItem

func (s SubItems) Len() int {
	return len([]*SubItem(s))
}

func (s SubItems) Less(i, j int) bool {
	return s[i].From < s[j].From
}

func (s SubItems) Swap(i, j int) {
	t := s[i]
	s[i] = s[j]
	s[j] = t
}

func linebreak(r rune) bool {
	return r == '\r' || r == '\n'
}

func Parse(str string, width, height float64) []*SubItem {
	lines := strings.FieldsFunc(str, linebreak)

	items := make([]*SubItem, 0)

	parseContent(&lines)

	// log.Print("head: ", head)
	for len(lines) > 0 {
		if ok, from, to := parseTime(lines[0]); ok {
			// println("line after parseTime:", lines[0])
			lines = lines[1:]

			usePos, pos, text := parsePosition(lines[0], width, height)
			lines[0] = text
			content := parseContent(&lines)
			// log.Print("content:", content)
			items = append(items, &SubItem{from, to, content, usePos, pos, SubItemExtra{0, 0}})
		} else {
			log.Println("parse time error:", lines[0])
			panic("parse error")
		}
	}

	sort.Sort(SubItems(items))
	return items
}

func parseContent(lines *[]string) []AttributedString {
	i := 0
	for ; i < len(*lines); i++ {
		line := (*lines)[i]
		// println("line:", line)

		// _, err := strconv.Atoi(line)
		// if err == nil {
		// 	break
		// }

		ok, _, _ := parseTime(line)
		if ok {
			break
		}
	}
	content := ""
	if i == 0 {
		*lines = nil
	} else if i == len(*lines) {
		content = strings.Join((*lines)[:i], "\n")
		*lines = nil
	} else {
		if i > 1 {
			content = strings.Join((*lines)[:i-1], "\n")
		}
		*lines = (*lines)[i:]
	}

	// log.Print("i:",i)
	// if i+1 < len(*lines) {
	// 	println("parseContent:", i)
	// 	content := ""
	// 	if i > 1 {
	// 		content = strings.Join((*lines)[:i-1], "\n")
	// 	}
	// 	*lines = (*lines)[i:]
	// } else {
	// 	// *lines = make([]string,0)
	// 	*lines = nil
	// }

	return parseAttributedString(content)
	// return content
}
func toColor(c string) uint {
	c = strings.ToLower(c)
	switch c {
	case "white":
		return 0xffffff
	case "black":
		return 0
		//...
	}

	defaultColor := uint(0)

	if c[0] == '#' {
		if len(c) == 4 {
			c = c[1:]
			bytes := []byte(c)
			b := make([]byte, 6)
			b[0] = bytes[0]
			b[1] = bytes[0]
			b[2] = bytes[1]
			b[3] = bytes[1]
			b[4] = bytes[2]
			b[5] = bytes[2]
			c = string(b)
		} else if len(c) == 7 {
			c = c[1:]
		}

		if len(c) == 6 {
			bytes, err := hex.DecodeString(c)
			if err != nil {
				return defaultColor
			}

			// println("test")
			// println(c)
			// println(bytes[0], bytes[1], bytes[2])
			return uint(bytes[2]) + (uint(bytes[1]) << 8) + (uint(bytes[0]) << 16)
		} else {
			return defaultColor
		}
	} else {
		return defaultColor
	}
}
func removePositionInfo(text string) string {
	regPos := regexp.MustCompile(`\{\\[^}]+\}`)
	return regPos.ReplaceAllString(text, "")
}

func parsePos(text string, width, height float64) Position {
	regPos := regexp.MustCompile(`\{\\pos\(([0-9]+)[.]?[0-9]*,([0-9]+)[.]?[0-9]*\)\}`)
	matches := regPos.FindStringSubmatch(text)
	if matches == nil {
		return Position{-1, -1}
	} else {
		x, _ := strconv.Atoi(matches[1])
		y, _ := strconv.Atoi(matches[2])

		//w:384,h:303 two magic numbers come from Baofeng player, don't why
		return Position{float64(x) / 384 * width, float64(y) / 303 * height}
	}
}
func parseAlign(text string) int {
	regPos := regexp.MustCompile(`^\{\\an?([1-9])\}`)
	matches := regPos.FindStringSubmatch(text)
	if matches == nil {
		return 2
	} else {
		i, _ := strconv.Atoi(matches[1])
		return i
	}
}

func parsePosition(text string, width, height float64) (int, Position, string) {
	return parseAlign(text), parsePos(text, width, height), removePositionInfo(text)
}
func parseTag(nodes []*html.Node, as AttributedString, res *[]AttributedString) {
	for _, n := range nodes {
		if n.Type == html.TextNode {
			as.Content = n.Data
			*res = append(*res, as)
			as.Content = ""
		} else if n.Type == html.ElementNode {
			savedas := as
			switch strings.ToLower(n.Data) {
			case "i":
				as.Style |= 1
				break
			case "b":
				as.Style |= 2
				break
			case "font":
				for _, attr := range n.Attr {
					if strings.ToLower(attr.Key) == "color" {
						as.Color = toColor(attr.Val)
					}
				}
				break
			}
			parseTag(n.Child, as, res)
			as = savedas
		}
	}
}
func parseAttributedString(text string) []AttributedString {
	res := make([]AttributedString, 0)

	r := bytes.NewReader([]byte(text))
	nodes, err := html.ParseFragment(r, nil)
	if err == nil {
		parseTag(nodes, AttributedString{
			Color: 0xffffff,
		}, &res)
	}

	return res
}

func parseTime(line string) (bool, time.Duration, time.Duration) {
	r := regexp.MustCompile(`([0-9]{2}):([0-9]{2}):([0-9]{2}),([0-9]{3})\s+-->\s+([0-9]{2}):([0-9]{2}):([0-9]{2}),([0-9]{3})`)
	if r.MatchString(line) {
		matches := r.FindStringSubmatch(line)
		// fmt.Printf("%v", matches)
		return true, convertTime(matches[1:5]), convertTime(matches[5:9])
	} else {
		// fmt.Printf("error")
		return false, 0, 0
	}
}

func convertTime(parts []string) time.Duration {
	for i, p := range parts {
		parts[i] = strings.TrimLeft(p, "0")
	}
	h, _ := strconv.Atoi(parts[0])
	m, _ := strconv.Atoi(parts[1])
	s, _ := strconv.Atoi(parts[2])
	ms, _ := strconv.Atoi(parts[3])

	return time.Duration(h)*time.Hour + time.Duration(m)*time.Minute + time.Duration(s)*time.Second + time.Duration(ms)*time.Millisecond
}

// func main() {
// 	srt, err := ioutil.ReadFile("a.srt")
// 	if err != nil {
// 		return
// 	}
// 	// log.Print(string(srt))
// 	// fmt.log.Print("Hello, playground")
// 	head, items := Parse(string(srt))
// 	log.Print("head:", head)
// 	for _, item := range items {
// 		fmt.Printf("item: %v\n", item)
// 	}
// }
