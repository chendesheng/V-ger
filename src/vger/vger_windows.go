package main

import (
	"download"
	// "native"
	"runtime"
	"website"
)

func main() {
	go download.Start()
	website.Run()
}
