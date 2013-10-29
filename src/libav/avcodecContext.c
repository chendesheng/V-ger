#include "_cgo_export.h"

int getBuffer(AVCodecContext *c, AVFrame *frame) {
	return goGetBufferCallback(c, frame);
}

void releaseBuffer(AVCodecContext *c, AVFrame *frame) {
	goReleaseBufferCallback(c, frame);
}


void SetGetBufferCallbackCB(AVCodecContext *c) {
	c->get_buffer = getBuffer;
}
void SetReleaseBufferCallbackCB(AVCodecContext *c) {
	c->release_buffer = releaseBuffer;
}

