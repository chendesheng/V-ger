package util

import (
	"regexp"
	"strings"
)

func CleanMovieName(name string) string {
	return cleanMovieName2(name)
}
func CleanMovieNameWithMaxLen(name string, maxLen int) string {
	name = cleanMovieName2(name)
	if len(name) < maxLen {
		return name
	} else {
		return name[:maxLen]
	}
}

var regCleanName *regexp.Regexp = regexp.MustCompile("(?i)720p|[[]720p[]]|x[.]264|BluRay|DTS|x264|1080p|H[.]264|AC3|[.]ENG|[.]BD|Rip|BRRip|H264|HDTV|-IMMERSE|-DIMENSION|xvid|[[]PublicHD[]]|[.]Rus|Chi_Eng|DD5[.]1|HR-HDTV|[.]AAC|[0-9]+x[0-9]+|blu-ray|Remux|dxva|dvdscr|WEB-DL|[_.]")

func cleanMovieName2(name string) string {
	name = cleanMovieName1(name)
	name = string(regCleanName.ReplaceAll([]byte(name), []byte("")))
	index := strings.LastIndex(name, "..")
	if index > 0 {
		name = name[:index]
	}
	name = strings.Replace(name, ".", " ", -1)
	name = strings.TrimSpace(name)

	return name
}
func cleanMovieName1(name string) string {
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
