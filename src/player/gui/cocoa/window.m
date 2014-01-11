#import "window.h"

@implementation Window
- (id)initWithTitle:(NSString*)title width:(int)w height:(int)h  {
	unsigned int styleMask = NSTitledWindowMask | NSClosableWindowMask 
		| NSMiniaturizableWindowMask | NSResizableWindowMask;

	NSLog(@"%dx%d", w, h);

    self = [super initWithContentRect:NSMakeRect(0, 0, w, h)
    	styleMask:styleMask
    	backing:NSBackingStoreBuffered
      	defer:NO];

    [self setTitle:title];
    // [self setContentAspectRatio:NSMakeSize(w, h)];
    self->customAspectRatio = NSMakeSize(w, h);
    [self setOpaque:YES];
    [self setHasShadow:YES];
    [self setContentMinSize:NSMakeSize(500, 500*h/w)];
    [self setAcceptsMouseMovedEvents:YES];
	[self setRestorable:NO];
    [self setCollectionBehavior:NSWindowCollectionBehaviorFullScreenPrimary];

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
	[[self contentView] setNeedsDisplay:b];
}
- (void)timerTick:(NSEvent *)event {
	onTimerTick((void*)self);
}
- (void)makeCurrentContext {
    [NSOpenGLContext clearCurrentContext];
    [[[self contentView] openGLContext] makeCurrentContext];
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
    onAudioMenuClicked(self, [audioMenuItem tag]);
}
- (void)subtitleMenuItemClick:(id)sender {
    NSLog(@"subtitleMenuItemClick");
    NSMenuItem* subtitleMenuItem = (NSMenuItem*)sender;
    if ([subtitleMenuItem state] == NSOnState) {
        [subtitleMenuItem setState:NSOffState];
        onSubtitleMenuClicked(self, [subtitleMenuItem tag], 0);
        return;
    }

    NSMenu *menu = [[subtitleMenuItem parentItem] submenu];
    int cnt = 0;
    for (int i = 0; i < [menu numberOfItems]; i++) {
        NSMenuItem *item = [menu itemAtIndex:i];
        if ([item state] == NSOnState) {
            cnt++;
        }
    }
    if (cnt == 2) {
        return;
    } else {
        [subtitleMenuItem setState:NSOnState];
        onSubtitleMenuClicked(self, [subtitleMenuItem tag], 1);
    }
}
@end