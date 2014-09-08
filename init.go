package main

import (

	// "cocoa/cookiejar"
	"log"
	"net/http"
	"net/http/cookiejar"
	"os"
	"path"
	"runtime"
	"time"
	"vger/dbHelper"
	"vger/download"
	"vger/filelock"
	"vger/logger"
	"vger/thunder"
	"vger/util"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU() - 1)
	err := os.Chdir(path.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}

	logbase := util.ReadConfig("log")
	logger.InitLog("[V'ger]", path.Join(logbase, "vger.log"))

	// http.DefaultClient.Jar = &cookiejar.SafariCookieJar{}
	jar, _ := cookiejar.New(nil)
	http.DefaultClient.Jar = jar

	util.SaveConfig("shutdown-after-finish", "false")

	//set timeout
	networkTimeout := time.Duration(util.ReadIntConfig("network-timeout")) * time.Second
	transport := http.DefaultTransport.(*http.Transport)
	transport.ResponseHeaderTimeout = networkTimeout
	transport.MaxIdleConnsPerHost = 3

	go func() {
		err := thunder.Login(nil)
		if err != nil {
			log.Print(err)
		}
	}()

	dbHelper.Init("sqlite3", path.Join(util.ReadConfig("dir"), "vger.db"))

	filelock.DefaultLock, _ = filelock.New("/tmp/vger.db.lock.txt")

	download.BaseDir = util.ReadConfig("dir")
	download.NetworkTimeout = networkTimeout
}
