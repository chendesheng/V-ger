package audio

import (
	"fmt"
	"log"
	"math"
	. "player/clock"
	. "player/libav"
	"sync/atomic"
	"time"
)

type Audio struct {
	c           *Clock
	codecCtx    *AVCodecContext
	resampleCtx AVAudioResampleContext
	stream      AVStream
	frame       AVFrame

	PacketChan  chan *AVPacket
	audioBuffer *sampleBuffer

	skipBytes int

	quit chan struct{}

	driver *portAudio

	silence []byte

	diffCnt       int
	diffCum       float64
	diffCoef      float64
	diffThreshold time.Duration

	Offset time.Duration

	chQuitDone chan struct{}
}

func NewAudio(c *Clock, volume float64) *Audio {
	a := &Audio{}

	resampleCtx := AVAudioResampleContext{}
	resampleCtx.Alloc()

	a.frame = AllocFrame()
	a.resampleCtx = resampleCtx
	a.PacketChan = make(chan *AVPacket, 500)
	a.c = c
	a.driver = &portAudio{volume: volume}
	a.audioBuffer = &sampleBuffer{}
	return a
}
func (a *Audio) StreamIndex() int {
	return a.stream.Index()
}

func (a *Audio) receivePacket() (*AVPacket, bool) {
	select {
	case packet, ok := <-a.PacketChan:
		return packet, ok
	case <-a.quit:
		close(a.chQuitDone)
		return nil, false
	}
}

func (a *Audio) AddOffset(dur time.Duration) {
	atomic.AddInt64((*int64)(&a.Offset), int64(dur))
}

func (a *Audio) GetOffset() time.Duration {
	return time.Duration(atomic.LoadInt64((*int64)(&a.Offset)))
}

func (a *Audio) sync(pts time.Duration) {
	now := a.c.GetTime()

	diff := pts - now
	avgDiff := time.Duration(math.Abs(float64(diff)))

	if avgDiff < a.diffThreshold {
		return
	} else if pts > now && pts-now < 5*time.Second {
		log.Print("wait audio:", (pts - now).String())
		a.c.WaitUntilWithQuit(pts, a.quit)
		return
	} else {
		log.Print("skip audio packet:", (now - pts).String())
		a.audioBuffer.cutByTime(now - pts)
		// log.Print("rest:", len(a.audioBuffer.buf))
		return
	}
}

func (a *Audio) getPts(packet *AVPacket) time.Duration {
	var pts time.Duration
	if packet.Pts() != AV_NOPTS_VALUE {
		pts = time.Duration(float64(packet.Pts()) * a.stream.Timebase().Q2D() * (float64(time.Second)))
	}
	return pts + a.GetOffset()
}

//decode one packet
func (a *Audio) decode(packet *AVPacket) {
	packetSize := packet.Size()
	pts := a.getPts(packet)

	//decode frame from this packet, there may be many frames in one packet
	for packetSize > 0 { //continue decode until packet is empty
		gotFrame, size := a.codecCtx.DecodeAudio(a.frame, packet)
		if size >= 0 {
			packetSize -= size
			if gotFrame {
				data := resampleFrame(a.resampleCtx, a.frame, a.codecCtx)
				if !data.IsNil() {
					pts = a.audioBuffer.append(&samples{data.Bytes(), pts})
					data.Free()
				}
			}
		} else {
			log.Print("audio decode error")
			break
		}
	}
}

func min2(a, b int) int {
	if a < b {
		return a
	} else {
		return b
	}
}

func (a *Audio) getSilence(size int) []byte {
	if len(a.silence) < size {
		a.silence = append(a.silence, make([]byte, size-len(a.silence))...)
	}

	return a.silence[:size]
}

func (a *Audio) Open(stream AVStream) error {
	a.stream = stream
	codecCtx := stream.Codec()
	a.codecCtx = &codecCtx
	a.audioBuffer.bytesPerSec = 2 * GetBytesPerSample(AV_SAMPLE_FMT_S16) * a.codecCtx.SampleRate()

	a.diffThreshold = 50 * time.Millisecond

	a.quit = make(chan struct{})

	decoder := codecCtx.FindDecoder()
	if decoder.IsNil() {
		return fmt.Errorf("Unsupported audio codec!")
	}
	errCode := codecCtx.Open(decoder)
	if errCode < 0 {
		return fmt.Errorf("Open decoder error:%d", errCode)
	}

	log.Print("open audio")
	return a.driver.Open(a.codecCtx.SampleRate(),
		func(length int) []byte {
			if a.c.WaitUntilRunning(a.quit) {
				close(a.chQuitDone)
				return a.getSilence(length)
			}

			for a.audioBuffer.empty() {
				if packet, ok := a.receivePacket(); ok {
					a.decode(packet)
					packet.Free()
				} else {
					//can't receive more packets
					// log.Print("No more audio packets.")
					return a.getSilence(length)
				}

				a.sync(a.audioBuffer.pts())
			}

			ret := a.audioBuffer.cut(length)
			return ret
		})
}

func (a *Audio) Close() {
	log.Print("close audio")

	a.chQuitDone = make(chan struct{})
	close(a.quit)
	<-a.chQuitDone

	// a.FlushBuffer()
	a.codecCtx.Close()
	a.driver.Close()
}

func (a *Audio) FlushBuffer() {
	for {
		select {
		case p := <-a.PacketChan:
			p.Free()
			break
		default:
			a.audioBuffer.clear()
			return
		}
	}
}

func (a *Audio) IncreaseVolume() float64 {
	return a.driver.IncreaseVolume()
}
func (a *Audio) DecreaseVolume() float64 {
	return a.driver.DecreaseVolume()
}
