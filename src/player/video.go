package main

import (
	. "player/clock"
	// "github.com/go-gl/gl"
	"log"
	"path/filepath"
	. "player/libav"
	// "player/glfw"
	// "runtime"
	"player/gui"
	// "sync"
	"time"
)

type video struct {
	formatCtx AVFormatContext
	codecCtx  *AVCodecContext
	swsCtx    SwsContext

	stream AVStream

	frame      AVFrame
	pictureRGB AVPicture

	videoPktPts uint64
	videoClock  float64

	width, height int

	window *gui.Window
	status string

	c *Clock
}

func (v *video) setup(formatCtx AVFormatContext, stream AVStream, filename string, start time.Duration) {
	codecCtx := stream.Codec()
	v.codecCtx = &codecCtx

	v.width, v.height = codecCtx.Width(), codecCtx.Height()

	decoder := codecCtx.FindDecoder()
	if decoder.IsNil() {
		println("Unsupported codec!!")
		return
	}
	errCode := codecCtx.Open(decoder)
	if errCode < 0 {
		log.Println("open decoder error code ", errCode)
		return
	}

	v.videoClock = float64(start / time.Second)

	codecCtx.SetGetBufferCallback(func(ctx *AVCodecContext, frame *AVFrame) int {
		ret := ctx.DefaultGetBuffer(frame)

		pts := AVObject{}
		pts.Malloc(8)

		pts.WriteUInt64(v.videoPktPts)
		frame.SetOpaque(pts)
		return ret
	})
	codecCtx.SetReleaseBufferCallback(func(ctx *AVCodecContext, frame *AVFrame) {
		if !frame.IsNil() {
			pts := frame.Opaque()
			pts.Free()
		}

		ctx.DefaultReleaseBuffer(frame)
	})
	// println("source pix format:", codecCtx.PixelFormat())

	numBytes := AVPictureGetSize(AV_PIX_FMT_RGB24, v.width, v.height)

	picFrame := AllocFrame()
	pictureRGB := picFrame.Picture()
	pictureRGBBuffer := AVObject{}
	pictureRGBBuffer.Malloc(numBytes)
	pictureRGB.Fill(pictureRGBBuffer, AV_PIX_FMT_RGB24, v.width, v.height)
	v.pictureRGB = pictureRGB

	v.formatCtx = formatCtx
	v.stream = stream
	v.frame = AllocFrame()
	v.videoPktPts = AV_NOPTS_VALUE
	// v.ch = make(chan picture)

	width := v.width
	if width%4 != 0 {
		/*
			It's a trick for some videos with weired width (like 1278x720), but don't known why it works.
			I got this trick from here:
				http://forum.doom9.org/showthread.php?t=169036
		*/
		width += 1
	}
	v.swsCtx = SwsGetContext(width, v.height, codecCtx.PixelFormat(),
		v.width, v.height, AV_PIX_FMT_RGB24, SWS_BICUBIC)

	v.window = gui.NewWindow(filepath.Base(filename), v.width, v.height)
}

func (v *video) decode(packet *AVPacket) {
	stream := v.stream
	codecCtx := v.codecCtx
	frame := v.frame
	pictureRGB := v.pictureRGB
	// b := time.Now()

	v.videoPktPts = packet.Pts()

	frameFinished := codecCtx.DecodeVideo(frame, packet)

	opaque := frame.Opaque()
	var pts float64
	if packet.Dts() == AV_NOPTS_VALUE &&
		!opaque.IsNil() && opaque.UInt64() != AV_NOPTS_VALUE {
		pts = float64(opaque.UInt64())
	} else if packet.Dts() != AV_NOPTS_VALUE {
		pts = float64(packet.Dts())
	} else {
		pts = 0
	}
	pts *= stream.Timebase().Q2D()
	// println("pts:", pts)
	if frameFinished {

		var frameDelay float64
		if pts != 0 {
			v.videoClock = pts
		} else {
			pts = v.videoClock
		}
		codec := stream.Codec()
		frameDelay = codec.Timebase().Q2D()
		frameDelay += float64(frame.RepeatPict()) * (frameDelay * 0.5)
		v.videoClock += frameDelay

		frame.Flip(v.height)
		v.swsCtx.Scale(frame, pictureRGB)

		t := time.Duration(pts * (float64(time.Second)))
		v.c.WaitUtil(t)

		v.window.ChanDraw <- pictureRGB.RGBBytes(v.width, v.height)
	}
}

func (v *video) play() {
	gui.PollEvents()
	return
}
