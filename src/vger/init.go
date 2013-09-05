package main

import (
	"log"
	"net/http"
	"net/http/cookiejar"
	"os"
	"runtime"
	"thunder"
	"util"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU() - 1)
	config := util.ReadAllConfigs()

	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	if logPath, ok := config["log"]; ok {
		f, err := os.OpenFile(logPath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			log.Fatal(err)
		}
		log.SetOutput(f)
	}
	log.Print("log initialized.")

	if http.DefaultClient.Jar == nil {
		jar, _ := cookiejar.New(nil)
		http.DefaultClient.Jar = jar
	}

	go func() {
		err := thunder.Login()
		if err != nil {
			log.Print(err)
		}
	}()
}
