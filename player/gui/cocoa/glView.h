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

        NSSize originalSize;

        ProgressView* progressView;
        BlurView* bvProgressView;

        SpinningView* spinningView;

        VolumeView* volumeView;
        BlurView* bvVolumeView;

        CGFloat _collisionOffsets[10];
        CGFloat _collisionDeltas[10];
        
        double _fontSize;
}

@property (nonatomic, strong, retain) NSDate* showCursorDeadline;

-(id)initWithFrame2:(NSRect)frame;
-(void)updatePorgressInfo:(NSString*)leftStr rightString:(NSString*)rightStr percent:(CGFloat)p;
-(void)updateBufferInfo:(NSString*)speed bufferPercent:(CGFloat)p;

-(TextView*)showSubtitle:(SubItem*)item;
-(void)hideSubtitle:(TextView*)tv;
-(void)refreshTexts;
// -(void)hideAllTexts;

-(void)setPlaybackViewHidden:(BOOL)b;
-(void)setCursorHidden:(BOOL)b;
-(void)setOriginalSize:(NSSize)size;
-(void)setSpinningHidden:(BOOL)b;
-(void)setVolume:(int)volume;
-(void)setVolumeHidden:(BOOL)b;
-(BOOL)isCursorHidden;
-(void)makeCurrentContext;
-(void)flushBuffer;
-(void)setFontSize:(double)sz;
-(void)timerTick:(NSEvent *)event;
@end
