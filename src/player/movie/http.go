package movie

import (
	"download"
	"fmt"
	"log"
	"os"

	// "path/filepath"
	. "player/libav"
	"time"
	// "task"
	// "time"
	"util"
)

func (m *Movie) openHttp(file string) (AVFormatContext, string, error) {
	download.NetworkTimeout = time.Duration(util.ReadIntConfig("network-timeout")) * time.Second
	download.BaseDir = util.ReadConfig("dir")

	m.chSpeed = make(chan float64)

	url, name, size, _, err := download.GetDownloadInfoN(file, nil, 10, false, m.quit)

	if err != nil {
		return AVFormatContext{}, "", err
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
				select {
				case <-time.After(20 * time.Millisecond):
				case <-m.quit:
					return AVERROR_INVALIDDATA
				}

				got += m.httpBuffer.Read(&buf, require-got)
				if got >= require || m.httpBuffer.IsFinish() {
					break
				} else {
					if time.Since(startWaitTime) > download.NetworkTimeout {
						pos := m.httpBuffer.CurrentPos()

						log.Print("Streamming timeout restart:", pos)

						startWaitTime = time.Now()
						select {
						case streaming.Restart() <- pos:
						case <-m.quit:
							return AVERROR_INVALIDDATA
						}
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
			select {
			case streaming.Restart() <- start:
			case <-m.quit:
				return pos
			}
		}
		return pos
	})

	ctx := NewAVFormatContext()
	ctx.SetPb(ioctx)

	m.httpBuffer.Seek(0, os.SEEK_SET)
	select {
	case streaming.Restart() <- 0:
	case <-m.quit:
		return ctx, name, fmt.Errorf("quit")
	}

	if err := ctx.OpenInput(file); err != nil {
		return ctx, name, err
	}

	log.Print("open http return")
	return ctx, name, nil
}
