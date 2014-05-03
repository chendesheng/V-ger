package sdl

//#include "stdlib.h"
//#include "string.h"
//#include "sdl2/sdl.h"
import "C"
import (
	"unsafe"
)

type Object struct {
	ptr     unsafe.Pointer
	current uintptr
}

func (obj *Object) IsNil() bool {
	return obj.ptr == nil
}

func (obj *Object) SetZero(length int) {
	// obj.current = uintptr(obj.ptr)
	C.SDL_memset(unsafe.Pointer(obj.current), 0, C.size_t(length))
}

func (obj *Object) Write(p []byte) {
	if len(p) == 0 {
		return
	}

	if obj.current == 0 {
		obj.current = uintptr(obj.ptr)
	}

	C.memcpy(unsafe.Pointer(obj.current), unsafe.Pointer(&p[0]), C.size_t(len(p)))

	obj.current += uintptr(len(p))
}
