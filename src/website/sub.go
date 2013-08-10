package website

import (
	"download"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	// "os/exec"
	"path"
	"regexp"
	"strings"
	"subtitles"
	"util"
)

func subtitlesSearchHandler(w http.ResponseWriter, r *http.Request) {
	print("subtitlesSearchHandler")
	movieName, _ := url.QueryUnescape(r.URL.String()[strings.LastIndex(r.URL.String(), "/")+1:])
	data := subtitles.SearchSubtitles(util.CleanMovieName(movieName))
	text, _ := json.Marshal(data)
	w.Write([]byte(text))
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

	if ok, err := subtitles.QuickDownload(url, path.Join(config["dir"], movieName+path.Ext(name))); !ok {
		w.Write([]byte(err.Error()))
		return
	}
	// else {
	// 	cmd := exec.Command("open", path.Join(config["dir"], movieName+path.Ext(name)))
	// 	cmd.Start()
	// 	// extractSubtitle(name, movieName)
	// 	cleanFileName(movieName)
	// }
}

func cleanFileName(movieName string) {
	fileInfoes, err := ioutil.ReadDir(path.Join(config["dir"], movieName))
	if err != nil {
		log.Print(err)
		return
	}

	for _, f := range fileInfoes {
		name := f.Name()

		regSubName := regexp.MustCompile(".*[.](srt|ass)$")

		if !regSubName.Match([]byte(name)) {
			continue
		}

		print(name)
	}

}
