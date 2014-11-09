package video

import (
	"errors"
	"fmt"
	"log"
	"math"
	"time"
	"vger/player/clock"
	"vger/player/libav"
)

type VideoRender interface {
	Draw([]byte)
	SendShowSpinning()
	SendHideSpinning(bool)
}
type SubRender interface {
	SeekPlayingSubs(time.Duration)
}
type VideoFrame struct {
	Pts time.Duration
	Img []byte
}
type Video struct {
	formatCtx   libav.AVFormatContext
	stream      libav.AVStream
	StreamIndex int
	codec       libav.AVCodecContext

	frame               libav.AVFrame
	imageData           *libav.AVObject
	currentPictureIndex int

	Width, Height int
	c             *clock.Clock

	// ChanDecoded chan *VideoFrame
	ChPackets chan *libav.AVPacket

	flushQuit  chan struct{}
	quit       chan struct{}
	chQuitDone chan struct{}
	r          VideoRender
	sr         SubRender

	global_pts uint64 //for avframe only

	numFaultyDts int
	numFaultyPts int

	lastPts float64
	lastDts float64

	chHold chan struct{}
	chEOF  chan struct{}
}

func (v *Video) SendEOF() {
	select {
	case v.chEOF <- struct{}{}:
	case <-time.After(20 * time.Millisecond):
	case <-v.quit:
	}
}

func (v *Video) setupCodec(codec libav.AVCodecContext) error {
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
	pts = libav.AV_NOPTS_VALUE

	if dts != libav.AV_NOPTS_VALUE {
		if dts <= v.lastDts {
			v.numFaultyDts += 1
		}
		v.lastDts = dts
	}
	if reorderedPts != libav.AV_NOPTS_VALUE {
		if reorderedPts <= v.lastPts {
			v.numFaultyPts += 1
		}
		v.lastPts = reorderedPts
	}
	if (v.numFaultyPts <= v.numFaultyDts || dts == libav.AV_NOPTS_VALUE) && reorderedPts != libav.AV_NOPTS_VALUE {
		pts = reorderedPts
	} else {
		pts = dts
	}

	return pts
}

func (v *Video) setupPictureRGB() {
	sz := libav.AVPictureGetSize(libav.AV_PIX_FMT_YUV420P, v.Width, v.Height)

	obj := libav.AVObject{}
	obj.Malloc(sz)
	v.imageData = &obj
}

func NewVideo(formatCtx libav.AVFormatContext, stream libav.AVStream, c *clock.Clock, r VideoRender, sr SubRender) (*Video, error) {
	v := &Video{}
	v.formatCtx = formatCtx
	v.stream = stream
	v.StreamIndex = stream.Index()
	v.global_pts = libav.AV_NOPTS_VALUE
	v.r = r
	v.sr = sr

	err := v.setupCodec(stream.Codec())
	if err != nil {
		return nil, err
	}

	v.Width, v.Height = v.codec.Width(), v.codec.Height()

	v.setupPictureRGB()
	v.frame = libav.AllocFrame()
	v.c = c

	v.ChPackets = make(chan *libav.AVPacket, 200)
	v.flushQuit = make(chan struct{})
	v.quit = make(chan struct{})
	v.chHold = make(chan struct{})
	v.chEOF = make(chan struct{})

	log.Print("new video success")
	return v, nil
}

func (v *Video) Decode(packet *libav.AVPacket) (bool, time.Duration) {
	if v.stream.Index() != packet.StreamIndex() {
		return false, 0
	}

	codec := v.codec
	frame := v.frame

	v.global_pts = packet.Pts()
	frameFinished := codec.DecodeVideo(frame, packet)

	if frameFinished {
		pts := v.guessCorrectPts(frame.Pts(), frame.Dts())
		if pts == libav.AV_NOPTS_VALUE {
			pts = 0
		}

		dur := time.Duration(float64(pts) * v.stream.Timebase().Q2D() * (float64(time.Second)))

		return true, dur
	}

	return false, 0
}

func (v *Video) SeekAccurate(t time.Duration) (time.Duration, []byte, error) {
	flags := libav.AVSEEK_FLAG_FRAME | libav.AVSEEK_FLAG_BACKWARD

	ctx := v.formatCtx
	err := ctx.SeekFrame(v.stream, t, flags)
	if err != nil {
		return t, nil, err
	}

	t1, img, err := v.DropFramesUtil(t)
	if err != nil {
		return t, nil, err
	} else {
		return t1, img, nil
	}
}

func (v *Video) Seek(t time.Duration) (time.Duration, []byte, error) {
	v.r.SendShowSpinning()
	defer v.r.SendHideSpinning(false)

	flags := libav.AVSEEK_FLAG_FRAME

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

func (v *Video) DecodeAndScale(packet *libav.AVPacket) (bool, time.Duration, []byte) {
	if v.stream.Index() != packet.StreamIndex() {
		return false, 0, nil
	}
	if frameFinished, pts := v.Decode(packet); frameFinished {
		frame := v.frame
		width, height := v.Width, v.Height
		pic := frame.Picture()
		obj := v.imageData
		pic.Layout(libav.AV_PIX_FMT_YUV420P, width, height, *obj)
		return true, pts, obj.Bytes()
	}

	return false, 0, nil
}

func (v *Video) FlushBuffer() {
	for {
		select {
		case packet := <-v.ChPackets:
			packet.Free()
			break
		default:
			close(v.flushQuit)
			v.flushQuit = make(chan struct{})
			return
		}
	}
}

func (v *Video) Play() {
	defer func() {
		if v.chQuitDone != nil {
			close(v.chQuitDone)
		}
	}()

	for {
		v.r.SendShowSpinning()

		t := time.Now()
		select {
		case packet := <-v.ChPackets:
			v.r.SendHideSpinning(false)

			if frameFinished, pts, img := v.DecodeAndScale(packet); frameFinished {

				d := time.Since(t)
				if d > 30*time.Millisecond {
					log.Print("long decode time:", d.String())
					v.c.AddTime(-d)
				}

				packet.Free()

				// log.Printf("playing:%s,%s", pts.String(), v.c.GetTime())
				select {
				case <-v.chHold:
					select {
					case <-v.chHold:
					case <-v.quit:
						return
					}
				case <-v.c.WaitRunning():
					v.r.SendShowSpinning()
					select {
					case <-v.chHold:
						v.r.SendHideSpinning(false)
						select {
						case <-v.chHold:
						case <-v.quit:
							return
						}
					case <-v.c.WaitUntil(pts):
						v.r.SendHideSpinning(true)
						v.sr.SeekPlayingSubs(v.c.GetTime())
						v.r.Draw(img)
					case <-v.flushQuit:
						v.r.SendHideSpinning(false)
						continue
					case <-v.quit:
						v.r.SendHideSpinning(false)
						return
					}
				case <-v.quit:
					return
				}
			}
		case <-v.chEOF:
			v.r.SendHideSpinning(false)
			select {
			case <-v.chHold:
				select {
				case <-v.chHold:
				case <-v.quit:
					return
				}
			case <-v.quit:
				return
			}
		case <-v.quit:
			v.r.SendHideSpinning(false)
			return
		}
	}
}

func (v *Video) Hold() {
	select {
	case v.chHold <- struct{}{}:
		v.FlushBuffer()
	case <-v.quit:
		if v.chQuitDone != nil {
			close(v.chQuitDone)
		}
	}
}

func (v *Video) Unhold() {
	v.FlushBuffer()
	select {
	case v.chHold <- struct{}{}:
	case <-v.quit:
		if v.chQuitDone != nil {
			close(v.chQuitDone)
		}
	}
}

func (v *Video) ReadOneFrame() (time.Duration, []byte, error) {
	packet := libav.AVPacket{}
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
			obj := v.imageData
			pic.Layout(libav.AV_PIX_FMT_YUV420P, width, height, *obj)
			return pts, obj.Bytes(), nil
		} else {
			packet.Free()
		}
	}

	return 0, nil, fmt.Errorf("read frame error: %x", errCode)
}
func (v *Video) DropFramesUtil(t time.Duration) (time.Duration, []byte, error) {
	packet := libav.AVPacket{}
	ctx := v.formatCtx
	width, height := v.Width, v.Height
	frame := v.frame
	for ctx.ReadFrame(&packet) >= 0 {
		if frameFinished, pts := v.Decode(&packet); frameFinished {

			// log.Print("pts:", pts.String())
			packet.Free()

			if t-pts < 0*time.Millisecond {
				pic := frame.Picture()
				obj := v.imageData
				pic.Layout(libav.AV_PIX_FMT_YUV420P, width, height, *obj)
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
	v.imageData.Free()

	v.codec.Close()
}
