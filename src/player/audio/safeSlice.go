package audio

import "sync"

type safeSlice struct {
	sync.Mutex
	buf []byte
}

func (ss *safeSlice) append(p []byte) {
	ss.Lock()
	defer ss.Unlock()

	ss.buf = append(ss.buf, p...)
}

func (ss *safeSlice) empty() bool {
	ss.Lock()
	defer ss.Unlock()

	return len(ss.buf) == 0
}

func (ss *safeSlice) clear() {
	ss.Lock()
	defer ss.Unlock()

	ss.buf = ss.buf[0:0]
}

func (ss *safeSlice) cut(length int) []byte {
	ss.Lock()
	defer ss.Unlock()

	retLen := min2(len(ss.buf), length)
	ret := ss.buf[:retLen]
	ss.buf = ss.buf[retLen:]
	return ret
}
