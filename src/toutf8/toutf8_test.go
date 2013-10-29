package toutf8

import (
	// "io/ioutil"
	// "os"
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

	res, err := ConverToUTF8("utf16le.txt")
	if err != nil {
		t.Error(err)
	}

	println(res)
	println(len(res))
}
