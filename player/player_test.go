package main

import (
	"testing"
	"vger/player/libav"
)

func TestDecodeWMV(t *testing.T) {
	ctx := libav.NewAVFormatContext()
	if err := ctx.OpenInput(""); err != nil {
		return
	}

	if err := ctx.FindStreamInfo(); err != nil {
		return
	}

	ctx.DumpFormat()

	audioStream := ctx.AudioStream()[0]
	avctx := audioStream.Codec()
	frame := libav.AllocFrame()
	codec := avctx.FindDecoder()
	if codec.IsNil() {
		return
	}
	errCode := avctx.Open(codec)
	if errCode < 0 {
		return
	}
	println(avctx.ChannelLayout())
	for {
		packet := libav.AVPacket{}
		resCode := ctx.ReadFrame(&packet)
		println("read frame:", resCode)
		if resCode >= 0 {
			if audioStream.Index() == packet.StreamIndex() {
				println("decode package")
				for packet.Size() > 0 {
					//		println("packet size:", rest)
					gotFrame, sz := avctx.DecodeAudio(frame, &packet)
					if sz < 0 {
						println("decode error")
						return
					} else {
						if gotFrame {
							println("got frame")
						}
						packet.DecodeSize(sz)
					}
				}
			} else {
			}
		}
		packet.Free()
	}
}
