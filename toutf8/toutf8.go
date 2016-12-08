package toutf8

/*
#cgo LDFLAGS: -lcharsetdetect -liconv
#include "charsetdetect.h"
*/
import "C"
import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"strings"
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
	c, err := Open("UTF-8", encoding)
	defer c.Close()
	if err != nil {
		return "", err
	}

	return c.ConvStr(s)
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

	c, err := Open("UTF-8", encoding)
	defer c.Close()
	if err != nil {
		return "", "", err
	}

	res, err := c.ConvStr(string(data))
	if err != nil {
		return "", "", err
	}


	return res, encoding, nil
}

func GB18030ToUTF8(text string) (res string, err error) {
	res = text

	var encoding string
	encoding, err = GuessEncoding([]byte(text))
	if err != nil {
		return
	}

	if encoding != "utf-8" && encoding != "ascii" {
		var utf8text string
		utf8text, err = ConvertToUTF8From(text, "gb18030")
		if err == nil {
			res = utf8text
		}
	}

	return
}
