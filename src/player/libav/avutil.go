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

const (
	AVSEEK_SIZE = 0x10000
)

/**
 * @defgroup channel_masks Audio channel masks
 * @{
 */
const (
	AV_CH_FRONT_LEFT            = 0x00000001
	AV_CH_FRONT_RIGHT           = 0x00000002
	AV_CH_FRONT_CENTER          = 0x00000004
	AV_CH_LOW_FREQUENCY         = 0x00000008
	AV_CH_BACK_LEFT             = 0x00000010
	AV_CH_BACK_RIGHT            = 0x00000020
	AV_CH_FRONT_LEFT_OF_CENTER  = 0x00000040
	AV_CH_FRONT_RIGHT_OF_CENTER = 0x00000080
	AV_CH_BACK_CENTER           = 0x00000100
	AV_CH_SIDE_LEFT             = 0x00000200
	AV_CH_SIDE_RIGHT            = 0x00000400
	AV_CH_TOP_CENTER            = 0x00000800
	AV_CH_TOP_FRONT_LEFT        = 0x00001000
	AV_CH_TOP_FRONT_CENTER      = 0x00002000
	AV_CH_TOP_FRONT_RIGHT       = 0x00004000
	AV_CH_TOP_BACK_LEFT         = 0x00008000
	AV_CH_TOP_BACK_CENTER       = 0x00010000
	AV_CH_TOP_BACK_RIGHT        = 0x00020000
	AV_CH_STEREO_LEFT           = 0x20000000 ///< Stereo downmix.
	AV_CH_STEREO_RIGHT          = 0x40000000 ///< See AV_CH_STEREO_LEFT.
	AV_CH_WIDE_LEFT             = 0x0000000080000000
	AV_CH_WIDE_RIGHT            = 0x0000000100000000
	AV_CH_SURROUND_DIRECT_LEFT  = 0x0000000200000000
	AV_CH_SURROUND_DIRECT_RIGHT = 0x0000000400000000
	AV_CH_LOW_FREQUENCY_2       = 0x0000000800000000

	/** Channel mask value used for AVCodecContext.request_channel_layout
	  to indicate that the user requests the channel order of the decoder output
	  to be the native codec channel order. */
	AV_CH_LAYOUT_NATIVE = 0x8000000000000000

	/**
	 * @}
	 * @defgroup channel_mask_c Audio channel convenience macros
	 * @{
	 * */
	AV_CH_LAYOUT_MONO              = (AV_CH_FRONT_CENTER)
	AV_CH_LAYOUT_STEREO            = (AV_CH_FRONT_LEFT | AV_CH_FRONT_RIGHT)
	AV_CH_LAYOUT_2POINT1           = (AV_CH_LAYOUT_STEREO | AV_CH_LOW_FREQUENCY)
	AV_CH_LAYOUT_2_1               = (AV_CH_LAYOUT_STEREO | AV_CH_BACK_CENTER)
	AV_CH_LAYOUT_SURROUND          = (AV_CH_LAYOUT_STEREO | AV_CH_FRONT_CENTER)
	AV_CH_LAYOUT_3POINT1           = (AV_CH_LAYOUT_SURROUND | AV_CH_LOW_FREQUENCY)
	AV_CH_LAYOUT_4POINT0           = (AV_CH_LAYOUT_SURROUND | AV_CH_BACK_CENTER)
	AV_CH_LAYOUT_4POINT1           = (AV_CH_LAYOUT_4POINT0 | AV_CH_LOW_FREQUENCY)
	AV_CH_LAYOUT_2_2               = (AV_CH_LAYOUT_STEREO | AV_CH_SIDE_LEFT | AV_CH_SIDE_RIGHT)
	AV_CH_LAYOUT_QUAD              = (AV_CH_LAYOUT_STEREO | AV_CH_BACK_LEFT | AV_CH_BACK_RIGHT)
	AV_CH_LAYOUT_5POINT0           = (AV_CH_LAYOUT_SURROUND | AV_CH_SIDE_LEFT | AV_CH_SIDE_RIGHT)
	AV_CH_LAYOUT_5POINT1           = (AV_CH_LAYOUT_5POINT0 | AV_CH_LOW_FREQUENCY)
	AV_CH_LAYOUT_5POINT0_BACK      = (AV_CH_LAYOUT_SURROUND | AV_CH_BACK_LEFT | AV_CH_BACK_RIGHT)
	AV_CH_LAYOUT_5POINT1_BACK      = (AV_CH_LAYOUT_5POINT0_BACK | AV_CH_LOW_FREQUENCY)
	AV_CH_LAYOUT_6POINT0           = (AV_CH_LAYOUT_5POINT0 | AV_CH_BACK_CENTER)
	AV_CH_LAYOUT_6POINT0_FRONT     = (AV_CH_LAYOUT_2_2 | AV_CH_FRONT_LEFT_OF_CENTER | AV_CH_FRONT_RIGHT_OF_CENTER)
	AV_CH_LAYOUT_HEXAGONAL         = (AV_CH_LAYOUT_5POINT0_BACK | AV_CH_BACK_CENTER)
	AV_CH_LAYOUT_6POINT1           = (AV_CH_LAYOUT_5POINT1 | AV_CH_BACK_CENTER)
	AV_CH_LAYOUT_6POINT1_BACK      = (AV_CH_LAYOUT_5POINT1_BACK | AV_CH_BACK_CENTER)
	AV_CH_LAYOUT_6POINT1_FRONT     = (AV_CH_LAYOUT_6POINT0_FRONT | AV_CH_LOW_FREQUENCY)
	AV_CH_LAYOUT_7POINT0           = (AV_CH_LAYOUT_5POINT0 | AV_CH_BACK_LEFT | AV_CH_BACK_RIGHT)
	AV_CH_LAYOUT_7POINT0_FRONT     = (AV_CH_LAYOUT_5POINT0 | AV_CH_FRONT_LEFT_OF_CENTER | AV_CH_FRONT_RIGHT_OF_CENTER)
	AV_CH_LAYOUT_7POINT1           = (AV_CH_LAYOUT_5POINT1 | AV_CH_BACK_LEFT | AV_CH_BACK_RIGHT)
	AV_CH_LAYOUT_7POINT1_WIDE      = (AV_CH_LAYOUT_5POINT1 | AV_CH_FRONT_LEFT_OF_CENTER | AV_CH_FRONT_RIGHT_OF_CENTER)
	AV_CH_LAYOUT_7POINT1_WIDE_BACK = (AV_CH_LAYOUT_5POINT1_BACK | AV_CH_FRONT_LEFT_OF_CENTER | AV_CH_FRONT_RIGHT_OF_CENTER)
	AV_CH_LAYOUT_OCTAGONAL         = (AV_CH_LAYOUT_5POINT0 | AV_CH_BACK_LEFT | AV_CH_BACK_CENTER | AV_CH_BACK_RIGHT)
	AV_CH_LAYOUT_STEREO_DOWNMIX    = (AV_CH_STEREO_LEFT | AV_CH_STEREO_RIGHT)
)
