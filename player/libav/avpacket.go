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
AVPacket* av_malloc_avpacket() {
	AVPacket* p = av_malloc(sizeof(AVPacket));
	memset(p, 0, sizeof(AVPacket));
	return p;
}
*/
import "C"
import "unsafe"

type AVPacket struct {
	ptr *C.AVPacket
}

func NewAVPacket() AVPacket {
	return AVPacket{C.av_malloc_avpacket()}
}

func (packet AVPacket) IsNil() bool {
	return packet.ptr == nil
}

func (packet AVPacket) Size() int {
	return int(packet.ptr.size)
}

func (packet AVPacket) Free() {
	C.av_free(unsafe.Pointer(packet.ptr))
}

func (packet AVPacket) FreePacket() {
	C.av_free_packet(packet.ptr)
}

func (packet AVPacket) Dts() uint64 {
	return uint64(packet.ptr.dts)
}

func (packet AVPacket) Pts() uint64 {
	return uint64(packet.ptr.pts)
}

func (packet AVPacket) StreamIndex() int {
	return int(packet.ptr.stream_index)
}

func (packet AVPacket) pointer() *C.AVPacket {
	return packet.ptr
}

func (packet AVPacket) DecodeSize(len int) {
	C.decode_size((*C.struct_AVPacket)(packet.ptr), C.int(len))
}
