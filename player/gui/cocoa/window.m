#import "window.h"
#include <stdlib.h>

@implementation Window
- (id)initWithWidth:(int)w height:(int)h  {
	unsigned int styleMask = NSTitledWindowMask | NSClosableWindowMask 
		| NSMiniaturizableWindowMask | NSResizableWindowMask;

    self = [super initWithContentRect:NSMakeRect(0,0,w,h-22)
    	styleMask:styleMask
    	backing:NSBackingStoreBuffered
      	defer:YES];

    self->customAspectRatio = NSMakeSize(w, h);
    [self setHasShadow:YES];
    [self setContentMinSize:NSMakeSize(200, 200*h/w)];
    [self setAcceptsMouseMovedEvents:YES];
	[self setRestorable:NO];
    [self setCollectionBehavior:NSWindowCollectionBehaviorFullScreenPrimary];
    [self setOpaque:YES];
    
    [self center];
    
    
    //Window > NSFrameView > NSOpenGLView > TitleView & ProgressView & TextView
    NSRect bounds = NSMakeRect(0, 0, w, h);

    NSView* fv = [[self contentView] superview];
    glView = [[GLView alloc] initWithFrame2:bounds];

    [fv addSubview:glView positioned:NSWindowBelow relativeTo:nil];
    fv.wantsLayer = YES;
    
    glView.frame = bounds;
    [glView setAutoresizingMask:NSViewWidthSizable|NSViewHeightSizable];
    
    TitleTextView* ttv = [[TitleTextView alloc] init];
    BlurView* tiv = [[BlurView alloc] initWithView:ttv frame:NSMakeRect(0,h-22,w,22)];
    [tiv setAutoresizingMask:NSViewWidthSizable|NSViewMinYMargin];

    //must add title view to glView
    [glView addSubview:tiv positioned:NSWindowAbove relativeTo:nil];
    self->titleTextView = ttv;
    self->bvTitleTextView = tiv;

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
- (void)updateRoundCorner {
    NSView* fv = [[self contentView] superview];
    fv.layer.cornerRadius=4.1;
    fv.layer.masksToBounds=YES;
}

-(void)close {
    [super close];
    [NSApp terminate:nil];
}
-(void)setFrame:(NSRect)frameRect display:(BOOL)flag {
    // Maintain round corner when resizing window
    // Remove this window's round corner disappear after resize, don't known why.
    [self updateRoundCorner];
    
    [super setFrame:frameRect display:flag];
}
-(void)setTitleHidden:(BOOL)b {
    NSView* fv = [self.contentView superview];
    if (b) {
        [bvTitleTextView setHidden:YES];
        for (NSView* v in [fv subviews]) {
            if (v != glView) {
                [v setHidden:YES];
            }
        }
    } else {
        [bvTitleTextView setHidden:NO];
        for (NSView* v in fv.subviews) {
            if (v != glView) {
                [v setHidden:NO];
            }
        }
    }
}

-(void)setTitle:(NSString *)title {
    titleTextView.title = title;
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
    }
    return YES;
}
@end



