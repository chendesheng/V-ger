package libav

/*
#include "libavcodec/avcodec.h"
#include "libavformat/avformat.h"
#include "libavutil/avutil.h"
#include "libavutil/mathematics.h"
#include <stdlib.h>
void decode_size(AVPacket* pkt, int len) {
	pkt->size -= len;
	pkt->data += len;
}
*/
import "C"
import "unsafe"

type AVPacket C.AVPacket

func (packet *AVPacket) Size() int {
	return int(packet.size)
}

func (packet *AVPacket) Free() {
	C.av_free_packet(packet.pointer())
}

func (packet *AVPacket) Dts() uint64 {
	return uint64(packet.dts)
}

func (packet *AVPacket) Pts() uint64 {
	return uint64(packet.pts)
}

func (packet *AVPacket) StreamIndex() int {
	return int(packet.stream_index)
}

func (packet *AVPacket) pointer() *C.AVPacket {
	return (*C.AVPacket)(unsafe.Pointer(packet))
}

func (packet *AVPacket) DecodeSize(len int) {
	C.decode_size((*C.struct_AVPacket)(packet), C.int(len))
}
