package website

import (
	"b1"
	"download"
	"fmt"
	"net/url"
	"os"
	"path"
	"regexp"
	"strings"
	"subtitles"
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
func getSubList(movieName string, filters []filter) []subtitles.Subtitle {
	for _, f := range filters {
		name := f(movieName)
		fmt.Printf("searching subtitles for \"%s\"...\n", name)
		subs := subtitles.SearchSubtitles(name)
		if len(subs) > 0 {
			return subs
		}
	}

	return make([]subtitles.Subtitle, 0)
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

func GetMovieSub(movieName string) {
	getMovieSub(movieName)
}
func getFileName(fullURL string) string {
	e := strings.Index(fullURL, "?")
	if e < 0 {
		e = len(fullURL)
	}
	name, _ := url.QueryUnescape(fullURL[strings.LastIndex(fullURL, `/`)+1 : e])
	return name
}
func getMovieSub(movieName string) {
	subs := getSubList(movieName, []filter{filterMovieName1, filterMovieName2})

	arr := make([]string, len(subs))
	for i, s := range subs {
		arr[i] = s.String()
	}
	i, _ := pick(arr, "no subtitle :(")
	if i != -1 {
		selectedSub := subs[i]
		url := selectedSub.URL
		fmt.Printf("download subtitle: %s", url)
		name := getFileName(url)
		if ok, err := subtitles.QuickDownload(url, path.Join(download.BaseDir, name)); !ok {
			print(err)
			return
		}

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

					subtitles.QuickDownload(f, fmt.Sprintf("%s%c%s.%s.srt", download.BaseDir, os.PathSeparator, movieName, category))

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

func extractSubtitle(name, movieName string) {
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

				subtitles.QuickDownload(f, fmt.Sprintf("%s%c%s.%s.srt", download.BaseDir, os.PathSeparator, movieName, category))

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
