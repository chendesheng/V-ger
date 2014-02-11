package shared

import (
	"bufio"
	"strings"
)

type LineScanner bufio.Scanner

func (ls *LineScanner) NextLine() string {
	scanner := (*bufio.Scanner)(ls)

	if scanner.Scan() {
		line := scanner.Text()
		line = strings.Trim(line, " \t\r\n")

		if len(line) == 0 || line[0] == ';' {
			return ls.NextLine()
		} else {
			return line
		}
	} else {
		return ""
	}
}
