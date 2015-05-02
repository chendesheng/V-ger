package ass

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"regexp"
	"runtime/debug"
	"sort"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
	. "vger/player/shared"
)

var regAn = regexp.MustCompile(`^an([0-9])`)
var regSvg = regexp.MustCompile("^([mlb] (-?[0-9]+ ?)+)+")
var regBold = regexp.MustCompile(`^b([0-9])+`)
var regColor = regexp.MustCompile(`^[0-9]?c&H([0-9a-fA-F]+)&`)
var regPos = regexp.MustCompile(`^pos\(([0-9]+)[.]?[0-9]*,([0-9]+)[.]?[0-9]*\)`)
var regBreak = regexp.MustCompile("(?i)\\\\n")

type parser struct {
	*LineScanner

	formats []string
	items   []*SubItem

	width  float64
	height float64
}

func Parse(r io.Reader, width, height float64) (items []*SubItem, err error) {
	defer func() {
		r := recover()
		if r != nil {
			err = r.(error)
			items = nil

			log.Printf("%s\n%s", err, string(debug.Stack()))
		}
	}()

	p := newParser(r, width, height)
	p.parse()

	items = p.items
	sort.Sort(SubItems(items))

	err = nil

	return
}

func newParser(r io.Reader, width float64, height float64) *parser {
	p := parser{}
	p.LineScanner = (*LineScanner)(bufio.NewScanner(r))
	p.width, p.height = width, height
	return &p
}
func (p *parser) parseSection(line string) {
	if line == "[Events]" {
		line = p.NextLine()
		if len(line) == 0 {
			return
		}

		if line[0] == '[' {
			p.parseSection(line)
		} else {
			p.formats = parseFormats(line)
			for {
				line = p.NextLine()
				if len(line) == 0 {
					return
				}

				d := p.parseDialogue(line)
				if d != nil {
					p.items = append(p.items, d)
				}
			}
		}
	}
}
func (p *parser) parse() {
	for {
		line := p.NextLine()
		if len(line) == 0 {
			return
		}

		if line[0] == '[' {
			p.parseSection(line)
		}
	}
}

func parseField(content string) (string, string) {
	var state int

	for i, c := range content {

		if c == '{' {
			state = 1
		}
		if c == '}' {
			state = 0
		}

		if state == 0 && c == ',' {
			return strings.Trim(content[:i], " \t\r"), content[i+1:]
		}
	}

	return content, ""
}

func parseTime(content string) time.Duration {
	var h, m, s, ms time.Duration
	fmt.Sscanf(content, "%d:%d:%d.%d", &h, &m, &s, &ms)

	return h*time.Hour + m*time.Minute + s*time.Second + ms*10*time.Millisecond
}

func (p *parser) parseText(text string, item *SubItem) {
	var state int
	part := make([]rune, 0)
	attrstart := 0
	attrend := 1

	var currentStyle int
	var currentColor uint
	var currentPos Position
	currentPos.X = -1
	currentPos.Y = -1
	currentColor = 0xffffff
	var currentAlign int
	currentAlign = 2

	res := make([]AttributedString, 0)

	for i, w := 0, 0; i < len(text); i += w {
		c, width := utf8.DecodeRuneInString(text[i:])
		// fmt.Printf("%#U starts at byte position %d\n", runeValue, i)

		if c == '{' {
			state = 1
			attrstart = i
		}

		if i < len(text)-1 && c == '\\' && text[i+1] == 'N' {
			part = append(part, '\n')
			i++
		} else if state == 0 {
			part = append(part, c)
		}

		if c == '}' {
			state = 0
			attrend = i

			if len(part) > 0 {
				as := AttributedString{}
				as.Color = currentColor
				// as.Position = currentPos
				as.Style = currentStyle
				as.Content = string(part)
				res = append(res, as)
				part = make([]rune, 0)

				// fmt.Printf("%v", as)
			}

			currentStyle, currentColor, currentPos, currentAlign = p.parseAttr(text[attrstart:attrend+1], currentStyle, currentColor, currentPos, currentAlign)
		}

		w = width
	}

	// println(string(part))
	as := AttributedString{}
	as.Color = currentColor
	// as.Position = currentPos
	as.Style = currentStyle
	as.Content = string(part)
	res = append(res, as)
	item.Content = res
	item.Position = currentPos
	item.PositionType = currentAlign

	if len(item.Content) >= 1 {
		item.Content[0].Content = dropSvgContent(item.Content[0].Content)
	}
	// fmt.Printf("%v", *item)
}
func dropSvgContent(text string) string {
	text = strings.TrimLeft(text, " \r\n\t")

	if regSvg.MatchString(text) {
		return ""
	} else {
		return text
	}
}
func (p *parser) parsePos(text string) (bool, Position) {
	matches := regPos.FindStringSubmatch(text)

	if matches != nil {
		x, _ := strconv.Atoi(matches[1])
		y, _ := strconv.Atoi(matches[2])
		return true, Position{float64(x) / 384 * p.width, float64(y) / 303 * p.height}
	} else {
		return false, Position{-1, -1}
	}
}
func parseAlignment(text string) int {
	matches := regAn.FindStringSubmatch(text)
	if matches != nil {
		i, _ := strconv.Atoi(matches[1])
		return i
	} else {
		return 2
	}
}
func parseBold(text string) (bool, int) {
	matches := regBold.FindStringSubmatch(text)
	if matches != nil {
		i, _ := strconv.Atoi(matches[1])
		return true, i
	} else {
		return false, 0
	}
}
func hex2int(c byte) uint {
	if c >= 'a' && c <= 'f' {
		return 10 + uint(c-'a')
	} else if c >= 'A' && c <= 'F' {
		return 10 + uint(c-'A')
	} else {
		return uint(c - '0')
	}
}

//Dialogue: 0,0:04:26.38,0:04:32.38,FlancyForBrBa,NTP,0000,0000,0000,,{\fade(255,0,255,0,2000,4000,5000)}{\fs16\3c&H18490B&}本字幕由 {\fs22\b1\3c&H1E1EB9&}YounFlancy{\b0}{\fs22\3c&H6C0D0A&}@{\fs22\3c&H9449B1&\c&HFF06FE&\b1}Newsmth{\b0} {\fs16\c&HFFFFFF&\3c&H18490B&}翻译制作\N{\rFlancyForBrBa}仅供交流学习，勿用于商业用途
func parseColor(text string) (bool, uint) {
	// println(text)
	matches := regColor.FindStringSubmatch(text)
	if matches != nil {
		hexStr := matches[1]
		if len(hexStr) < 6 {
			paddingZeroes := ""
			for i := 0; i < 6-len(hexStr); i++ {
				paddingZeroes = paddingZeroes + "0"
			}

			hexStr = paddingZeroes + hexStr
		}

		b1, b2, g1, g2, r1, r2 := hex2int(hexStr[0]), hex2int(hexStr[1]), hex2int(hexStr[2]), hex2int(hexStr[3]), hex2int(hexStr[4]), hex2int(hexStr[5])

		return true, (r1 << 20) + (r2 << 16) + (g1 << 12) + (g2 << 8) + (b1 << 4) + (b2 << 0)
	} else {
		return false, 0xffffff
	}
}
func (p *parser) parseAttr(text string, style int, color uint, pos Position, an int) (int, uint, Position, int) {
	attrs := strings.Split(strings.Trim(text, "{} \r\n\t"), "\\")
	for _, a := range attrs {
		if strings.HasPrefix(a, "pos") {
			if ok, tmpPos := p.parsePos(a); ok {
				pos = tmpPos
			}
		} else if strings.HasPrefix(a, "an") {
			an = parseAlignment(a)
		} else if ok, weight := parseBold(a); ok {
			if weight > 0 {
				style = style | 2
			} else {
				style = style & (^2)
			}
		} else if a == "i0" {
			style = style & (^1)
		} else if a == "i1" {
			style = style | 1
		} else if ok, col := parseColor(a); ok {
			color = col
		}
	}

	return style, color, pos, an
}

func parseLine(line string) (string, string) {
	i := strings.Index(line, ":")
	if i < 0 {
		panic(fmt.Errorf("Parse line error: expect ':'. \n%s", line))
	}
	title := line[:i]
	content := line[i+1:]
	content = regBreak.ReplaceAllString(content, "\n")
	return title, content
}

func parseFormats(line string) []string {
	title, content := parseLine(line)
	if title != "Format" {
		// panic(fmt.Errorf("Expect title 'Format'."))
		return nil
	}

	formats := strings.Split(content, ",")
	for i, _ := range formats {
		formats[i] = strings.Trim(formats[i], " \t\r")
	}
	return formats
}

func (p *parser) parseDialogue(line string) *SubItem {
	defer func() {
		err := recover()
		if err != nil {
			log.Printf("%s\n%s", err, string(debug.Stack()))
		}
	}()

	title, content := parseLine(line)
	if title != "Dialogue" {
		panic(fmt.Errorf("Expect title 'Dialogue'."))
	}

	var item SubItem
	for _, format := range p.formats {
		var field string
		field, content = parseField(content)
		switch format {
		case "Text": //always the last one so it can contains comma
			p.parseText(field+content, &item)
			break
		case "Start":
			item.From = parseTime(field)
			break
		case "End":
			item.To = parseTime(field)
			break
		}
	}
	return &item
}
