#import "popupView.h"
@implementation PopupView

- (id)initWithFrame:(NSRect)frame {
    NSLog(@"popupView initWithFrame");
    self = [super initWithFrame:frame];    

    // [self setHasVerticalScroller:YES];
    // [self setHasHorizontalScroller:YES];
    // [self setBorderType:NSNoBorder];
    // [self setDrawsBackground:NO];
    // [[self enclosingScrollView] setDrawsBackground:NO];
    // [self setBackgroundColor:[NSColor clearColor]];

    // NSView *testView = [[NSView alloc] initWithFrame:NSMakeRect(0,0,1000,1000)];
    // [self setContentView:testView];

    return self;
}
// -(void)drawRect:(NSRect)dirtyRect
// {
//     [[NSColor clearColor] setFill];
//     NSRectFill(dirtyRect);
// }
@end