#import "glView.h"
#import <OpenGL/gl.h>
#import <objc/runtime.h>

@implementation GLView : NSOpenGLView

- (id)initWithFrame2:(NSRect)frame {
    self = [super init];
    if (self) {
        originalSize = frame.size;
        // Initialization code here.
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

        textView = [[TextView alloc] init];
        [self addSubview:textView];

        textView2 = [[TextView alloc] init];
        [self addSubview:textView2];

        [self updateTrackingAreas];
    }
    
    return self;
}

- (void)prepareOpenGL{
    // [[self openGLContext] makeCurrentContext];
    GLint swapInt = 1;
    [[self openGLContext] setValues:&swapInt forParameter:NSOpenGLCPSwapInterval];
}

-(void)dealloc {
    [trackingArea release];
    [textView release];
    [textView2 release];
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
        [win toggleFullScreen:nil];
        // return;
    }

    [self hideCursor];
    [self hideProgress];
}

-(BOOL)mouseDownCanMoveWindow {
    return YES;
}
-(void)hideCursor {
    currentCursor = noneCursor;
}
-(void)showCursor {
    currentCursor = [NSCursor arrowCursor];
}
-(void)hideProgress {
    [titleView setHidden:YES];
    [progressView setHidden:YES];

    NSView* target = [self superview];
    for (NSView* v in [frameView subviews]) {
        if (v != target) {
            [v setHidden:YES];
        }
    }
}
-(void)showProgress {
    [progressView setHidden:NO];
    [self showTitle];
}

-(void)showTitle {
    [titleView setHidden:NO];

    NSView* target = [self superview];
    for (NSView* v in [frameView subviews]) {
        if (v != target) {
            [v setHidden:NO];
        }
    }    
}
- (void)mouseMoved:(NSEvent *)event {
     NSPoint mouse = [NSEvent mouseLocation];
    if ([NSWindow windowNumberAtPoint:mouse belowWindowWithWindowNumber:0] == [self window].windowNumber) {
        onMouseMove();
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

- (void)keyUp:(NSEvent *)event {
}

-(void)showProgress:(char*)left right:(char*)right percent:(double)percent {
    ProgressView* pv = progressView;
    
    [pv->leftString autorelease];

    if (strlen(left) == 0) {
        pv->leftString = @"00:00:00";
    } else {
        pv->leftString = [[NSString stringWithUTF8String:left] retain];
    }

    [pv->rightString autorelease];

    if (strlen(right) == 0) {
        pv->rightString = @"00:00:00";
    } else {
        pv->rightString = [[NSString stringWithUTF8String:right] retain];
    }

    pv->percent = percent;
    
    [pv setNeedsDisplay:YES];
}
-(void)showBufferInfo:(char*)speed bufferPercent:(double)percent {
    ProgressView* pv = progressView;

    [pv->speedString autorelease];
    pv->speedString = [[NSString stringWithUTF8String:speed] retain];
    pv->percent2 = percent;

    [pv setNeedsDisplay:YES];
}
-(void)setProgressView:(ProgressView*)pv {
    progressView = pv;
}

-(TextView*)showText:(SubItem*)item {
    int align = item->align;

    TextView* tv = nil;
    if (item->x < 0 && item->y < 0) {
        if (align == 2) {
            tv = textView;
        } else if (align == 10) {
            tv = textView2;
        }
    }

    if (tv == nil) {
        tv = [[TextView alloc] init];
        [self addSubview:tv positioned:NSWindowBelow relativeTo:progressView];
    }

    [tv setText:item->texts length:item->length];
    tv->x = item->x;
    tv->y = item->y;
    tv->align = item->align;
    [self updateTextViewPosition:tv];
    [tv setHidden:NO];
    return tv;
}
-(void)updateTextViewPosition:(TextView*)tv {
    int align = tv->align;

    BOOL secondSub = NO;
    if (align == 10) {
        align = 2;
        secondSub = YES;
    }

    int xalign = (align-1)%3;   //0-left, 1-center, 2-right
    int yalign = (align-1)/3;   //0-bottom, 1-middle, 2-top

    NSSize wsz = [[self window] frame].size;
        
    CGFloat PADDING = 30;
    CGFloat GAP = 5.0;  //5 pixes gap between first subtitle and second subtitle

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

        //handle default position
        if (align == 2) {
            if (secondSub && [textView isHidden] == NO) {
                //current is second sub and first sub is visible
                NSSize sz1 = [textView sizeForWidth:(wsz.width-2*PADDING) height:FLT_MAX];
                if (sz1.height > 0) {
                    y += sz1.height + GAP;
                }
            } else if ([textView2 isHidden] == NO && sz.height > 0) {
                //current is first sub and second sub is visible
                NSRect rt = [textView2 frame];
                rt.origin.y += sz.height + GAP;
                [textView2 setFrame:rt];
            }
        }
    }

    //NSLog(@"%f %f %f %f %f %f", x, y, sz.width, sz.height, x-0.5*xalign*sz.width, y-0.5*yalign*sz.height);
    [tv setFrame:NSMakeRect(x-0.5*xalign*sz.width, y-0.5*yalign*sz.height, sz.width, sz.height)];
}
-(void)hideText:(TextView*)tv {
    if (tv == textView) {
        [tv setText:NULL length:0];
        [tv setHidden:YES];
    } else if (tv == textView2) {
        [tv setText:NULL length:0];
        [tv setHidden:YES];
    } else {
        [tv removeFromSuperview];
        [tv release];
    }
}

- (void)cursorUpdate:(NSEvent *)event {
    NSCursor* cur = currentCursor;
    [cur set];
}

- (void)drawRect:(NSRect)dirtyRect {
    onDraw((void*)[self window]);
    [[self openGLContext] flushBuffer];
}
-(void)setStartupView:(StartupView*)sv {
    startupView = sv;
}
-(void)hideStartupView {
    [startupView setHidden:YES];
}
-(void)showStartupView {
    [startupView setHidden:NO];

    [self setNeedsDisplay:YES];
}
- (void)scrollWheel:(NSEvent *)event
{
    onMouseWheel([event deltaY]);
}
- (void)setOriginalSize:(NSSize)size {
    originalSize = size;
}
- (void)showAllTexts {
    for (NSView* v in [self subviews]) {
        if ([v isKindOfClass:[textView class]] && [v isHidden]==NO) {
        // NSLog(@"begin update %d", [(TextView*)v getSubItem]->align);
        [self updateTextViewPosition:(TextView*)v];
        // NSLog(@"end update %d", [(TextView*)v getSubItem]->align);

            // [v setHidden:NO];
        }
    }
}
// - (void)hideAllTexts {
//     for (NSView* v in [self subviews]) {
//         if ([v isKindOfClass:[TextView class]]) {
//             [v setHidden:YES];
//         }
//     }
// }
@end