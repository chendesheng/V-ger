package libav

//#include "libswscale/swscale.h"
import "C"

type SwsContext struct {
	ptr *C.struct_SwsContext
}

// var cachedContext SwsContext

// func SwsGetCachedContext(w int, h int, pixFmt int, outw int, outh int, outPixFmt int, flags int) SwsContext {
// 	cachedContext.ptr = C.sws_getCachedContext(cachedContext.ptr, C.int(w), C.int(h), int32(pixFmt), C.int(outw),
// 		C.int(outh), int32(outPixFmt), C.int(flags), nil, nil, nil)

// 	return cachedContext
// }

func SwsGetContext(w, h int, pixFmt int, outw, outh int, outPixFmt int, flags int) SwsContext {
	return SwsContext{ptr: C.sws_getContext(C.int(w), C.int(h), int32(pixFmt), C.int(outw),
		C.int(outh), int32(outPixFmt), C.int(flags), nil, nil, nil)}
}

func (ctx *SwsContext) Scale(frame1 AVFrame, frame2 AVPicture) {
	C.sws_scale(ctx.ptr, (**C.uint8_t)(&frame1.ptr.data[0]), (*C.int)(&frame1.ptr.linesize[0]),
		0, frame1.ptr.height, (**C.uint8_t)(&frame2.ptr.data[0]), (*C.int)(&frame2.ptr.linesize[0]))
}

const (
	SWS_FAST_BILINEAR = 1 << iota
	SWS_BILINEAR
	SWS_BICUBIC
	SWS_X
	SWS_POINT
	SWS_AREA
	SWS_BICUBLIN
	SWS_GAUSS
	SWS_SINC
	SWS_LANCZOS
	SWS_SPLINE
)
