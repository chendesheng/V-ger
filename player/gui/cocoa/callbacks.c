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

	inline int onMenuClick(int type, int tag) {
		return goOnMenuClick(type, tag);
	}
	inline int isPlaying() {
		return goIsPlaying();
	}
	inline void getAllSubtitleNames(void*** names, int* length) {
		goGetSubtitles(names, length);   
	}
	inline void getPlayingSubtitles(int* firstSub, int* secondSub) {
		goGetPlayingSubtitles(firstSub, secondSub);
	}
	inline void getAllAudioTracks(void*** names, int* length) {
		goGetAllAudioTracks(names, length);
	}
	inline int getPlayingAudioTrack() {
		return goGetPlayingAudioTrack();
	}
	inline int isSearchingSubtitle() {
		return goIsSearchingSubtitle();
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

	inline int onMenuClick(int type, int tag) {
        printf("onMenuClick: %d %d\n", type, tag);
        return 0;
	}
    inline int isPlaying() {
        return 0;
    }
inline void getAllSubtitleNames(void*** names, int* length) {}
inline void getPlayingSubtitles(int* firstSub, int* secondSub) {}
inline int getPlayingAudioTrack(){return -1;}
inline void getAllAudioTracks(void*** names, int* length) {}
inline int isSearchingSubtitle() { return 0; }
#endif