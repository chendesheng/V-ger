#ifndef GUI_H
#define GUI_H
#include <stdint.h>

typedef struct AttributedString {
	char* str;
	int style;
	unsigned int color;
} AttributedString;

typedef struct SubItem {
	AttributedString* texts;
	int length;
	int align;
	double x;
	double y;
} SubItem;

typedef struct CSize {
	int width;
	int height;
} CSize;
CSize getScreenSize();
void alert(void*, char*);

void initialize();
void pollEvents();
void* newWindow(char*, int, int);
void showWindow(void*);
CSize getWindowSize(void*);
void setWindowSize(void*, int, int);
void setWindowTitle(void*, char* title);
void setWindowTitleWithRepresentedFilename(void* wptr, char* title);

void refreshWindowContent(void*);
void initWindowCurrentContext(void*);
void toggleFullScreen(void*);
int isFullScreen(void*);
void closeWindow(void*);

void updatePlaybackInfo(void*,char*,char*,double);
void updateBufferInfo(void*, char*, double);

void* showSubtitle(void*, SubItem*);
void hideSubtitle(void*, long);

//including playback panel, window title, cursor
void setControlsVisible(void*, int, int);
void setSpinningVisible(void*, int);

void setVolume(void*, int);
void setVolumeVisible(void*, int);

void addRecentOpenedFile(char*);

void allowDisplaySleep();
void preventDisplaySleep();

//callbacks
int onMenuClick(int, int);
void onOpenOpenPanel();
void onCloseOpenPanel(char*);
void onDraw();
void onTimerTick();
int onKeyDown(int);
void onMouseWheel(double);
void onPlaybackChange(int, double);
int onOpenFile(const char*);
void onWillTerminate();
void onFullScreen(int);
void onWillSleep();
void onDidWake();

int isPlaying();
void getAllSubtitleNames(void***, int*);
void getPlayingSubtitles(int*, int*);

void getAllAudioTracks(void***, int*);
int getPlayingAudioTrack();
int isSearchingSubtitle();

#define MENU_AUDIO 0
#define	MENU_SUBTITLE 1
#define MENU_SEARCH_SUBTITLE 2
#define MENU_PLAY 3
#define MENU_SEEK 4
#define MENU_VOLUME 5
#define MENU_SYNC_SUBTITLE 6
#define MENU_SYNC_AUDIO 7

#define WILL_ENTER_FULL_SCREEN 0
#define DID_ENTER_FULL_SCREEN 1
#define WILL_EXIT_FULL_SCREEN 2
#define DID_EXIT_FULL_SCREEN 3

void flushBuffer(void* ptr);
void makeCurrentContext(void* ptr);

void setSubFontSize(void*, double);
void onImportSubtitle(char*);
#endif
