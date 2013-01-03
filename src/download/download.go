package download

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"time"
)

type Task struct {
	URL       string
	Size      int64
	Name      string
	Path      string
	StartDate string
}
type Config struct {
	BaseDir  string
	TaskList []Task
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
func BeginDownload(url string, name string) {
	globalConfig = readConfig()

	if DownloadClient == nil {
		DownloadClient = http.DefaultClient
	}

	currentTask, currentTaskIfNew := getOrNewTask(url, name)

	DownloadClient.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		if currentTaskIfNew {
			return nil
		}
		pos, _ := readProgress(currentTask)
		if pos > 0 {
			req.Header.Add("Range", fmt.Sprintf("bytes=%d-", pos))
		}
		return nil
	}

	resp := sendGet(currentTask.URL)
	defer func() {
		resp.Body.Close()
		DownloadClient.CheckRedirect = nil
	}()

	if currentTaskIfNew {
		currentTask.Name, currentTask.Size = getFileInfo(resp)
		if name != "" {
			currentTask.Name = name
		}

		currentTask.StartDate = time.Now().String()
		currentTask.Path = fmt.Sprintf("%s%c%s", globalConfig.BaseDir, os.PathSeparator, currentTask.Name)

		globalConfig.TaskList = append(globalConfig.TaskList, currentTask)

		saveConfig(&globalConfig)

		fmt.Printf("New Task: %s    %d\n", currentTask.Name, currentTask.Size)
	}
	size := currentTask.Size

	//write file
	f, err := os.OpenFile(currentTask.Path, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	pos, elapsedTime := readProgress(currentTask)
	f.Seek(pos, 0)
	//io.Copy(f, resp.Body)
	bytes := make([]byte, 40000)
	total := pos
	part := int64(0)
	parts := [5]int64{0, 0, 0, 0, 0}
	checkTimes := [5]time.Time{time.Now(), time.Now(), time.Now(), time.Now(), time.Now()}
	cnt := 0

	beginTime := time.Now()

	percentage := float64(total) / float64(size) * 100
	totalElapsedTime := elapsedTime
	speed := float64(0)
	for {
		readLen, _ := resp.Body.Read(bytes)
		if readLen == 0 {
			removeTask(currentTask)
			fmt.Printf("\nIt's done! %.2f KB/s on average.\n", float64(size)*float64(time.Second)/float64(elapsedTime)/1024)
			return
		}
		f.Write(bytes[:readLen])

		total += int64(readLen)
		part += int64(readLen)

		if time.Since(checkTimes[cnt]) > time.Second {
			totalElapsedTime = time.Since(beginTime) + elapsedTime
			saveProgress(currentTask, total, totalElapsedTime)
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

			printProgress(percentage, speed, totalElapsedTime, est)

			// if speed > 200 {
			// 	time.Sleep(time.Duration(float64(sum)*float64(time.Second)/speed/1024))
			// }
		}
	}
}

// func StopDownload(url string) {

// }

// func GetTaskProgress(name string) Progress {

// }
func GetTasks() []Task {
	globalConfig = readConfig()
	return globalConfig.TaskList
}

var globalConfig Config

func getFileInfo(resp *http.Response) (name string, size int64) {
	contentDisposition := resp.Header["Content-Disposition"][0]
	regexFile, err := regexp.Compile(`filename="([^"]+)"`)
	if err != nil {
		log.Fatal(err)
	}
	name = regexFile.FindStringSubmatch(contentDisposition)[1]

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
func readProgress(currentTask Task) (size int64, elapsedTime time.Duration) {
	path := currentTask.Path + ".json"

	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return 0, 0
	}

	data := ProgressInfo{}
	err = json.Unmarshal(bytes, &data)
	if err != nil {
		return 0, 0
	}

	//fmt.Println(string(bytes))
	return data.Size, data.ElapsedTime
}
func saveProgress(currentTask Task, size int64, elapsedTime time.Duration) {
	data, _ := json.Marshal(ProgressInfo{
		Size:        size,
		ElapsedTime: elapsedTime,
	})

	path := currentTask.Path + ".json"
	ioutil.WriteFile(path, data, 0666)
}
func readConfig() Config {
	config := Config{}
	bytes, err := ioutil.ReadFile("config.json")
	if err != nil {
		log.Fatal(err)
	}
	json.Unmarshal(bytes, &config)
	return config
}
func saveConfig(config *Config) {
	data, _ := json.Marshal(config)
	ioutil.WriteFile("config.json", data, 0666)
}
func printProgress(percentage float64, speed float64, elapsedTime time.Duration, est time.Duration) {
	fmt.Printf("\r%.2f%%    %.2f KB/s    %s    Est. %s    ", percentage, speed, elapsedTime.String(), est.String())
}

func getOrNewTask(url string, name string) (Task, bool) {
	for _, t := range globalConfig.TaskList {
		if url == t.URL {
			return t, false
		}
	}

	t := Task{URL: url, Name: name}
	return t, true
}

func removeTask(task Task) {
	config := readConfig()

	for i, t := range config.TaskList {
		if task.URL == t.URL {
			config.TaskList = append(config.TaskList[:i], config.TaskList[i+1:]...)
			break
		}
	}

	saveConfig(&config)
}

// func main() {
// 	var url string
// 	if len(os.Args) > 1 {
// 		url = os.Args[1]
// 		currentTask, currentTaskIfNew = getOrNewTask(url)
// 	} else {
// 		if len(globalConfig.TaskList) == 0 {
// 			fmt.Println("no task yet.")
// 			return
// 		}
// 		for i, t := range globalConfig.TaskList {
// 			fmt.Printf("[%d] %s  %s\n", i+1, t.Name, t.StartDate)
// 		}
// 		i := 0
// 		_, err := fmt.Scanf("%d", &i)
// 		if err != nil {
// 			log.Fatal(err)
// 			return
// 		}
// 		i--
// 		if i >= 0 && i < len(globalConfig.TaskList) {
// 			currentTask = globalConfig.TaskList[i]
// 		} else {
// 			fmt.Println("pick wrong number.")
// 			return
// 		}
// 	}

// 	req, err := http.NewRequest("GET", currentTask.URL, nil)
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	req.Header.Add("Cookie", "gdriveid=5120E7CE422D1E3F34D7ED1501A1C86A")

// 	client := &http.Client{CheckRedirect: redirect}
// 	resp, err := client.Do(req)
// 	defer resp.Body.Close()

// 	if currentTaskIfNew {
// 		currentTask.Name, currentTask.Size = getFileInfo(resp)
// 		currentTask.StartDate = time.Now().String()
// 		currentTask.Path = fmt.Sprintf("%s%c%s", globalConfig.BaseDir, os.PathSeparator, currentTask.Name)

// 		globalConfig.TaskList = append(globalConfig.TaskList, currentTask)
// 		saveConfig(&globalConfig)

// 		fmt.Printf("New Task: %s    %d\n", currentTask.Name, currentTask.Size)
// 	}
// 	size := currentTask.Size

// 	//write file
// 	f, err := os.OpenFile(currentTask.Path, os.O_RDWR|os.O_CREATE, 0666)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer f.Close()
// 	pos, elapsedTime := readProgress()
// 	f.Seek(pos, 0)
// 	//io.Copy(f, resp.Body)
// 	bytes := make([]byte, 40000)
// 	total := pos
// 	part := int64(0)
// 	parts := [5]int64{0, 0, 0, 0, 0}
// 	checkTimes := [5]time.Time{time.Now(), time.Now(), time.Now(), time.Now(), time.Now()}
// 	cnt := 0

// 	beginTime := time.Now()

// 	percentage := float64(total) / float64(size) * 100
// 	totalElapsedTime := elapsedTime
// 	speed := float64(0)
// 	for {
// 		readLen, _ := resp.Body.Read(bytes)
// 		if readLen == 0 {
// 			removeTask(currentTask)
// 			fmt.Printf("\nIt's done! %.2f KB/s on average.\n", float64(size)*float64(time.Second)/float64(elapsedTime)/1024)
// 			return
// 		}
// 		f.Write(bytes[:readLen])

// 		total += int64(readLen)
// 		part += int64(readLen)

// 		if time.Since(checkTimes[cnt]) > time.Second {
// 			totalElapsedTime = time.Since(beginTime) + elapsedTime
// 			saveProgress(currentTask, total, totalElapsedTime)
// 			percentage = float64(total) / float64(size) * 100

// 			cnt++
// 			cnt = cnt % 5

// 			sinceLastCheck := time.Since(checkTimes[cnt])

// 			checkTimes[cnt] = time.Now()
// 			parts[cnt] = part
// 			part = 0

// 			sum := int64(0)
// 			for _, p := range parts {
// 				sum += p
// 			}
// 			speed = float64(sum) * float64(time.Second) / float64(sinceLastCheck) / 1024
// 			est := time.Duration(float64((size-total))/speed) * time.Millisecond

// 			printProgress(percentage, speed, totalElapsedTime, est)

// 			// if speed > 200 {
// 			// 	time.Sleep(time.Duration(float64(sum)*float64(time.Second)/speed/1024))
// 			// }
// 		}
// 	}
// }
