package libav

/*
#include "libavutil/channel_layout.h"
#include "libavutil/samplefmt.h"
#include "stdlib.h"
#include "libavresample/avresample.h"
*/
import "C"
import (
	"unsafe"
)

const (
	AVMEDIA_TYPE_UNKNOWN = iota - 1
	AVMEDIA_TYPE_VIDEO
	AVMEDIA_TYPE_AUDIO
	AVMEDIA_TYPE_DATA
	AVMEDIA_TYPE_SUBTITLE
	AVMEDIA_TYPE_ATTACHMENT
	AVMEDIA_TYPE_NB
)

const (
	AV_NOPTS_VALUE = 0x8000000000000000
)

func GetChannelLayout(name string) uint64 {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	return uint64(C.av_get_channel_layout(cname))
}

func GetChannelLayoutNbChannels(layout uint64) int {
	return int(C.av_get_channel_layout_nb_channels(C.uint64_t(layout)))
}

func GetBytesPerSample(fmt int) int {
	return int(C.av_get_bytes_per_sample(int32(fmt)))
}

const (
	AV_SAMPLE_FMT_NONE = iota - 1
	AV_SAMPLE_FMT_U8
	AV_SAMPLE_FMT_S16
	AV_SAMPLE_FMT_S32

	AV_SAMPLE_FMT_FLT
	AV_SAMPLE_FMT_DBL
	AV_SAMPLE_FMT_U8P
	AV_SAMPLE_FMT_S16P

	AV_SAMPLE_FMT_S32P
	AV_SAMPLE_FMT_FLTP
	AV_SAMPLE_FMT_DBLP
	AV_SAMPLE_FMT_NB
)

func AVSampleGetBufferSize(channels int, samples int, fmt int) (int, int) {
	var linesize C.int
	sz := int(C.av_samples_get_buffer_size(&linesize, C.int(channels), C.int(samples), int32(fmt), 0))

	return sz, int(linesize)
}

const (
	AV_TIME_BASE = 1000000
)

const (
	AV_PIX_FMT_NONE = iota - 1
	AV_PIX_FMT_YUV420P
	AV_PIX_FMT_YUYV422
	AV_PIX_FMT_RGB24
)
