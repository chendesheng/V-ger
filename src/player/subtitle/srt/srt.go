package srt

import (
	"encoding/hex"
	"fmt"
	"io"
	// "fmt"
	// "log"
	"sort"
	// "io/ioutil"
	"bufio"
	"bytes"
	"github.com/peterbourgon/html"
	"regexp"
	"strconv"
	"strings"
	"time"

	. "player/shared"
)

type parser struct {
	*LineScanner

	items []*SubItem
}

func newParser(r io.Reader) *parser {
	p := parser{}
	p.LineScanner = (*LineScanner)(bufio.NewScanner(r))
	return &p
}
func (p *parser) appendItem(item SubItem, text string, width, height float64) {
	if item.To != 0 { //item is not inited yet
		text = strings.Trim(text, "\n")

		// drop last line
		i := strings.LastIndex(text, "\n")
		if i >= 0 {
			text = text[:i]
		}
		item.PositionType, item.Position, text = parsePosition(text, width, height)

		item.Content = parseAttributedString(text)
		if len(item.Content) >= 1 {
			item.Content[0].Content = dropSvgContent(item.Content[0].Content)
		}

		p.items = append(p.items, &item)
	}
}
func dropSvgContent(text string) string {
	text = strings.TrimLeft(text, " \r\n\t")
	regSvg := regexp.MustCompile("^([mlb] ([0-9]+ ?)+)+")

	if regSvg.MatchString(text) {
		return ""
	} else {
		return text
	}
}

func (p *parser) parse(width, height float64) {
	var item SubItem
	var text string

	for {
		line := p.NextLine()
		if len(line) == 0 {
			p.appendItem(item, fmt.Sprintf("%s\n0", text), width, height)
			return
		}

		if ok, from, to := parseTime(line); ok {
			p.appendItem(item, text, width, height)

			//for next item
			text = ""
			item.From = from
			item.To = to
			// item.SubItemExtra = SubItemExtra{0, 0}

		} else {
			text = fmt.Sprintf("%s\n%s", text, line)
		}
	}

	return
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
	r := regexp.MustCompile(`([0-9]{1,2}):([0-9]{1,2}):([0-9]{1,2})[,.]([0-9]{1,3})\s+-->\s+([0-9]{1,2}):([0-9]{1,2}):([0-9]{2})[,.]([0-9]{1,3})`)
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

func Parse(r io.Reader, width, height float64) (items []*SubItem, err error) {
	defer func() {
		r := recover()
		if r != nil {
			err = r.(error)
			items = nil
		}
	}()

	p := newParser(r)
	p.parse(width, height)

	items = p.items
	sort.Sort(SubItems(items))

	err = nil

	return
}
