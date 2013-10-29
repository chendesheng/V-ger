package sdl

//#include "stdlib.h"
//#include "string.h"
import "C"
import (
	"unsafe"
)

type Object struct {
	ptr unsafe.Pointer

	pos  uintptr
	size int
}

// func (obj *Object) Free() {
// 	if obj.IsNil() {
// 		return
// 	}
// }

func (obj *Object) SetSize(sz int) {
	obj.size = sz
}

// func (obj *Object) Size() int {
// 	return obj.size
// }

func (obj *Object) IsNil() bool {
	return obj.ptr == nil
}

func (obj *Object) Write(bytes []byte) {
	if obj.pos == 0 {
		obj.pos = uintptr(obj.ptr)
	}
	C.memcpy(unsafe.Pointer(obj.pos), unsafe.Pointer(&bytes[0]), C.size_t(len(bytes)))
}

func (obj *Object) Offset(offset int) {
	if obj.pos == 0 {
		obj.pos = uintptr(obj.ptr)
	}
	obj.pos += uintptr(offset)
}

// func (obj *Object) WriteUInt64(data uint64) {
// 	C.memcpy(obj.ptr, unsafe.Pointer(&data), C.size_t(unsafe.Sizeof(data)))
// }

// func (obj *Object) UInt64() uint64 {
// 	return uint64(*(*C.uint64_t)(obj.ptr))
// }

// func (obj *Object) SetOptInt(name string, value int64, searchFlags int) int {
// 	cname := C.CString(name)
// 	defer C.free(unsafe.Pointer(cname))
// 	return int(C.av_opt_set_int(obj.ptr, cname, C.int64_t(value), C.int(searchFlags)))
// }

// func (obj *Object) Bytes() []byte {
// 	var bytes []byte
// 	header := (*reflect.SliceHeader)(unsafe.Pointer(&bytes))
// 	if obj.size > 0 {
// 		header.Len = obj.size
// 	} else {
// 		header.Len = (1 << 31)
// 	}
// 	header.Cap = header.Len
// 	header.Data = uintptr(obj.ptr)

// 	return bytes
// }
