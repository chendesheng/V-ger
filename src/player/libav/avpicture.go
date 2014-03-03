package libav

//#include "libavcodec/avcodec.h"
//#include "libavutil/mem.h"
//#include <stdlib.h>
import "C"
import (
	"fmt"
	"os"
	"unsafe"
)

type AVPicture struct {
	ptr *C.AVPicture
	buf AVObject
}

func (picture *AVPicture) IsNil() bool {
	return picture.ptr == nil
}

func (picture *AVPicture) Fill(buffer AVObject, pixFmt int, width, height int) AVObject {
	picture.buf = buffer
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

func (picture *AVPicture) Data() AVObject {
	return AVObject{ptr: unsafe.Pointer(&picture.ptr.data[0])}
}

func (picture *AVPicture) SaveToPPMFile(file string, width, height int) []byte {
	f, _ := os.OpenFile(file, os.O_CREATE|os.O_RDWR, 0666)
	defer f.Close()

	f.WriteString(fmt.Sprintf("P6\n%d %d\n255\n", width, height))

	// obj := picture.Data()
	// obj.size = buf

	bytes := picture.buf.Bytes()
	linesize := picture.ptr.linesize[0]

	outBytes := make([]byte, 0, len(bytes))

	for y := 0; y < height; y++ {
		start := int(C.int(y) * linesize)

		outBytes = append(outBytes, bytes[start:start+width*3]...)
	}

	f.Write(outBytes)
	return outBytes
}

func (picture *AVPicture) RGBBytes(width, height int) []byte {
	bytes := picture.buf.Bytes()
	linesize := picture.ptr.linesize[0]

	outBytes := make([]byte, 0, len(bytes))

	for y := 0; y < height; y++ {
		start := int(C.int(y) * linesize)
		outBytes = append(outBytes, bytes[start:start+width*3]...)
	}

	return outBytes
}

func (picture *AVPicture) Frame() AVFrame {
	return AVFrame{ptr: (*C.AVFrame)(unsafe.Pointer(picture.ptr))}
}
