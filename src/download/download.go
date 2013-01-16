package download

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"time"
)

type Config struct {
	BaseDir string
}
type ProgressInfo struct {
	Size        int64
	ElapsedTime time.Duration
}

var DownloadClient *http.Client

func sendGet(url string) *http.Response {
	resp, err := DownloadClient.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	return resp
}
func addRangeHeader(req *http.Request, pos int64) {
	if pos <= 0 {
		return
	}
	req.Header.Add("Range", fmt.Sprintf("bytes=%d-", pos))
}
func createRequest(t *Task) *http.Request {
	req := new(http.Request)
	req.Method = "GET"
	req.URL, _ = url.Parse(t.URL)
	req.Header = make(http.Header)
	return req
}
func openOrCreateFileRW(path string) *os.File {
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		log.Fatal(err)
	}
	return f
}

func writeDownload(t *Task, resp *http.Response) {
	f := openOrCreateFileRW(t.Path)
	defer f.Close()

	size, total, elapsedTime := t.Size, t.DownloadedSize, t.ElapsedTime
	f.Seek(total, 0)

	buffer := make([]byte, 40000)
	part := int64(0)
	parts := [5]int64{0, 0, 0, 0, 0}
	checkTimes := [5]time.Time{time.Now(), time.Now(), time.Now(), time.Now(), time.Now()}
	cnt := 0
	percentage := float64(total) / float64(size) * 100
	speed := float64(0) // average speed of recent 5 seconds
	for {
		readLen, _ := resp.Body.Read(buffer)
		if readLen == 0 {
			return
		}
		f.Write(buffer[:readLen])

		total += int64(readLen)
		part += int64(readLen)

		if time.Since(checkTimes[cnt]) > time.Second {
			t.DownloadedSize = total
			elapsedTime += time.Second
			t.ElapsedTime = elapsedTime
			saveTask(t)

			percentage = float64(total) / float64(size) * 100

			cnt++
			cnt = cnt % 5

			sinceLastCheck := time.Since(checkTimes[cnt])

			checkTimes[cnt] = time.Now()
			parts[cnt] = part
			part = 0

			sum := int64(0)
			for _, p := range parts {
				sum += p
			}
			speed = float64(sum) * float64(time.Second) / float64(sinceLastCheck) / 1024
			est := time.Duration(float64((size-total))/speed) * time.Millisecond

			printProgress(percentage, speed, elapsedTime, est)

			// if speed > 200 {
			// 	time.Sleep(time.Duration(float64(sum)*float64(time.Second)/speed/1024))
			// }
		}
	}
}
func BeginDownload(url string, name string) {
	globalConfig = readConfig()

	if DownloadClient == nil {
		DownloadClient = http.DefaultClient
	}

	currentTask := getOrNewTask(url, name)

	DownloadClient.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		currentTask.URL = req.URL.String()

		prevReq := via[len(via)-1]
		rangeVal := prevReq.Header.Get("Range")
		if rangeVal != "" {
			req.Header.Add("Range", rangeVal)
		}
		return nil
	}

	req := createRequest(&currentTask)
	addRangeHeader(req, currentTask.DownloadedSize)
	// fmt.Println(req.URL)
	resp, err := DownloadClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		resp.Body.Close()
		DownloadClient.CheckRedirect = nil
	}()
	currentTask.Name = name

	if currentTask.isNew {
		var name string
		name, currentTask.Size = getFileInfo(resp)
		if currentTask.Name == "" {
			currentTask.Name = name
		}

		currentTask.StartDate = time.Now().String()
		currentTask.Path = fmt.Sprintf("%s%c%s", globalConfig.BaseDir, os.PathSeparator, currentTask.Name)

		saveTask(&currentTask)

		fmt.Printf("New Task: %s    %d\n", currentTask.Name, currentTask.Size)
	}
	writeDownload(&currentTask, resp)
	removeTask(currentTask.Name)
	fmt.Printf("It's done!\n\n")
}

// func StopDownload(url string) {

// }

// func GetTaskProgress(name string) Progress {

// }
func GetTasks() []Task {
	return getTasks()
}

var globalConfig Config

func getFileInfo(resp *http.Response) (name string, size int64) {
	if len(resp.Header["Content-Disposition"]) > 0 {
		contentDisposition := resp.Header["Content-Disposition"][0]
		regexFile, err := regexp.Compile(`filename="([^"]+)"`)
		if err != nil {
			log.Fatal(err)
		}
		name = regexFile.FindStringSubmatch(contentDisposition)[1]
	}

	if cr := resp.Header["Content-Range"]; len(cr) > 0 {
		regexSize, err := regexp.Compile(`/(\d+)`)
		if err != nil {
			log.Fatal(err)
		}
		sizeStr := regexSize.FindStringSubmatch(cr[0])[1]
		size, _ = strconv.ParseInt(sizeStr, 10, 64)
	} else {
		size, _ = strconv.ParseInt(resp.Header["Content-Length"][0], 10, 64)
	}

	return
}

func writeJson(path string, object interface{}) {
	data, err := json.Marshal(object)
	if err != nil {
		log.Fatal(err)
	}

	ioutil.WriteFile(path, data, 0666)
}
func readJson(path string, object interface{}) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("read json: ")
	log.Println(path)
	log.Println(string(data))

	json.Unmarshal(data, &object)
}
func readConfig() Config {
	config := Config{}
	readJson("config.json", &config)
	return config
}
func saveConfig(config Config) {
	writeJson("config.json", config)
}
func printProgress(percentage float64, speed float64, elapsedTime time.Duration, est time.Duration) {
	fmt.Printf("\r%.2f%%    %.2f KB/s    %s    Est. %s    ", percentage, speed, elapsedTime.String(), est.String())
}
