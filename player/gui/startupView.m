#import "startupView.h"

@implementation StartupView

- (id)initWithFrame:(NSRect)frame
{
    return [super initWithFrame:frame];
}

- (void)drawRect:(NSRect)dirtyRect
{
    [[NSColor blackColor] setFill];
    NSRectFill(dirtyRect);
}

@end
