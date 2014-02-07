package toutf8

import (
	// "fmt"
	"io/ioutil"
	"os"
	// "strings"
	"testing"
)

// func TestGB18030ToUTF8(t *testing.T) {
// 	// data, err := ioutil.ReadFile("gb18030.txt")
// 	// if err != nil {
// 	// 	t.Error(err)
// 	// 	return
// 	// }

// 	res, err := ConverToUTF8("gb18030.txt")
// 	if err != nil {
// 		t.Error(err)
// 	}

// 	println(res)
// }

func TestGuessEncoding(t *testing.T) {
	// os.OpenFile("utf16le.txt", flag, perm)
	// data, err := ioutil.ReadFile("gb18030.txt")
	// if err != nil {
	// 	t.Error(err)
	// 	return
	// }

	// res, err := ConverToUTF8("utf16le.txt")
	// if err != nil {
	// 	t.Error(err)
	// }

	// println(res)
	// println(len(res))

	// infoes, _ := ioutil.ReadDir("/Volumes/Data/Downloads/Video/Girls")
	// for _, f := range infoes {
	// 	println(f.Name())
	// 	if strings.HasPrefix(f.Name(), "Icon") {
	// 		if f.IsDir() {
	// 			println("yes")
	// 		}
	// 		println("size:", f.Size())
	// r, err := os.Open("/Volumes/Data/Downloads/Video/Girls/Icon\r/..namedfork/rsrc")

	// if err != nil {
	// 	println(err)
	// }

	// bytes, err := ioutil.ReadAll(r)
	// if err != nil {
	// 	println(err)
	// }
	// println(len(bytes))
	// r.Close()
	bytes, _ := ioutil.ReadFile("/Volumes/Data/Downloads/Video/Rake/Icon\r/..namedfork/rsrc")
	// fmt.Printf("%v", bytes)
	// bytes, _ := ioutil.ReadFile("/Volumes/Data/Downloads/Video/Rake/b.jpg")
	// ioutil.WriteFile("/Volumes/Data/Downloads/Video/Rake/Icon\r/..namedfork/rsrc", bytes, os.ModeDevice)
	ioutil.WriteFile("/Volumes/Data/Downloads/Video/Rake/a.jpg", bytes, os.ModeType)
	// f, err := os.OpenFile("/Volumes/Data/Downloads/Video/Rake/Icon\r/..namedfork/rsrc", os.O_RDWR)
	// f.Write(bytes)
	// 	}
	// }
}
