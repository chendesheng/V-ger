#include "_cgo_export.h"

int readFunc(void* ptr, uint8_t* buf, int bufSize) {
	return goReadFunc(ptr, buf, bufSize);
}

int64_t seekFunc(void* ptr, int64_t pos, int whence) {
	return goSeekFunc(ptr, pos, whence);
}


AVIOContext* new_io_context(unsigned char* buf, int bufSize, void* userData) {
	return avio_alloc_context(buf, bufSize, 0, userData, readFunc, 0, seekFunc);
}