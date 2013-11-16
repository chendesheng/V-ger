package main

import (
	. "player/clock"
	// "github.com/go-gl/gl"
	. "libav"
	"log"
	"path/filepath"
	// "player/glfw"
	// "runtime"
	"player/gui"
	"sync"
	"time"
)

type video struct {
	sync.RWMutex

	formatCtx AVFormatContext
	codecCtx  *AVCodecContext
	swsCtx    SwsContext

	stream AVStream

	frame      AVFrame
	pictureRGB AVPicture

	videoPktPts uint64
	videoClock  float64

	width, height int

	ch chan picture

	window *gui.Window
	status string

	pic picture

	c *Clock
}
type picture struct {
	bytes []byte
	pts   time.Duration
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
	v.ch = make(chan picture)

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

	//run in main thread, safe to operate ui elements
	v.window.FuncDraw = append(v.window.FuncDraw, func() {
		v.draw()
	})
	v.window.FuncKeyDown = append(v.window.FuncKeyDown, func(keycode int) {
		switch keycode {
		case gui.KEY_SPACE:
			v.c.Toggle()
			break
		case gui.KEY_LEFT:
			println("key left pressed")
			v.c.StartSeekTo(-10 * time.Second)
			break
		case gui.KEY_RIGHT:
			println("key right pressed")
			v.c.StartSeekTo(10 * time.Second)
			break
		case gui.KEY_UP:
			v.c.StartSeekTo(time.Minute)
			break
		case gui.KEY_DOWN:
			v.c.StartSeekTo(-time.Minute)
			break
		}
	})
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
		// println(time.Since(b).String())
		// b = time.Now()

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

		// println("pixel format:", codecCtx.PixelFormat())
		// println("width:", v.width, "height:", v.height)

		v.swsCtx.Scale(frame, pictureRGB)

		// println("after scale")

		// tmp := make([]byte, v.pictureSize)
		// copy(tmp, pictureRGB.DataAt(0)[:v.pictureSize])

		// b := time.Now()
		// obj := pictureRGB.Layout(AV_PIX_FMT_RGB24, v.width, v.height)

		// pictureRGB.SaveToPPMFile("a.ppm", v.width, v.height)

		pic := picture{pictureRGB.RGBBytes(v.width, v.height), time.Duration(pts * (float64(time.Second)))}
		v.setPic(pic)
		v.c.WaitUtil(pic.pts)
		v.window.PostEvent(gui.Event{gui.Draw, nil})
	}
}

func (v *video) setPic(pic picture) {
	v.Lock()
	defer v.Unlock()
	// v.pic.Free()
	v.pic = pic
}

func (v *video) getPic() picture {
	v.RLock()
	defer v.RUnlock()

	return v.pic
}

func (v *video) draw() {
	pic := v.getPic()
	if len(pic.bytes) > 0 {
		v.window.Draw(pic.bytes, v.width, v.height)
	} else {
		println("DrawClear")
		// v.window.DrawClear(v.width, v.height)
	}
}

func (v *video) play() {
	// go func() {
	// 	// v.window.PostEvent(Event{Draw, nil})
	// 	// skip := 0
	// 	for pic := range v.ch {
	// 		v.setPic(pic)
	// 		// now := v.c.GetTime()
	// 		// if now-pic.pts > 5*time.Millisecond && skip < 10 {
	// 		// 	skip++
	// 		// 	continue
	// 		// }
	// 		// skip = 0
	// 		v.c.WaitUtil(pic.pts)
	// 		v.window.PostEvent(Event{Draw, nil})
	// 	}
	// }()

	gui.PollEvents()
	return
}
