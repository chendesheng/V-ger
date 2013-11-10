#ifndef GUI_H
#define GUI_H

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
void showText(void*,SubItem*,int,double,double);


// typedef void (* TimerTickFunc)(void*);
// typedef void (* DrawFunc)(void*);

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