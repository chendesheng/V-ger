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
	"path"
	"subtitles"
	"task"
	"unicode/utf8"
	"util"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket" // "regexp"
)

func subtitlesSearchHandler(w http.ResponseWriter, r *http.Request) {
	ws, err := websocket.Upgrade(w, r, nil, 1024, 1024)
	if err != nil {
		writeError(w, err)
		return
	}

	vars := mux.Vars(r)
	movie := vars["movie"]

	log.Printf("search subtitle for '%s'", movie)

	result := make(chan subtitles.Subtitle)
	var t *task.Task
	var url string
	if t, _ = task.GetTask(movie); t != nil {
		url = t.URL
	}

	var search = util.CleanMovieName(movie)
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
	vars := mux.Vars(r)
	movie := vars["movie"]

	input, _ := ioutil.ReadAll(r.Body)

	arg := make(map[string]string)

	err := json.Unmarshal(input, &arg)
	if err != nil {
		writeError(w, err)
		return
	}

	url := arg["url"]
	_, name, _, data, err := download.GetDownloadInfo(url, true)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	if name == "content" {
		name = arg["name"] + ".srt"
	}

	subFileDir := path.Join(util.ReadConfig("dir"), "subs", movie)
	err, _ = util.MakeSurePathExists(subFileDir)
	if err != nil {
		writeError(w, err)
		return
	}

	subFile := path.Join(subFileDir, name)

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
