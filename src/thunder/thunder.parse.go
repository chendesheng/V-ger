package thunder

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	// "os"
	"regexp"
	"strconv"
	"strings"
)

type ThunderTask struct {
	Name        string
	DownloadURL string
	PlayURL     string
	Size        string
	Percent     int
}

func (t *ThunderTask) String() string {
	return fmt.Sprintf("%s  %s %d%%\n", t.Name, t.Size, t.Percent)
}
func getError(text string) error {
	regexAlert := regexp.MustCompile(`alert\('(.+)'\)`)
	result := regexAlert.FindStringSubmatch(text)
	if len(result) > 0 {
		return errors.New(fmt.Sprintf("Thunder server error:%s.\n", result[1]))
	}

	return errors.New(fmt.Sprintln("Unknown thunder server error"))
}
func parseUrlQueryResult(text string) (cid string, tsize string, btname string, size string, findex string) {
	regexUrlQuery := regexp.MustCompile(`queryUrl\((-?[0-9]*),'([^']*)','([^']*)','([^']*)','([^']*)',new Array\((.+)\),new Array\(([^)]*)\),new Array\(([^)]*)\),new Array\(([^)]*)\),new Array\(([^)]*)\),new Array\(([^)]*)\),'[^']*','[^']*'\)`)

	args := make([]string, 0, 10)

	text = strings.Replace(text, "\\'", "", -1) //ensure no ' in side ''
	matches := regexUrlQuery.FindStringSubmatch(text)

	if matches == nil {
		panic(fmt.Errorf("Parse unexpected response: %s", text))
	}

	for _, s := range matches[1:] {
		args = append(args, strings.TrimSpace(s))
	}

	sizeList := strings.Split(args[7], ",")
	trimStringSlice(&sizeList, " '")

	selectionList := make([]string, 0, 10)
	selectionSizeList := make([]string, 0, 10)
	for i, s := range strings.Split(args[8], ",") {
		if strings.Trim(s, " '") == "1" {
			selectionList = append(selectionList, strconv.Itoa(i))
			selectionSizeList = append(selectionSizeList, sizeList[i])
		}
	}

	cid = args[1]
	tsize = args[2]
	btname = args[3]
	size = strings.Join(selectionSizeList, "_")
	findex = strings.Join(selectionList, "_")

	return
}
func parseBtTaskList(text string) ([]ThunderTask, error) {
	regexUrl, _ := regexp.Compile(`"Record":(\[.+\])`)

	result := regexUrl.FindStringSubmatch(text)
	if len(result) == 0 {
		return nil, getError(text)
	}

	jsonStr := result[1]
	log.Println(jsonStr)

	var r []interface{}
	json.Unmarshal([]byte(jsonStr), &r)

	res := make([]ThunderTask, 0, len(r))

	for _, item := range r {
		t := item.(map[string]interface{})
		percent := t["percent"].(float64)
		res = append(res, ThunderTask{
			Name:        t["title"].(string),
			Size:        t["size"].(string),
			Percent:     int(percent),
			DownloadURL: t["downurl"].(string),
		})
	}

	return res, nil
}

func parseNewlyCreateTask(text string) map[string]interface{} {
	regexUrl := regexp.MustCompile(`("id":"[0-9]*").*("filesize":"[^"]*").*("cid":"[^"]*").*("taskname":"[^"]*").*("lixian_url":"[^"]*")`)

	if matches := regexUrl.FindStringSubmatch(text); matches != nil {

		jsonStr := "{" + strings.Join(matches[1:], ",") + "}"
		log.Println(jsonStr)

		var r interface{}
		json.Unmarshal([]byte(jsonStr), &r)
		return r.(map[string]interface{})
	} else {
		log.Println("parseNewlyCreateTask")
		panic(fmt.Errorf("Parse unexpected response: %s", text))
	}
}
func parseTaskCheck(text string) (cid string, gcid string, size string, t string) {
	args := parseJsFuncArgs(text)

	cid = args[0]
	gcid = args[1]
	size = args[2]
	t = args[4]

	return
}
func parseJsFuncArgs(text string) []string {
	regex, _ := regexp.Compile(`\(([^)]+)\)`)
	args := strings.Split(regex.FindStringSubmatch(text)[1], ",")

	trimStringSlice(&args, " '")

	return args
}
func trimStringSlice(strs *[]string, cutset string) {
	for i, s := range *strs {
		(*strs)[i] = strings.Trim(s, cutset)
	}
}
func parseUploadTorrentResutl(text string) (map[string]interface{}, error) {
	i := strings.Index(text, "var btResult =")
	j := strings.LastIndex(text, ";")
	s := 14
	if i == -1 {
		i = strings.Index(text, "edit_bt_list(")
		j = strings.LastIndex(text, "}")
		j = strings.LastIndex(text[:j], "}") + 1
		s = 13
	}
	if i == -1 {
		return nil, errors.New("Unknown upload .torrent file result.")
	}
	text = text[i+s : j]
	fmt.Println(text)
	res := make(map[string]interface{})
	json.Unmarshal([]byte(text), &res)
	return res, nil
}
