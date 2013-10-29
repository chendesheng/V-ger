package libav

//#include "libavcodec/avcodec.h"
//#include "libavutil/mem.h"
//#include <stdlib.h>
import "C"
import (
// "unsafe"
)

type AVPicture struct {
	ptr *C.AVPicture
}

func (picture *AVPicture) IsNil() bool {
	return picture.ptr == nil
}

func (picture *AVPicture) Fill(buffer AVObject, pixFmt int, width, height int) AVObject {
	C.avpicture_fill(picture.ptr, (*C.uint8_t)(buffer.ptr), int32(pixFmt), C.int(width), C.int(height))
	return buffer
}

func (picture *AVPicture) Alloc(pixFmt int, width, height int) {
	picture.ptr = &C.AVPicture{}
	C.avpicture_alloc(picture.ptr, int32(pixFmt), C.int(width), C.int(height))
}

func AVPictureGetSize(pixFmt int, width, height int) int {
	return int(C.avpicture_get_size(int32(pixFmt), C.int(width), C.int(height)))
}

// func (picture *AVPicture) DataAt(i int) []byte {
// 	obj := AVObject{ptr: unsafe.Pointer(picture.ptr.data[i])}
// 	return obj.Bytes()
// }

func (picture *AVPicture) Layout(pixFmt int, width, height int) AVObject {
	size := AVPictureGetSize(pixFmt, width, height)
	dest := C.malloc(C.size_t(size))
	C.avpicture_layout(picture.ptr, int32(pixFmt), C.int(width), C.int(height), (*_Ctype_unsignedchar)(dest), C.int(size))

	return AVObject{ptr: dest, size: size}
}
