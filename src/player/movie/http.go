package movie

import (
	"download"
	"log"
	"os"
	. "player/libav"
	"task"
	"time"
	"util"
)

func (m *Movie) openHttp(file string) (AVFormatContext, string) {
	_, name, size, err := download.GetDownloadInfo(file)
	if err != nil {
		log.Fatal(err)
	}
	t, err := task.GetTask(name)
	if err != nil {
		// log.Fatal(err)
		t = &task.Task{}
		t.Name = name
		t.Size = size
		t.StartTime = time.Now().Unix()
		t.Status = "Stopped"
		t.URL = file
		task.SaveTask(t)
	}

	mbuf := &util.Buffer{}

	// currentPos := int64(0)
	buf := AVObject{}
	buf.Malloc(1024 * 32)
	ioctx := NewAVIOContext(buf, func(buf AVObject) int {
		if buf.Size() == 0 {
			return 0
		}

		size := buf.Size()
		// currentPos := mbuf.CurrentPos

		// println("read: currentPos11:", mbuf.CurrentPos)
		// for size > 0 && currentPos+int64(size) < t.Size {
		for {

			if bytes, err := mbuf.Read(size); err == nil {
				buf.Write(bytes)
				// size -= len(bytes)
				currentPos := mbuf.GetCurrentPos()
				currentPos += int64(len(bytes))
				mbuf.SetCurrentPos(currentPos)
				println("readfunc:", currentPos, len(bytes))
				return len(bytes)
			}

			time.Sleep(50 * time.Millisecond)
		}

		// println("read: currentPos:", mbuf.CurrentPos)
		// return buf.Size() - size
	}, func(pos int64, whence int) int64 {
		println("seekfunc:", pos, whence)

		// download.Play(t, w, from, to)
		// t, _ = task.GetTask(t.Name)
		currentPos := mbuf.GetCurrentPos()
		switch whence {
		case os.SEEK_SET:
			currentPos = pos
			break
		case os.SEEK_CUR:
			currentPos += pos
			break
		case os.SEEK_END:
			currentPos = t.Size + pos
			break
		// case AVSEEK_SIZE:
		default:
			return t.Size
		}

		if currentPos > t.Size {
			currentPos = t.Size
			return currentPos
		}

		if currentPos < 0 {
			return -1
		}
		mbuf.SetCurrentPos(currentPos)
		mbuf.ClearData()
		go download.Play(t, mbuf, currentPos, t.Size)
		return currentPos
	})

	ctx := NewAVFormatContext()
	ctx.SetPb(ioctx)

	download.Play(t, mbuf, 0, 1024*32)

	if bytes, err := mbuf.Read(1024 * 32); err == nil {
		println("bytes:", len(bytes))
		pd := NewAVProbeData()

		obj := AVObject{}
		obj.Malloc(len(bytes))
		obj.Write(bytes)

		pd.SetBuffer(obj)
		pd.SetFileName("")

		ctx.SetInputFormat(pd.InputFormat())

		mbuf.ClearData()

		go download.Play(t, mbuf, 0, t.Size)

		ctx.OpenInput("")
	}

	return ctx, name
}
