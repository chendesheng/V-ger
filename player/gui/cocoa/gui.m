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

	//create memu bar
	id menubar = [[NSMenu new] autorelease];
    id appMenuItem = [[NSMenuItem new] autorelease];
    [menubar addItem:appMenuItem];
    [NSApp setMainMenu:menubar];
    id appMenu = [[NSMenu new] autorelease];

    NSMenuItem *searchSubtitleMenuItem = [[[NSMenuItem alloc] initWithTitle:@"Search Subtitle"
        action:@selector(searchSubtitleMenuItemClick:) keyEquivalent:@"s"] autorelease];
    [searchSubtitleMenuItem setKeyEquivalentModifierMask:0];
    [searchSubtitleMenuItem setTarget: appDelegate];
    [appMenu addItem:searchSubtitleMenuItem];

    NSMenuItem *openFileMenuItem = [[[NSMenuItem alloc] initWithTitle:@"Open..."
        action:@selector(openFileMenuItemClick:) keyEquivalent:@"o"] autorelease];
    [openFileMenuItem setTarget: appDelegate];
    [appMenu addItem:openFileMenuItem];

    id appName = [[NSProcessInfo processInfo] processName];
    id quitTitle = [@"Quit " stringByAppendingString:appName];
    id quitMenuItem = [[[NSMenuItem alloc] initWithTitle:quitTitle
        action:@selector(terminate:) keyEquivalent:@"q"] autorelease];
    [appMenu addItem:quitMenuItem];


    [appMenuItem setSubmenu:appMenu];

    [[NSAppleEventManager sharedAppleEventManager] setEventHandler:appDelegate andSelector:@selector(handleAppleEvent:withReplyEvent:) forEventClass:kInternetEventClass andEventID:kAEGetURL];
}
NSMenuItem* getTopMenuByTitle(NSString* title) {
    NSMenu* menubar = [NSApp mainMenu];
    NSArray* menus = [menubar itemArray];
    for (NSMenuItem* menu in menus) {
        if ([menu title] == title) {
            return menu;
        }
    }

    return nil;
}
void hideMenuNSString(NSString* title) {
    @autoreleasepool {
        NSMenuItem* item = getTopMenuByTitle(title);
        if (item != nil) {
            [[NSApp mainMenu] removeItem:item];
        }
    }
}
void hideSubtitleMenu() {
    hideMenuNSString(@"Subtitle");
}
void hideAudioMenu() {
    hideMenuNSString(@"Audio");
}
void initAudioMenu(void* wptr, char** names, int32_t* tags, int len, int selected) {
    hideAudioMenu();

    @autoreleasepool {
        if (len > 0) {
            NSWindow* w = (NSWindow*)wptr;

            NSMenu *menubar = [NSApp mainMenu];
            NSMenuItem* audioMenuItem = [[NSMenuItem new] autorelease];
            [audioMenuItem setTitle:@"Audio"];
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
    }
}

void initSubtitleMenu(void* wptr, char** names, int32_t* tags, int len, int32_t selected1, int32_t selected2) {
    hideSubtitleMenu();
    
    @autoreleasepool {
        if (len > 0) {
            NSWindow* w = (NSWindow*)wptr;

            NSMenu* menubar = [NSApp mainMenu];
            NSMenuItem* subtitleMenuItem = [[NSMenuItem new] autorelease];
            [subtitleMenuItem setTitle:@"Subtitle"];
            [menubar addItem:subtitleMenuItem];
            NSMenu* subtitleMenu = [[NSMenu alloc] initWithTitle:@"Subtitle"];

            for (int i = 0; i < len; i++) {
                char* name = names[i];
                int tag = tags[i];
                NSMenuItem* item = [[NSMenuItem alloc] initWithTitle:[NSString stringWithUTF8String:name] 
                    action:@selector(subtitleMenuItemClick:) keyEquivalent:@""];
                [item autorelease];
                [item setTarget: w];
                [item setTag: tag];
                [subtitleMenu addItem:item];

                if (tag == selected1 || tag == selected2) {
                    [item setState: NSOnState];
                }
            }

            [subtitleMenuItem setSubmenu:subtitleMenu];
        }
    }
}
void setSubtitleMenuItem(int t1, int t2) {
    @autoreleasepool {
        NSMenuItem* menu = getTopMenuByTitle(@"Subtitle");
        for (NSMenuItem* item in [[menu submenu] itemArray]) {
            int tag = (int)[item tag];
            if (tag == t1 || tag == t2) {
                [item setState:NSOnState];
            } else {
                [item setState:NSOffState];
            }
        }
    
    }
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
void makeWindowCurrentContext(void*ptr) {
    Window* w = (Window*)ptr;
    [w makeCurrentContext];
}
void pollEvents() {
    [NSApp run];
}
void refreshWindowContent(void*wptr) {
	Window* w = (Window*)wptr;
    [w->glView setNeedsDisplay:YES];
}

int getWindowWidth(void* ptr) {
    Window* w = (Window*)ptr;
    return (int)([w->glView frame].size.width);
}
int getWindowHeight(void* ptr) {
    Window* w = (Window*)ptr;
    return (int)([w->glView frame].size.height);
}
void showWindowProgress(void* ptr, char* left, char* right, double percent) {
    Window* w = (Window*)ptr;
    
    NSString* leftStr;
    if (strlen(left) == 0) {
        leftStr = @"00:00:00";
    } else {
        leftStr = [[NSString stringWithUTF8String:left] retain];
    }
    NSString* rightStr;
    if (strlen(left) == 0) {
        rightStr = @"00:00:00";
    } else {
        rightStr = [[NSString stringWithUTF8String:left] retain];
    }
    [w->glView updatePorgressInfo:leftStr rightString:rightStr percent:percent];
}
void showWindowBufferInfo(void* ptr, char* speed, double percent) {
    Window* w = (Window*)ptr;
    NSString* str;
    if (strlen(speed) == 0) {
        str = @"";
    } else {
        str = [[NSString stringWithUTF8String:speed] retain];
    }
    [w->glView updateBufferInfo:str bufferPercent:percent];
}
void* showText(void* ptr, SubItem* item) {
    Window* w = (Window*)ptr;
    return [w->glView showText:item];
}
void hideText(void* ptrWin, void* ptrText) {
    Window* w = (Window*)ptrWin;
    [w->glView hideText:ptrText];
}
void windowHideStartupView(void* ptr) {
    Window* w = (Window*)ptr;
    [w->glView hideStartupView];
}
void windowShowStartupView(void* ptr) {
    Window* w = (Window*)ptr;
    [w->glView showStartupView];
}
void showSpinning(void* ptr) {
    Window* w = (Window*)ptr;
    [w->glView setSpinningHidden:NO];
}
void hideSpinning(void* ptr) {
    Window* w = (Window*)ptr;
    [w->glView setSpinningHidden:YES];
}
void windowToggleFullScreen(void* ptr) {
    Window* w = (Window*)ptr;
    [w toggleFullScreen:nil];
}

void hideCursor(void* ptr) {
    Window* w = (Window*)ptr;
    [w->glView hideCursor];
    [w->glView hideProgress];
    [w setTitleHidden:YES];
}

void showCursor(void* ptr) {
    Window* w = (Window*)ptr;
    [w->glView showCursor];
    [w->glView showProgress];
    [w setTitleHidden:NO];
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

void setVolumeDisplay(void* wptr, int show) {
    Window* w = (Window*)wptr;
    [w->glView setVolumeHidden:(show==0)];
}

void alert(void* wptr, char* str) {
    Window* w = (Window*)wptr;
    [w setDelegate:nil];  //remove delegate prevent hide title bar
    [w->glView showProgress];

    NSAlert* alert = [[NSAlert alloc] init];
    [alert setMessageText:[NSString stringWithUTF8String:str]];
    [alert setAlertStyle:NSCriticalAlertStyle];

    [alert beginSheetModalForWindow:w modalDelegate:w didEndSelector:@selector(close) contextInfo:nil];
}

void closeWindow(void* wptr) {
    Window* w = (Window*)wptr;
    [w close];
}