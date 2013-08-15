package main

import (
	"download"
	// "native"
	"runtime"
	"website"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU() - 1)
	go download.Start()
	website.Run()
	// native.Start()
}
