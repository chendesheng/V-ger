package libav

//#include "libavutil/mem.h"
//#include "libavutil/opt.h"
//#include <string.h>
//#include <stdlib.h>
import "C"
import (
	"reflect"
	"unsafe"
)

type AVObject struct {
	ptr  unsafe.Pointer
	ptr2 uintptr
	size int
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
	var bytes []byte
	header := (*reflect.SliceHeader)(unsafe.Pointer(&bytes))
	if obj.size > 0 {
		header.Len = obj.size
	} else {
		header.Len = (1 << 31)
	}
	header.Cap = header.Len
	header.Data = uintptr(obj.ptr)

	return bytes
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

	if obj.ptr2 == 0 {
		obj.ptr2 = uintptr(obj.ptr)
	}

	C.memcpy(unsafe.Pointer(obj.ptr2), unsafe.Pointer(&p[0]), C.size_t(len(p)))

	obj.ptr2 = obj.ptr2 + uintptr(len(p))

	return len(p), nil
}
