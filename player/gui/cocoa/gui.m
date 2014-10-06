#include "gui.h"

#import <Cocoa/Cocoa.h>
#import "window.h"
#import "windowDelegate.h"
#import "glView.h"
#import "textView.h"
#import "blurView.h"
#import "progressView.h"
#import "startupView.h"
#import "volumeView.h"
// #import "titleTextView.h"
#import "app.h"

void initialize() {
    if (NSApp)
        return;

	[NSApplication sharedApplication];

    [NSApp setActivationPolicy:NSApplicationActivationPolicyRegular];

    Application *appDelegate = [[Application alloc] init];
    [NSApp setDelegate:appDelegate];

    [[NSBundle mainBundle] loadNibNamed:@"MainMenu" owner:NSApp topLevelObjects:nil];
    
    [[NSAppleEventManager sharedAppleEventManager] setEventHandler:appDelegate andSelector:@selector(handleAppleEvent:withReplyEvent:) forEventClass:kInternetEventClass andEventID:kAEGetURL];
}

void setWindowTitle(void* wptr, char* title) {
    Window* w = (Window*)wptr;

    NSString* str = [NSString stringWithUTF8String:title];
    [w setTitle:str];
}

void setWindowSize(void* wptr, int width, int height) {
    Window* w = (Window*)wptr;

    NSRect frame = [w frame];
    frame.origin.y -= (height - frame.size.height)/2;
    frame.origin.x -= (width - frame.size.width)/2;
    frame.size = NSMakeSize(width, height);

    w->customAspectRatio = NSMakeSize(width, height);
    [w->glView setOriginalSize:NSMakeSize(width, height)];

    [w setFrame:frame display:YES animate:YES];

    [w makeKeyAndOrderFront:nil];
}

void* newWindow(char* title, int width, int height) {
   	@autoreleasepool {

    	initialize();

    	Window* w = [[Window alloc] initWithWidth:width height:height];
        setWindowTitle(w, title);
        
        WindowDelegate* wd = (WindowDelegate*)[[WindowDelegate alloc] init];
    	[w setDelegate:(id)wd];

        [w makeFirstResponder:w->glView];

        NSTimer *renderTimer = [NSTimer timerWithTimeInterval:1.0/100.0 
                                target:[NSApp delegate]
                              selector:@selector(timerTick:)
                              userInfo:nil
                               repeats:YES];

        [[NSRunLoop currentRunLoop] addTimer:renderTimer
                                    forMode:NSDefaultRunLoopMode];
        [[NSRunLoop currentRunLoop] addTimer:renderTimer
                                    forMode:NSEventTrackingRunLoopMode]; //Ensure timer fires during resize

    	return w;
    }
}

void showWindow(void* ptr) {
	[NSApp activateIgnoringOtherApps:YES];

	Window* w = (Window*)ptr;
	[w makeKeyAndOrderFront:nil];
}
void initWindowCurrentContext(void*ptr) {
    Window* w = (Window*)ptr;
    [w makeCurrentContext];
}
void pollEvents() {
    [NSApp run];
    // NSApplicationMain(0, NULL);
}
void refreshWindowContent(void*wptr) {
	Window* w = (Window*)wptr;
    [w->glView setNeedsDisplay:YES];
}

CSize getWindowSize(void* ptr) {
    Window* w = (Window*)ptr;
    CSize sz;
    sz.width = (int)([w->glView frame].size.width);
    sz.height = (int)([w->glView frame].size.height);
    return sz;
}

void updatePlaybackInfo(void* ptr, char* left, char* right, double percent) {
    Window* w = (Window*)ptr;
    
    NSString* leftStr;
    if (strlen(left) == 0) {
        leftStr = @"00:00:00";
    } else {
        leftStr = [[NSString stringWithUTF8String:left] retain];
    }
    NSString* rightStr;
    if (strlen(right) == 0) {
        rightStr = @"00:00:00";
    } else {
        rightStr = [[NSString stringWithUTF8String:right] retain];
    }
    [w->glView updatePorgressInfo:leftStr rightString:rightStr percent:percent];
}
void updateBufferInfo(void* ptr, char* speed, double percent) {
    Window* w = (Window*)ptr;
    NSString* str;
    if (strlen(speed) == 0) {
        str = @"";
    } else {
        str = [[NSString stringWithUTF8String:speed] retain];
    }
    [w->glView updateBufferInfo:str bufferPercent:percent];
}
void* showSubtitle(void* ptr, SubItem* item) {
    Window* w = (Window*)ptr;
    return [w->glView showSubtitle:item];
}
void hideSubtitle(void* ptrWin, void* ptrText) {
    Window* w = (Window*)ptrWin;
    [w->glView hideSubtitle:ptrText];
}
void setStartupViewVisible(void* ptr, int b) {
    Window* w = (Window*)ptr;
    [w->glView setStartupViewHidden:(b==0)];
}
void setSpinningVisible(void* ptr, int b) {
    Window* w = (Window*)ptr;
    [w->glView setSpinningHidden:(b==0)];
}
void toggleFullScreen(void* ptr) {
    Window* w = (Window*)ptr;
    [w toggleFullScreen:nil];
}

void setControlsVisible(void* ptr, int b) {
    Window* w = (Window*)ptr;

    BOOL hidden = (b==0);
    [w->glView setCursorHidden:hidden];
    [w->glView setPlaybackViewHidden:hidden];
    [w setTitleHidden:hidden];
}

CSize getScreenSize() {
    NSSize sz = [[NSScreen mainScreen] frame].size;
    CSize csz;
    csz.width = (int)sz.width;
    csz.height = (int)sz.height;
    return csz;
}

void setVolume(void* wptr, int volume) {
    Window* w = (Window*)wptr;
    [w->glView setVolume:volume];
}

void setVolumeVisible(void* wptr, int b) {
    Window* w = (Window*)wptr;
    [w->glView setVolumeHidden:(b==0)];
}

void alert(void* wptr, char* str) {
    Window* w = (Window*)wptr;
    [w setDelegate:nil];  //remove delegate prevent hide title bar
    [w->glView setPlaybackViewHidden:NO];

    NSAlert* alert = [[NSAlert alloc] init];
    [alert setMessageText:[NSString stringWithUTF8String:str]];
    [alert setAlertStyle:NSCriticalAlertStyle];

    [alert beginSheetModalForWindow:w modalDelegate:w didEndSelector:@selector(close) contextInfo:nil];
}

void closeWindow(void* wptr) {
    Window* w = (Window*)wptr;
    [w close];
}

void addRecentOpenedFile(char* str) {
    NSString* filename = [NSString stringWithUTF8String:str];
    [[NSDocumentController sharedDocumentController] noteNewRecentDocumentURL:[NSURL fileURLWithPath:filename]];
}