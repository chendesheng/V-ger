package libav

/*
#include "libavresample/avresample.h"
*/
import "C"
import (
	"unsafe"
)

type AVAudioResampleContext struct {
	ptr *C.AVAudioResampleContext
}

func (ctx *AVAudioResampleContext) Alloc() {
	ctx.ptr = C.avresample_alloc_context()
}

func (ctx *AVAudioResampleContext) Object() AVObject {
	return AVObject{ptr: unsafe.Pointer(ctx.ptr)}
}

func (ctx *AVAudioResampleContext) Open() int {
	return int(C.avresample_open(ctx.ptr))
}

func (ctx *AVAudioResampleContext) Convert(out AVObject, linesize int, samples int, data unsafe.Pointer,
	inlinesze int, insamples int) int {

	// tmp_out = av_realloc(is->audio_buf1, out_size);
	return int(C.avresample_convert(ctx.ptr, (**C.uint8_t)(unsafe.Pointer(&out.ptr)),
		C.int(linesize), C.int(samples), (**C.uint8_t)(data), C.int(inlinesze), C.int(insamples)))
}

func (ctx *AVAudioResampleContext) Close() {
	C.avresample_close(ctx.ptr)
}
