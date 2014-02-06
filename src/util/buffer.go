package util

import (
	"io"
	"sync"
)

type Buffer struct {
	sync.RWMutex
	data [][]byte

	currentPos int64
}

func (buf *Buffer) GetCurrentPos() int64 {
	buf.RLock()
	defer buf.RUnlock()

	return buf.currentPos
}

func (buf *Buffer) SetCurrentPos(pos int64) {
	buf.Lock()
	defer buf.Unlock()

	buf.currentPos = pos
}

func (buf *Buffer) Write(p []byte) (n int, err error) {
	buf.Lock()
	defer buf.Unlock()

	buf.data = append(buf.data, p)
	err = nil
	n = len(p)
	// println("write:", len(p), len(buf.data), buf.data)
	return
}

func (buf *Buffer) largerThan(size int) bool {
	for _, bytes := range buf.data {
		size -= len(bytes)

		if size <= 0 {
			return true
		}

	}

	return false
}

func (buf *Buffer) Read(size int) ([]byte, error) {
	buf.RLock()
	defer buf.RUnlock()

	res := make([]byte, 0)

	// if buf.largerThan(size) {
	for size > 0 && len(buf.data) > 0 {
		bytes := buf.data[0]
		if len(bytes) > size {
			res = append(res, bytes[:size]...)

			buf.data[0] = bytes[size:]
			size = 0
		} else {
			size -= len(bytes)
			res = append(res, bytes...)
			buf.data = buf.data[1:]
		}
	}

	if len(res) == 0 {
		return nil, io.EOF
	} else {
		return res, nil
	}

	// } else {
	// return res, io.EOF
	// }
}

func (buf *Buffer) ClearData() {
	buf.Lock()
	defer buf.Unlock()

	println("buffer ClearData")
	buf.data = nil
}
