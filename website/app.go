package website

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/gorilla/mux"
	// _ "net/http/pprof"
	"bytes"
	"os"
	"runtime"
	"runtime/debug"
	"runtime/pprof"
	"time"
)

func appStatusHandler(w http.ResponseWriter, r *http.Request) {
	n := runtime.NumGoroutine()

	w.Header()["content-type"] = []string{"text/html"}
	w.Write([]byte(`
		<html>
			<head>
				<link href="/assets/style.css" rel="stylesheet" type="text/css">
				<style>
					body {
						padding: 20px 80px;
					}
				</style>
			</head>
		<body>`))

	w.Write([]byte("<a href=\"/app/shutdown\">Shutdown</a><br><br>"))

	w.Write([]byte(fmt.Sprintf("# of goruntines: %d.<br><br>", n)))

	buffer := new(bytes.Buffer)

	pprof.Lookup("goroutine").WriteTo(buffer, 2)
	buffer.Write([]byte("\n"))
	pprof.Lookup("heap").WriteTo(buffer, 2)
	buffer.Write([]byte("\n"))
	pprof.Lookup("block").WriteTo(buffer, 2)

	for {
		if line, err := buffer.ReadBytes('\n'); err == nil {

			w.Write(line)
			w.Write([]byte("<br>"))
		} else {
			break
		}
	}

	w.Write([]byte("</body></html>"))
}

func appShutdownHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("bye"))
	go func() {
		time.Sleep(time.Second)
		os.Exit(0)
	}()
}

func appGCHandler(w http.ResponseWriter, r *http.Request) {
	debug.FreeOSMemory()
}

func appCookieHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	u, _ := url.Parse("http://" + vars["domain"])
	cookies := http.DefaultClient.Jar.Cookies(u)
	writeJson(w, cookies)
}
