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


void onAudioMenuClicked(void* wptr, int tag) {
	goOnAudioMenuClicked(wptr, tag);
}
// void onSubtitleChanged(char* name1, char* name2) {
// 	goOnSubtitleChanged(name1, name2);
// }