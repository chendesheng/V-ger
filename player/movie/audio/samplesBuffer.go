package audio

import (
	"sync"
	"time"
)

type samples struct {
	data []byte
	pts  time.Duration
}

type sampleBuffer struct {
	sync.Mutex
	buf         []*samples
	bytesPerSec int
}

func (sb *sampleBuffer) append(s *samples) time.Duration {
	sb.Lock()
	defer sb.Unlock()

	sb.buf = append(sb.buf, s)
	return s.pts + sb.bytesToDur(len(s.data))
}

func (sb *sampleBuffer) empty() bool {
	sb.Lock()
	defer sb.Unlock()

	return len(sb.buf) == 0 || len(sb.buf[0].data) == 0
}

func (sb *sampleBuffer) clear() {
	sb.Lock()
	defer sb.Unlock()

	sb.buf = sb.buf[0:0]
}

func (sb *sampleBuffer) bytesToDur(bytes int) time.Duration {
	return time.Duration(float64(bytes) / float64(sb.bytesPerSec) * float64(time.Second))
}

func (sb *sampleBuffer) durToBytes(dur time.Duration) int {
	return int(float64(dur) / float64(time.Second) * float64(sb.bytesPerSec))
}

func (sb *sampleBuffer) cut(length int) (ret []byte) {
	sb.Lock()
	defer sb.Unlock()

	for len(sb.buf) > 0 && length > 0 {
		s := sb.buf[0]
		l := len(s.data)

		if length < l {
			ret = append(ret, s.data[:length]...)
			s.data = s.data[length:]
			s.pts += sb.bytesToDur(length)
			break
		} else {
			ret = append(ret, s.data[:l]...)
			length -= l
			sb.buf = sb.buf[1:]
		}
	}

	return
}

func (sb *sampleBuffer) pts() time.Duration {
	if sb.empty() {
		return 0
	} else {
		sb.Lock()
		defer sb.Unlock()

		return sb.buf[0].pts
	}
}

func (sb *sampleBuffer) cutByTime(dur time.Duration) {
	sb.cut(sb.durToBytes(dur) / 4 * 4)
}
