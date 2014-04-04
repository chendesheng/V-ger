package movie

import (
	"download"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	// "path/filepath"
	. "player/libav"
	"task"
	"time"
	// "task"
	// "time"
	"util"
)

func downloadBytes(url string, from int64, size int, filesize int64) []byte {
	to := from + int64(size)
	if to > filesize {
		to = 0
	}

	println("request:", from, to)
	req := download.CreateDownloadRequest(url, from, to-1)
	resp, _ := http.DefaultClient.Do(req)

	data, _ := ioutil.ReadAll(resp.Body)
	println("get:", len(data), from, size)
	return data
}

var mbuf *buffer

func (m *Movie) openHttp(file string) (AVFormatContext, string) {
	download.NetworkTimeout = 30 * time.Second
	download.BaseDir = util.ReadConfig("dir")

	_, name, size, err := download.GetDownloadInfo(file)

	if err != nil {
		log.Fatal(err)
	}

	t, err := task.GetTask(name)
	println("download info:", size, t.Size)
	if err != nil {
		t = &task.Task{}
		t.Name = name
		t.Size = size
		t.StartTime = time.Now().Unix()
		t.Status = "Playing"
		t.URL = file
		task.SaveTask(t)
	} else {
		t.Status = "Playing"
		task.SaveTask(t)
	}

	mbuf = NewBuffer(size)

	buf := AVObject{}
	buf.Malloc(1024 * 64)
	ioctx := NewAVIOContext(buf, func(buf AVObject) int {
		if buf.Size() == 0 {
			return 0
		}

		return mbuf.Read(&buf, int64(buf.Size()))
	}, func(offset int64, whence int) int64 {
		println("seek:", offset, whence)
		if whence == AVSEEK_SIZE {
			return mbuf.size
		}

		pos, start := mbuf.Seek(offset, whence)
		if start >= 0 && start < size {
			go func() {
				t, err := task.GetTask(name)
				if err != nil {
					log.Fatal(err)
				}

				download.QuitAndDownload(t, mbuf, start)
			}()
		}
		return pos
	})

	ctx := NewAVFormatContext()
	ctx.SetPb(ioctx)

	if ctx.IsNil() {
		println("avformatcontext is nil1")
	}

	// pd := NewAVProbeData()

	go download.QuitAndDownload(t, mbuf, 0)
	// mbuf.Read(&buf, 1024*64)
	mbuf.Seek(0, os.SEEK_SET)

	// pd.SetBuffer(buf)
	// pd.SetFileName(name)

	// ctx.SetInputFormat(pd.InputFormat())

	if ctx.IsNil() {
		println("avformatcontext is nil12")
	}

	ctx.OpenInput(name)

	println("open http return")
	if ctx.IsNil() {
		println("avformatcontext is nil")
	}
	return ctx, name
}
