package movie

import (
	"download"
	"log"
	"os"

	// "path/filepath"
	. "player/libav"
	"time"
	// "task"
	// "time"
	"util"
)

func (m *Movie) openHttp(file string) (AVFormatContext, string) {
	download.NetworkTimeout = time.Duration(util.ReadIntConfig("network-timeout")) * time.Second
	download.BaseDir = util.ReadConfig("dir")

	m.chSpeed = make(chan float64)

	url, name, size, err := download.GetDownloadInfoN(file, 3, m.quit)

	if err != nil {
		log.Fatal(err)
	}

	m.httpBuffer = NewBuffer(size)

	buf := AVObject{}
	buf.Malloc(1024 * 64)

	streaming := download.StartStreaming(url, size, m.httpBuffer, m)
	m.streaming = streaming

	ioctx := NewAVIOContext(buf, func(buf AVObject) int {
		if buf.Size() == 0 {
			return 0
		}

		if m.httpBuffer.CurrentPos() >= size {
			return AVERROR_EOF
		}

		require := int64(buf.Size())
		got := m.httpBuffer.Read(&buf, require)
		if got < require && !m.httpBuffer.IsFinish() {
			startWaitTime := time.Now()

			if m.c != nil {
				t := m.c.GetTime()
				defer m.c.SetTime(t)
			}

			for {
				time.Sleep(20 * time.Millisecond)
				got += m.httpBuffer.Read(&buf, require-got)
				if got >= require || m.httpBuffer.IsFinish() {
					break
				} else {
					if time.Since(startWaitTime) > download.NetworkTimeout {
						pos := m.httpBuffer.CurrentPos()

						log.Print("Streamming timeout restart:", pos)

						startWaitTime = time.Now()
						streaming.Restart(pos)
					}
				}
			}
		}

		return int(got)
	}, func(offset int64, whence int) int64 {
		if whence == AVSEEK_SIZE {
			return m.httpBuffer.size
		}

		pos, start := m.httpBuffer.Seek(offset, whence)
		if start >= 0 && start < size {
			m.w.SendShowSpinning()
			streaming.Restart(start)
		}
		return pos
	})

	ctx := NewAVFormatContext()
	ctx.SetPb(ioctx)

	m.httpBuffer.Seek(0, os.SEEK_SET)
	streaming.Restart(0)

	ctx.OpenInput(name)

	log.Print("open http return")
	return ctx, name
}
