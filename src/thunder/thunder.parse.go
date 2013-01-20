package thunder

import (
	"encoding/json"
	"fmt"
	"log"
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
func parseUrlQueryResult(text string) (cid string, tsize string, btname string, size string, findex string) {
	regexUrlQuery, _ := regexp.Compile(`queryUrl\((-?[0-9]*),'([^']*)','([^']*)','([^']*)','([^']*)',new Array\((.+)\),new Array\(([^)]*)\),new Array\(([^)]*)\),new Array\(([^)]*)\),new Array\(([^)]*)\),new Array\(([^)]*)\),'[^']*'\)`)

	args := make([]string, 0, 10)

	for _, s := range regexUrlQuery.FindStringSubmatch(text)[1:] {
		args = append(args, strings.TrimSpace(s))
	}
	log.Println(args)

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
func parseBtTaskList(text string) []ThunderTask {
	regexUrl, _ := regexp.Compile(`"Record":(\[.+\])`)

	jsonStr := regexUrl.FindStringSubmatch(text)[1]
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

	return res
}

func parseNewlyCreateTask(text string) map[string]interface{} {
	regexUrl, _ := regexp.Compile(`("id":"[0-9]*").*("filesize":"[^"]*").*("cid":"[^"]*").*("taskname":"[^"]*").*("lixian_url":"[^"]*")`)

	jsonStr := "{" + strings.Join(regexUrl.FindStringSubmatch(text)[1:], ",") + "}"
	log.Println(jsonStr)

	var r interface{}
	json.Unmarshal([]byte(jsonStr), &r)

	return r.(map[string]interface{})
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
