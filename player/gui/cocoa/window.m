#import "window.h"

@implementation Window
- (id)initWithTitle:(NSString*)title width:(int)w height:(int)h  {
	unsigned int styleMask = NSTitledWindowMask | NSClosableWindowMask 
		| NSMiniaturizableWindowMask | NSResizableWindowMask;

    // NSRect rt = [super contentRectForFrameRect:NSMakeRect(0,0,w,h)];
    self = [super initWithContentRect:NSMakeRect(0,0,w,h-22)
    	styleMask:styleMask
    	backing:NSBackingStoreBuffered
      	defer:YES];

    [self setTitle:title];
    // [self setContentAspectRatio:NSMakeSize(w, h)];
    self->customAspectRatio = NSMakeSize(w, h);
    [self setHasShadow:YES];
    [self setContentMinSize:NSMakeSize(500, 500*h/w)];
    [self setAcceptsMouseMovedEvents:YES];
	[self setRestorable:NO];
    [self setCollectionBehavior:NSWindowCollectionBehaviorFullScreenPrimary];

    // [self setBackgroundColor:[NSColor clearColor]];
    [self setOpaque:YES];

    [self center];

    return self;
}

- (BOOL)canBecomeKeyWindow {
    return YES;
}
- (BOOL)isMovableByWindowBackground {
    return YES;
}
- (void)setContentViewNeedsDisplay:(BOOL)b {
	[self->glView setNeedsDisplay:b];
}
- (void)timerTick:(NSEvent *)event {
    NSAutoreleasePool * pool = [[NSAutoreleasePool alloc] init];
	onTimerTick((void*)self);
    [pool drain];
}
- (void)makeCurrentContext {
    [NSOpenGLContext clearCurrentContext];
    [[self->glView openGLContext] makeCurrentContext];
}

- (void)audioMenuItemClick:(id)sender {
    // NSLog(sender);
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
    onMenuClicked(MENU_AUDIO, [audioMenuItem tag]);
}
- (void)subtitleMenuItemClick:(id)sender {
    // NSLog(@"subtitleMenuItemClick");
    // NSMenuItem* subtitleMenuItem = (NSMenuItem*)sender;
    // if ([subtitleMenuItem state] == NSOnState) {
    //     [subtitleMenuItem setState:NSOffState];
    //     onSubtitleMenuClicked(self, [subtitleMenuItem tag], 0);
    //     return;
    // }

    // NSMenu *menu = [[subtitleMenuItem parentItem] submenu];
    // int cnt = 0;
    // for (int i = 0; i < [menu numberOfItems]; i++) {
    //     NSMenuItem *item = [menu itemAtIndex:i];
    //     if ([item state] == NSOnState) {
    //         cnt++;
    //     }
    // }
    // if (cnt == 2) {
    //     return;
    // } else {
    //     [subtitleMenuItem setState:NSOnState];
    //     onSubtitleMenuClicked(self, [subtitleMenuItem tag], 1);
    // }
    NSMenuItem* subtitleMenuItem = (NSMenuItem*)sender;
    if ([subtitleMenuItem state] == NSOnState)
        onMenuClicked(MENU_SUBTITLE, [subtitleMenuItem tag]);
    else
        onMenuClicked(MENU_SUBTITLE, [subtitleMenuItem tag]);
}
- (void)updateRoundCorner {
    NSView* rv = [[self contentView] superview];
    rv.layer.cornerRadius=4.1;
    rv.layer.masksToBounds=YES;
}

-(void)close {
    [super close];
    [NSApp terminate:nil];
}
@end