package main

import (
	"b1"
	"download"
	"fmt"
	"os"
	"regexp"
	"shooter"
	"strings"
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
	index := strings.LastIndex(name, ".")
	if index > 0 {
		name = name[:index]
	}
	index = strings.LastIndex(name, "-")
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
func filterCategory(category string) string {
	if strings.Contains(category, "·±Ìå&Ó¢ÎÄ") {
		category = "cht&eng"
	} else if strings.Contains(category, "¼òÌå&Ó¢ÎÄ") {
		category = "chs&eng"
	} else if strings.Contains(category, "Ó¢ÎÄ") {
		category = "eng"
	} else if strings.Contains(category, "¼òÌå") {
		category = "chs"
	} else if strings.Contains(category, "·±Ìå") {
		category = "cht"
	}

	return category
}
func getMovieSub(movieName string) {
	subs, _ := getSubList(movieName, []filter{filterMovieName1, filterMovieName2})

	arr := make([]string, len(subs))
	for i, s := range subs {
		arr[i] = s.String()
	}
	i := pick(arr, "no subtitle :(")
	if i != -1 {
		selectedSub := subs[i]
		url, name := shooter.GetDownloadUrl(selectedSub.URL)
		fmt.Printf("download subtitle: %s from %s", name, url)
		download.BeginDownload(url, name, 0)

		if strings.HasSuffix(name, ".rar") || strings.HasSuffix(name, ".zip") {
			fileurls := b1.Extract(download.GetFilePath(name))
			count := 0
			for _, f := range fileurls {
				if strings.HasSuffix(f, ".srt") || strings.HasSuffix(f, ".ass") {
					fmt.Println(f)

					temp := f[:len(f)-4]
					index := strings.LastIndex(temp, ".")
					category := fmt.Sprint(count)
					if index > 0 {
						category = temp[index+1:]
						fmt.Println(category)
						category = filterCategory(category)
						if strings.Contains(category, "cht") {
							continue
						}
					}

					download.BeginDownload(f, fmt.Sprintf("%s.%s.srt", movieName, category), 0)
					count++
				}
			}
			if count > 0 {
				os.Remove(download.GetFilePath(name))
			}
		}
		if strings.HasSuffix(name, ".srt") {
			os.Rename(download.GetFilePath(name), download.GetFilePath(movieName+".srt"))
		}
		if strings.HasSuffix(name, ".ass") {
			os.Rename(download.GetFilePath(name), download.GetFilePath(movieName+".ass"))
		}
	}
}
