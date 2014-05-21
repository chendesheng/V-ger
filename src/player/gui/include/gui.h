#ifndef GUI_H
#define GUI_H
#include <stdint.h>

typedef struct SubItem {
	char* str;
	int style;
	unsigned int color;
} SubItem;
typedef struct CSize {
	int width;
	int height;
} CSize;

void test();
void* newWindow(char* title, int width, int height);
void* newDialog(char* title, int width, int height);

void initialize();
void pollEvents();
void showWindow(void* wptr);
void refreshWindowContent(void*wptr);
void makeWindowCurrentContext(void*wptr);
int getWindowWidth(void*);
int getWindowHeight(void*);
void showWindowProgress(void*,char*,char*,double, double, char*);
void* showText(void*,SubItem*,int,int,double,double);
void hideText(void* ptrWin, void* ptrText);
void windowHideStartupView(void*);
void windowShowStartupView(void*);
void windowToggleFullScreen(void* wptr);
void setWindowSize(void* wptr, int width, int height);
void setWindowTitle(void* wptr, char* title);

void initialize();

void initAudioMenu(void* wptr, char** name, int32_t*, int, int);
void hideAudioMenu();

void initSubtitleMenu(void* wptr, char** name, int32_t*, int, int32_t, int32_t);
void hideSubtitleMenu();

// void setSubtitles(void* wptr, char** name, int, int, int);
// typedef void (* TimerTickFunc)(void*);
// typedef void (* DrawFunc)(void*);

void onAudioMenuClicked(void*, int);
void onSubtitleMenuClicked(void*, int, int);


void onDraw(void*);
void onTimerTick(void*);
void onKeyDown(void*, int);
void onMouseWheel(void*, double);
void onProgressChanged(void* wptr, int typ, double position);
void onFullscreenChanged(void* wptr, int b);

int onOpenFile(const char* file);
void onWillTerminate();
void onSearchSubtitleMenuItemClick();
// void setText(void* wptr, SubItem* items, int len);
// void setCallbackKeyDown();
// void setCallbackMouseDown();
// void setCallbackFullSreenChanged();
CSize getScreenSize();

void onOpenOpenPanel();
void onCloseOpenPanel(char* filename);
void hideCursor(void*);
void showCursor(void*);
#endif