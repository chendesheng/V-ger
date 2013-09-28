package main

import (
	"cocoa"
	"download"
	"website"
)

func main() {
	go download.Start()
	go website.Run()
	cocoa.Start()
}
