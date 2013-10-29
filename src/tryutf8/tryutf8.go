package main

import (
	// "fmt"
	"io/ioutil"
	"unicode/utf8"
)

func main() {
	data, _ := ioutil.ReadFile("a.srt")
	println(string(data))
	out := make([]byte, 0, len(data))
	for len(data) > 0 {
		r, l := utf8.DecodeRune(data)
		// fmt.Printf("0x%x\n", r)
		if r != utf8.RuneError {
			if string(data[:l]) != "\uC080" {
				out = append(out, data[:l]...)
				print(string(data[:l]))
			}
		} else {
			// println("error")
		}
		data = data[l:]
	}
	// println(string(out))
	// ioutil.WriteFile("b.srt", out, 0666)

}
