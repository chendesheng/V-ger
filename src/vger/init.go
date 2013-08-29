package main

import (
	"log"
	"os"
	"util"
	"runtime"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU() - 1)

	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	if logPath := util.ReadConfig("log"); logPath != "" {
		f, err := os.OpenFile(logPath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			log.Fatal(err)
		}
		log.SetOutput(f)
	}
}
