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
int getWindowWidth(void*);
int getWindowHeight(void*);
void setWindowSize(void*, int, int);
void setWindowTitle(void*, char* title);

void refreshWindowContent(void*);
void makeWindowCurrentContext(void*);
void windowHideStartupView(void*);
void windowShowStartupView(void*);
void windowToggleFullScreen(void*);
void closeWindow(void*);

void showWindowProgress(void*,char*,char*,double);
void showWindowBufferInfo(void*, char*, double);

void* showText(void*,SubItem*);
void hideText(void*, void*);

void initAudioMenu(void*, char**, int32_t*, int, int);
void hideAudioMenu();
void initSubtitleMenu(void*, char**, int32_t*, int, int32_t, int32_t);
void hideSubtitleMenu();
void setSubtitleMenuItem(int,int);
void hideCursor(void*);
void showCursor(void*);
void showSpinning(void*);
void hideSpinning(void*);
void setVolume(void*, int);
void setVolumeDisplay(void*, int);

//callbacks
void onMouseMove();
void onMenuClicked(int, int);
void onOpenOpenPanel();
void onCloseOpenPanel(char*);
void onDraw();
void onTimerTick();
int onKeyDown(int);
void onMouseWheel(double);
void onProgressChanged(int, double);
void onFullscreenChanged(int);
int onOpenFile(const char*);
void onWillTerminate();


#define MENU_AUDIO 0
#define	MENU_SUBTITLE 1
#define MENU_SEARCH_SUBTITLE 2

#endif