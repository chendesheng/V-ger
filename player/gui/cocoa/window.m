#import "window.h"
#include <stdlib.h>

@implementation Window
- (id)initWithWidth:(int)w height:(int)h  {
	unsigned int styleMask = NSTitledWindowMask | NSClosableWindowMask 
		| NSMiniaturizableWindowMask | NSResizableWindowMask | NSFullSizeContentViewWindowMask;

    CGFloat screenh = [[NSScreen mainScreen] frame].size.height;
    self = [super initWithContentRect:NSMakeRect(50, screenh - h + 22 - 150, w, h-22)
    	styleMask:styleMask
    	backing:NSBackingStoreBuffered
      	defer:YES];

    [self setHasShadow:YES];
    [self setContentMinSize:NSMakeSize(300, 300*h/w)];
    [self setAcceptsMouseMovedEvents:YES];
	[self setRestorable:NO];
    [self setCollectionBehavior:NSWindowCollectionBehaviorFullScreenPrimary];
    [self setOpaque:YES];

    NSRect bounds = NSMakeRect(0, 0, w, h);
    glView = [[GLView alloc] initWithFrame2:bounds];
    self.contentView = glView;
    self.aspectRatio = bounds.size;

    glView.frame = bounds;
    [glView setAutoresizingMask:NSViewWidthSizable|NSViewHeightSizable];
    
    return self;
}

- (BOOL)canBecomeKeyWindow {
    return YES;
}
- (BOOL)isMovableByWindowBackground {
    return YES;
}
- (void)makeCurrentContext {
    [NSOpenGLContext clearCurrentContext];
    [[self->glView openGLContext] makeCurrentContext];
}

- (void)audioMenuItemClick:(id)sender {
    NSMenuItem* audioMenuItem = (NSMenuItem*)sender;
    if ([audioMenuItem state] == NSOnState) {
        return;
    }

    NSMenu *menu = [[audioMenuItem parentItem] submenu];
    for (int i = 0; i < [menu numberOfItems]; i++) {
        NSMenuItem *item = [menu itemAtIndex:i];
        [item setState:NSOffState];
    }

    [audioMenuItem setState:NSOnState];
    onMenuClick(MENU_AUDIO, (int)[audioMenuItem tag]);
}

- (void)subtitleMenuItemClick:(id)sender {
    NSMenuItem* subtitleMenuItem = (NSMenuItem*)sender;
    if ([subtitleMenuItem state] == NSOnState)
        onMenuClick(MENU_SUBTITLE, (int)[subtitleMenuItem tag]);
    else
        onMenuClick(MENU_SUBTITLE, (int)[subtitleMenuItem tag]);
}

-(void)close {
    [super close];
    [NSApp terminate:nil];
}

-(void)setTitleHidden:(BOOL)b {
    if ([self isFullScreen]) {
        return;
    }
    //self.titlebarAppearsTransparent = b;

    NSView* fv = [self.contentView superview];
    for (NSView* v in fv.subviews) {
        if (v != glView) {
            [v setHidden:b];
        }
    }
}

-(void)playPause:(id)sender {
    onMenuClick(MENU_PLAY, 0);
}

-(void)selectSubtitle:(id)sender {
}

-(void)selectSubtitleItem:(id)sender {
    NSMenuItem* item = (NSMenuItem*)sender;
    onMenuClick(MENU_SUBTITLE, (int)item.tag);
}

-(void)selectAudio:(id)sender {
}

-(void)selectAudioItem:(id)sender {
    NSMenuItem* item = (NSMenuItem*)sender;
    onMenuClick(MENU_AUDIO, (int)item.tag);
}

-(void)seekBackwardBySubtitle:(id)sender {
    onMenuClick(MENU_SEEK, 0);
}

-(void)seekForwardBySubtitle:(id)sender {
    onMenuClick(MENU_SEEK, 1);
}

-(void)seekBackwardByTime:(id)sender {
    onMenuClick(MENU_SEEK, 2);
}

-(void)seekForwardByTime:(id)sender {
    onMenuClick(MENU_SEEK, 3);
}

-(void)increaseVolume:(id)sender {
    onMenuClick(MENU_VOLUME, 1);
}

-(void)decreaseVolume:(id)sender {
    onMenuClick(MENU_VOLUME, -1);
}

-(void)pullMainSubtitle:(id)sender {
    onMenuClick(MENU_SYNC_SUBTITLE, 0);
}

-(void)pushMainSubtitle:(id)sender {
    onMenuClick(MENU_SYNC_SUBTITLE, 1);
}

-(void)pullSecondSubtitle:(id)sender {
    onMenuClick(MENU_SYNC_SUBTITLE, 2);
}

-(void)pushSecondSubtitle:(id)sender {
    onMenuClick(MENU_SYNC_SUBTITLE, 3);
}

-(void)pullAudio:(id)sender {
    onMenuClick(MENU_SYNC_AUDIO, -1);
}

-(void)pushAudio:(id)sender {
    onMenuClick(MENU_SYNC_AUDIO, 1);
}

-(void)searchSubtitle:(id)sender {
    onMenuClick(MENU_SEARCH_SUBTITLE, 0);
}

- (BOOL)validateMenuItem:(NSMenuItem *)item {
    if ([item action] == @selector(playPause:)) {
        if (isPlaying()) {
            item.title = @"Pause";
        } else {
            item.title = @"Play";
        }
    } else if ([item action] == @selector(selectSubtitle:)) {
        NSMenu* menu = item.submenu;
        [menu removeAllItems];

        char** names;
        int length;
        int firstSub, secondSub;
        getAllSubtitleNames((void***)&names, &length);  //read subtitles very time open the menu
        getPlayingSubtitles(&firstSub, &secondSub);
        if (length == 0) {
            [menu addItemWithTitle:@"None" action:nil keyEquivalent:@""];
        } else {
            for (int i = 0; i < length; i++) {
                NSMenuItem * submenuItem = [menu addItemWithTitle:[NSString stringWithUTF8String:names[i]] action:@selector(selectSubtitleItem:) keyEquivalent:@""];
                submenuItem.tag = i;
                if (i == firstSub || i == secondSub) {
                    submenuItem.state = NSOnState;
                } else {
                    submenuItem.state = NSOffState;
                }
                free(names[i]);
            }
        }
    } else if ([item action] == @selector(selectAudio:)) {
        NSMenu* menu = item.submenu;
        [menu removeAllItems];

        char** names;
        int length;
        getAllAudioTracks((void***)&names, &length);

        int selected = getPlayingAudioTrack();

        if (length == 0) {
            [menu addItemWithTitle:@"None" action:nil keyEquivalent:@""];
        } else {
            for (int i = 0; i < length; i++) {
                NSMenuItem * submenuItem = [menu addItemWithTitle:[NSString stringWithUTF8String:names[i]] action:@selector(selectAudioItem:) keyEquivalent:@""];
                submenuItem.tag = i;
                if (i == selected) {
                    submenuItem.state = NSOnState;
                } else {
                    submenuItem.state = NSOffState;
                }
                free(names[i]);
            }
        }
    } else if ([item action] == @selector(seekForwardBySubtitle:) || 
        [item action] == @selector(seekBackwardBySubtitle:) ||
        [item action] == @selector(pullMainSubtitle:) ||
        [item action] == @selector(pushMainSubtitle:)) {
        int firstSub, secondSub;
        getPlayingSubtitles(&firstSub, &secondSub);
        if (firstSub == -1) {
            return NO;
        }
    } else if ([item action] == @selector(pullSecondSubtitle:) ||
        [item action] == @selector(pushSecondSubtitle:)) {
        int firstSub, secondSub;
        getPlayingSubtitles(&firstSub, &secondSub);
        if (secondSub == -1) {
            return NO;
        }
    } else if ([item action] == @selector(increaseVolume:) || 
        [item action] == @selector(decreaseVolume:) ||
        [item action] == @selector(pullAudio:) ||
        [item action] == @selector(pushAudio:)) {
        if (getPlayingAudioTrack() == -1) {
            return NO;
        }
    } else if ([item action] == @selector(searchSubtitle:)) {
        if (isSearchingSubtitle() != 0) {
            item.title = @"Stop Search";
        } else {
            item.title = @"Search Online";
        }
    } else if ([item action] == @selector(alwaysOnTop:)) {
        if ([self isFullScreen]) {
            item.state = NSOffState;
            return NO;
        } else {
            if (self.level == NSFloatingWindowLevel) {
                item.state = NSOnState;
            } else {
                item.state = NSOffState;
            }
        }
    }
    return YES;
}

-(void)alwaysOnTop:(id)sender {
    if (self.level == NSNormalWindowLevel) {
        self.level = NSFloatingWindowLevel;
    } else {
        self.level = NSNormalWindowLevel;
    }
}

-(void)open:(id)sender {
    setControlsVisible(self, 1, 0);
    onOpenOpenPanel();

    NSOpenPanel *panel = [NSOpenPanel openPanel];
    [panel setCanChooseDirectories:NO];
    [panel setAllowsMultipleSelection:NO];
    [panel beginSheetModalForWindow: self completionHandler:^(NSInteger result){
        setControlsVisible(self, 1, 1);

        if(result == NSFileHandlingPanelOKButton){
            NSString* filename = [[panel URL] path];
            char* cfilename = (char*)[filename UTF8String];
            onCloseOpenPanel(cfilename);
        } else {
            onCloseOpenPanel("");
            return;
        }
    }];
}

- (BOOL)isFullScreen {
    return (([self styleMask] & NSFullScreenWindowMask) == NSFullScreenWindowMask);
}

- (void)magnifyWithEvent:(NSEvent *)event {
    NSLog(@"magnifyWithEvent %f", [event magnification]);
    if ([event magnification] < 0 && [self isFullScreen]) {
        [self toggleFullScreen:nil];
    } else if ([event magnification] > 0 && ![self isFullScreen]) {
        [self toggleFullScreen:nil];
    }
}

- (void)fatal:(NSString *)message {
    NSAlert* alert = [[NSAlert alloc] init];
    [alert setMessageText:message];
    [alert setAlertStyle:NSCriticalAlertStyle];

    [[NSNotificationCenter defaultCenter] addObserver:self
        selector:@selector(openPanelDidClose:)
        name:NSWindowDidEndSheetNotification
        object:self];

    [alert beginSheetModalForWindow:self completionHandler: ^(NSInteger result){
        self->isFatalHappen = YES;
    }];
}

- (void)openPanelDidClose:(NSNotification *)notification {
    if (self->isFatalHappen == YES) {
        self->isFatalHappen = NO;
        [self close];
    }

    [[NSNotificationCenter defaultCenter] removeObserver:self
        name:NSWindowDidEndSheetNotification
        object:self];
}
@end



