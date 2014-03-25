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

	driver *sdlAudio

	silence []byte
}

func NewAudio(c *Clock, volume byte) *Audio {
	a := &Audio{}

	resampleCtx := AVAudioResampleContext{}
	resampleCtx.Alloc()

	a.frame = AllocFrame()
	a.resampleCtx = resampleCtx
	a.PacketChan = make(chan *AVPacket, 200)
	a.c = c
	a.driver = &sdlAudio{volume: volume}

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
	pts := time.Duration(float64(packet.Dts()) * a.stream.Timebase().Q2D() * float64(time.Second))
	now := a.c.GetSeekTime()

	if time.Duration(math.Abs(float64(pts-now))) < 50*time.Millisecond {
		return true
	} else if pts > now {
		return !a.c.WaitUtilWithQuit(pts, a.quit)
	} else {
		log.Print("skip audio packet:", pts.String())
		return false
	}
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
				data := resampleFrame(a.resampleCtx, a.frame, a.codecCtx.Channels())
				if !data.IsNil() {
					defer data.Free()
					fn(data.Bytes())
				}
			}
		} else {
			break
		}
	}
}
func (a *Audio) Open(stream AVStream) error {
	a.stream = stream
	codecCtx := stream.Codec()
	a.codecCtx = &codecCtx

	a.quit = make(chan bool)

	decoder := codecCtx.FindDecoder()
	if decoder.IsNil() {
		return fmt.Errorf("Unsupported audio codec!")
	}
	errCode := codecCtx.Open(decoder)
	if errCode < 0 {
		return fmt.Errorf("Open decoder error:%d", errCode)
	}
	println("open audio driver")
	return a.driver.Open(a.codecCtx.Channels(), a.codecCtx.SampleRate(),
		func(length int) []byte {
			if a.c.WaitUtilRunning(a.quit) {
				return nil
			}

			for len(a.audioBuffer) < length {
				if packet, ok := a.receivePacket(); ok {
					if a.sync(packet) {
						a.decode(packet, func(bytes []byte) {
							a.audioBuffer = append(a.audioBuffer, bytes...)
						})
					}
					packet.Free()
				} else {
					//can't receive more packets
					log.Print("No more audio packets.")
					return nil
				}
			}
			retlen := int(math.Min(float64(len(a.audioBuffer)), float64(length)))
			ret := a.audioBuffer[:retlen]
			a.audioBuffer = a.audioBuffer[retlen:]
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
	log.Print("audio flush buffer")
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

func (a *Audio) IncreaseVolume() byte {
	return a.driver.IncreaseVolume()
}
func (a *Audio) DecreaseVolume() byte {
	return a.driver.DecreaseVolume()
}
