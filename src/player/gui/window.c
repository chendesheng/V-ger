#include "_cgo_export.h"
#include "gui.h"

void onDraw(void* wptr) {
	goOnDraw(wptr);
}
void onTimerTick(void* wptr) {
	goOnTimerTick(wptr);
}

int onKeyDown(void* wptr, int key) {
	return goOnKeyDown(wptr, key);
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
void onSubtitleMenuClicked(void* wptr, int tag) {
	goOnSubtitleMenuClicked(wptr, tag);
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

void onMouseWheel(void* wptr, double deltaY) {
	goOnMouseWheel(wptr, deltaY);
}

void onMouseMove(void* wptr) {
	goOnMouseMove(wptr);
}