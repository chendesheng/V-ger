package libav

/*
#include "libavcodec/avcodec.h"
#include "libavformat/avformat.h"
#include "libavutil/avutil.h"
#include "libavutil/mathematics.h"
#include <stdlib.h>
*/
import "C"

type AVCodecContext struct {
	ptr *C.AVCodecContext
}

func (ctx *AVCodecContext) IsNil() bool {
	return ctx.ptr == nil
}

func (ctx *AVCodecContext) FindDecoder() AVCodec {
	return AVCodec{ptr: C.avcodec_find_decoder(ctx.ptr.codec_id)}
}

func (ctx *AVCodecContext) Open(codec AVCodec) int {
	return int(C.avcodec_open2(ctx.ptr, codec.ptr, nil))
}

func (ctx *AVCodecContext) SampleRate() int {
	return int(ctx.ptr.sample_rate)
}

func (ctx *AVCodecContext) SampleFormat() int {
	return int(ctx.ptr.sample_fmt)
}

func (ctx *AVCodecContext) Channels() int {
	return int(ctx.ptr.channels)
}

func (ctx *AVCodecContext) DecodeAudio(frame AVFrame, packet AVPacket) (bool, int) {
	var gotFrame C.int
	sz := int(C.avcodec_decode_audio4(ctx.ptr, frame.ptr, &gotFrame, packet.pointer()))
	return gotFrame != 0, sz
}

func (ctx *AVCodecContext) Width() int {
	return int(ctx.ptr.width)
}
func (ctx *AVCodecContext) Height() int {
	return int(ctx.ptr.height)
}

func (ctx *AVCodecContext) PixelFormat() int {
	return int(ctx.ptr.pix_fmt)
}

func (ctx *AVCodecContext) DecodeVideo(frame AVFrame, packet AVPacket) bool {
	// println("decode video ", frame.ptr)
	// println("decode video ", ctx.ptr)
	var gotFrame C.int
	C.avcodec_decode_video2(ctx.ptr, frame.ptr, &gotFrame, packet.pointer())
	return gotFrame != 0
}

func (ctx *AVCodecContext) Timebase() AVRational {
	return AVRational(ctx.ptr.time_base)
}

func (ctx *AVCodecContext) FlushBuffer() {
	C.avcodec_flush_buffers(ctx.ptr)
}

func (ctx *AVCodecContext) Close() {
	C.avcodec_close(ctx.ptr)
}
