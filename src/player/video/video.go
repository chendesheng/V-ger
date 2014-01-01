package video

import (
	"errors"
	"fmt"
	. "player/clock"
	// "github.com/go-gl/gl"
	"log"
	// "path/filepath"
	. "player/libav"
	// "player/glfw"
	// "runtime"
	// "player/gui"
	// "sync"
	"time"
)

type VideoRender interface {
	// SendSetSize(int, int)
	SendDrawImage([]byte)
}
type VideoFrame struct {
	Pts time.Duration
	Img []byte
}
type Video struct {
	formatCtx   AVFormatContext
	stream      AVStream
	StreamIndex int
	codec       AVCodecContext
	swsCtx      SwsContext

	//buffers
	frame      AVFrame
	pictureRGB AVPicture

	// videoPktPts uint64
	//for calcuting pts if video file not supply it
	// videoClock time.Duration

	Width, Height int
	c             *Clock

	ChanPacket  chan *AVPacket
	ChanDecoded chan *VideoFrame
	ChanFlush   chan bool
	r           VideoRender
}

// func myGetBuffer(ctx *AVCodecContext, frame *AVFrame) int {
// 	log.Println("SetGetBufferCallback")
// 	log.Println("SetGetBufferCallback ", globalVideo.videoPktPts)
// 	ret := ctx.DefaultGetBuffer(frame)

// 	pts := AVObject{}
// 	pts.Malloc(8)

// 	pts.WriteUInt64(globalVideo.videoPktPts)
// 	frame.SetOpaque(pts)
// 	return ret
// }
// func myReleaseBuffer(ctx *AVCodecContext, frame *AVFrame) {
// 	if !frame.IsNil() {
// 		pts := frame.Opaque()
// 		pts.Free()
// 	}

// 	ctx.DefaultReleaseBuffer(frame)
// }

func (v *Video) setupCodec(codec AVCodecContext) error {
	v.codec = codec

	decoder := codec.FindDecoder()
	if decoder.IsNil() {
		return errors.New("Unsupported codec!!")
	}

	errCode := codec.Open(decoder)
	if errCode < 0 {
		return fmt.Errorf("open decoder error code %s", errCode)
	}

	// codec.SetGetBufferCallback(myGetBuffer)
	// codec.SetReleaseBufferCallback(myReleaseBuffer)

	return nil
}

func (v *Video) setupPictureRGB() {
	numBytes := AVPictureGetSize(AV_PIX_FMT_RGB24, v.Width, v.Height)
	picFrame := AllocFrame()
	pictureRGB := picFrame.Picture()
	pictureRGBBuffer := AVObject{}
	pictureRGBBuffer.Malloc(numBytes)
	pictureRGB.Fill(pictureRGBBuffer, AV_PIX_FMT_RGB24, v.Width, v.Height)
	v.pictureRGB = pictureRGB
}

func (v *Video) setupSwsContext() {
	width := v.Width
	if width%4 != 0 {
		/*
			It's a trick for some videos with weired width (like 1278x720), but don't known why it works.
			I got this trick from here:
				http://forum.doom9.org/showthread.php?t=169036
		*/
		width += 1
	}
	v.swsCtx = SwsGetContext(width, v.Height, v.codec.PixelFormat(),
		v.Width, v.Height, AV_PIX_FMT_RGB24, SWS_BICUBIC)
}

func NewVideo(formatCtx AVFormatContext, stream AVStream, c *Clock) (*Video, error) {
	v := &Video{}
	// globalVideo = v
	v.formatCtx = formatCtx
	v.stream = stream
	v.StreamIndex = stream.Index()

	err := v.setupCodec(stream.Codec())
	if err != nil {
		return nil, err
	}

	v.Width, v.Height = v.codec.Width(), v.codec.Height()

	v.setupPictureRGB()
	v.frame = AllocFrame()

	// v.videoPktPts = AV_NOPTS_VALUE

	v.setupSwsContext()

	v.c = c
	// v.videoClock = start
	// r.SendSetSize(v.Width, v.Height)
	// v.r = r

	v.ChanPacket = make(chan *AVPacket, 10)
	v.ChanDecoded = make(chan *VideoFrame, 10)
	v.ChanFlush = make(chan bool)

	log.Print("new video success")
	return v, nil
}

func (v *Video) Decode(packet *AVPacket) (bool, time.Duration) {
	if v.stream.Index() != packet.StreamIndex() {
		return false, 0
	}

	stream := v.stream
	codec := v.codec
	frame := v.frame
	// pictureRGB := v.pictureRGB
	// b := time.Now()

	// v.videoPktPts = packet.Pts()

	frameFinished := codec.DecodeVideo(frame, packet)

	// opaque := frame.Opaque()
	// var pts time.Duration
	// if packet.Dts() == AV_NOPTS_VALUE &&
	// 	!opaque.IsNil() && opaque.UInt64() != AV_NOPTS_VALUE {
	// 	pts = time.Duration(float64(opaque.UInt64()) * stream.Timebase().Q2D() * (float64(time.Second)))
	// } else if packet.Dts() != AV_NOPTS_VALUE {
	// 	pts = time.Duration(float64(packet.Dts()) * stream.Timebase().Q2D() * (float64(time.Second)))
	// } else {
	// 	pts = 0
	// }

	// println("pts:", pts)
	if frameFinished {
		//TODO: get pts in more safe way
		var pts time.Duration
		if packet.Dts() != AV_NOPTS_VALUE {
			pts = time.Duration(float64(packet.Dts()) * stream.Timebase().Q2D() * (float64(time.Second)))
		}

		// var frameDelay float64
		// if pts != 0 {
		// 	v.videoClock = pts
		// } else {
		// 	pts = v.videoClock
		// }
		// codec := stream.Codec()
		// frameDelay = codec.Timebase().Q2D()
		// frameDelay += float64(frame.RepeatPict()) * (frameDelay * 0.5)
		// v.videoClock += time.Duration(frameDelay * float64(time.Second))

		// frame.Flip(v.height)
		// v.swsCtx.Scale(frame, pictureRGB)

		return true, pts
	}

	return false, 0
}

//small seek
func (v *Video) SeekOffset(t time.Duration) time.Duration {
	flags := AVSEEK_FLAG_FRAME
	if t < v.c.GetTime() {
		flags |= AVSEEK_FLAG_BACKWARD
	}
	ctx := v.formatCtx
	ctx.SeekFrame(v.stream, t, flags)

	timeAfterSeek, _ := v.DropFramesUtil(t)
	if timeAfterSeek > t+time.Second {
		ctx.SeekFrame(v.stream, t, flags|AVSEEK_FLAG_BACKWARD)
		timeAfterSeek, _ = v.DropFramesUtil(t)
	}
	return timeAfterSeek
}

func (v *Video) Seek(t time.Duration) (time.Duration, []byte) {
	// log.Print("video seek ", t.String())

	flags := AVSEEK_FLAG_FRAME
	// if t < v.c.GetTime() {
	// 	flags |= AVSEEK_FLAG_BACKWARD
	// }
	ctx := v.formatCtx
	ctx.SeekFrame(v.stream, t, flags)
	// v.videoClock = t

	timeAfterSeek, img := v.DropFramesUtil(t)
	return timeAfterSeek, img
}

func (v *Video) DecodeAndScale(packet *AVPacket) (bool, time.Duration, []byte) {
	if v.stream.Index() != packet.StreamIndex() {
		return false, 0, nil
	}

	frame := v.frame
	pictureRGB := v.pictureRGB
	swsCtx := v.swsCtx
	width, height := v.Width, v.Height
	// c := v.c
	// r := v.r

	if frameFinished, pts := v.Decode(packet); frameFinished {
		frame.Flip(height)
		swsCtx.Scale(frame, pictureRGB)

		return true, pts, pictureRGB.RGBBytes(width, height)
	}

	return false, 0, nil
}

func (v *Video) FlushBuffer() {
	// close(v.ChanPacket)
	// v.ChanPacket = make(chan *AVPacket, 200)
	// v.ChanFlush <- true
	for {
		select {
		case <-v.ChanDecoded:
			break
		default:
			return
		}
	}

}

func (v *Video) Play() {
	for data := range v.ChanDecoded {
		v.c.WaitUtil(data.Pts)
		log.Printf("playing:%s,%s", data.Pts.String(), v.c.GetTime())
		v.r.SendDrawImage(data.Img)

		v.c.WaitUtilRunning()
	}
}
func (v *Video) SetRender(r VideoRender) {
	v.r = r
}

// func readOneFrame(ctx AVFormatContext, stream AVStream, frame AVFrame) (time.Duration, bool) {
// 	packet := AVPacket{}
// 	codec := stream.Codec()

// 	for ctx.ReadFrame(&packet) >= 0 {
// 		if packet.StreamIndex() == stream.Index() {
// 			if codec.DecodeVideo(frame, &packet) {
// 				tmp := packet.Pts()
// 				if tmp == AV_NOPTS_VALUE {
// 					tmp = 0
// 				}

// 				pts := time.Duration(float64(tmp) * stream.Timebase().Q2D() * float64(time.Second))
// 				println("pts:", pts.String())
// 				packet.Free()

// 				return pts, true
// 			} else {
// 				packet.Free()
// 			}
// 		}
// 	}

// 	return 0, false
// }

func (v *Video) DropFramesUtil(t time.Duration) (time.Duration, []byte) {
	packet := AVPacket{}
	ctx := v.formatCtx
	width, height := v.Width, v.Height
	frame := v.frame
	pictureRGB := v.pictureRGB
	swsCtx := v.swsCtx

	for ctx.ReadFrame(&packet) >= 0 {
		if frameFinished, pts := v.Decode(&packet); frameFinished {

			// println("pts:", pts.String())
			packet.Free()

			if t-pts < 0*time.Millisecond {
				frame.Flip(height)
				swsCtx.Scale(frame, pictureRGB)

				return pts, pictureRGB.RGBBytes(width, height)
			}
		} else {
			packet.Free()
		}
	}

	return 0, nil
}
