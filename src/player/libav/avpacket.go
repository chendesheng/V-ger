package libav

/*
#include "libavcodec/avcodec.h"
#include "libavformat/avformat.h"
#include "libavutil/avutil.h"
#include "libavutil/mathematics.h"
#include <stdlib.h>
*/
import "C"

type AVPacket struct {
	cAVPacket C.AVPacket
}

func (packet *AVPacket) Size() int {
	return int(packet.cAVPacket.size)
}

func (packet *AVPacket) Free() {
	C.av_free_packet(&packet.cAVPacket)
}

func (packet *AVPacket) Dts() uint64 {
	return uint64(packet.cAVPacket.dts)
}

func (packet *AVPacket) Pts() uint64 {
	return uint64(packet.cAVPacket.pts)
}

func (packet *AVPacket) StreamIndex() int {
	return int(packet.cAVPacket.stream_index)
}

func (packet *AVPacket) Dup() {
	C.av_dup_packet(&packet.cAVPacket)
}
