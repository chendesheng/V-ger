package main

import (
	"github.com/nightlyone/lockfile"
	"log"
	"net/http"
	"net/http/cookiejar"
	"os"
	"path"
	"player/shared"
	"runtime"
	"subscribe"
	"task"
	"thunder"
	"time"
	"util"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU() - 1)

	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	if logPath := util.ReadConfig("log"); logPath != "" {
		f, err := os.OpenFile(logPath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0666)
		if err == nil {
			log.SetOutput(f)
			log.Print("log initialized.")
		} else {
			log.Print(err)
		}
	}

	if http.DefaultClient.Jar == nil {
		jar, _ := cookiejar.New(nil)
		http.DefaultClient.Jar = jar
	}

	util.SaveConfig("shutdown-after-finish", "false")

	//set timeout
	networkTimeout := time.Duration(util.ReadIntConfig("network-timeout")) * time.Second
	transport := http.DefaultTransport.(*http.Transport)
	transport.ResponseHeaderTimeout = networkTimeout
	transport.MaxIdleConnsPerHost = 3

	go func() {
		err := thunder.Login()
		if err != nil {
			log.Print(err)
		}
	}()

	task.TaskDir = path.Join(util.ReadConfig("dir"), "vger.db")
	subscribe.DbPath = task.TaskDir
	shared.DbFile = task.TaskDir

	//only block when file locked by another process
	lockfile.DefaultLock, _ = lockfile.New("/tmp/vger.db.lock.txt")
}
