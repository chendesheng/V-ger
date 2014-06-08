package movie

import (
	"download"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	// "path/filepath"
	. "player/libav"
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
func max(a, b int64) int64 {
	if a > b {
		return a
	}

	return b
}

func (m *Movie) openHttp(file string) (AVFormatContext, string) {
	download.NetworkTimeout = time.Duration(util.ReadIntConfig("network-timeout")) * time.Second
	download.BaseDir = util.ReadConfig("dir")

	url, name, size, err := download.GetDownloadInfoN(file, 3, m.quit)

	if err != nil {
		log.Fatal(err)
	}

	m.httpBuffer = NewBuffer(size)

	buf := AVObject{}
	buf.Malloc(1024 * 64)
	ioctx := NewAVIOContext(buf, func(buf AVObject) int {
		if buf.Size() == 0 {
			return 0
		}
		if m.c != nil {
			t := m.c.GetTime()
			defer m.c.SetTime(t)
		}

		require := int64(buf.Size())

		got := m.httpBuffer.Read(&buf, require)

		if got < require && !m.httpBuffer.IsFinish() {
			if m.c != nil {
				m.c.Pause()
				defer m.c.Resume()
			}

			for got < require && !m.httpBuffer.IsFinish() {
				time.Sleep(20 * time.Millisecond)
				got += m.httpBuffer.Read(&buf, require-got)
			}
		}

		return int(got)
	}, func(offset int64, whence int) int64 {
		println("seek:", offset, whence)
		if whence == AVSEEK_SIZE {
			return m.httpBuffer.size
		}

		pos, start := m.httpBuffer.Seek(offset, whence)
		if start >= 0 && start < size {
			go download.Streaming(url, size, m.httpBuffer, start, m)
		}
		return pos
	})

	ctx := NewAVFormatContext()
	ctx.SetPb(ioctx)

	go download.Streaming(url, size, m.httpBuffer, 0, m)
	m.httpBuffer.Seek(0, os.SEEK_SET)

	ctx.OpenInput(name)

	println("open http return")
	return ctx, name
}
