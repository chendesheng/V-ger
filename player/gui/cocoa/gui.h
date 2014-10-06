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

void refreshWindowContent(void*);
void initWindowCurrentContext(void*);
void setStartupViewVisible(void*, int);
void toggleFullScreen(void*);
void closeWindow(void*);

void updatePlaybackInfo(void*,char*,char*,double);
void updateBufferInfo(void*, char*, double);

void* showSubtitle(void*, SubItem*);
void hideSubtitle(void*, void*);

void initAudioMenu(void*, char**, int32_t*, int, int);
void hideAudioMenu();
void initSubtitleMenu(void*, char**, int32_t*, int, int32_t, int32_t);
void hideSubtitleMenu();
void selectSubtitleMenu(int,int);

//including playback panel, window title, cursor
void setControlsVisible(void*, int);
void setSpinningVisible(void*, int);

void setVolume(void*, int);
void setVolumeVisible(void*, int);

void addRecentOpenedFile(char*);

//callbacks
void onMouseMove();
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

int isPlaying();


#define MENU_AUDIO 0
#define	MENU_SUBTITLE 1
#define MENU_SEARCH_SUBTITLE 2
#define MENU_PLAY 3

#endif