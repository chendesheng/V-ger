package main

import (
	"bytes"
	// "fmt"
	"io/ioutil"
	. "player/gui"
	"strconv"
	"strings"
)

func readPPMFile(file string) (int, int, []byte) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		println(err.Error())
	}
	println(string(data[:16]))

	l1 := bytes.IndexByte(data, '\n')
	println("l1", l1)
	l2 := bytes.IndexByte(data[l1+1:], '\n') + l1 + 1
	println("l2", l2)
	l3 := bytes.IndexByte(data[l2+1:], '\n') + l2 + 1
	println("l3", l3)

	str := string(data[l1+1 : l2])
	strs := strings.Split(str, " ")
	w, _ := strconv.Atoi(strs[0])
	h, _ := strconv.Atoi(strs[1])
	return w, h, data[l3+1:]
}
func main() {
	width, height, img := readPPMFile("b.ppm")
	// fmt.Printf("%v", img)
	// linesize := width * 3
	// for i := 0; i < height*linesize; i += 3 {
	// 	img[i] = 20
	// 	img[i+1] = 255
	// 	img[i+2] = 20
	// }

	w := NewWindow("b.ppm", width, height)
	// w.ShowText(strs, withPosition, x, y)
	// w.ChanDraw <- img
	w.RefreshContent(img)
	PollEvents()
}
