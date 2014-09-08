package libav

/*
#cgo darwin LDFLAGS: -L/usr/local/lib/ -lavcodec -lavformat -lavutil -lswscale -lavresample -lbz2 -framework Foundation -lz -framework CoreVideo -framework VideoDecodeAcceleration
#include "libavformat/avformat.h"
*/
import "C"

func init() {
	C.av_register_all()
}
