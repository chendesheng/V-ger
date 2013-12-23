#include "gui.h"

#import <Cocoa/Cocoa.h>
#import "window.h"
#import "windowDelegate.h"
#import "glView.h"
#import "textView.h"
#import "blurView.h"
#import "progressView.h"
#import "startupView.h"

void initialize() {
    if (NSApp)
        return;

	[NSApplication sharedApplication];

    [NSApp setActivationPolicy:NSApplicationActivationPolicyRegular];

	NSLog(@"initialized");
	//create memu bar
	id menubar = [[NSMenu new] autorelease];
    id appMenuItem = [[NSMenuItem new] autorelease];
    [menubar addItem:appMenuItem];
    [NSApp setMainMenu:menubar];
    id appMenu = [[NSMenu new] autorelease];
    id appName = [[NSProcessInfo processInfo] processName];
    id quitTitle = [@"Quit " stringByAppendingString:appName];
    id quitMenuItem = [[[NSMenuItem alloc] initWithTitle:quitTitle
        action:@selector(terminate:) keyEquivalent:@"q"] autorelease];
    [appMenu addItem:quitMenuItem];
    [appMenuItem setSubmenu:appMenu];
}

void initAudioMenu(void* wptr, char** names, int32_t* tags, int len, int selected) {
    NSWindow* w = (NSWindow*)wptr;

    NSMenu *menubar = [NSApp mainMenu];
    NSMenuItem* audioMenuItem = [[NSMenuItem new] autorelease];
    [menubar addItem:audioMenuItem];
    NSMenu* audioMenu = [[NSMenu alloc] initWithTitle:@"Audio"];

    for (int i = 0; i < len; i++) {
        char* name = names[i];
        int tag = tags[i];
        NSMenuItem* item = [[NSMenuItem alloc] initWithTitle:[NSString stringWithUTF8String:name] 
            action:@selector(audioMenuItemClick:) keyEquivalent:@""];
        [item setTarget: w];
        [item setTag: tag];
        [audioMenu addItem:item];

        if (tag == selected) {
            [item setState: NSOnState];
        }
    }

    [audioMenuItem setSubmenu:audioMenu];
}

void initSubtitleMenu(void* wptr, char** names, int32_t* tags, int len, int selected) {
    NSWindow* w = (NSWindow*)wptr;

    NSMenu *menubar = [NSApp mainMenu];
    NSMenuItem* subtitleMenuItem = [[NSMenuItem new] autorelease];
    [menubar addItem:subtitleMenuItem];
    NSMenu* subtitleMenu = [[NSMenu alloc] initWithTitle:@"Subtitle"];

    for (int i = 0; i < len; i++) {
        char* name = names[i];
        int tag = tags[i];
        NSMenuItem* item = [[NSMenuItem alloc] initWithTitle:[NSString stringWithUTF8String:name] 
            action:@selector(subtitleMenuItemClick:) keyEquivalent:@""];
        [item setTarget: w];
        [item setTag: tag];
        [subtitleMenu addItem:item];

        if (tag == selected) {
            [item setState: NSOnState];
        }
    }

    [subtitleMenuItem setSubmenu:subtitleMenu];
}

void* newWindow(char* title, int width, int height) {
	NSAutoreleasePool * pool = [[NSAutoreleasePool alloc] init];

	initialize();

	Window* w = [[Window alloc] initWithTitle:[NSString stringWithUTF8String:title]
		width:width height:height];
	
    WindowDelegate* wd = [[WindowDelegate alloc] init];
    wd->window = w;
	[w setDelegate:wd];

	GLView* v = [[GLView alloc] initWithFrame2: [w frame]];
	[w setContentView:v];


    TextView* tv = [[TextView alloc] initWithFrame:NSMakeRect(0, 30, width, 0)];
    [v addSubview:tv];
    [tv setAutoresizingMask:NSViewWidthSizable];
    [v setTextView:tv];

    TextView* tv2 = [[TextView alloc] initWithFrame:NSMakeRect(0, 30, width, 0)];
    [v addSubview:tv2];
    [tv2 setAutoresizingMask:NSViewWidthSizable];
    [v setTextView2:tv2];

    BlurView* bv = [[BlurView alloc] initWithFrame:NSMakeRect(0,0,width,30)];
    [v addSubview:bv];
    [bv setAutoresizingMask:NSViewWidthSizable];

    ProgressView* pv = [[ProgressView alloc] initWithFrame:[bv frame]];
    [bv addSubview:pv];
    [pv setAutoresizingMask:NSViewWidthSizable|NSViewHeightSizable];

    [v setProgressView:pv];

    StartupView* sv = [[StartupView alloc] initWithFrame:[v frame]];
    [v addSubview:sv];

    [v setStartupView:sv];


    NSTimer *renderTimer = [NSTimer timerWithTimeInterval:1.0/60.0 
                            target:w
                          selector:@selector(timerTick:)
                          userInfo:nil
                           repeats:YES];

    [[NSRunLoop currentRunLoop] addTimer:renderTimer
                                forMode:NSDefaultRunLoopMode];
    [[NSRunLoop currentRunLoop] addTimer:renderTimer
                                forMode:NSEventTrackingRunLoopMode]; //Ensure timer fires during resize

	[pool drain];

	return w;
}

void showWindow(void* ptr) {
	[NSApp activateIgnoringOtherApps:YES];

	Window* w = (Window*)ptr;
	[w makeKeyAndOrderFront:nil];
}
void makeWindowCurrentContext(void*ptr) {
    Window* w = (Window*)ptr;
    [w makeCurrentContext];
}
void pollEvents() {
    NSAutoreleasePool* pool = [[NSAutoreleasePool alloc] init];
    [NSApp finishLaunching];

    while(YES) {
	    [pool drain];
		pool = [[NSAutoreleasePool alloc] init];

	    NSEvent* event = [NSApp nextEventMatchingMask:NSAnyEventMask
	                                        untilDate:[NSDate distantFuture]
	                                        inMode:NSDefaultRunLoopMode
	                                        dequeue:YES];
	    [NSApp sendEvent:event];
	}
	// [NSApp activateIgnoringOtherApps:YES];
    //[NSApp run];

    [pool drain];
}
void refreshWindowContent(void*wptr) {
	Window* w = (Window*)wptr;
	[w setContentViewNeedsDisplay:YES];
}

int getWindowWidth(void* ptr) {
    Window* w = (Window*)ptr;
    return (int)([[w contentView] frame].size.width);
}
int getWindowHeight(void* ptr) {
    Window* w = (Window*)ptr;
    return (int)([[w contentView] frame].size.height);
}
void showWindowProgress(void* ptr, char* left, char* right, double percent) {
    Window* w = (Window*)ptr;
    [[w contentView] showProgress:left right:right percent:percent];
}
void* showText(void* ptr, SubItem* items, int length, int position, double x, double y) {
    Window* w = (Window*)ptr;
    return [[w contentView] showText:items length:length position:position x:x y:y];
}
void hideText(void* ptrWin, void* ptrText) {
    Window* w = (Window*)ptrWin;
    [[w contentView] hideText:ptrText];
}
void windowHideStartupView(void* ptr) {
    Window* w = (Window*)ptr;
    [[w contentView] hideStartupView];
}
void windowToggleFullScreen(void* ptr) {
    Window* w = (Window*)ptr;
    [w toggleFullScreen:nil];
}