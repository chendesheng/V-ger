package toutf8

/*
#cgo LDFLAGS: -lcharsetdetect -liconv
#include "charsetdetect.h"
*/
import "C"
import (
	"fmt"
	"github.com/qiniu/iconv"
	// "io"
	"io"
	"io/ioutil"
	"log"
	// "os"
	"strings"
	// "os/exec"
	"unsafe"
)

type CharsetDetect struct {
	ptr C.csd_t
}

func (csd *CharsetDetect) Open() {
	csd.ptr = C.csd_open()
}

// Feeds some more data to the character set detector. Returns 0 if it
// needs more data to come to a conclusion and a positive number if it has enough to say what
// the character set is. Returns a negative number if there is an error.
func (csd *CharsetDetect) Consider(data []byte) int {
	return int(C.csd_consider(csd.ptr, (*_Ctype_char)(unsafe.Pointer(&data[0])), C.int(len(data))))
}

func (csd *CharsetDetect) Close() string {
	return C.GoString(C.csd_close(csd.ptr))
}

func GuessEncoding(data []byte) (string, error) {
	cd := CharsetDetect{}
	cd.Open()

	// loop:
	for i := 0; i < len(data); i += 512 {
		// buf := make([]byte, 512)
		e := i + 512
		if e > len(data) {
			e = len(data)
		}

		code := cd.Consider(data[i:e])
		switch {
		case code < 0:
			return "", fmt.Errorf("Detect encode failed: code %d", code)
		}
	}

	return strings.ToLower(cd.Close()), nil
}
func ConvertToUTF8From(s string, encoding string) (string, error) {
	c, err := iconv.Open("UTF-8", encoding)
	defer c.Close()
	if err != nil {
		return "", err
	}

	res, _, err := c.Conv([]byte(s), make([]byte, 512))
	if err != nil {
		return "", err
	}

	return string(res), nil
}
func ConverToUTF8(r io.Reader) (string, string, error) {
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return "", "", err
	}

	encoding, err := GuessEncoding(data)
	if err != nil {
		return "", "", err
	}
	log.Print("encodeing:", encoding)

	if encoding == "utf-8" {
		return string(data), encoding, nil
	}

	c, err := iconv.Open("UTF-8", encoding)
	defer c.Close()
	if err != nil {
		return "", "", err
	}

	res, _, err := c.Conv(data, make([]byte, 512))
	if err != nil {
		return "", "", err
	}

	return string(res), encoding, nil
}

// func GB18030ToUTF8(filename string) (string, error) {
// 	f, err := os.OpenFile(filename, os.O_RDONLY, 0666)
// 	defer f.Close()
// 	if err != nil {
// 		return "", err
// 	}

// 	c, err := iconv.Open("utf-8", "gb18030")
// 	defer c.Close()
// 	if err != nil {
// 		return "", err
// 	}

// 	data, err := ioutil.ReadAll(f)
// 	if err != nil {
// 		return "", err
// 	}

// 	res, _, err := c.Conv(data, make([]byte, 512))
// 	if err != nil {
// 		return "", err
// 	}

// 	return string(res), nil
// }
