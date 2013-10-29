package libav

/*
#cgo darwin LDFLAGS: -lavresample -lswscale -lbz2 -framework Foundation -lz -framework CoreVideo -framework VideoDecodeAcceleration
#include "libavformat/avformat.h"
*/
import "C"

func init() {
	C.av_register_all()
}
