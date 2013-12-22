#ifndef GUI_H
#define GUI_H
#include <stdint.h>

typedef struct SubItem {
	char* str;
	int style;
	unsigned int color;
} SubItem;

void test();
void* newWindow(char* title, int width, int height);
void initialize();
void pollEvents();
void showWindow(void* wptr);
void refreshWindowContent(void*wptr);
void makeWindowCurrentContext(void*wptr);
int getWindowWidth(void*);
int getWindowHeight(void*);
void showWindowProgress(void*,char*,char*,double);
void* showText(void*,SubItem*,int,int,double,double);
void hideText(void* ptrWin, void* ptrText);
void windowHideStartupView(void*);
void windowToggleFullScreen(void* wptr);

void initAudioMenu(void* wptr, char** name, int32_t*, int, int);
// void setSubtitles(void* wptr, char** name, int, int, int);
// typedef void (* TimerTickFunc)(void*);
// typedef void (* DrawFunc)(void*);

void onAudioMenuClicked(void*, int);
// void onSubtitleChanged(char*, char*);

void onDraw(void*);
void onTimerTick(void*);
void onKeyDown(void*, int);
void onProgressChanged(void* wptr, int typ, double position);
void onFullscreenChanged(void* wptr, int b);



// void setText(void* wptr, SubItem* items, int len);
// void setCallbackKeyDown();
// void setCallbackMouseDown();
// void setCallbackFullSreenChanged();
#endif