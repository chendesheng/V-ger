package main

import (
	// "fmt"
	"log"
	// "math"
	// "math"
	. "player/clock"

	// glfw "github.com/go-gl/glfw3"
	. "player/libav"
	"player/sdl"
	"time"
	// "unsafe"
)

func init() {
	if sdl.Init(sdl.SDL_INIT_AUDIO) {
		println("could not init sdl: ", sdl.GetError())
		return
	}
}

type audio struct {
	streams []AVStream

	codecCtx *AVCodecContext
	stream   AVStream
	frame    AVFrame

	resampleCtx AVAudioResampleContext

	ch chan *AVPacket

	audioBuffer []byte

	c *Clock

	skipBytes int

	quit       chan bool
	quitFinish chan bool
	volume     byte

	audioSpec sdl.AudioSpec

	silence []byte
}

func (a *audio) setCurrentStream(index int) {
	var stream AVStream
	for _, s := range a.streams {
		if s.Index() == index {
			stream = s
		}
	}
	codecCtx := stream.Codec()
	a.codecCtx = &codecCtx

	decoder := codecCtx.FindDecoder()
	if decoder.IsNil() {
		println("Unsupported audio codec!")
		return
	}
	errCode := codecCtx.Open(decoder)
	if errCode < 0 {
		log.Println("open decoder error code ", errCode)
		return
	}
	// println("audio simple rate:", codecCtx.SampleRate()) //samples per second
	// println("audio simple format:", codecCtx.SampleFormat())
	// println("audio channels:", codecCtx.Channels()) //bytes per samples

	resampleCtx := AVAudioResampleContext{}
	resampleCtx.Alloc()

	a.stream = stream
	a.frame = AllocFrame()
	a.resampleCtx = resampleCtx

	if a.ch != nil {
		a.flushBuffer()
	} else {
		a.ch = make(chan *AVPacket, 200)
	}

	a.initsdl()
}
func abs(i int64) int64 {
	if i < i {
		return -i
	}

	return i
}
func (a *audio) flushBuffer() {
	log.Print("audio flush buffer")
	for {
		select {
		case p := <-a.ch:
			p.Free()
			break
		default:
			a.audioBuffer = a.audioBuffer[0:0]
			// codec := a.stream.Codec()
			// codec.FlushBuffer()
			return
		}
	}
	// close(a.ch)
	// a.ch = make(chan *AVPacket, 200)
}
func dropHalfSamples(buffer []byte, bytesPerSample int) []byte {
	j := 0
	for i := 0; i < len(buffer); i += bytesPerSample {
		if (i/bytesPerSample)%3 == 0 {
			for k := 0; k < bytesPerSample; k++ {
				buffer[j+k] = buffer[i+k]
			}
			j += bytesPerSample
		}
	}
	return buffer[:j-bytesPerSample]
}
func (a *audio) skipBuffer(samples int) int {
	// if samples > 64 {
	// 	samples = 64
	// }

	bytesPerSample := GetBytesPerSample(a.codecCtx.SampleFormat())
	buffersamples := len(a.audioBuffer) / bytesPerSample

	if samples >= buffersamples {
		a.audioBuffer = a.audioBuffer[buffersamples*bytesPerSample:] //dropHalfSamples(a.audioBuffer, bytesPerSample) // a.audioBuffer[buffersamples*bytesPerSample:]
	}

	return buffersamples - len(a.audioBuffer)/bytesPerSample
}

func (a *audio) initsdl() {
	codecCtx := a.codecCtx
	a.quit = make(chan bool)

	layout := GetChannelLayout("stereo")
	if codecCtx.Channels() == 1 {
		layout = GetChannelLayout("mono")
	}

	var desired sdl.AudioSpec
	desired.Init()
	desired.SetFreq(codecCtx.SampleRate())
	desired.SetFormat(sdl.AUDIO_S16LSB)
	desired.SetChannels(uint8(GetChannelLayoutNbChannels(layout)))
	desired.SetSilence(0)
	desired.SetSamples(4096) //audio buffer size
	// bytesPerSample := GetBytesPerSample(a.codecCtx.SampleFormat())

	// log.Print(GetBytesPerSample(a.codecCtx.SampleFormat()))

	desired.SetCallback(func(userdata sdl.Object, stream sdl.Object, length int) {
		if a.c.WaitUtilRunning(a.quit) {
			a.codecCtx.Close()
			close(a.quitFinish)
			log.Printf("Audio close quit return")
			return
		}

		if length > len(a.silence) {
			a.silence = make([]byte, length)
		}

		stream.Write(a.silence[:length])

		codecCtx := a.codecCtx
		frame := a.frame

		for length > 0 {
			if len(a.audioBuffer) > 0 { //try get from buffer first
				if length < len(a.audioBuffer) { //buffer has enough data
					// stream.Write(volume(a.audioBuffer[:length], a.volume))
					sdl.MixAudioFormat(stream, a.audioBuffer[:length], a.audioSpec.Format(), int(float64(a.volume)/100*sdl.MIX_MAXVOLUME))
					// fmt.Print(stream.Bytes(200))

					stream.Offset(length)
					a.audioBuffer = a.audioBuffer[length:]
					return
				} else {
					// stream.Write(volume(a.audioBuffer, a.volume))
					// println(a.audioSpec.Format())
					sdl.MixAudioFormat(stream, a.audioBuffer, a.audioSpec.Format(), int(float64(a.volume)/100*sdl.MIX_MAXVOLUME))
					// fmt.Print(stream.Bytes(200))

					stream.Offset(len(a.audioBuffer))
					length -= len(a.audioBuffer)
					a.audioBuffer = a.audioBuffer[0:0]
				}
			} else {
				//buffer is empty, fill buffer
				//fill one packet data in one time
				var packet *AVPacket
			receivePackage:
				for {
					var ok bool
					select {
					case packet, ok = <-a.ch:

						if !ok {
							//write silence audio
							stream.Write(a.silence[:length])
							return
						}

						pts := time.Duration(float64(packet.Dts()) * a.stream.Timebase().Q2D() * float64(time.Second))
						now := a.c.GetSeekTime()

						if now < pts+time.Second && now > pts-time.Second && (now > pts+100*time.Millisecond || now < pts-100*time.Millisecond) {
							if pts > now {
								if a.c.WaitUtilWithQuit(pts, a.quit) {
									a.codecCtx.Close()
									close(a.quitFinish)
									log.Printf("Audio close quit return1")
									return
								}
								// diff := float64(pts-now) / float64(time.Second)
								// silenceLen := int(diff*float64(a.codecCtx.SampleRate())) * a.codecCtx.Channels() * GetBytesPerSample(a.codecCtx.SampleFormat())
								// println("silenceLen:", silenceLen)
								// a.audioBuffer = append(a.audioBuffer, make([]byte, silenceLen)...)
								break receivePackage
							} else {
								// diff := float64(now-pts) / float64(time.Second)
								diffsamples := int(0.1*float64(a.codecCtx.SampleRate())) * a.codecCtx.Channels()
								log.Print("set diff samples:", diffsamples)

								log.Print("skip packet:", pts.String())
								packet.Free()
								// pts = time.Duration(float64(packet.Dts()) * a.stream.Timebase().Q2D() * float64(time.Second))
								// packet, ok = <-a.ch
								// if !ok {
								// 	a.audioBuffer = append(a.audioBuffer, make([]byte, length)...)
								// 	continue
								// }
							}
						} else {
							break receivePackage
						}
					case <-a.quit:
						a.codecCtx.Close()
						close(a.quitFinish)
						log.Printf("Audio close quit return2")
						return
					}
				}

				// println("audio buffer:", len(a.audioBuffer))

				packetSize := packet.Size()
				//decode frame from this packet, there may be many frames in one packet
				for packetSize > 0 { //continue decode until packet is empty
					gotFrame, size := codecCtx.DecodeAudio(frame, packet)
					if size < 0 {
						//skip error frame
						log.Print("error audio frame")
						break
					} else {
						packetSize -= size
						// println("decode audio rest size: ", packetSize)
						if gotFrame {
							data := a.resampleFrame(frame)
							bytes := data.Bytes()

							a.audioBuffer = append(a.audioBuffer, bytes...)
							data.Free()
						}
					}
				}
				packet.Free()

				// println("audio buffer after:", len(a.audioBuffer))

				// if len(a.audioBuffer) > 0 && diffsamples > 0 { //
				// 	log.Printf("sync audio: diffsamples=%d, buffersize=%d,", diffsamples, len(a.audioBuffer))
				// 	diffsamples -= a.skipBuffer(diffsamples)
				// 	log.Printf("after sync: diffsamples=%d, buffersize=%d,", diffsamples, len(a.audioBuffer))
				// }
			}
		}
	})

	res, obtained := sdl.OpenAudio(desired)
	if res < 0 {
		println("sdl open audio error: ", sdl.GetError())
	}
	if obtained.IsNil() {
		println("sdl get nil obtained audio spec")
	}
	a.audioSpec = obtained

	println("sdl play")
	sdl.PauseAudio(0)
}

func (a *audio) resampleFrame(frame AVFrame) AVObject {
	resampleCtx := a.resampleCtx

	resampleCtxObj := resampleCtx.Object()
	resampleCtxObj.SetOptInt("in_channel_layout", int64(frame.ChannelLayout()), 0)
	resampleCtxObj.SetOptInt("in_sample_fmt", int64(frame.Format()), 0)
	resampleCtxObj.SetOptInt("in_sample_rate", int64(frame.SampleRate()), 0)
	resampleCtxObj.SetOptInt("out_channel_layout", int64(GetChannelLayout("stereo")), 0)
	resampleCtxObj.SetOptInt("out_sample_fmt", AV_SAMPLE_FMT_S16, 0)
	resampleCtxObj.SetOptInt("out_sample_rate", int64(frame.SampleRate()), 0)

	if resampleCtx.Open() < 0 {
		println("error initializing libavresample")
		return AVObject{}
	}
	defer resampleCtx.Close()

	osize := GetBytesPerSample(AV_SAMPLE_FMT_S16)
	outSize, outLinesize := AVSampleGetBufferSize(a.codecCtx.Channels(), frame.NbSamples(), frame.Format())
	// println("frame data size:", outSize)

	tmpOut := AVObject{}
	tmpOut.Malloc(outSize)
	// tmpOut := make([]byte, outSize)
	outSamples := resampleCtx.Convert(tmpOut, outLinesize, frame.NbSamples(),
		frame.Data(), frame.Linesize(0), frame.NbSamples())
	if outSamples < 0 {
		println("avresample_convert() failed")
		return AVObject{}
	}
	tmpOut.SetSize(outSamples * osize * 2)
	return tmpOut //must free after copy to buffer
}

func (a *audio) getAudioDelay(packet *AVPacket) {
	var t float64
	if packet.Pts() == AV_NOPTS_VALUE {
		t = 0
	} else {
		t = float64(packet.Pts())
	}

	t *= a.stream.Timebase().Q2D()
	pts := time.Duration(t * float64(time.Second))

	now := a.c.GetTime()
	if now < pts+time.Second && now > pts-time.Second {
		a.c.SetTime(pts)
	}
}

func (a *audio) decode(packet *AVPacket) {
	if a.stream.Index() == packet.StreamIndex() {
		packet.Dup()
		p := *packet
		a.ch <- &p
	}
}

func (a *audio) Close() {
	a.flushBuffer()
	a.quitFinish = make(chan bool)
	close(a.quit)
	<-a.quitFinish

	sdl.CloseAudio()
}

func (a *audio) IncreaseVolume() byte {
	a.volume++

	if a.volume > 100 {
		a.volume = 100
	}
	return a.volume
}
func (a *audio) DecreaseVolume() byte {
	a.volume--

	if a.volume < 0 {
		a.volume = 0
	}
	return a.volume
}

// func (a *audio) realtimeClock() time.Duration {
// 	return a.audioClock + time.Now().Sub(a.audioClockTime)
// }

type sample struct {
	data AVObject
	pts  time.Duration
}
