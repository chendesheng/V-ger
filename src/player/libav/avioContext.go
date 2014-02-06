package libav

/*
#include "libavformat/avio.h"
int readFunc(void* ptr, uint8_t* buf, int bufSize);
int64_t seekFunc(void* ptr, int64_t pos, int whence);
AVIOContext* new_io_context(unsigned char* buf, int bufSize, void* userData);
*/
import "C"
import (
	"unsafe"
)

type AVIOContext struct {
	ptr *C.AVIOContext
}

var cbReadFunc func(AVObject) int
var cbSeekFunc func(pos int64, whence int) int64

func NewAVIOContext(buffer AVObject, readfn func(AVObject) int, seekfn func(pos int64, whence int) int64) AVIOContext {
	ptr := C.new_io_context(
		(*_Ctype_unsignedchar)(buffer.ptr),
		C.int(buffer.size),
		nil)

	cbReadFunc = readfn
	cbSeekFunc = seekfn

	return AVIOContext{ptr: ptr}
}

//export goReadFunc
func goReadFunc(ptr unsafe.Pointer, buf unsafe.Pointer, bufSize C.int) C.int {
	return C.int(cbReadFunc(AVObject{ptr: buf, size: int(bufSize)}))
}

//export goSeekFunc
func goSeekFunc(ptr unsafe.Pointer, pos C.int64_t, whence C.int) C.int64_t {
	return C.int64_t(cbSeekFunc(int64(pos), int(whence)))
}
