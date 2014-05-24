package main

import (
	// "cocoa/cookiejar"
	"dbHelper"
	"download"
	"filelock"
	"logger"
	"net/http"
	"net/http/cookiejar"
	"os"
	// "os"
	"log"
	"path"
	"runtime"
	"thunder"
	"time"
	"util"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU() - 1)
	err := os.Chdir(path.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}

	// log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	// if logPath := util.ReadConfig("log"); logPath != "" {
	// 	f, err := os.OpenFile(logPath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0666)
	// 	if err == nil {
	// 		log.SetOutput(f)
	// 		log.Print("log initialized.")
	// 	} else {
	// 		log.Print(err)
	// 	}
	// }

	logger.InitLog("V'ger", util.ReadConfig("log"))

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
		err := thunder.Login()
		if err != nil {
			log.Print(err)
		}
	}()

	dbHelper.Init("sqlite3", path.Join(util.ReadConfig("dir"), "vger.db"))

	filelock.DefaultLock, _ = filelock.New("/tmp/vger.db.lock.txt")

	download.BaseDir = util.ReadConfig("dir")
	download.NetworkTimeout = time.Duration(util.ReadIntConfig("network-timeout")) * time.Second
}
