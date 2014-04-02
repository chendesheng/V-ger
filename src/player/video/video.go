package video

import (
	"errors"
	"fmt"
	// "io/ioutil"
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
	frame AVFrame
	// pictureRGBs         [8]AVPicture
	pictureObjects      [8]*AVObject
	currentPictureIndex int

	Width, Height int
	c             *Clock

	// ChanPacket  chan *AVPacket
	ChanDecoded chan *VideoFrame
	ChanFlush   chan bool
	flushQuit   chan bool
	quit        chan bool
	r           VideoRender
}

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

	return nil
}

func (v *Video) setupPictureRGB() {
	for i, _ := range v.pictureObjects {
		obj := AVObject{}
		obj.Malloc(v.Width * v.Height * 2)
		println("setup picture objects", obj.Size())
		v.pictureObjects[i] = &obj
	}
	// for i, _ := range v.pictureRGBs {
	// 	numBytes := AVPictureGetSize(AV_PIX_FMT_RGB24, v.Width, v.Height)
	// 	picFrame := AllocFrame()
	// 	pictureRGB := picFrame.Picture()
	// 	pictureRGBBuffer := AVObject{}
	// 	pictureRGBBuffer.Malloc(numBytes)
	// 	pictureRGB.Fill(pictureRGBBuffer, AV_PIX_FMT_RGB24, v.Width, v.Height)

	// 	v.pictureRGBs[i] = pictureRGB
	// }
}

func (v *Video) getPictureObject() *AVObject {
	obj := v.pictureObjects[v.currentPictureIndex]
	v.currentPictureIndex++
	v.currentPictureIndex = v.currentPictureIndex % len(v.pictureObjects)
	return obj
}

// func (v *Video) getPictureRGB() AVPicture {
// 	pic := v.pictureRGBs[v.currentPictureIndex]
// 	v.currentPictureIndex++
// 	v.currentPictureIndex = v.currentPictureIndex % len(v.pictureRGBs)
// 	return pic
// }

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

	println("setupSwsContext", v.Width, v.Height, v.codec.PixelFormat())
	v.swsCtx = SwsGetContext(width, v.Height, AV_PIX_FMT_YUV420P,
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

	v.ChanDecoded = make(chan *VideoFrame, 10)
	v.ChanFlush = make(chan bool)
	v.flushQuit = make(chan bool)
	v.quit = make(chan bool)

	log.Print("new video success")
	return v, nil
}

func (v *Video) Decode(packet *AVPacket) (bool, time.Duration) {
	if v.stream.Index() != packet.StreamIndex() {
		return false, 0
	}

	codec := v.codec
	frame := v.frame

	frameFinished := codec.DecodeVideo(frame, packet)

	if frameFinished {
		//TODO: get pts in more accurate way
		var pts time.Duration
		if packet.Dts() != AV_NOPTS_VALUE {
			pts = time.Duration(float64(packet.Dts()) * v.stream.Timebase().Q2D() * (float64(time.Second)))
		} else if packet.Pts() != AV_NOPTS_VALUE {
			pts = time.Duration(float64(packet.Pts()) * v.stream.Timebase().Q2D() * (float64(time.Second)))
		}

		return true, pts
	}

	return false, 0
}

//small seek
func (v *Video) SeekOffset(t time.Duration) (time.Duration, []byte, error) {
	flags := AVSEEK_FLAG_FRAME
	if t < v.c.GetTime() {
		flags |= AVSEEK_FLAG_BACKWARD
	}
	ctx := v.formatCtx
	err := ctx.SeekFrame(v.stream, t, flags)
	if err != nil {
		return t, nil, err
	}

	timeAfterSeek, img, err := v.DropFramesUtil(t)
	if timeAfterSeek > t+time.Second {
		err = ctx.SeekFrame(v.stream, t, flags|AVSEEK_FLAG_BACKWARD)
		if err != nil {
			return t, nil, err
		}
		timeAfterSeek, _, err = v.DropFramesUtil(t)
	}
	return timeAfterSeek, img, err
}

func (v *Video) Seek(t time.Duration) (time.Duration, []byte, error) {
	// log.Print("video seek ", t.String())

	flags := AVSEEK_FLAG_FRAME
	// if t < v.c.GetTime() {
	// 	flags |= AVSEEK_FLAG_BACKWARD
	// }
	ctx := v.formatCtx
	err := ctx.SeekFrame(v.stream, t, flags)
	if err != nil {
		return t, nil, err
	}

	// return t, nil, nil
	timeAfterSeek, img, err := v.DropFramesUtil(t)
	return timeAfterSeek, img, err
}

func (v *Video) DecodeAndScale(packet *AVPacket) (bool, time.Duration, []byte) {
	if v.stream.Index() != packet.StreamIndex() {
		return false, 0, nil
	}

	if frameFinished, pts := v.Decode(packet); frameFinished {
		frame := v.frame
		// pictureRGB := v.getPictureRGB()
		// swsCtx := v.swsCtx
		width, height := v.Width, v.Height

		// frame.Flip(height)
		// swsCtx.Scale(frame, pictureRGB)

		// return true, pts, pictureRGB.RGBBytes(width, height)

		pic := frame.Picture()
		obj := v.getPictureObject()
		pic.Layout(AV_PIX_FMT_YUV420P, width, height, *obj)
		return true, pts, obj.Bytes()
	}

	return false, 0, nil
}

func (v *Video) FlushBuffer() {
	log.Print("video flush buffer")
	for {
		select {
		case <-v.ChanDecoded:
			break
		default:
			close(v.flushQuit)
			v.flushQuit = make(chan bool)
			return
		}
	}
}

func (v *Video) Play() {
	for {
		select {
		case data := <-v.ChanDecoded:
			if v.c.WaitUtilWithQuit(data.Pts, v.flushQuit) {
				continue
			}

			// log.Printf("playing:%s,%s", data.Pts.String(), v.c.GetTime())

			v.r.SendDrawImage(data.Img)

			if v.c.WaitUtilRunning(v.quit) {
				return
			}
			break
		case <-v.quit:
			return
		}
	}
}
func (v *Video) SetRender(r VideoRender) {
	v.r = r
}

func (v *Video) DropFramesUtil(t time.Duration) (time.Duration, []byte, error) {
	packet := AVPacket{}
	ctx := v.formatCtx
	width, height := v.Width, v.Height
	frame := v.frame
	// pictureRGB := v.getPictureRGB()
	// swsCtx := v.swsCtx

	for ctx.ReadFrame(&packet) >= 0 {
		if frameFinished, pts := v.Decode(&packet); frameFinished {

			// println("pts:", pts.String())
			packet.Free()

			if t-pts < 0*time.Millisecond {
				// frame.Flip(height)
				// swsCtx.Scale(frame, pictureRGB)

				pic := frame.Picture()
				obj := v.getPictureObject()
				pic.Layout(AV_PIX_FMT_YUV420P, width, height, *obj)
				return pts, obj.Bytes(), nil

				// pd := frame.DataObject()
				// pd.SetSize(width*height + width*height/2)
				// pd.Bytes()
				// picYUV.SaveToPPMFile("test.yuv", width, height)
				// ioutil.WriteFile("test.yuv", picYUV.RGBBytes(width, height), 0666)
				// println(len(pd.Bytes()))
				// println(width, height)
				// ioutil.WriteFile("test.yuv", obj.Bytes(), 0666)
				// log.Fatal("yes")

				// return pts, pictureRGB.RGBBytes(width, height), nil
			}
		} else {
			packet.Free()
		}
	}

	return t, nil, errors.New("drop frame error")
}

func (v *Video) Close() {
	v.FlushBuffer()
	close(v.quit)

	v.swsCtx.Free()
	v.frame.Free()

	// for _, pic := range v.pictureRGBs {
	// 	pic.FreeBuffer()
	// 	f := pic.Frame()
	// 	f.Free()
	// }

	for _, obj := range v.pictureObjects {
		obj.Free()
	}

	v.codec.Close()
}
