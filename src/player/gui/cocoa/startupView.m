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
    [[NSColor whiteColor] setFill];
    NSRectFill(dirtyRect);

    [super drawRect:dirtyRect];
}

@end
