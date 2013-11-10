#import <Cocoa/Cocoa.h>
#import "gui.h"
#import "progressView.h"
#import "textView.h"
@interface GLView : NSOpenGLView 

{
    NSTrackingArea* trackingArea;
    NSCursor* noneCursor;
    ProgressView* progressView;
    TextView* textView;
    NSCursor* currentCursor;
}

-(id)initWithFrame2:(NSRect)frame;
-(void)showProgress:(char*)left right:(char*)right percent:(double)percent;
-(void)setProgressView:(ProgressView*)pv;
-(void)showText:(SubItem*)items length:(int)length x:(double)x y:(double)y;
-(void)setTextView:(TextView*)tv;
@end
