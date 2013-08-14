package main

import (
	"download"
	"native"
	"runtime"
	"website"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU() - 1)
	go website.Run()
	go download.Start()
	native.Start()
}
