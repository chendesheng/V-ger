package main

import (
	"log"
	"net/http"
	"net/http/cookiejar"
	"os"
	"runtime"
	"thunder"
	"time"
	"util"
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
	log.Print("log initialized.")

	if http.DefaultClient.Jar == nil {
		jar, _ := cookiejar.New(nil)
		http.DefaultClient.Jar = jar
	}

	//set timeout
	networkTimeout := time.Duration(util.ReadIntConfig("network-timeout")) * time.Second
	http.DefaultTransport.(*http.Transport).ResponseHeaderTimeout = networkTimeout

	go func() {
		err := thunder.Login()
		if err != nil {
			log.Print(err)
		}
	}()
}
