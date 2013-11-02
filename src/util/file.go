package util

import (
	"io/ioutil"
	"os"
	"path"
	// "strings"
)

func MakeSurePathExists(path string) error {
	if _, err := ioutil.ReadDir(path); os.IsNotExist(err) {
		return os.Mkdir(path, 0777)
	}

	return nil
}

func IsPathExists(path string) bool {
	if _, err := ioutil.ReadDir(path); os.IsNotExist(err) {
		return false
	} else {
		return true
	}
}

func CheckExt(filename string, exts ...string) bool {
	ext := path.Ext(filename)[1:]
	println("ext:", ext)

	for _, e := range exts {
		if e == ext {
			return true
		}
	}

	return false
}
