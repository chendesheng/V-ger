#import "window.h"

@implementation Window
- (id)initWithTitle:(NSString*)title width:(int)w height:(int)h  {
	unsigned int styleMask = NSTitledWindowMask | NSClosableWindowMask 
		| NSMiniaturizableWindowMask | NSResizableWindowMask;

    self = [super initWithContentRect:NSMakeRect(0,0,w,h-22)
    	styleMask:styleMask
    	backing:NSBackingStoreBuffered
      	defer:YES];

    [self setTitle:title];
    self->customAspectRatio = NSMakeSize(w, h);
    [self setHasShadow:YES];
    [self setContentMinSize:NSMakeSize(500, 500*h/w)];
    [self setAcceptsMouseMovedEvents:YES];
	[self setRestorable:NO];
    [self setCollectionBehavior:NSWindowCollectionBehaviorFullScreenPrimary];

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
    onMenuClicked(MENU_AUDIO, (int)[audioMenuItem tag]);
}
- (void)subtitleMenuItemClick:(id)sender {
    NSMenuItem* subtitleMenuItem = (NSMenuItem*)sender;
    if ([subtitleMenuItem state] == NSOnState)
        onMenuClicked(MENU_SUBTITLE, (int)[subtitleMenuItem tag]);
    else
        onMenuClicked(MENU_SUBTITLE, (int)[subtitleMenuItem tag]);
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