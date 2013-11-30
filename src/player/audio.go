package main

import (
	"log"
	// "math"
	. "player/clock"

	// glfw "github.com/go-gl/glfw3"
	. "player/libav"
	"player/sdl"
	"time"
	// "unsafe"
)

type audio struct {
	formatCtx AVFormatContext
	codecCtx  *AVCodecContext
	stream    AVStream

	frame AVFrame

	audioPktPts uint64
	// audioClock     time.Duration
	// audioClockTime time.Time

	resampleCtx AVAudioResampleContext

	ch chan *AVPacket

	audioBuffer []byte

	c *Clock

	skipBytes int
}

func (a *audio) setup(formatCtx AVFormatContext, stream AVStream) {

	codecCtx := stream.Codec()
	a.codecCtx = &codecCtx

	// codecCtx.SetGetBufferCallback(func(ctx *AVCodecContext, frame *AVFrame) int {
	// 	ret := ctx.DefaultGetBuffer(frame)

	// 	pts := AVObject{}
	// 	pts.Malloc(8)

	// 	pts.WriteUInt64(a.audioPktPts)
	// 	frame.SetOpaque(pts)
	// 	return ret
	// })
	// codecCtx.SetReleaseBufferCallback(func(ctx *AVCodecContext, frame *AVFrame) {
	// 	if !frame.IsNil() {
	// 		pts := frame.Opaque()
	// 		pts.Free()
	// 	}

	// 	ctx.DefaultReleaseBuffer(frame)
	// })

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
	println("audio simple rate:", codecCtx.SimpleRate()) //samples per second
	println("audio simple format:", codecCtx.SimpleFormat())
	println("audio channels:", codecCtx.Channels()) //bytes per samples

	resampleCtx := AVAudioResampleContext{}
	resampleCtx.Alloc()

	a.formatCtx = formatCtx
	a.stream = stream
	a.frame = AllocFrame()
	a.audioPktPts = AV_NOPTS_VALUE
	a.resampleCtx = resampleCtx
	a.ch = make(chan *AVPacket, 200)

	a.initsdl()
}
func abs(i int64) int64 {
	if i < i {
		return -i
	}

	return i
}
func (a *audio) flushBuffer() {
	for {
		select {
		case pkt := <-a.ch:
			pkt.Free()
			println("skip package")
		default:
			println("flush return")
			return
		}
	}
	codec := a.stream.Codec()
	codec.FlushBuffer()
}
func (a *audio) initsdl() {
	codecCtx := a.codecCtx

	if sdl.Init(sdl.SDL_INIT_AUDIO) {
		println("could not init sdl: ", sdl.GetError())
		return
	}

	layout := GetChannelLayout("stereo")
	if codecCtx.Channels() == 1 {
		layout = GetChannelLayout("mono")
	}

	var desired sdl.AudioSpec
	desired.Init()
	desired.SetFreq(codecCtx.SimpleRate())
	desired.SetFormat(sdl.AUDIO_S16LSB)
	desired.SetChannels(uint8(GetChannelLayoutNbChannels(layout)))
	desired.SetSilence(0)
	desired.SetSamples(4096) //audio buffer size

	desired.SetCallback(func(userdata sdl.Object, stream sdl.Object, length int) {
		// println(time.Duration(glfw.GetTime() * float64(time.Second)).String())

		// startTime := time.Duration(0)
		// delayBytes := 0

		codecCtx := a.codecCtx
		frame := a.frame

		for length > 0 {
			if len(a.audioBuffer) > 0 { //try get from buffer first
				// println("read buffer:", len(a.audioBuffer))
				if length < len(a.audioBuffer) { //buffer has enough data
					stream.Write(a.audioBuffer[:length])
					stream.Offset(length)
					a.audioBuffer = a.audioBuffer[length:]
					return
				} else {
					stream.Write(a.audioBuffer)
					stream.Offset(len(a.audioBuffer))
					length -= len(a.audioBuffer)
					// println("offset ", len(a.audioBuffer))
					a.audioBuffer = make([]byte, 0) //make a new block, the last one will be released by gc.
				}
			} else {
				//buffer is empty, fill buffer
				//fill one packet data in one time
				packet, ok := <-a.ch
				// if startTime == 0 {
				// 	startTime = time.Duration(float64(packet.Pts()) * a.stream.Timebase().Q2D() * float64(time.Second))
				// }
				if !ok {
					//write silence audio
					stream.Write(make([]byte, length))
					return
				}

				pts := time.Duration(float64(packet.Pts()) * a.stream.Timebase().Q2D() * float64(time.Second))
				now := a.c.GetTime()
				diffsize := 0
				if now < pts+time.Second && now > pts-time.Second &&
					(now > pts+100*time.Millisecond || now < pts-100*time.Millisecond) {
					println((pts - now).String())
					if pts > now {
						a.c.WaitUtil(pts)
						// diff := float64(pts-now) / float64(time.Second)
						// size := int(diff * float64(a.codecCtx.SimpleRate()*
						// 	a.codecCtx.Channels()*GetBytesPerSample(a.codecCtx.SimpleFormat())))
						// println("wait:", size)

						// zeroes := make([]byte, size)
						// for i, _ := range zeroes {
						// 	zeroes[i] = 0
						// }
						// a.audioBuffer = append(a.audioBuffer, zeroes...)
						// continue
					} else {
						diff := float64(now-pts) / float64(time.Second) / 2
						diffsize = int(diff * float64(a.codecCtx.SimpleRate()*
							a.codecCtx.Channels()*GetBytesPerSample(a.codecCtx.SimpleFormat())))

						// println("diffsize:", diffsize)

						if diffsize >= len(a.audioBuffer) {
							diffsize -= len(a.audioBuffer)
							a.audioBuffer = make([]byte, 0) //skip all
						} else {
							diffsize = 0
							a.audioBuffer = a.audioBuffer[diffsize:]
						}
					}
				}

				packetSize := packet.Size()
				//decode frame from this packet, there may be many frames in one packet
				// firstFrame := true
				for packetSize > 0 { //continue decode until packet is empty
					a.c.WaitUtilRunning()

					gotFrame, size := codecCtx.DecodeAudio(frame, packet)

					// println("decode audio read size: ", size)
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
							if diffsize > 0 {
								// println("diffsize2:", diffsize)
								if diffsize < len(bytes) {
									bytes = bytes[diffsize:]
									diffsize = 0
								} else {
									diffsize -= len(bytes)
									bytes = make([]byte, 0)
								}
								println("diffsize3:", diffsize)
							}
							//TODO: we can avoid this copy by use AVObject as buffer directly
							a.audioBuffer = append(a.audioBuffer, bytes...)
							data.Free()
						}
					}
				}
				packet.Free()
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

// func (a *audio) realtimeClock() time.Duration {
// 	return a.audioClock + time.Now().Sub(a.audioClockTime)
// }

type sample struct {
	data AVObject
	pts  time.Duration
}
