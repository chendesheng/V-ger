package util

import (
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"
)

func MakeSurePathExists(path string) (error, bool) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return os.Mkdir(path, 0777), false
	}

	return nil, true
}

func IsPathExists(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	} else {
		return true
	}
}

func CheckExt(filename string, exts ...string) bool {
	ext := path.Ext(filename)[1:]

	for _, e := range exts {
		if e == ext {
			return true
		}
	}

	return false
}

func GetFileSize(file string) int64 {
	if f, err := os.OpenFile(file, os.O_RDONLY, 0666); err == nil {
		defer f.Close()
		if info, err := f.Stat(); err == nil {
			return info.Size()
		}

	}
	return 0
}

func EmulateFiles(dir string, fn func(string), exts ...string) {
	f, err := os.Open(dir)
	defer f.Close()

	if err != nil {
		log.Print(err)
		return
	}

	files, err := f.Readdir(-1)
	if err != nil {
		log.Print(err)
		return
	}

	for _, file := range files {
		name := file.Name()
		if file.IsDir() {
			EmulateFiles(path.Join(dir, name), fn, exts...)
		} else if CheckExt(name, exts...) {
			fn(path.Join(dir, name))
		}
	}
}

func extractOneFile(unarPath, filename string) {
	dir := path.Dir(filename)
	cmd := exec.Command(unarPath, filename, "-f", "-o", dir)

	if err := cmd.Run(); err != nil {
		log.Print(err)
	} else {
		os.Remove(filename)
	}
}
func Extract(unarPath string, filename string) {
	if CheckExt(filename, "rar", "zip") {
		extractOneFile(unarPath, filename)

		dir := filename[:len(filename)-len(path.Ext(filename))]

		infoes, err := ioutil.ReadDir(dir)
		if err != nil {
			log.Print(err)
			return
		}

		for _, f := range infoes {
			filename := strings.ToLower(path.Join(dir, f.Name()))
			extractOneFile(unarPath, filename)
		}
	}
}
