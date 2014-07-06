package website

import (

	// "bufio"
	"bytes"
	"download"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"path"
	"task"
	"github.com/gorilla/websocket"

	// "regexp"
	"strings"
	"subtitles"
	"unicode/utf8"
	"util"
)

func subtitlesSearchHandler(w http.ResponseWriter, r *http.Request) {
	ws, err := websocket.Upgrade(w, r, nil, 1024, 1024)
	if err != nil {
		writeError(w, err)
		return
	}

	movieName, _ := url.QueryUnescape(r.URL.String()[strings.LastIndex(r.URL.String(), "/")+1:])
	log.Printf("search subtitle for '%s'", movieName)

	result := make(chan subtitles.Subtitle)
	var t *task.Task
	var url string
	if t, _ = task.GetTask(movieName); t != nil {
		url = t.URL
	}

	var search = util.CleanMovieName(movieName)
	if t != nil && len(t.Subscribe) != 0 && t.Season > 0 {
		search = fmt.Sprintf("%s s%2de%2d", t.Subscribe, t.Season, t.Episode)
	}

	go subtitles.SearchSubtitles(search, url, result, nil)

	for s := range result {
		log.Printf("%v", s)
		ws.WriteJSON(s)
	}
}

func subtitlesDownloadHandler(w http.ResponseWriter, r *http.Request) {
	movieName, _ := url.QueryUnescape(r.URL.String()[strings.LastIndex(r.URL.String(), "/")+1:])
	input, _ := ioutil.ReadAll(r.Body)

	arg := make(map[string]string)

	println(string(input))
	json.Unmarshal(input, &arg)
	println(arg["url"], arg["name"])

	url := arg["url"]
	url, name, _, err := download.GetDownloadInfo(url)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	if name == "content" {
		name = arg["name"] + ".srt"
	}

	subFileDir := path.Join(util.ReadConfig("dir"), "subs", movieName)
	util.MakeSurePathExists(subFileDir)

	subFile := path.Join(subFileDir, name)

	data, err := subtitles.QuickDownload(url)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	if util.CheckExt(name, ".rar", ".zip") {
		ioutil.WriteFile(subFile, data, 0666)
		util.Extract("./unar", subFile)
	} else {
		data = bytes.Replace(data, []byte{'+'}, []byte{' '}, -1)

		spaceBytes := make([]byte, 4)
		n := utf8.EncodeRune(spaceBytes, '　')
		spaceBytes = spaceBytes[:n]
		data = bytes.Replace(data, spaceBytes, []byte{' '}, -1)

		data = bytes.Replace(data, []byte{'\\', 'N'}, []byte{'\n'}, -1)

		ioutil.WriteFile(subFile, data, 0666)
	}
}
