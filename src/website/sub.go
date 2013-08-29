package website

import (
	"code.google.com/p/go.net/websocket"
	"download"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	// "os/exec"
	"path"
	// "regexp"
	"strings"
	"subtitles"
	"util"
)

func subtitlesSearchHandler(ws *websocket.Conn) {
	r := ws.Request()
	movieName, _ := url.QueryUnescape(r.URL.String()[strings.LastIndex(r.URL.String(), "/")+1:])
	log.Printf("search subtitle for '%s'", movieName)

	result := make(chan subtitles.Subtitle)
	go subtitles.SearchSubtitles(util.CleanMovieName(movieName), result)

	for s := range result {
		log.Printf("%v", s)
		text, _ := json.Marshal(s)
		io.WriteString(ws, string(text))
	}
}

func subtitlesDownloadHandler(w http.ResponseWriter, r *http.Request) {
	movieName, _ := url.QueryUnescape(r.URL.String()[strings.LastIndex(r.URL.String(), "/")+1:])
	input, _ := ioutil.ReadAll(r.Body)
	url := string(input)
	url, name, _, err := download.GetDownloadInfo(url)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	if ok, err := subtitles.QuickDownload(url, path.Join(util.ReadConfig("dir"), movieName+path.Ext(name))); !ok {
		w.Write([]byte(err.Error()))
		return
	}
}
