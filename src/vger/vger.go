package main

import (
	"native"
	"runtime"
	"website"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU() - 1)
	go website.Run()
	native.Start()
}
