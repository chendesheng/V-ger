package sdl

//#include "stdlib.h"
//#include "string.h"
//#include "sdl2/sdl.h"
import "C"
import (
	"unsafe"
)

type Object struct {
	ptr unsafe.Pointer
}

func (obj *Object) IsNil() bool {
	return obj.ptr == nil
}

func (obj *Object) SetZero(length int) {
	C.SDL_memset(obj.ptr, 0, C.size_t(length))
}
