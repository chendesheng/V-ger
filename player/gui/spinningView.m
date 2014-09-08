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

        _blurView = [[BlurView alloc] initWithFrame:rt];
        _blurView.wantsLayer = YES;
        _blurView.layer.masksToBounds = YES;
        _blurView.layer.cornerRadius = 4.1;

        [self addSubview:_blurView];
        [self addSubview:_progressIndicator positioned:NSWindowAbove relativeTo:nil];
    }
    
    return self;
}
@end
