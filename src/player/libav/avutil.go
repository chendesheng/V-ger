package libav

/*
#include "libavutil/channel_layout.h"
#include "libavutil/samplefmt.h"
#include "stdlib.h"
#include "libavresample/avresample.h"
#include "libavutil/dict.h"
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

type AVDictionary struct {
	ptr *C.AVDictionary
}
type AVDictionaryEntry struct {
	ptr *C.AVDictionaryEntry
}

func (e *AVDictionaryEntry) Key() string {
	return C.GoString(e.ptr.key)
}

func (e *AVDictionaryEntry) Value() string {
	return C.GoString(e.ptr.value)
}
func (e *AVDictionaryEntry) IsNil() bool {
	return e.ptr == nil
}

func (m *AVDictionary) AVDictGet(key string, prev AVDictionaryEntry, flags int) AVDictionaryEntry {
	ckey := C.CString(key)
	defer C.free(unsafe.Pointer(ckey))

	return AVDictionaryEntry{ptr: C.av_dict_get(m.ptr, ckey, prev.ptr, C.int(flags))}
}
func (m *AVDictionary) GetValue(key string) string {
	e := m.AVDictGet(key, AVDictionaryEntry{}, 2)
	return e.Value()
}
func (m *AVDictionary) Map() map[string]string {
	tag := AVDictionaryEntry{}
	res := make(map[string]string)
	for {
		tag = m.AVDictGet("", tag, 2)
		if !tag.IsNil() {
			res[tag.Key()] = tag.Value()
		} else {
			break
		}
	}
	return res
}

const (
	AV_DICT_MATCH_CASE      = 1 << iota
	AV_DICT_IGNORE_SUFFIX   //2
	AV_DICT_DONT_STRDUP_KEY //4   /**< Take ownership of a key that's been
	//          allocated with av_malloc() and children. */
	AV_DICT_DONT_STRDUP_VAL //8   *< Take ownership of a value that's been
	//            allocated with av_malloc() and chilren.
	AV_DICT_DONT_OVERWRITE //16   ///< Don't overwrite existing entries.
	AV_DICT_APPEND         //32   /**< If the entry already exists, append to it.  Note that no
	//     delimiter is added, the strings are simply concatenated. */
)
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
