package website

import (
	// "bufio"
	"bytes"
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
	"task"
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
	var t *task.Task
	var url string
	if t, _ = task.GetTask(movieName); t != nil {
		url = t.URL
	}
	go subtitles.SearchSubtitles(util.CleanMovieName(movieName), url, result)

	for s := range result {
		log.Printf("%v", s)
		text, _ := json.Marshal(s)
		io.WriteString(ws, string(text))
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
		name = arg["name"]
	}

	subFileDir := path.Join(util.ReadConfig("dir"), "subs", movieName)
	util.MakeSurePathExists(subFileDir)

	subFile := path.Join(subFileDir, name)

	data, err := subtitles.QuickDownload(url)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	data = bytes.Replace(data, []byte{'+'}, []byte{' '}, -1)

	ioutil.WriteFile(subFile, data, 0666)

	if util.CheckExt(name, "rar", "zip") {
		cmd := exec.Command("./unar", subFile, "-f", "-o", subFileDir)

		if err := cmd.Run(); err != nil {
			w.Write([]byte(err.Error()))
		} else {
			os.Remove(subFile)
		}
	}
}
