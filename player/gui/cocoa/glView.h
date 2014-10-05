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
    NSCursor* currentCursor;
    
    TextView* textView;
    TextView* textView2;
    
    StartupView* startupView;
    
    NSSize originalSize;

    ProgressView* progressView;
    BlurView* bvProgressView;
    
    SpinningView* spinningView;
    
    VolumeView* volumeView;
    BlurView* bvVolumeView;
}

-(id)initWithFrame2:(NSRect)frame;
-(void)updatePorgressInfo:(NSString*)leftStr rightString:(NSString*)rightStr percent:(CGFloat)p;
-(void)updateBufferInfo:(NSString*)speed bufferPercent:(CGFloat)p;

-(TextView*)showSubtitle:(SubItem*)item;
-(void)hideSubtitle:(TextView*)tv;
-(void)showAllTexts;
// -(void)hideAllTexts;

-(void)setStartupViewHidden:(BOOL)b;
-(void)setPlaybackViewHidden:(BOOL)b;
-(void)setCursorHidden:(BOOL)b;
-(void)setOriginalSize:(NSSize)size;
- (void)setSpinningHidden:(BOOL)b;
- (void)setVolume:(int)volume;
- (void)setVolumeHidden:(BOOL)b;
@end