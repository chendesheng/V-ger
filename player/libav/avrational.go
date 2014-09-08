package libav

//#include "libavformat/avformat.h"
//#include <stdlib.h>
import "C"

type AVRational C.AVRational

// func (r *AVRational) IsNil() bool {
// 	return r.ptr == nil
// }

func (r AVRational) Q2D() float64 {
	return float64(C.av_q2d(C.AVRational(r)))
}
