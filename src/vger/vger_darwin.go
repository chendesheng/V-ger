package main

import (
	"download"
	// "native"
	"website"
)

func main() {
	go download.Start()
	website.Run()
	// native.Start()
}
