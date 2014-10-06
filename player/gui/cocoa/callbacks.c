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

	inline void onPlaybackChange(int typ, double position) {
		goOnPlaybackChange(typ, position);
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
	inline int onMenuClick(int type, int tag) {
		return goOnMenuClick(type, tag);
	}
	inline int isPlaying() {
		return goIsPlaying();
	}
	inline void getSubtitles(void*** names, int* length, int* firstSub, int* secondSub) {
		goGetSubtitles(names, length, firstSub, secondSub);   
	}
	void getAudioes(void*** names, int* length, int* selected) {
		goGetAudioes(names, length, selected);
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

	inline void onPlaybackChange(int typ, double position) {
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
        //showCursor(w);
	}
	inline int onMenuClick(int type, int tag) {
        printf("onMenuClick: %d %d\n", type, tag);
        return 0;
	}
    inline int isPlaying() {
        return 0;
    }
inline void getSubtitles(void*** names, int* length, int* firstSub, int* secondSub) {
    
}
void getAudioes(void** names, int* length, int* selected) {

}
#endif