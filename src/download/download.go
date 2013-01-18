package download

import (
	"bytes"
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

type Downloader struct {
	task *Task
}

func (d *Downloader) writeDownload(resp *http.Response) {
	t := d.task
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
		}
	}
}
func (d *Downloader) doDownload(from int64, blockSize int64, f *os.File, connChans []chan int64, buffers [][]byte) int64 {
	to := from + blockSize - 1
	t := d.task
	size := t.Size

	var i int
	for i = 0; i < 6; i++ {

		if to > size {
			to = -1
		}

		go func(pBuffer *[]byte, connChan chan int64, from, to int64) {
			d.downloadBlock(from, to, pBuffer)
			f.WriteAt(*pBuffer, from)
			connChan <- int64(len(*pBuffer))
		}(&buffers[i], connChans[i], from, to)

		if to == -1 {
			i++
			break
		} else {
			from, to = to+1, to+blockSize
		}

	}

	total := int64(0)
	for i--; i >= 0; i-- {
		length := <-connChans[i]
		total += length
	}
	return total
}
func (d *Downloader) download() {
	t := d.task

	f := openOrCreateFileRW(t.Path)
	if t.isNew {
		f.Truncate(t.Size)
	}

	defer f.Close()

	blockSize := int64(400 * 1024)

	size, total, elapsedTime := t.Size, t.DownloadedSize, t.ElapsedTime

	part := int64(0)
	parts := [5]int64{0, 0, 0, 0, 0}
	checkTimes := [5]time.Time{time.Now(), time.Now(), time.Now(), time.Now(), time.Now()}
	cnt := 0
	percentage := float64(total) / float64(size) * 100
	speed := float64(0) // average speed of recent 5 seconds

	connChans := make([]chan int64, 10)
	buffers := make([][]byte, 10)
	for i := 0; i < 10; i++ {
		connChans[i] = make(chan int64)
		buffers[i] = make([]byte, blockSize)
	}

	for total < size {
		length := d.doDownload(total, blockSize, f, connChans, buffers)

		total += length
		part += length

		if time.Since(checkTimes[cnt]) > time.Second || total == size {
			t.DownloadedSize = total
			elapsedTime += time.Since(checkTimes[cnt])
			t.ElapsedTime = elapsedTime
			saveTask(t)

			percentage = float64(total) / float64(size) * 100

			cnt++
			cnt = cnt % 5

			sinceLastCheck := time.Since(checkTimes[cnt])

			checkTimes[cnt] = time.Now()
			parts[cnt] = part
			part = 0

			//sum up download size of recent 5 seconds
			sum := int64(0)
			for _, p := range parts {
				sum += p
			}
			speed = float64(sum) * float64(time.Second) / float64(sinceLastCheck) / 1024
			est := time.Duration(float64((size-total))/speed) * time.Millisecond

			printProgress(percentage, speed, elapsedTime, est)
		}
	}
}
func BeginDownload(url, name string) {
	globalConfig = readConfig()

	if DownloadClient == nil {
		DownloadClient = http.DefaultClient
	}
	d := Downloader{}
	d.initDownload(url, name)
	// req := createRequest(d.task)
	// resp, err := DownloadClient.Do(req)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// d.writeDownload(resp)
	d.download()
	d.done()
}
func (d *Downloader) initDownload(url, name string) {
	t := getOrNewTask(url, name)
	d.task = &t

	req := createRequest(d.task)
	DownloadClient.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		t.URL = req.URL.String()
		return nil
	}

	resp, err := DownloadClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	t.Name = name

	if t.isNew {
		var name string
		name, t.Size = getFileInfo(resp)
		if t.Name == "" {
			t.Name = name
		}

		t.StartDate = time.Now().String()
		t.Path = fmt.Sprintf("%s%c%s", globalConfig.BaseDir, os.PathSeparator, t.Name)

		saveTask(&t)

		fmt.Printf("New Task: %s\t%d\n", t.Name, t.Size)
	}
}
func (d *Downloader) downloadBlock(from, to int64, block *[]byte) {
	buffer := bytes.NewBuffer(*block)
	buffer.Truncate(0)

	req := createRequest(d.task)
	addRangeHeader(req, from, to)
	resp, err := DownloadClient.Do(req)

	defer resp.Body.Close()
	if err != nil {
		log.Fatal(err)
	}
	part := make([]byte, 2000)
	total := int64(0)
	for {
		readLen, err := resp.Body.Read(part)
		if readLen == 0 {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		buffer.Write(part[:readLen])
		total += int64(readLen)
	}

	*block = buffer.Bytes()[:total]
}
func (d *Downloader) done() {
	removeTask(d.task.Name)
	fmt.Printf("\nIt's done!\n\n")
}

var DownloadClient *http.Client

func sendGet(url string) *http.Response {
	resp, err := DownloadClient.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	return resp
}
func addRangeHeader(req *http.Request, from, to int64) {
	if from == to || (from <= 0 && to < 0) {
		return
	}
	if to < 0 {
		req.Header.Add("Range", fmt.Sprintf("bytes=%d-", from))
	} else {
		req.Header.Add("Range", fmt.Sprintf("bytes=%d-%d", from, to))
	}
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

// func StopDownload(url string) {

// }

// func GetTaskProgress(name string) Progress {

// }

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
func printProgress(percentage float64, speed float64, elapsedTime time.Duration, est time.Duration) {
	fmt.Printf("\r%.2f%%    %.2f KB/s    %s    Est. %s    ", percentage, speed, elapsedTime.String(), est.String())
}
