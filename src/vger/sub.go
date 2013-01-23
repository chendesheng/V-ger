package main

import (
	"regexp"
	"shooter"
	"strings"
	"fmt"
	// "os"
	"download"
)

type filter func(name string) string

func filterMovieName2(name string) string {
	name = filterMovieName1(name)
	reg, _ := regexp.Compile("(?i)720p|x[.]264|BluRay|DTS|x264|1080p|H[.]264|AC3|[.]ENG|[.]BD|Rip|H264|HDTV|-IMMERSE|-DIMENSION|xvid|[[]PublicHD[]]|[.]Rus|Chi_Eng|DD5[.]1|HR-HDTV|[.]AAC|[0-9]+x[0-9]+|blu-ray|Remux|dxva|dvdscr")
	name = string(reg.ReplaceAll([]byte(name), []byte("")))
	name = strings.Replace(name, ".", " ", -1)
	name = strings.TrimSpace(name)
	return name
}
func filterMovieName1(name string) string {
	name = name[:strings.LastIndex(name, ".")]
	index := strings.LastIndex(name, "-")
	if index > 0 {
		name = name[:index]
	}
	return name
}
func getSubList(movieName string, filters []filter) ([]shooter.Subtitle, string) {
	for _, f := range filters {
		name := f(movieName)
		fmt.Printf("searching subtitles for \"%s\"...\n", name)
		subs := shooter.SearchSubtitles(name)
		if len(subs) > 0 {
			return subs, name
		}
	}

	return make([]shooter.Subtitle, 0), movieName
}

func getMovieSub(movieName string) {
	subs, movieName := getSubList(movieName, []filter{filterMovieName1, filterMovieName2})

	arr := make([]string, len(subs))
	for i, s := range subs {
		arr[i] = s.String()
	}
	i := pick(arr, "no subtitle :(")
	if i != -1 {
		selectedSub := subs[i]
		url, name := shooter.GetDownloadUrl(selectedSub.URL)
		fmt.Printf("download subtitle: %s from %s", name, url)
		download.BeginDownload(url, name)
	}
}
