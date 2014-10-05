#import "window.h"

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
    [self setContentMinSize:NSMakeSize(100, 100*h/w)];
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
@end



