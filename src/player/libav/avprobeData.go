package libav

/*
#include "libavformat/avformat.h"
*/
import "C"
import (
	"unsafe"
)

type AVProbeData struct {
	ptr *C.AVProbeData
}

func NewAVProbeData() AVProbeData {
	return AVProbeData{ptr: &C.AVProbeData{}}
}

func (p *AVProbeData) SetBuffer(buf AVObject) {
	p.ptr.buf = (*_Ctype_unsignedchar)(buf.ptr)
	p.ptr.buf_size = C.int(buf.size)
}

func (p *AVProbeData) SetFileName(name string) {
	cname := C.CString(name)
	p.ptr.filename = cname
}

func (p *AVProbeData) InputFormat() unsafe.Pointer {
	return unsafe.Pointer(C.av_probe_input_format(p.ptr, 1))
}
