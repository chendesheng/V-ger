#import "glView.h"
#import <OpenGL/gl.h>

@implementation GLView : NSOpenGLView

- (id)initWithFrame2:(NSRect)frame {
    NSLog(@"glView initWithFrame2");
    self = [super init];
    if (self) {
        // Initialization code here.
        trackingArea = nil;

        [self updateTrackingAreas];

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
        self->noneCursor = [[NSCursor alloc] initWithImage:data
                                                  hotSpot:NSZeroPoint];
        [data release];
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
    NSLog(@"mouseDown");

    self->currentCursor = self->noneCursor;
    [self cursorUpdate:nil];
}

-(BOOL)mouseDownCanMoveWindow {
    return YES;
}
- (void)mouseMoved:(NSEvent *)event {
    self->currentCursor = [NSCursor arrowCursor];
    [self cursorUpdate:nil];
}

- (void)viewDidChangeBackingProperties {
    NSLog(@"viewDidChangeBackingProperties");
}

- (void)updateTrackingAreas {
    NSLog(@"updateTrackingAreas");
    if (trackingArea != nil) {
        [self removeTrackingArea:trackingArea];
        [trackingArea release];
    }

    NSTrackingAreaOptions options = NSTrackingMouseEnteredAndExited |
                                    NSTrackingActiveInKeyWindow |
                                    NSTrackingCursorUpdate |
                                    NSTrackingInVisibleRect;

    trackingArea = [[NSTrackingArea alloc] initWithRect:[self bounds]
                                                options:options
                                                  owner:self
                                               userInfo:nil];

    [self addTrackingArea:trackingArea];
    [super updateTrackingAreas];

    NSLog(@"endupdateTrackingAreas");
}

- (void)keyDown:(NSEvent *)event {
    NSLog(@"keyDown");
    onKeyDown([self window], [event keyCode]);
}

- (void)keyUp:(NSEvent *)event {
    NSLog(@"keyUp");
}

-(void)showProgress:(char*)left right:(char*)right percent:(double)percent {
    ProgressView* pv = self->progressView;
    
    [pv->leftString autorelease];
    pv->leftString = [[NSString stringWithUTF8String:left] retain];

    [pv->rightString autorelease];
    pv->rightString = [[NSString stringWithUTF8String:right] retain];

    pv->percent = percent;
    
    [pv setNeedsDisplay:YES];
}
-(void)setProgressView:(ProgressView*)pv {
    self->progressView = pv;
}
-(void)setTextView:(TextView*)tv {
    self->textView = tv;
}
-(void)showText:(SubItem*)items length:(int)length x:(double)x y:(double)y {
    TextView* tv = self->textView;

    [tv setText:items length:length];
}
- (void)cursorUpdate:(NSEvent *)event {
    // setModeCursor(window, window->cursorMode);
    NSCursor* cur = self->currentCursor;
    [cur set];

    if (cur == [NSCursor arrowCursor]) {
        [self->progressView setHidden:NO];
        
        NSRect frame = [self->textView frame];
        frame.origin.y = 60;
        [self->textView setFrame:frame];
        [self->textView setNeedsDisplay:YES];
    }
    else {
        [self->progressView setHidden:YES];
        
        NSRect frame = [self->textView frame];
        frame.origin.y = 20;
        [self->textView setFrame:frame];
        [self->textView setNeedsDisplay:YES];
    }
}
// - (void)scrollWheel:(NSEvent *)event {
//     double deltaX, deltaY;

// #if MAC_OS_X_VERSION_MAX_ALLOWED >= 1070
//     if (floor(NSAppKitVersionNumber) >= NSAppKitVersionNumber10_7)
//     {
//         deltaX = [event scrollingDeltaX];
//         deltaY = [event scrollingDeltaY];

//         if ([event hasPreciseScrollingDeltas])
//         {
//             deltaX *= 0.1;
//             deltaY *= 0.1;
//         }
//     }
//     else
// #endif /*MAC_OS_X_VERSION_MAX_ALLOWED*/
//     {
//         deltaX = [event deltaX];
//         deltaY = [event deltaY];
//     }

//     if (fabs(deltaX) > 0.0 || fabs(deltaY) > 0.0)
//         _glfwInputScroll(window, deltaX, deltaY);
// }

- (void)drawRect:(NSRect)dirtyRect {
    onDraw((void*)[self window]);
    [[self openGLContext] flushBuffer];
}
-(void)setStartupView:(StartupView*)sv {
    self->startupView = sv;
}
-(void)hideStartupView {
    [self->startupView setHidden:YES];
}
@end