#import "glView.h"
#import <OpenGL/gl.h>
#import <objc/runtime.h>

@implementation GLView : NSOpenGLView

- (id)initWithFrame2:(NSRect)frame {
        self = [super init];
        if (self) {
                originalSize = frame.size;

                trackingArea = nil;

                NSOpenGLPixelFormatAttribute attrs[] = {

                        // Specifying "NoRecovery" gives us a context that cannot fall back to the software renderer.  This makes the View-based context a compatible with the layer-backed context, enabling us to use the "shareContext" feature to share textures, display lists, and other OpenGL objects between the two.
                        NSOpenGLPFANoRecovery, // Enable automatic use of OpenGL "share" contexts.

                        NSOpenGLPFAColorSize, 24,
                        NSOpenGLPFAAlphaSize, 8,
                        NSOpenGLPFADepthSize, 16,
                        NSOpenGLPFADoubleBuffer,
                        NSOpenGLPFAAccelerated,
                        0
                };
                NSOpenGLPixelFormat* pixelFormat = [[NSOpenGLPixelFormat alloc] initWithAttributes:attrs];
                [super initWithFrame:frame pixelFormat:pixelFormat];
                [pixelFormat release];

                [self setWantsLayer:YES];

                NSImage* data = [[NSImage alloc] initWithSize:NSMakeSize(1, 1)];
                noneCursor = [[NSCursor alloc] initWithImage:data
                                                     hotSpot:NSZeroPoint];
                [data release];

                [self updateTrackingAreas];

                CGFloat width = frame.size.width;
                CGFloat height = frame.size.height;

                progressView = [[ProgressView alloc] init];
                bvProgressView = [[BlurView alloc] initWithView:progressView frame:NSMakeRect(0,0,width,22)];
                [bvProgressView setAutoresizingMask:NSViewWidthSizable|NSViewMaxYMargin];
                [self addSubview:bvProgressView positioned:NSWindowAbove relativeTo:nil];

                spinningView = [[SpinningView alloc] initWithFrame:NSMakeRect((width-50)/2, (height-50)/2, 50, 50)];
                [spinningView setAutoresizingMask:NSViewMinXMargin|NSViewMaxXMargin|NSViewMinYMargin|NSViewMaxYMargin];
                [self addSubview:spinningView positioned:NSWindowAbove relativeTo:nil];
                [spinningView setHidden:YES];

                volumeView = [[VolumeView alloc] init];
                bvVolumeView = [[BlurView alloc] initWithView:volumeView frame:NSMakeRect((width-120)/2, (height-120)/2, 120, 120)];
                [bvVolumeView setAutoresizingMask:NSViewMinXMargin|NSViewMaxXMargin|NSViewMinYMargin|NSViewMaxYMargin];
                [bvVolumeView setCornerRadius:4.1];
                [self addSubview:bvVolumeView positioned:NSWindowAbove relativeTo:spinningView];
                [bvVolumeView setHidden:YES];

                self.showCursorDeadline = [NSDate distantFuture];
                self->_fontSize = 25.0;
        }

        return self;
}

- (void)prepareOpenGL{
        GLint swapInt = 1;
        [[self openGLContext] setValues:&swapInt forParameter:NSOpenGLCPSwapInterval];
}
- (void)makeCurrentContext {
        [[self openGLContext] makeCurrentContext];
}
- (void)flushBuffer {
        [[self openGLContext] flushBuffer];
}
-(void)dealloc {
        [trackingArea release];
        [super dealloc];
}

- (BOOL)isOpaque {
        return YES;
}

- (BOOL)canBecomeKeyView {
        return YES;
}

- (BOOL)acceptsFirstResponder {
        return YES;
}

- (void)mouseDown:(NSEvent *)event {
        if (event.clickCount == 2) {
                toggleFullScreen(self.window);
        }

        if ([self.showCursorDeadline compare:[NSDate distantFuture]] == NSOrderedAscending) {
                setControlsVisible(self.window, 0, 0);
        }
}

-(BOOL)mouseDownCanMoveWindow {
        return YES;
}
-(void)setCursorHidden:(BOOL)b {
        if (b) {
                currentCursor = noneCursor;
        } else {
                currentCursor = [NSCursor arrowCursor];
        }
}
-(BOOL)isCursorHidden {
        return currentCursor == noneCursor;
}
-(void)setPlaybackViewHidden:(BOOL)b {
     [progressView setHidden:b];
}

- (void)mouseMoved:(NSEvent *)event {
        NSPoint mouse = [NSEvent mouseLocation];
        if ([NSWindow windowNumberAtPoint:mouse belowWindowWithWindowNumber:0] == [self window].windowNumber
      && [self.showCursorDeadline compare:[NSDate distantFuture]] == NSOrderedAscending) {
                setControlsVisible(self.window, 1, 1);
        }
}

- (void)timerTick:(NSEvent *)event {
        @autoreleasepool {
                onTimerTick();
                if ([self.showCursorDeadline compare:[NSDate date]] == NSOrderedAscending) {
                        setControlsVisible(self.window, 0, 0);
                }
        }
}

- (void)updateTrackingAreas {
        if (trackingArea != nil) {
                [self removeTrackingArea:trackingArea];
                [trackingArea release];
        }

        NSTrackingAreaOptions options = NSTrackingMouseMoved |
                NSTrackingActiveInKeyWindow |
                NSTrackingCursorUpdate |
                NSTrackingInVisibleRect;

        trackingArea = [[NSTrackingArea alloc] initWithRect:[self bounds]
                                                    options:options
                                                      owner:self
                                                   userInfo:nil];

        [self addTrackingArea:trackingArea];
        [super updateTrackingAreas];
}

- (void)keyDown:(NSEvent *)event {
        if (!onKeyDown([event keyCode])) {
                [super keyDown:event];
        }
}

- (void)updatePorgressInfo:(NSString*)leftString rightString:(NSString*)rightString percent:(CGFloat)percent {
    [progressView updatePorgressInfo:leftString rightString:rightString percent:percent];
}
-(void)updateBufferInfo:(NSString*)speed bufferPercent:(CGFloat)percent {
    [progressView updateBufferInfo:speed bufferPercent:percent];
}

-(TextView*)showSubtitle:(SubItem*)item {
        TextView *tv = [self getOrCreateTextView];
        [tv setText:item->texts length:item->length];
        tv->x = item->x;
        tv->y = item->y;
        tv->align = item->align;
        [tv setHidden:NO];
        [self refreshTexts];
        return tv;
}
// get an unused textView or create one
-(TextView*)getOrCreateTextView {
        for (NSView* v in [self subviews]) {
                if ([NSStringFromClass([v class]) isEqualToString:@"TextView"] && [v isHidden] == YES) {
                        return (TextView*)v;
                }
        }

        TextView *tv = [[TextView alloc] initWithFrameAndSize:NSMakeRect(0,0,0,0) fontSize:self->_fontSize];
        [tv setVerticallyResizable:NO]; // this is required or setFrame won't do right
        [self addSubview:tv positioned:NSWindowBelow relativeTo:progressView];
        [tv release];
        return tv;
}
-(void)updateTextViewPosition:(TextView*)tv {
        int align = tv->align;

        int xalign = (align-1)%3;   //0-left, 1-center, 2-right
        int yalign = (align-1)/3;   //0-bottom, 1-middle, 2-top

        NSSize wsz = [[self window] frame].size;

        CGFloat PADDING = 25;
        CGFloat GAP = 5.0;  //5 pixes space between collisioned texts

        CGFloat x;
        CGFloat y;

        NSSize sz = [tv sizeForWidth:(wsz.width-2*PADDING) height:FLT_MAX];

        if (tv->x >= 0 && tv->x >= 0) {
                x = tv->x;
                y = tv->y;
                // get x y from subtitle file, need scale to current view
                x = x/originalSize.width * wsz.width;
                y = wsz.height - y/originalSize.height * wsz.height;
        } else {
                //no x y, set position by align
                x = 0.5*xalign*wsz.width + (1-xalign)*PADDING;
                y = 0.5*yalign*wsz.height + (1-yalign)*PADDING;

                // handle position collision
                y += _collisionOffsets[align];
                _collisionOffsets[align] += (sz.height + GAP) * _collisionDeltas[align];
        }

        // NSLog(@"set:%f,%f,%f,%f",x-0.5*xalign*sz.width, y-0.5*yalign*sz.height, sz.width, sz.height);
        [tv setFrame:NSMakeRect(x-0.5*xalign*sz.width, y-0.5*yalign*sz.height, sz.width, sz.height)];
        // NSLog(@"after:%f %f %f %f", tv.frame.origin.x, tv.frame.origin.y, tv.frame.size.width, tv.frame.size.height);
}
-(void)hideSubtitle:(TextView*)tv {
        [tv setText:NULL length:0];
        [tv setHidden:YES];

        [self refreshTexts];
}

- (void)cursorUpdate:(NSEvent *)event {
        NSCursor* cur = currentCursor;
        [cur set];
}

- (void)scrollWheel:(NSEvent *)event
{
        onMouseWheel([event deltaY]);
}
- (void)setOriginalSize:(NSSize)size {
        originalSize = size;
}
- (void)refreshTexts {

        for (int i = 0; i < 10; i++) {
                self->_collisionOffsets[i] = 0;
                if (i < 7) self->_collisionDeltas[i] = 1;
                else self->_collisionDeltas[i] = -1;
        }

        // NSLog(@"subview length: %ld", [self subviews].count);
        for (NSView* v in [self subviews]) {
                if ([NSStringFromClass([v class]) isEqualToString:@"TextView"] && [v isHidden]==NO) {
                        // NSLog(@"begin update %@", v);
                        [self updateTextViewPosition:(TextView*)v ];
                        // NSLog(@"end update %@", v);

                        // [v setHidden:NO];
                }
        }
}
- (void)setFontSize:(double)sz {
        self->_fontSize = sz;

        for (NSView* v in [self subviews]) {
                if ([NSStringFromClass([v class]) isEqualToString:@"TextView"]) {
                        [(TextView*)v setFontSize:sz];
                }
        }

        [self refreshTexts];
}
- (void)setSpinningHidden:(BOOL)b {
  [spinningView setHidden:b];
}
- (void)setVolume:(int)volume {
    [volumeView setVolume:volume];
}
- (void)setVolumeHidden:(BOOL)b {
        BlurView* bv = bvVolumeView;
        [bv setHidden:b];
        if (!b) {
                NSSize sz = self.frame.size;
                [bv setFrame:NSMakeRect((sz.width-120)/2, (sz.height-120)/2, 120, 120)];
                [volumeView setNeedsDisplay:YES];
        }
}
@end

