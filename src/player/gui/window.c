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
void onSubtitleMenuClicked(void* wptr, int tag, int showOrHide) {
	goOnSubtitleMenuClicked(wptr, tag, showOrHide);
}

int onOpenFile(const char* file) {
	return goOnOpenFile((void*)file);
}

void onWillTerminate() {
	goOnWillTerminate();
}

void onSearchSubtitleMenuItemClick() {
	goOnSearchSubtitleMenuItemClick();
}

void onOpenOpenPanel() {
	goOnOpenOpenPanel();
}
void onCloseOpenPanel(char* filename) {
	goOnCloseOpenPanel(filename);
}
