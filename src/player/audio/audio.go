package audio

import (
	"fmt"
	"io"
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
	a.quit = make(chan bool)

	return a
}

func (a *Audio) receivePacket() (*AVPacket, bool) {
	var packet *AVPacket
	var ok bool
	for {
		select {
		case packet, ok = <-a.PacketChan:
			if !ok {
				return nil, false
			}

			pts := time.Duration(float64(packet.Dts()) * a.stream.Timebase().Q2D() * float64(time.Second))
			now := a.c.GetSeekTime()

			if now < pts+time.Second && now > pts-time.Second && (now > pts+100*time.Millisecond || now < pts-100*time.Millisecond) {
				if pts > now {
					println("packet pts:", pts.String())
					if a.c.WaitUtilWithQuit(pts, a.quit) {
						packet.Free()
						return nil, false
					} else {
						return packet, true
					}
				}
			} else {
				return packet, true
			}
		case <-a.quit:
			return nil, false
		}
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
	codecCtx := stream.Codec()
	a.codecCtx = &codecCtx

	decoder := codecCtx.FindDecoder()
	if decoder.IsNil() {
		return fmt.Errorf("Unsupported audio codec!")
	}
	errCode := codecCtx.Open(decoder)
	if errCode < 0 {
		return fmt.Errorf("open decoder error code")
	}
	println("open audio driver")
	a.stream = stream
	return a.driver.Open(a.codecCtx.Channels(), a.codecCtx.SampleRate(), func(w io.Writer, length int) {
		defer func() {
			select {
			case <-a.quit:
				log.Printf("Audio close quit return")
				close(a.quitFinish)
			default:
				break
			}
		}()

		if a.c.WaitUtilRunning(a.quit) {
			return
		}

		for length > 0 {
			if len(a.audioBuffer) > 0 {
				writelen, _ := w.Write(a.audioBuffer[:int(math.Min(float64(len(a.audioBuffer)), float64(length)))])
				a.audioBuffer = a.audioBuffer[writelen:]
				length -= writelen
			} else {
				if packet, ok := a.receivePacket(); ok {
					a.decode(packet, func(bytes []byte) {
						a.audioBuffer = append(a.audioBuffer, bytes...)
					})
					packet.Free()
				} else {
					//can't receive more packets
					log.Print("No more audio packets.")
					return
				}
			}
		}
	})
}
func (a *Audio) Close() {
	a.FlushBuffer()
	a.quitFinish = make(chan bool)
	close(a.quit)
	<-a.quitFinish

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
