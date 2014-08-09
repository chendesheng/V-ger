package libav

/*
#include "libavcodec/avcodec.h"
#include "libavformat/avformat.h"
#include "libavutil/avutil.h"
#include "libavutil/mathematics.h"
#include "libavutil/dict.h"

#include <stdlib.h>
*/
import "C"
import (
	"errors"
	"fmt"
	"log"
	"math"
	"reflect"
	// "sync"
	"time"
	"unsafe"
)

type AVFormatContext struct {
	ptr *C.AVFormatContext
}

// var frameLock sync.Mutex = sync.Mutex{}
func NewAVFormatContext() AVFormatContext {
	ptr := C.avformat_alloc_context()
	return AVFormatContext{ptr: ptr}
}

func (ctx AVFormatContext) OpenInput(filename string) error {
	cfilename := C.CString(filename)
	defer C.free(unsafe.Pointer(cfilename))

	if code := int(C.avformat_open_input(&ctx.ptr, cfilename, nil, nil)); code != 0 {
		// ctx.ptr = nil
		return fmt.Errorf("open input error: %x", code)
	}

	return nil
}
func NetworkInit() {
	C.avformat_network_init()
}

func (ctx AVFormatContext) IsNil() bool {
	return ctx.ptr == nil
}

func (ctx AVFormatContext) DumpFormat() {
	C.av_dump_format(ctx.ptr, 0, &(ctx.ptr.filename[0]), 0)
}

func (ctx AVFormatContext) FindStreamInfo() error {
	if code := int(C.avformat_find_stream_info(ctx.ptr, nil)); code < 0 {
		return errors.New(fmt.Sprint("Find stream info error: ", code))
	} else {
		return nil
	}
}

func (ctx AVFormatContext) SetInputFormat(f unsafe.Pointer) {
	ctx.ptr.iformat = (*_Ctype_struct_AVInputFormat)(f)
}

func (ctx AVFormatContext) VideoStream() AVStream {
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
func (ctx AVFormatContext) AudioStream() []AVStream {
	var streams []*C.AVStream
	header := (*reflect.SliceHeader)(unsafe.Pointer(&streams))
	header.Len = int(ctx.ptr.nb_streams)
	header.Cap = header.Len
	header.Data = uintptr(unsafe.Pointer(ctx.ptr.streams))

	res := make([]AVStream, 0)
	for i := 0; i < len(streams); i++ {
		// stream := (ctx.ptr.streams)[i]
		stream := streams[i]
		if int(stream.codec.codec_type) == AVMEDIA_TYPE_AUDIO {
			res = append(res, AVStream{ptr: stream})
		}
	}
	return res
}
func (ctx AVFormatContext) Stream(i int) AVStream {
	var streams []*C.AVStream
	header := (*reflect.SliceHeader)(unsafe.Pointer(&streams))
	header.Len = int(ctx.ptr.nb_streams)
	header.Cap = header.Len
	header.Data = uintptr(unsafe.Pointer(ctx.ptr.streams))

	return AVStream{ptr: streams[i]}
}
func (ctx *AVFormatContext) ReadFrame(packet *AVPacket) int {
	// frameLock.Lock()
	// defer frameLock.Unlock()
	return int(C.av_read_frame(ctx.ptr, &packet.cAVPacket))
}
func (ctx AVFormatContext) SeekFrame2(stream AVStream, t time.Duration, flags int) error {
	var timeBase C.AVRational
	timeBase.num = 1
	timeBase.den = C.AV_TIME_BASE

	seek_pos := t / time.Second * C.AV_TIME_BASE

	var timebaseq C.AVRational
	timebaseq.num = 1
	timebaseq.den = C.AV_TIME_BASE

	seek_target := C.av_rescale_q(C.int64_t(seek_pos), timebaseq, timeBase)
	res := C.av_seek_frame(ctx.ptr, -1, C.int64_t(seek_target), C.int(flags))
	if res < 0 {
		return fmt.Errorf("Seek frame error:", res)
	}

	//this is required! otherwise will get history data after seeking
	C.avcodec_flush_buffers(stream.Codec().ptr)
	return nil
}
func (ctx AVFormatContext) SeekFrame(stream AVStream, t time.Duration, flags int) error {
	// frameLock.Lock()
	// defer frameLock.Unlock()

	// timeBase := stream.ptr.time_base

	// seek_pos := t / time.Second * C.AV_TIME_BASE
	seek_pos := float64(t) / float64(time.Second) * AV_TIME_BASE

	// var timebaseq C.AVRational
	// timebaseq.num = 1
	// timebaseq.den = C.AV_TIME_BASE

	// seek_target := C.av_rescale_q(C.int64_t(seek_pos), timebaseq, timeBase)
	res := C.av_seek_frame(ctx.ptr, -1, C.int64_t(seek_pos), C.int(flags))
	if res < 0 {
		return fmt.Errorf("Seek frame error:", res)
	}

	//this is required! otherwise will get history data after seeking
	C.avcodec_flush_buffers(stream.Codec().ptr)
	return nil
}
func (ctx AVFormatContext) SeekFile(t time.Duration, flags int) error {
	// frameLock.Lock()
	// defer frameLock.Unlock()
	// timeBase := stream.ptr.time_base

	seek_target := float64(t) / float64(time.Second) * AV_TIME_BASE
	if (flags & AVSEEK_FLAG_BYTE) == AVSEEK_FLAG_BYTE {
		seek_target = float64(t)
	}

	// seek_target := C.av_rescale(C.int64_t(t/time.Millisecond), C.int64_t(timeBase.den), C.int64_t(timeBase.num)) / 1000

	ret := int(C.avformat_seek_file(ctx.ptr, -1, C.int64_t(math.MinInt64), C.int64_t(seek_target), C.int64_t(math.MaxInt64), C.int(flags)))
	if ret < 0 {
		return fmt.Errorf("Seek file error:", ret)
	}

	return nil
}

func (ctx AVFormatContext) Duration() uint64 {
	return uint64(ctx.ptr.duration)
}

func (ctx AVFormatContext) Duration2() time.Duration {
	var duration time.Duration
	if ctx.Duration() != AV_NOPTS_VALUE {
		duration = time.Duration(float64(ctx.Duration()) / AV_TIME_BASE * float64(time.Second))
	} else {
		// duration = 2 * time.Hour
		log.Fatal("Can't get video duration.")
	}
	return duration
}

func (ctx AVFormatContext) StartTime() time.Duration {
	return time.Duration((float64(ctx.ptr.start_time) / float64(AV_TIME_BASE)) * float64(time.Second))
}

func (ctx AVFormatContext) SetPb(pb AVIOContext) {
	ctx.ptr.pb = pb.ptr
	ctx.ptr.flags |= AVFMT_FLAG_CUSTOM_IO
}

func (ctx AVFormatContext) Close() {
	C.avformat_close_input(&ctx.ptr)
}

// func (ctx *AVFormatContext) FindStreamInfo(count int) []map[string]string {
// 	var dices *C.AVDictionary
// 	C.avformat_find_stream_info(ctx.ptr, &dices)
// 	// entries := dices.elems
// 	infoes := make([](map[string]string), 0)
// 	for i := 0; i < count; i++ {
// 		var d C.AVDictionary = dices[i]
// 		info := make(map[string]string)
// 		infoes = append(infoes, info)

// 		elems := make([]*C.AVDictionaryEntry, 0)
// 		header := (*reflect.SliceHeader)(unsafe.Pointer(&elems))
// 		header.Len = C.int(d.count)
// 		header.Cap = C.int(d.count)
// 		header.Data = d.elems

// 		// header := reflect.SliceHeader{d.elems, C.int(d.count), C.int(d.count)}
// 		for j := 0; j < len(elems); j++ {
// 			ele := elems[0]
// 			info[C.GoString(ele.key)] = C.GoString(ele.val)
// 		}
// 	}

// 	println("%v", infoes)

// 	return infoes
// }

const (
	AVSEEK_FLAG_BACKWARD = 1 << iota ///< seek backward
	AVSEEK_FLAG_BYTE                 ///< seeking based on position in bytes
	AVSEEK_FLAG_ANY                  ///< seek to any frame, even non-keyframes
	AVSEEK_FLAG_FRAME                ///< seeking based on frame number
)

const (
	AVFMT_FLAG_CUSTOM_IO = 0x0080
)
