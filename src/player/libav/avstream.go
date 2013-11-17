package libav

//#include "libavformat/avformat.h"
//#include <stdlib.h>
import "C"

type AVStream struct {
	ptr *C.AVStream
}

func (stream *AVStream) IsNil() bool {
	return stream.ptr == nil
}

func (stream *AVStream) Codec() AVCodecContext {
	return AVCodecContext{ptr: stream.ptr.codec}
}

func (stream *AVStream) Index() int {
	return int(stream.ptr.index)
}
func (stream *AVStream) Timebase() AVRational {
	return AVRational(stream.ptr.time_base)
}
func (stream *AVStream) Duration() int64 {
	return int64(stream.ptr.duration)
}
