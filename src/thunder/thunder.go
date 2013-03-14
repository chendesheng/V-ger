package thunder

import (
	"bytes"
	// "errors"
	"fmt"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"
	// "encoding/json"
	// "regexp"
	// "io"
	// "os"
)

var Client *http.Client

func NewTask(taskURL string) ([]ThunderTask, error) {
	log.Println("thunder new task: ", taskURL)

	taskType := getTaskType(taskURL)
	userId := getCookieValue("userid")

	if taskType == 4 {
		if err := btTaskCommit(userId, taskURL); err != nil {
			return nil, err
		}
	} else if taskType == 1 {
		torrent, err := quickDownload(taskURL)
		if err != nil {
			return nil, err
		}

		text, err := uploadTorrentFile(torrent)
		if err != nil {
			return nil, err
		}

		result, err := parseUploadTorrentResutl(text)
		if err != nil {
			return nil, err
		}
		ret_value := result["ret_value"].(float64)
		if ret_value == 0 {
			return nil, fmt.Errorf("Upload torrent file: Can't find files.")
		}

		btsize := int64(result["btsize"].(float64))
		infoid := result["infoid"].(string)
		ftitle := result["ftitle"].(string)

		filelist := result["filelist"].([]interface{})
		selectionList := make([]string, 0)
		sizelist := make([]string, 0)
		for _, f := range filelist {
			item := f.(map[string]interface{})
			if item["valid"].(float64) == 1 {
				selectionList = append(selectionList, item["findex"].(string))
				sizelist = append(sizelist, item["subsize"].(string))
			}
		}

		findex := strings.Join(selectionList, "_")
		size := strings.Join(sizelist, "_")

		sendPost("http://dynamic.cloud.vip.xunlei.com/interface/bt_task_commit",
			&url.Values{
				"callback": {"jsonp"},
				"t":        {time.Now().String()},
			},
			&url.Values{
				"uid":        {userId},
				"cid":        {infoid},
				"tsize":      {fmt.Sprint(btsize)},
				"goldbean":   {"0"},
				"silverbean": {"0"},
				"btname":     {ftitle},
				"size":       {size},
				"findex":     {findex},
				"o_page":     {"task"},
				"o_taskid":   {"0"},
				"class_id":   {"0"},
			})
	} else {
		if err := taskCommit(userId, taskURL, taskType); err != nil {
			return nil, err
		}
	}
	return getNewlyCreateTask(userId)
}
func taskCommit(userId string, taskURL string, taskType int) error {
	text := sendGet("http://dynamic.cloud.vip.xunlei.com/interface/task_check",
		&url.Values{
			"callback": {"fun"},
			"url":      {taskURL},
		})

	cid, gcid, size, t := parseTaskCheck(text)
	if cid == "" {
		return fmt.Errorf("Commit task error, try again later")
	}

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

	return nil
}
func uploadTorrentFile(torrent []byte) (string, error) {
	url := "http://dynamic.cloud.vip.xunlei.com/interface/torrent_upload"
	resp, err := postFile("a.torrent", torrent, url)
	if err == nil {
		defer resp.Body.Close()
		bytes, _ := ioutil.ReadAll(resp.Body)
		text := string(bytes)
		return text, nil
	}

	return "", err
}
func btTaskCommit(userId string, taskURL string) error {

	text := sendGet("http://dynamic.cloud.vip.xunlei.com/interface/url_query", &url.Values{
		"u":        {taskURL},
		"callback": {"queryUrl"},
	})

	cid, tsize, btname, size, findex := parseUrlQueryResult(text)

	if cid == "" {
		return fmt.Errorf("Commit bt task error, try again later.")
	}

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

	return nil
}
func getNewlyCreateTask(userId string) ([]ThunderTask, error) {
	text := sendGet("http://dynamic.cloud.vip.xunlei.com/interface/showtask_unfresh",
		&url.Values{
			"callback": {"jsonp1"},
			"t":        {time.Now().String()},
			"type_id":  {"4"},
			"page":     {"1"},
			"tasknum":  {"1"},
		})

	info := parseNewlyCreateTask(text)

	if info["lixian_url"] != "" {
		return []ThunderTask{
			ThunderTask{
				Name:        info["taskname"].(string),
				DownloadURL: info["lixian_url"].(string),
				Size:        info["filesize"].(string),
				Percent:     100,
			},
		}, nil
	}

	return getBtTaskList(userId, info["id"].(string), info["cid"].(string))
}
func getBtTaskList(userId string, id string, cid string) ([]ThunderTask, error) {
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
	} else if strings.Index(url, ".torrent") != -1 {
		return 1
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
	log.Println(text)
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

	text := readBody(resp)

	log.Println("Request Get:", url)
	log.Println("Response:", text)

	return text
}
func readBody(resp *http.Response) string {
	bytes, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		log.Fatal(err)
	}

	dumpBytes, _ := httputil.DumpResponse(resp, true)
	log.Println(string(dumpBytes))

	text := string(bytes)
	return text
}

//download small files like .torrent or .srt file
func quickDownload(url string) ([]byte, error) {
	resp, err := Client.Get(url)
	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	if err != nil {
		return nil, err
	}

	return data, nil
}

func postFile(filename string, filebytes []byte, target_url string) (*http.Response, error) {
	fmt.Println("filename:", filename)
	fmt.Println("target_url:", target_url)

	buffer := bytes.NewBufferString("")
	writer := multipart.NewWriter(buffer)
	w, _ := writer.CreateFormFile("filepath", filename)
	w.Write(filebytes)
	writer.WriteField("random", "136282211134691729.1585377371")
	writer.WriteField("interfrom", "task")
	writer.Close()

	resp, err := Client.Post(target_url, writer.FormDataContentType(), buffer)

	return resp, err
}
