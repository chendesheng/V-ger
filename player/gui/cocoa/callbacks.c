#include "_cgo_export.h"
#include "gui.h"

inline void onDraw(void* wptr) {
	goOnDraw(wptr);
}
inline void onTimerTick(void* wptr) {
	goOnTimerTick(wptr);
}

inline int onKeyDown(int key) {
	return goOnKeyDown(key);
}

inline void onProgressChanged(int typ, double position) {
	goOnProgressChanged(typ, position);
}

inline void onFullscreenChanged(int b) {
	goOnFullscreenChanged(b);
}

inline int onOpenFile(const char* file) {
	return goOnOpenFile((void*)file);
}

inline void onWillTerminate() {
	goOnWillTerminate();
}


inline void onOpenOpenPanel() {
	goOnOpenOpenPanel();
}
inline void onCloseOpenPanel(char* filename) {
	goOnCloseOpenPanel(filename);
}

inline void onMouseWheel(double deltaY) {
	goOnMouseWheel(deltaY);
}

inline void onMouseMove() {
	goOnMouseMove();
}
inline void onMenuClicked(int type, int tag) {
	goOnMenuClicked(type, tag);
}