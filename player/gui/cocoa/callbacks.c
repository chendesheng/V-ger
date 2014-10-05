#include "gui.h"

#ifndef COCOA_TEST
	#include "_cgo_export.h"
	inline void onDraw() {
		goOnDraw();
	}
	inline void onTimerTick() {
		goOnTimerTick();
	}

	inline int onKeyDown(int key) {
		return goOnKeyDown(key);
	}

	inline void onProgressChange(int typ, double position) {
		goOnProgressChange(typ, position);
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
	inline void onMenuClick(int type, int tag) {
		goOnMenuClick(type, tag);
	}
#else
#include <stdio.h>
    inline void onDraw() {
	//	goOnDraw(wptr);
	}
	inline void onTimerTick() {
//        printf("onTimerTick");
	}

	inline int onKeyDown(int key) {
//        printf("onKeyDown:%x\n", key);
//        windowHideStartupView(w);
//        setVolumeDisplay(w, 1);
	    return 0;
	}

	inline void onProgressChange(int typ, double position) {
	}

	inline void onFullscreenChanged(int b) {
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

	inline void onMouseMove(void* w) {
        showCursor(w);
	}
	inline void onMenuClick(int type, int tag) {
        printf("onMenuClick: %d %d\n", type, tag);
	}
#endif