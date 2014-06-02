#import "startupView.h"

@implementation StartupView

- (id)initWithFrame:(NSRect)frame
{
    self = [super initWithFrame:frame];
    if (self) {
        _progressIndicator = [[NSProgressIndicator alloc] initWithFrame:frame];
        [_progressIndicator setStyle:1];

        [_progressIndicator startAnimation:nil];

        [self addSubview:_progressIndicator];
	    [_progressIndicator setAutoresizingMask:NSViewWidthSizable|NSViewHeightSizable];
    }
    
    return self;
}

- (void)drawRect:(NSRect)dirtyRect
{
    [[NSColor blackColor] setFill];
    NSRectFill(dirtyRect);

    [[NSColor whiteColor] setFill];
    NSRect frame = [self bounds];
    NSRect rt = NSMakeRect((frame.size.width - 50)/2, (frame.size.height - 50)/2, 50, 50);
    // NSRectFill(rt);

    NSShadow* shadow = [[NSShadow alloc] init];
    [shadow setShadowBlurRadius:5.0];
    [shadow setShadowColor:[NSColor blackColor]];
    [shadow set];

    NSBezierPath *path =
      [NSBezierPath bezierPathWithRoundedRect:rt
                                      xRadius:3.0f
                                      yRadius:3.0f];
    [path fill];

    [super drawRect:dirtyRect];
}

@end
