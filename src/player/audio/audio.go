package audio

import (
	"fmt"
	"log"
	"math"
	. "player/clock"
	. "player/libav"
	"time"
)

type Audio struct {
	c           *Clock
	codecCtx    *AVCodecContext
	resampleCtx AVAudioResampleContext
	stream      AVStream
	frame       AVFrame

	PacketChan  chan *AVPacket
	audioBuffer []byte

	skipBytes int

	quit       chan bool
	quitFinish chan bool

	driver *portAudio

	silence []byte

	diffCnt       int
	diffCum       float64
	diffCoef      float64
	diffThreshold time.Duration

	decodePkt     *AVPacket
	decodePktSize int
}

func NewAudio(c *Clock, volume float64) *Audio {
	a := &Audio{}

	resampleCtx := AVAudioResampleContext{}
	resampleCtx.Alloc()

	a.frame = AllocFrame()
	a.resampleCtx = resampleCtx
	a.PacketChan = make(chan *AVPacket, 200)
	a.c = c
	a.driver = &portAudio{volume: volume}

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
		return nil, false
	}
}

//drop packet if return false
func (a *Audio) sync(packet *AVPacket) bool {
	// println("sync audio")
	var pts time.Duration
	if packet.Pts() != AV_NOPTS_VALUE {
		pts = time.Duration(float64(packet.Pts()) * a.stream.Timebase().Q2D() * (float64(time.Second)))
	}

	now := a.c.GetTime()

	diff := pts - now
	avgDiff := time.Duration(math.Abs(float64(diff)))

	if avgDiff < a.diffThreshold {
		return true
	} else if pts > now && pts-now < 500*time.Millisecond {
		return !a.c.WaitUtilWithQuit(pts, a.quit)
	} else {
		log.Print("skip audio packet:", (now - pts).String())
		return false
	}
}

func (a *Audio) getClock() time.Duration {
	var pts time.Duration
	if a.decodePkt.Pts() != AV_NOPTS_VALUE {
		pts = time.Duration(float64(a.decodePkt.Pts()) * a.stream.Timebase().Q2D() * (float64(time.Second)))
	}
	return pts
}

//decode one packet
func (a *Audio) decode(packet *AVPacket, fn func([]byte)) {
	packetSize := packet.Size()
	//decode frame from this packet, there may be many frames in one packet
	for packetSize > 0 { //continue decode until packet is empty
		gotFrame, size := a.codecCtx.DecodeAudio(a.frame, packet)
		if size >= 0 {
			packetSize -= size
			if gotFrame {
				data := resampleFrame(a.resampleCtx, a.frame, a.codecCtx)
				if !data.IsNil() {
					defer data.Free()
					fn(data.Bytes())
				}
			}
		} else {
			println("audio decode error")
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

	a.diffThreshold = 50 * time.Millisecond

	a.quit = make(chan bool)

	decoder := codecCtx.FindDecoder()
	if decoder.IsNil() {
		return fmt.Errorf("Unsupported audio codec!")
	}
	errCode := codecCtx.Open(decoder)
	if errCode < 0 {
		return fmt.Errorf("Open decoder error:%d", errCode)
	}

	println("open audio")
	return a.driver.Open(a.codecCtx.Channels(), a.codecCtx.SampleRate(),
		func(length int) []byte {
			if a.c.WaitUtilRunning(a.quit) {
				return a.getSilence(length)
			}

			for len(a.audioBuffer) == 0 {
				if packet, ok := a.receivePacket(); ok {
					a.decode(packet, func(bytes []byte) {
						if a.sync(packet) {
							a.audioBuffer = append(a.audioBuffer, bytes...)
						}
					})

					packet.Free()
				} else {
					//can't receive more packets
					// log.Print("No more audio packets.")
					return a.getSilence(length)
				}
			}

			retLen := min2(len(a.audioBuffer), length)
			ret := a.audioBuffer[:retLen]
			a.audioBuffer = a.audioBuffer[retLen:]
			return ret
		})
}

func (a *Audio) Close() {
	log.Print("close audio")

	close(a.quit)

	a.FlushBuffer()
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
			a.audioBuffer = a.audioBuffer[0:0]
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
