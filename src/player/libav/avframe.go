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
	"unsafe"
)

type AVFrame struct {
	ptr *C.AVFrame
}

func (frame *AVFrame) SetOpaque(obj AVObject) {
	frame.ptr.opaque = obj.ptr
}

func (frame *AVFrame) Opaque() AVObject {
	return AVObject{ptr: frame.ptr.opaque}
}

func (frame *AVFrame) IsNil() bool {
	return frame.ptr == nil
}

func AllocFrame() AVFrame {
	frame := AVFrame{ptr: C.avcodec_alloc_frame()}
	println("alloc frame ", frame.ptr)
	return frame
}

func (frame *AVFrame) ChannelLayout() int64 {
	return int64(frame.ptr.channel_layout)
}

func (frame *AVFrame) NbSamples() int {
	return int(frame.ptr.nb_samples)
}

func (frame *AVFrame) Format() int {
	return int(frame.ptr.format)
}

func (frame *AVFrame) SampleRate() int {
	return int(frame.ptr.sample_rate)
}

func (frame *AVFrame) Data() unsafe.Pointer {
	return unsafe.Pointer(&frame.ptr.data[0])
}

func (frame *AVFrame) Linesize(i int) int {
	return int(frame.ptr.linesize[i])
}

func (frame *AVFrame) Picture() AVPicture {
	return AVPicture{ptr: (*C.AVPicture)(unsafe.Pointer(frame.ptr))}
}

func (frame *AVFrame) RepeatPict() int {
	return int(frame.ptr.repeat_pict)
}

func (frame *AVFrame) Flip(height int) {
	// for i := 0; i < 3; i++ {
	// 	frame.ptr.data[i] = (*C.uint8_t)(unsafe.Pointer(uintptr(unsafe.Pointer(frame.ptr.data[i])) + uintptr(frame.ptr.linesize[i]*(frame.ptr.height-1))))
	// 	frame.ptr.linesize[i] = -frame.ptr.linesize[i]
	// }

	pic := frame.ptr
	nDivisor := 0
	nMaxLineSize := 0
	// find max linesize
	for i := 0; i < 4; i++ {
		if int(pic.linesize[i]) > nMaxLineSize {
			nMaxLineSize = int(pic.linesize[i])
		}
	}
	if pic.linesize[0] != 0 {
		nDivisor = (nMaxLineSize / int(pic.linesize[0]))
		if nDivisor == 0 {
			nDivisor = 1
		}
		pic.data[0] = (*C.uint8_t)(unsafe.Pointer(uintptr(unsafe.Pointer(pic.data[0])) + uintptr((int(pic.linesize[0]) * ((height / nDivisor) - 1)))))
	}
	if pic.linesize[1] != 0 {
		nDivisor = (nMaxLineSize / int(pic.linesize[1]))
		if nDivisor == 0 {
			nDivisor = 1
		}
		pic.data[1] = (*C.uint8_t)(unsafe.Pointer(uintptr(unsafe.Pointer(pic.data[1])) + uintptr((int(pic.linesize[1]) * ((height / nDivisor) - 1)))))
	}
	if pic.linesize[2] != 0 {
		nDivisor = (nMaxLineSize / int(pic.linesize[2]))
		if nDivisor == 0 {
			nDivisor = 1
		}
		pic.data[2] = (*C.uint8_t)(unsafe.Pointer(uintptr(unsafe.Pointer(pic.data[2])) + uintptr((int(pic.linesize[2]) * ((height / nDivisor) - 1)))))
	}
	if pic.linesize[3] != 0 {
		nDivisor = (nMaxLineSize / int(pic.linesize[3]))
		if nDivisor == 0 {
			nDivisor = 1
		}
		pic.data[3] = (*C.uint8_t)(unsafe.Pointer(uintptr(unsafe.Pointer(pic.data[3])) + uintptr((int(pic.linesize[3]) * ((height / nDivisor) - 1)))))
	}
	pic.linesize[0] *= -1
	pic.linesize[1] *= -1
	pic.linesize[2] *= -1
	pic.linesize[3] *= -1
}
