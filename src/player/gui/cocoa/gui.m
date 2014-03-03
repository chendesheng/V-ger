#include "gui.h"

#import <Cocoa/Cocoa.h>
#import "window.h"
#import "windowDelegate.h"
#import "glView.h"
#import "textView.h"
#import "blurView.h"
#import "progressView.h"
#import "popupView.h"
#import "subtitleView.h"
#import "startupView.h"
#import "app.h"

void initialize() {
    if (NSApp)
        return;

	[NSApplication sharedApplication];

    [NSApp setActivationPolicy:NSApplicationActivationPolicyRegular];

    Application *appDelegate = [[Application alloc] init];
    [NSApp setDelegate:appDelegate];

	NSLog(@"initialized");
	//create memu bar
	id menubar = [[NSMenu new] autorelease];
    id appMenuItem = [[NSMenuItem new] autorelease];
    [menubar addItem:appMenuItem];
    [NSApp setMainMenu:menubar];
    id appMenu = [[NSMenu new] autorelease];

    NSMenuItem *searchSubtitleMenuItem = [[[NSMenuItem alloc] initWithTitle:@"Search Subtitle"
        action:@selector(searchSubtitleMenuItemClick:) keyEquivalent:@""] autorelease];
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
}
void hideMenuNSString(NSString* title) {
    NSLog(@"remove subtitle menu %@", title);

    NSMenu* menubar = [NSApp mainMenu];
    NSArray* menus = [menubar itemArray];
    for (NSMenuItem* menu in menus) {
        NSLog(@"compare %@ to %@", [menu title], title);
        if ([menu title] == title) {
            NSLog(@"remove menu item");
            [menubar removeItem:menu];
            break;
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

void initSubtitleMenu(void* wptr, char** names, int32_t* tags, int len, int32_t selected1, int32_t selected2) {
    NSAutoreleasePool * pool = [[NSAutoreleasePool alloc] init];

    hideSubtitleMenu();

    NSLog(@"len:%d", len);

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
            [item setTarget: w];
            [item setTag: tag];
            [subtitleMenu addItem:item];

            if (tag == selected1) {
                [item setState: NSOnState];
            }

            if (tag == selected2) {
                [item setState: NSOnState];
            }
        }

        [subtitleMenuItem setSubmenu:subtitleMenu];
    }

    [pool drain];
}

void setWindowTitle(void* wptr, char* title) {
    Window* w = (Window*)wptr;
    [w setTitle:[NSString stringWithUTF8String:title]];
}

void setWindowSize(void* wptr, int width, int height) {
    Window* w = (Window*)wptr;

    NSRect frame = [w frame];
    frame.origin.y -= (height - frame.size.height)/2;
    frame.origin.x -= (width - frame.size.width)/2;
    frame.size = NSMakeSize(width, height);

    [w setFrame:frame display:YES animate:YES];
}

void* newWindow(char* title, int width, int height) {
	NSAutoreleasePool * pool = [[NSAutoreleasePool alloc] init];

	initialize();

	Window* w = [[Window alloc] initWithTitle:[NSString stringWithUTF8String:title]
		width:width height:height];
	
    WindowDelegate* wd = (WindowDelegate*)[[WindowDelegate alloc] init];
    wd->window = w;
	[w setDelegate:(id)wd];

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

    // BlurView* bvPopup = [[BlurView alloc] initWithFrame:NSMakeRect(200,40,400,500)];
    // [v addSubview:bvPopup];
    // [bvPopup setAutoresizingMask:NSViewWidthSizable];

    // PopupView* ppv = [[PopupView alloc] initWithFrame:NSMakeRect(0,0,400,500)];
    // [bvPopup addSubview:ppv];
    // [ppv setAutoresizingMask:NSViewWidthSizable|NSViewHeightSizable];


    StartupView* sv = [[StartupView alloc] initWithFrame:[v frame]];
    [v addSubview:sv];
    [sv setAutoresizingMask:NSViewWidthSizable|NSViewHeightSizable];

    [v setStartupView:sv];


    NSTimer *renderTimer = [NSTimer timerWithTimeInterval:1.0/100.0 
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
void showWindowProgress(void* ptr, char* left, char* right, double percent, double percent2) {
    Window* w = (Window*)ptr;
    [[w contentView] showProgress:left right:right percent:percent percent2:percent2];
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
void windowShowStartupView(void* ptr) {
    Window* w = (Window*)ptr;
    [[w contentView] showStartupView];
}
void windowToggleFullScreen(void* ptr) {
    Window* w = (Window*)ptr;
    [w toggleFullScreen:nil];
}

void *newDialog(char* title, int width, int height) {
    NSAutoreleasePool * pool = [[NSAutoreleasePool alloc] init];

    NSPanel *dialog = [[NSPanel alloc] initWithContentRect:NSMakeRect(200.0, 200.0, 300, 200)
        styleMask:NSHUDWindowMask | NSClosableWindowMask | NSTitledWindowMask | NSUtilityWindowMask | NSResizableWindowMask
          backing:NSBackingStoreBuffered
            defer:YES];

    [dialog makeKeyAndOrderFront:nil];

    [dialog setTitle:[NSString stringWithUTF8String:title]];
    
    SubtitleView* sv = [[SubtitleView alloc] initWithFrame:NSMakeRect(0,0,dialog.frame.size.width,dialog.frame.size.height)];
    [dialog setContentView:sv];
    [sv setAutoresizingMask:NSViewWidthSizable|NSViewHeightSizable];


    [pool drain];

    return dialog;
}

CSize getScreenSize() {
    NSSize sz = [[NSScreen mainScreen] frame].size;
    CSize csz;
    csz.width = (int)sz.width;
    csz.height = (int)sz.height;
    return csz;
}