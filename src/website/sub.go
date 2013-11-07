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
	"os"
	"os/exec"
	"path"
	// "regexp"
	"strings"
	"subtitles"
	"util"
)

func init() {
	util.MakeSurePathExists(path.Join(util.ReadConfig("dir"), "subs"))
}

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

	subFileDir := path.Join(util.ReadConfig("dir"), "subs", movieName)
	util.MakeSurePathExists(subFileDir)

	subFile := path.Join(subFileDir, name)
	if ok, err := subtitles.QuickDownload(url, subFile); !ok {
		w.Write([]byte(err.Error()))
		return
	}
	println(name)
	if util.CheckExt(name, "rar", "zip") {
		cmd := exec.Command("./unar", subFile, "-f", "-o", subFileDir)

		if err := cmd.Run(); err != nil {
			w.Write([]byte(err.Error()))
		} else {
			os.Remove(subFile)
		}
	}
}
