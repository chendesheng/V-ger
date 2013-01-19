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

func (d *Downloader) beginSampleDownload() chan []byte {
	req := createRequest(d.task)
	resp, err := DownloadClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	buffer := make([]byte, 40000)
	output := make(chan []byte)
	for {
		readLen, _ := resp.Body.Read(buffer)
		if readLen == 0 {
			break
		}
		output <- buffer[:readLen]
	}

	return output
}
func (d *Downloader) sampleDownload(resp *http.Response) {
	output := d.beginSampleDownload()
	d.writeOutput(output)
}

func (d *Downloader) pipeDownloadBlocks(from int64, blockSize int, size int64, output chan []byte, readyOutput chan bool, cntControl chan bool) {
	cntControl <- true // block if full
	if from >= size {
		<-readyOutput
		close(output)
	}

	complete := make(chan bool)

	to := from + int64(blockSize)
	if to > size {
		to = size
	}

	go func() {
		block := d.downloadBlock(from, to-1)

		<-readyOutput
		output <- block
		complete <- true
		<-cntControl
	}()

	d.pipeDownloadBlocks(to, blockSize, size, output, complete, cntControl)
}
func (d *Downloader) beginPipeDownload() chan []byte {
	t := d.task

	blockCnt := 6
	blockSize := 200 * 1024

	output := make(chan []byte, blockCnt)
	cntControl := make(chan bool, blockCnt)

	readyOutput := make(chan bool)
	go d.pipeDownloadBlocks(t.DownloadedSize, blockSize, t.Size, output, readyOutput, cntControl)
	readyOutput <- true

	return output
}
func (d *Downloader) pipeDownload() {
	output := d.beginPipeDownload()
	d.writeOutput(output)
}

func (d *Downloader) downloadBlock(from, to int64) []byte {
	req := createRequest(d.task)
	addRangeHeader(req, from, to)

	resp, err := DownloadClient.Do(req)
	defer resp.Body.Close()
	if err != nil {
		log.Fatal(err)
	}

	buffer := bytes.NewBuffer(make([]byte, 0, to-from+1))
	_, err = buffer.ReadFrom(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	return buffer.Bytes()
}
func (d *Downloader) progress() chan int64 {
	t := d.task
	downloadBlockChan := make(chan int64)
	go func() {
		size, total, elapsedTime := t.Size, t.DownloadedSize, t.ElapsedTime

		part := int64(0)
		parts := [5]int64{0, 0, 0, 0, 0}
		checkTimes := [5]time.Time{time.Now(), time.Now(), time.Now(), time.Now(), time.Now()}
		cnt := 0
		percentage := float64(total) / float64(size) * 100
		speed := float64(0) // average speed of recent 5 seconds

		for length := range downloadBlockChan {
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
	}()

	return downloadBlockChan
}
func (d *Downloader) writeOutput(output chan []byte) {
	t := d.task

	f := openOrCreateFileRW(t.Path)
	defer f.Close()
	f.Seek(t.DownloadedSize, 0)

	progressChan := d.progress()

	for b := range output {
		f.Write(b)

		progressChan <- int64(len(b))
	}
	close(progressChan)
}
func BeginDownload(url, name string) {
	globalConfig = readConfig()

	if DownloadClient == nil {
		DownloadClient = http.DefaultClient
	}
	d := Downloader{}
	d.initDownload(url, name)

	// d.sampleDownload(resp)
	// d.download()
	d.pipeDownload()
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
