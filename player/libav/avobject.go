package libav

//#include "libavutil/mem.h"
//#include "libavutil/opt.h"
//#include <string.h>
//#include <stdlib.h>
//void copybytes(void* p, int offset, void* source, int len) {
//	memcpy(p+offset, source, len);
//}
import "C"
import "unsafe"

type AVObject struct {
	ptr    unsafe.Pointer
	length int
	size   int
}

func (obj *AVObject) Malloc(sz int) {
	obj.ptr = C.av_malloc(C.size_t(sz))
	obj.size = sz
}

func (obj *AVObject) Free() {
	if obj.IsNil() {
		return
	}

	C.av_free(obj.ptr)
	obj.ptr = nil
	obj.size = 0
}

func (obj *AVObject) SetSize(sz int) {
	obj.size = sz
}
func (obj *AVObject) Size() int {
	return obj.size
}

func (obj *AVObject) IsNil() bool {
	return obj.ptr == nil
}

func (obj *AVObject) WriteUInt64(data uint64) {
	C.memcpy(obj.ptr, unsafe.Pointer(&data), C.size_t(unsafe.Sizeof(data)))
}

func (obj *AVObject) UInt64() uint64 {
	return uint64(*(*C.uint64_t)(obj.ptr))
}

func (obj *AVObject) SetOptInt(name string, value int64, searchFlags int) int {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	return int(C.av_opt_set_int(obj.ptr, cname, C.int64_t(value), C.int(searchFlags)))
}

func (obj *AVObject) Bytes() []byte {
	if obj.size <= 0 {
		return nil
	}

	return (*[1 << 30]byte)(obj.ptr)[:obj.size:obj.size]
}

func (obj *AVObject) Copy() []byte {
	data := obj.Bytes()
	res := make([]byte, len(data))
	copy(res, data)
	return res
}

//can only write once
func (obj *AVObject) Write(p []byte) (int, error) {
	if len(p) == 0 {
		return 0, nil
	}
	if obj.length+len(p) > obj.size {
		panic("write out of range")
	}

	C.copybytes(obj.ptr, C.int(obj.length), unsafe.Pointer(&p[0]), C.int(len(p)))
	obj.length += len(p)

	return len(p), nil
}
