#import <Cocoa/Cocoa.h>
#import "gui.h"
#import "progressView.h"
#import "textView.h"
#import "startupView.h"
#import "blurView.h"
#import "titleTextView.h"
#import "spinningView.h"
#import "volumeView.h"

@interface GLView : NSOpenGLView {
    NSTrackingArea* trackingArea;
    NSCursor* noneCursor;
    TextView* textView;
    TextView* textView2;
    NSCursor* currentCursor;
    StartupView* startupView;
    NSSize originalSize;
@public
    BlurView* blurView;
    ProgressView* progressView;
    BlurView* titleView;
    TitleTextView* titleTextView;
    NSView* frameView;
    NSWindow* win;
    SpinningView* spinningView;
    VolumeView* volumeView;
    BlurView* volumeView2;
}

-(id)initWithFrame2:(NSRect)frame;
-(void)showProgress:(char*)left right:(char*)right percent:(double)percent;
-(void)setProgressView:(ProgressView*)pv;

-(TextView*)showText:(SubItem*)item;
-(void)hideText:(TextView*)tv;
-(void)showAllTexts;
// -(void)hideAllTexts;

-(void)setStartupView:(StartupView*)sv;
-(void)hideStartupView;
-(void)showStartupView;

-(void)showProgress;
-(void)hideProgress;
-(void)showBufferInfo:(char*)speed bufferPercent:(double)percent;
-(void)showCursor;
-(void)hideCursor;

-(void)setOriginalSize:(NSSize)size;
@end