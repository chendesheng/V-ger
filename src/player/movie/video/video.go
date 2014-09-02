package video

import (
	"errors"
	"fmt"
	"log"
	"math"
	. "player/clock"
	. "player/libav"
	"time"
)

type VideoRender interface {
	SendDrawImage([]byte)
	SendShowSpinning()
	SendHideSpinning()
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

	frame               AVFrame
	pictureObjects      [40]*AVObject
	currentPictureIndex int

	Width, Height int
	c             *Clock

	ChanDecoded chan *VideoFrame
	flushQuit   chan struct{}
	quit        chan struct{}
	chQuitDone  chan struct{}
	r           VideoRender

	global_pts uint64 //for avframe only

	numFaultyDts int
	numFaultyPts int

	lastPts float64
	lastDts float64
}

func (v *Video) setupCodec(codec AVCodecContext) error {
	log.Print("setupCodec")

	v.codec = codec

	decoder := codec.FindDecoder()
	if decoder.IsNil() {
		return errors.New("Unsupported codec!!")
	}

	errCode := codec.Open(decoder)
	if errCode < 0 {
		return fmt.Errorf("open decoder error code %d", errCode)
	}

	v.lastPts = math.MinInt64
	v.lastDts = math.MinInt64

	return nil
}

//copy from avplay.c
func (v *Video) guessCorrectPts(reorderedPts float64, dts float64) (pts float64) {
	pts = AV_NOPTS_VALUE

	if dts != AV_NOPTS_VALUE {
		if dts <= v.lastDts {
			v.numFaultyDts += 1
		}
		v.lastDts = dts
	}
	if reorderedPts != AV_NOPTS_VALUE {
		if reorderedPts <= v.lastPts {
			v.numFaultyPts += 1
		}
		v.lastPts = reorderedPts
	}
	if (v.numFaultyPts <= v.numFaultyDts || dts == AV_NOPTS_VALUE) && reorderedPts != AV_NOPTS_VALUE {
		pts = reorderedPts
	} else {
		pts = dts
	}

	return pts
}

func (v *Video) setupPictureRGB() {
	sz := AVPictureGetSize(AV_PIX_FMT_YUV420P, v.Width, v.Height)
	for i, _ := range v.pictureObjects {
		obj := AVObject{}
		obj.Malloc(sz)
		v.pictureObjects[i] = &obj
	}
}

func (v *Video) getPictureObject() *AVObject {
	obj := v.pictureObjects[v.currentPictureIndex]
	v.currentPictureIndex++
	v.currentPictureIndex = v.currentPictureIndex % len(v.pictureObjects)
	return obj
}

func NewVideo(formatCtx AVFormatContext, stream AVStream, c *Clock, r VideoRender) (*Video, error) {
	v := &Video{}
	v.formatCtx = formatCtx
	v.stream = stream
	v.StreamIndex = stream.Index()
	v.global_pts = AV_NOPTS_VALUE
	v.r = r

	err := v.setupCodec(stream.Codec())
	if err != nil {
		return nil, err
	}

	v.Width, v.Height = v.codec.Width(), v.codec.Height()

	v.setupPictureRGB()
	v.frame = AllocFrame()
	v.c = c

	v.ChanDecoded = make(chan *VideoFrame, 20)
	v.flushQuit = make(chan struct{})
	v.quit = make(chan struct{})

	log.Print("new video success")
	return v, nil
}

func (v *Video) Decode(packet *AVPacket) (bool, time.Duration) {
	if v.stream.Index() != packet.StreamIndex() {
		return false, 0
	}

	codec := v.codec
	frame := v.frame

	v.global_pts = packet.Pts()
	frameFinished := codec.DecodeVideo(frame, packet)

	if frameFinished {
		pts := v.guessCorrectPts(frame.Pts(), frame.Dts())
		if pts == AV_NOPTS_VALUE {
			pts = 0
		}

		dur := time.Duration(float64(pts) * v.stream.Timebase().Q2D() * (float64(time.Second)))

		return true, dur
	}

	return false, 0
}

func (v *Video) SeekAccurate(t time.Duration) (time.Duration, []byte, error) {
	flags := AVSEEK_FLAG_FRAME | AVSEEK_FLAG_BACKWARD

	ctx := v.formatCtx
	err := ctx.SeekFrame(v.stream, t, flags)
	if err != nil {
		return t, nil, err
	}

	return v.DropFramesUtil(t)
}

func (v *Video) Seek(t time.Duration) (time.Duration, []byte, error) {
	v.r.SendShowSpinning()
	defer v.r.SendHideSpinning()

	flags := AVSEEK_FLAG_FRAME

	ctx := v.formatCtx
	err := ctx.SeekFrame(v.stream, t, flags)

	if err != nil {
		return t, nil, err
	}

	t1, img, err := v.ReadOneFrame()
	if err != nil {
		return t, nil, err
	} else {
		return t1, img, nil
	}
}

func (v *Video) DecodeAndScale(packet *AVPacket) (bool, time.Duration, []byte) {
	if v.stream.Index() != packet.StreamIndex() {
		return false, 0, nil
	}
	if frameFinished, pts := v.Decode(packet); frameFinished {
		frame := v.frame
		width, height := v.Width, v.Height
		pic := frame.Picture()
		obj := v.getPictureObject()
		pic.Layout(AV_PIX_FMT_YUV420P, width, height, *obj)
		return true, pts, obj.Bytes()
	}

	return false, 0, nil
}

func (v *Video) FlushBuffer() {
	for {
		select {
		case <-v.ChanDecoded:
			break
		default:
			close(v.flushQuit)
			v.flushQuit = make(chan struct{})
			return
		}
	}
}

func (v *Video) Play() {
	for {
		select {
		case data := <-v.ChanDecoded:
			// log.Printf("playing:%s,%s", data.Pts.String(), v.c.GetTime())

			if v.c.WaitUntilWithQuit(data.Pts, v.flushQuit) {
				continue
			}

			v.r.SendDrawImage(data.Img)

			if v.c.WaitUntilRunning(v.quit) {
				return
			}
			break
		case <-v.quit:
			close(v.chQuitDone)
			return
		}
	}
}

func (v *Video) ReadOneFrame() (time.Duration, []byte, error) {
	packet := AVPacket{}
	ctx := v.formatCtx
	width, height := v.Width, v.Height
	frame := v.frame

	errCode := 0
	for {
		errCode = ctx.ReadFrame(&packet)
		if errCode < 0 {
			break
		}

		if frameFinished, pts := v.Decode(&packet); frameFinished {
			packet.Free()

			pic := frame.Picture()
			obj := v.getPictureObject()
			pic.Layout(AV_PIX_FMT_YUV420P, width, height, *obj)
			return pts, obj.Bytes(), nil
		} else {
			packet.Free()
		}
	}

	return 0, nil, fmt.Errorf("read frame error: %x", errCode)
}
func (v *Video) DropFramesUtil(t time.Duration) (time.Duration, []byte, error) {
	packet := AVPacket{}
	ctx := v.formatCtx
	width, height := v.Width, v.Height
	frame := v.frame
	for ctx.ReadFrame(&packet) >= 0 {
		if frameFinished, pts := v.Decode(&packet); frameFinished {

			// log.Print("pts:", pts.String())
			packet.Free()

			if t-pts < 0*time.Millisecond {
				pic := frame.Picture()
				obj := v.getPictureObject()
				pic.Layout(AV_PIX_FMT_YUV420P, width, height, *obj)
				return pts, obj.Bytes(), nil
			}
		} else {
			packet.Free()
		}
	}

	return t, nil, errors.New("drop frame error")
}

func (v *Video) Close() {
	log.Print("close video")

	v.chQuitDone = make(chan struct{})
	close(v.quit)
	<-v.chQuitDone

	v.frame.Free()

	for _, obj := range v.pictureObjects {
		obj.Free()
	}

	v.codec.Close()
}
