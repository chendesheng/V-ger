package main

import (
	"cocoa"
	"download"
	"subscribe"
	"website"
)

func main() {
	go download.Start()
	go website.Run()
	go subscribe.Monitor()
	cocoa.Start()
}
