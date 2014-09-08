package srt

import (
	"encoding/hex"
	"io"
	"io/ioutil"
	"log"
	// "fmt"
	// "log"
	"bytes"
	"regexp"
	"runtime/debug"
	"sort" // "io/ioutil"
	"strconv"
	"strings"
	"time"
	"github.com/peterbourgon/html"

	. "player/shared"
)

var regAn = regexp.MustCompile(`^\{\\an?([1-9])\}`)
var regSvg = regexp.MustCompile("^([mlb] ([0-9]+ ?)+)+")
var regPos = regexp.MustCompile(`\{\\pos\(([0-9]+)[.]?[0-9]*,([0-9]+)[.]?[0-9]*\)\}`)
var regBreak = regexp.MustCompile("(?i)\\\\n")
var regSubItem = regexp.MustCompile(`[0-9]+(?:\r\n|\r|\n)([0-9]{1,2}):([0-9]{1,2}):([0-9]{1,2})[.,]([0-9]{1,3}).*-->.*([0-9]{1,2}):([0-9]{1,2}):([0-9]{1,2})[.,]([0-9]{1,3})(?:\r\n|\r|\n)((?:.*(?:\r\n|\r|\n))*?)\s*(?:\r\n|\r|\n)`)

func dropSvgContent(text string) string {
	text = strings.TrimLeft(text, " \r\n\t")

	if regSvg.MatchString(text) {
		return ""
	} else {
		return text
	}
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
	matches := regAn.FindStringSubmatch(text)
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
			log.Print(err.Error())
			log.Print(string(debug.Stack()))
		}
	}()

	bytes, _ := ioutil.ReadAll(r)
	text := string(bytes) + "\n\n" //make sure the last item has enough \n to match the reg exp

	matches := regSubItem.FindAllStringSubmatch(text, -1)

	items = make([]*SubItem, 0, len(matches))

	for _, item := range matches {
		s := &SubItem{}
		s.From = convertTime(item[1:5])
		s.To = convertTime(item[5:9])

		text := strings.TrimSpace(item[9])
		text = regBreak.ReplaceAllString(text, "\n")

		s.PositionType, s.Position, text = parsePosition(text, width, height)
		s.Content = parseAttributedString(text)
		if len(s.Content) >= 1 {
			s.Content[0].Content = dropSvgContent(s.Content[0].Content)
		}

		if !s.IsEmpty() {
			items = append(items, s)
		}
	}

	sort.Sort(SubItems(items))

	err = nil

	return
}
