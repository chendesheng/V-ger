package libav

/*
#include "libavcodec/avcodec.h"
#include "libavformat/avformat.h"
#include "libavutil/avutil.h"
#include "libavutil/mathematics.h"
#include <stdlib.h>
*/
import "C"
import (
	"errors"
	"reflect"
	"time"
	"unsafe"
)

type AVFormatContext struct {
	ptr *C.AVFormatContext
}

func (ctx *AVFormatContext) OpenInput(filename string) {
	cfilename := C.CString(filename)
	defer C.free(unsafe.Pointer(cfilename))

	if int(C.avformat_open_input(&ctx.ptr, cfilename, nil, nil)) != 0 {
		// ctx.ptr = nil
	}
}

func (ctx AVFormatContext) IsNil() bool {
	return ctx.ptr == nil
}

func (ctx *AVFormatContext) DumpFormat() {
	C.av_dump_format(ctx.ptr, 0, &(ctx.ptr.filename[0]), 0)
}

func (ctx *AVFormatContext) FindStreamInfo() error {
	if int(C.avformat_find_stream_info(ctx.ptr, nil)) < 0 {
		return errors.New("Fins stream info error")
	} else {
		return nil
	}
}

func (ctx *AVFormatContext) VideoStream() AVStream {
	var streams []*C.AVStream
	header := (*reflect.SliceHeader)(unsafe.Pointer(&streams))
	header.Len = int(ctx.ptr.nb_streams)
	header.Cap = header.Len
	header.Data = uintptr(unsafe.Pointer(ctx.ptr.streams))

	for i := 0; i < len(streams); i++ {
		// stream := (ctx.ptr.streams)[i]
		stream := streams[i]
		if int(stream.codec.codec_type) == AVMEDIA_TYPE_VIDEO {
			return AVStream{ptr: stream}
		}
	}

	return AVStream{ptr: nil}
}
func (ctx *AVFormatContext) AudioStream() AVStream {
	var streams []*C.AVStream
	header := (*reflect.SliceHeader)(unsafe.Pointer(&streams))
	header.Len = int(ctx.ptr.nb_streams)
	header.Cap = header.Len
	header.Data = uintptr(unsafe.Pointer(ctx.ptr.streams))

	for i := 0; i < len(streams); i++ {
		// stream := (ctx.ptr.streams)[i]
		stream := streams[i]
		if int(stream.codec.codec_type) == AVMEDIA_TYPE_AUDIO {
			return AVStream{ptr: stream}
		}
	}

	return AVStream{ptr: nil}
}
func (ctx *AVFormatContext) Stream(i int) AVStream {
	var streams []*C.AVStream
	header := (*reflect.SliceHeader)(unsafe.Pointer(&streams))
	header.Len = int(ctx.ptr.nb_streams)
	header.Cap = header.Len
	header.Data = uintptr(unsafe.Pointer(ctx.ptr.streams))

	return AVStream{ptr: streams[i]}
}
func (ctx *AVFormatContext) ReadFrame(packet *AVPacket) int {
	return int(C.av_read_frame(ctx.ptr, &packet.cAVPacket))
}

func (ctx *AVFormatContext) SeekFrame(stream AVStream, t time.Duration, flags int) {
	timeBase := stream.ptr.time_base

	seek_pos := t / time.Second * C.AV_TIME_BASE

	var timebaseq C.AVRational
	timebaseq.num = 1
	timebaseq.den = C.AV_TIME_BASE

	seek_target := C.av_rescale_q(C.int64_t(seek_pos), timebaseq, timeBase)
	C.av_seek_frame(ctx.ptr, C.int(stream.Index()), C.int64_t(seek_target), C.int(flags))

	C.avcodec_flush_buffers(stream.Codec().ptr)
}
func (ctx *AVFormatContext) SeekFile(stream AVStream, t time.Duration, flags int) {
	timeBase := stream.ptr.time_base

	seek_target := C.av_rescale(C.int64_t(t/time.Millisecond), C.int64_t(timeBase.den), C.int64_t(timeBase.num)) / 1000
	C.avformat_seek_file(ctx.ptr, C.int(stream.Index()), 0, C.int64_t(seek_target), C.int64_t(seek_target), C.int(flags))

	C.avcodec_flush_buffers(stream.Codec().ptr)
}

func (ctx *AVFormatContext) Duration() int64 {
	return int64(ctx.ptr.duration)
}
