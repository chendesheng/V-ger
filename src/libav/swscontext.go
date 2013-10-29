package libav

//#include "libswscale/swscale.h"
import "C"

type SwsContext struct {
	ptr *C.struct_SwsContext
}

func SwsGetCachedContext(w int, h int, pixFmt int, outw int, outh int, outPixFmt int, flags int) SwsContext {
	return SwsContext{ptr: C.sws_getCachedContext(nil, C.int(w), C.int(h), int32(pixFmt), C.int(outw),
		C.int(outh), int32(outPixFmt), C.int(flags), nil, nil, nil)}
}

func (ctx *SwsContext) Scale(frame1 AVFrame, frame2 AVPicture) {
	C.sws_scale(ctx.ptr, (**C.uint8_t)(&frame1.ptr.data[0]), (*C.int)(&frame1.ptr.linesize[0]),
		0, frame1.ptr.height, (**C.uint8_t)(&frame2.ptr.data[0]), (*C.int)(&frame2.ptr.linesize[0]))
}

const (
	SWS_BICUBIC = 4
)
