package libav

/*
#include "libavcodec/avcodec.h"
#include "libavformat/avformat.h"
#include "libavutil/avutil.h"
#include "libavutil/mathematics.h"
#include <stdlib.h>

void SetGetBufferCallbackCB(AVCodecContext*);
void SetReleaseBufferCallbackCB(AVCodecContext*);
*/
import "C"
import (
	"reflect"
	"unsafe"
)

type AVCodecContext struct {
	ptr *C.AVCodecContext

	getBufferFunc     GetBufferFunc
	releaseBufferFunc ReleaseBufferFuc
}

type GetBufferFunc func(*AVCodecContext, *AVFrame) int
type ReleaseBufferFuc func(*AVCodecContext, *AVFrame)

func (ctx *AVCodecContext) IsNil() bool {
	return ctx.ptr == nil
}
func (ctx *AVCodecContext) DefaultGetBuffer(frame *AVFrame) int {
	return int(C.avcodec_default_get_buffer(ctx.ptr, frame.ptr))
}
func (ctx *AVCodecContext) DefaultReleaseBuffer(frame *AVFrame) {
	C.avcodec_default_release_buffer(ctx.ptr, frame.ptr)
}

func (ctx *AVCodecContext) SetGetBufferCallback(fn GetBufferFunc) {
	println("aaaaaaaaaaaaa")
	// println("set opaque ", ctx.ptr)
	ctx.ptr.opaque = unsafe.Pointer(ctx)
	ctx.getBufferFunc = fn
	C.SetGetBufferCallbackCB(ctx.ptr)
}
func (ctx *AVCodecContext) SetReleaseBufferCallback(fn ReleaseBufferFuc) {
	ctx.ptr.opaque = unsafe.Pointer(ctx)
	ctx.releaseBufferFunc = fn
	C.SetReleaseBufferCallbackCB(ctx.ptr)
}

//export goGetBufferCallback
func goGetBufferCallback(ctxptr unsafe.Pointer, frameptr unsafe.Pointer) int {
	// println(ctxptr)
	// println(frameptr)
	// println(((*C.AVCodecContext)(ctxptr)).opaque)
	ctx := (*AVCodecContext)(unsafe.Pointer(reflect.NewAt(reflect.TypeOf(AVCodecContext{}), ((*C.AVCodecContext)(ctxptr)).opaque).Pointer()))
	// println("ctx.ptr ", ctx.ptr)
	// println(frameptr)

	return ctx.getBufferFunc(ctx, &AVFrame{ptr: (*C.AVFrame)(frameptr)})
}

//export goReleaseBufferCallback
func goReleaseBufferCallback(ctxptr unsafe.Pointer, frameptr unsafe.Pointer) {
	ctx := (*AVCodecContext)(unsafe.Pointer(reflect.NewAt(reflect.TypeOf(AVCodecContext{}), ((*C.AVCodecContext)(ctxptr)).opaque).Pointer()))
	ctx.releaseBufferFunc(ctx, &AVFrame{ptr: (*C.AVFrame)(frameptr)})
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

func (ctx *AVCodecContext) DecodeAudio(frame AVFrame, packet *AVPacket) (bool, int) {
	var gotFrame C.int
	sz := int(C.avcodec_decode_audio4(ctx.ptr, frame.ptr, &gotFrame, &packet.cAVPacket))
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

func (ctx *AVCodecContext) DecodeVideo(frame AVFrame, packet *AVPacket) bool {
	// println("decode video ", frame.ptr)
	// println("decode video ", ctx.ptr)
	var gotFrame C.int
	C.avcodec_decode_video2(ctx.ptr, frame.ptr, &gotFrame, &packet.cAVPacket)
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
