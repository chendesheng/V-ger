package util

import (
	"io/ioutil"
	"os"
	// "path"
	"strings"
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

func CheckExt(filename string, ext string) bool {
	return filename[strings.LastIndex(filename, "."):] == ("." + ext)
}
