#include "gui.h"

#ifndef COCOA_TEST
	#include "_cgo_export.h"
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
#else
	inline void onDraw(void* wptr) {
	//	goOnDraw(wptr);
	}
	inline void onTimerTick(void* wptr) {
	//	goOnTimerTick(wptr);
	}

	inline int onKeyDown(int key) {
	//	return goOnKeyDown(key);
	    return 0;
	}

	inline void onProgressChanged(int typ, double position) {
	//	goOnProgressChanged(typ, position);
	}

	inline void onFullscreenChanged(int b) {
	//	goOnFullscreenChanged(b);
	}

	inline int onOpenFile(const char* file) {
	//	return goOnOpenFile((void*)file);
	    return 0;
	}

	inline void onWillTerminate() {
	//	goOnWillTerminate();
	}


	inline void onOpenOpenPanel() {
	//	goOnOpenOpenPanel();
	}
	inline void onCloseOpenPanel(char* filename) {
	//	goOnCloseOpenPanel(filename);
	}

	inline void onMouseWheel(double deltaY) {
	//	goOnMouseWheel(deltaY);
	}

	inline void onMouseMove() {
	//	goOnMouseMove();
	}
	inline void onMenuClicked(int type, int tag) {
	//	goOnMenuClicked(type, tag);
	}
#endif