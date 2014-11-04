#import "spinningView.h"

@implementation SpinningView

- (id)initWithFrame:(NSRect)frame
{
    self = [super initWithFrame:frame];
    if (self) {

        NSRect rt = NSMakeRect(0, 0, frame.size.width, frame.size.height);

        _progressIndicator = [[NSProgressIndicator alloc] initWithFrame:rt];
        [_progressIndicator setStyle:1];
        [_progressIndicator startAnimation:nil];
        [_progressIndicator setAutoresizingMask:NSViewWidthSizable|NSViewHeightSizable];

        _blurView = [[BlurView alloc] initWithView:_progressIndicator frame:rt];
        [_blurView setCornerRadius:4.1];

        [self addSubview:_blurView];
    }
    
    return self;
}
@end
