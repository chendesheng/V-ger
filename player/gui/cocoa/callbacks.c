#include "_cgo_export.h"
#include "gui.h"

void onDraw(void* wptr) {
	goOnDraw(wptr);
}
void onTimerTick(void* wptr) {
	goOnTimerTick(wptr);
}

int onKeyDown(int key) {
	return goOnKeyDown(key);
}

void onProgressChanged(int typ, double position) {
	goOnProgressChanged(typ, position);
}

void onFullscreenChanged(int b) {
	goOnFullscreenChanged(b);
}

int onOpenFile(const char* file) {
	return goOnOpenFile((void*)file);
}

void onWillTerminate() {
	goOnWillTerminate();
}


void onOpenOpenPanel() {
	goOnOpenOpenPanel();
}
void onCloseOpenPanel(char* filename) {
	goOnCloseOpenPanel(filename);
}

void onMouseWheel(double deltaY) {
	goOnMouseWheel(deltaY);
}

void onMouseMove() {
	goOnMouseMove();
}
void onMenuClicked(int type, int tag) {
	goOnMenuClicked(type, tag);
}