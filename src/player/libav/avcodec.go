package libav

/*
#include "libavcodec/avcodec.h"
#include "libavformat/avformat.h"
#include "libavutil/avutil.h"
#include "libavutil/mathematics.h"
#include <stdlib.h>
*/
import "C"

type AVCodec struct {
	ptr *C.AVCodec
}

func (codec *AVCodec) IsNil() bool {
	return codec.ptr == nil
}
