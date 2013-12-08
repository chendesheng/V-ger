#import <Cocoa/Cocoa.h>
#import "gui.h"
#import "progressView.h"
#import "textView.h"
#import "startupView.h"
@interface GLView : NSOpenGLView 

{
    NSTrackingArea* trackingArea;
    NSCursor* noneCursor;
    ProgressView* progressView;
    TextView* textView;
    NSCursor* currentCursor;
    StartupView* startupView;
    NSSize originalSize;
}

-(id)initWithFrame2:(NSRect)frame;
-(void)showProgress:(char*)left right:(char*)right percent:(double)percent;
-(void)setProgressView:(ProgressView*)pv;

-(TextView*)showText:(SubItem*)items length:(int)length position:(int)position x:(double)x y:(double)y;
-(void)hideText:(TextView*)tv;

-(void)setTextView:(TextView*)tv;
-(void)setStartupView:(StartupView*)sv;
-(void)hideStartupView;
@end
