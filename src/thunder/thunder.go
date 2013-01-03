package thunder

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	// "net/http/httputil"
	"net/url"
	"strings"
	"time"
	// "encoding/json"
	// "regexp"
	// "io"
	// "os"
)

var Client *http.Client

// func pipe() {
// 	f, err := os.OpenFile("test.log", os.O_APPEND|os.O_CREATE, 0666)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	r, w := io.Pipe()

// }

func NewTask(taskURL string) []ThunderTask {
	taskType := getTaskType(taskURL)
	userId := getCookieValue("userid")

	if taskType == 4 {
		btTaskCommit(userId, taskURL)
	} else {
		taskCommit(userId, taskURL, taskType)
	}
	return getNewlyCreateTask(userId)
}
func taskCommit(userId string, taskURL string, taskType int) {
	text := sendGet("http://dynamic.cloud.vip.xunlei.com/interface/task_check",
		&url.Values{
			"callback": {"fun"},
			"url":      {taskURL},
		})

	cid, gcid, size, t := parseTaskCheck(text)

	sendGet("http://dynamic.cloud.vip.xunlei.com/interface/task_commit",
		&url.Values{
			"callback":   {"ret_task"},
			"uid":        {userId},
			"cid":        {cid},
			"gcid":       {gcid},
			"size":       {size},
			"goldbean":   {"0"},
			"silverbean": {"0"},
			"t":          {t},
			"url":        {taskURL},
			"type":       {fmt.Sprintf("%d", taskType)},
			"o_page":     {"history"},
			"o_taskid":   {"0"},
			"class_id":   {"0"},
			"database":   {"undefined"},
			"time":       {time.Now().String()},
		})
}
func btTaskCommit(userId string, taskURL string) {

	text := sendGet("http://dynamic.cloud.vip.xunlei.com/interface/url_query", &url.Values{
		"u":        {taskURL},
		"callback": {"queryUrl"},
	})

	cid, tsize, btname, size, findex := parseUrlQueryResult(text)

	sendPost("http://dynamic.cloud.vip.xunlei.com/interface/bt_task_commit",
		&url.Values{
			"callback": {"jsonp"},
			"t":        {time.Now().String()},
		},
		&url.Values{
			"uid":        {userId},
			"cid":        {cid},
			"tsize":      {tsize},
			"goldbean":   {"0"},
			"silverbean": {"0"},
			"btname":     {btname},
			"size":       {size},
			"findex":     {findex},
			"o_page":     {"task"},
			"o_taskid":   {"0"},
			"class_id":   {"0"},
		})
}
func getNewlyCreateTask(userId string) []ThunderTask {
	text := sendGet("http://dynamic.cloud.vip.xunlei.com/interface/showtask_unfresh",
		&url.Values{
			"callback": {"jsonp1"},
			"t":        {time.Now().String()},
			"type_id":  {"4"},
			"page":     {"1"},
			"tasknum":  {"1"},
		})

	// regexUrl, _ := regexp.Compile(`("id":"[0-9]*").*("cid":"[^"]*").*("lixian_url":"[^"]*")`)
	// finalJson := "{" + strings.Join(regexUrl.FindStringSubmatch(text)[1:], ",") + "}"
	// fmt.Println(finalJson)
	// var r interface{}
	// json.Unmarshal([]byte(finalJson), &r)
	// info := r.(map[string]interface{})
	// fmt.Println(info["lixian_url"])

	info := parseNewlyCreateTask(text)

	if info["lixian_url"] != "" {
		file := ThunderTask{
			Name:        info["taskname"].(string),
			DownloadURL: info["lixian_url"].(string),
			Size:        info["filesize"].(string),
		}
		files := [1]ThunderTask{file}

		return files[0:]
	}

	return getBtTaskList(userId, info["id"].(string), info["cid"].(string))
}
func getBtTaskList(userId string, id string, cid string) []ThunderTask {
	text := sendGet("http://dynamic.cloud.vip.xunlei.com/interface/fill_bt_list",
		&url.Values{
			"uid":      {userId},
			"callback": {"fill_bt_list"},
			"t":        {time.Now().String()},
			"tid":      {id},
			"infoid":   {cid},
			"p":        {"1"},
		})
	return parseBtTaskList(text)
}

func getCookieValue(name string) string {
	url, _ := url.Parse("http://vip.lixian.xunlei.com")
	for _, c := range Client.Jar.Cookies(url) {
		if c.Name == name {
			return c.Value
		}
	}

	return ""
}
func getTaskType(url string) int {
	if strings.Index(url, "magnet:") != -1 {
		return 4
	} else if strings.Index(url, "ed2k://") != -1 {
		return 2
	}
	return 0
}
func sendPost(url string, params *url.Values, data *url.Values) string {
	if params != nil {
		url = url + "?" + params.Encode()
	}
	resp, err := Client.PostForm(url, *data)
	if err != nil {
		log.Fatal(err)
	}

	text := readBody(resp)
	fmt.Println(text)
	return text
}
func sendGet(url string, params *url.Values) string {
	if params != nil {
		url = url + "?" + params.Encode()
	}
	resp, err := Client.Get(url)
	if err != nil {
		log.Fatal(err)
	}

	// dumpBytes, _ := httputil.DumpResponse(resp, true)
	// fmt.Println(string(dumpBytes))

	text := readBody(resp)
	return text
}
func readBody(resp *http.Response) string {
	bytes, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		log.Fatal(err)
	}

	// dumpBytes, _ := httputil.DumpResponse(resp, true)
	// fmt.Println(string(dumpBytes))

	text := string(bytes)
	return text
}
