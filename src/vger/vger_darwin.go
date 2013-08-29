package main

import (
	"download"
	"native"
	"website"
)

func main() {
	go download.Start()
	go website.Run()
	native.Start()
}
