package main

import (
	"bytes"
	. "player/shared"
	// "time"
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

	str := AttributedString{"hello", 0, 0}
	strs := make([]AttributedString, 0)
	strs = append(strs, str)

	h1 := w.ShowText(&SubItem{0, 0, strs, 10, Position{50, 700}})

	str = AttributedString{"Hello World1", 0, 0xffffff}
	strs = make([]AttributedString, 0)
	strs = append(strs, str)

	h := w.ShowText(&SubItem{0, 0, strs, 1, Position{0, 0}})

	str = AttributedString{"Hello World2", 0, 0xffffff}
	strs = make([]AttributedString, 0)
	strs = append(strs, str)
	h = w.ShowText(&SubItem{0, 0, strs, 2, Position{0, 0}})

	str = AttributedString{"Hello World3", 0, 0xffffff}
	strs = make([]AttributedString, 0)
	strs = append(strs, str)
	h = w.ShowText(&SubItem{0, 0, strs, 3, Position{0, 0}})

	str = AttributedString{"Hello World4", 0, 0xffffff}
	strs = make([]AttributedString, 0)
	strs = append(strs, str)
	h = w.ShowText(&SubItem{0, 0, strs, 4, Position{0, 0}})

	str = AttributedString{"Hello World5", 0, 0xffffff}
	strs = make([]AttributedString, 0)
	strs = append(strs, str)
	h = w.ShowText(&SubItem{0, 0, strs, 5, Position{0, 0}})

	str = AttributedString{"Hello World6", 0, 0xffffff}
	strs = make([]AttributedString, 0)
	strs = append(strs, str)
	h = w.ShowText(&SubItem{0, 0, strs, 6, Position{0, 0}})

	str = AttributedString{"Hello World7", 0, 0xffffff}
	strs = make([]AttributedString, 0)
	strs = append(strs, str)
	h = w.ShowText(&SubItem{0, 0, strs, 7, Position{0, 0}})

	str = AttributedString{"Hello World9", 0, 0xffffff}
	strs = make([]AttributedString, 0)
	strs = append(strs, str)
	h = w.ShowText(&SubItem{0, 0, strs, 9, Position{0, 0}})

	str = AttributedString{"Hello World8", 0, 0xffffff}
	strs = make([]AttributedString, 0)
	strs = append(strs, str)
	h3 := w.ShowText(&SubItem{0, 0, strs, 8, Position{0, 0}})

	// go func() {
	// time.Sleep(3 * time.Second)
	// w.HideText(h)
	// }()
	println(h)
	println(h1)
	println(h3)
	w.FuncKeyDown = append(w.FuncKeyDown, func(key int) {
		if key == KEY_0 {
			w.HideText(h3)
		}
	})

	audioStreamNames := make([]string, 0)
	audioStreamNames = append(audioStreamNames, "a1")
	audioStreamNames = append(audioStreamNames, "a2")
	indexes := make([]int32, 0)
	indexes = append(indexes, 1)
	indexes = append(indexes, 2)
	w.InitAudioMenu(audioStreamNames, indexes, 2)
	w.FuncAudioMenuClicked = append(w.FuncAudioMenuClicked, func(i int) {
		println(i)
	})

	PollEvents()
}
