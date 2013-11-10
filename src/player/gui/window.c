#include "_cgo_export.h"
#include "gui.h"

void onDraw(void* wptr) {
	goOnDraw(wptr);
}
void onTimerTick(void* wptr) {
	goOnTimerTick(wptr);
}

void onKeyDown(void* wptr, int key) {
	goOnKeyDown(wptr, key);
}

void onProgressChanged(void* wptr, int typ, double position) {
	goOnProgressChanged(wptr, typ, position);
}

void onFullscreenChanged(void* wptr, int b) {
	goOnFullscreenChanged(wptr, b);
}