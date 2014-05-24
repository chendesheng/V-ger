package main

import (
	"cocoa"
	"download"
	"flag"
	"os/exec"
	"util"
	// "subscribe"
	"website"
)

var debug *bool = flag.Bool("debug", false, "debug")

func main() {
	flag.Parse()

	go download.Start()
	go website.Run()
	// go subscribe.Monitor()
	if *debug {
		go func() {
			server := util.ReadConfig("server")
			cmd := exec.Command("open", "http://"+server)
			cmd.Run()
		}()
	}

	cocoa.Start()
}
