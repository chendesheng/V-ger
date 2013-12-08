#import "glView.h"
#import <OpenGL/gl.h>

@implementation GLView : NSOpenGLView

- (id)initWithFrame2:(NSRect)frame {
    NSLog(@"glView initWithFrame2");
    self = [super init];
    if (self) {
        self->originalSize = frame.size;
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
-(TextView*)showText:(SubItem*)items length:(int)length position:(int)position x:(double)x y:(double)y {
    double padding = 35;

    if (position == 0 || position == 2) {
        TextView* tv = self->textView;

        [tv setText:items length:length];

        double w = [self frame].size.width;
        NSSize sz = [tv sizeForWidth:w height:FLT_MAX];
        [tv setFrame:NSMakeRect(0, padding, w, sz.height)];

        return tv;
    } else if (position == 1) {
        TextView* tv = [[TextView alloc] init];
        [self addSubview:tv positioned:NSWindowBelow relativeTo:self->progressView];
        [tv setText:items length:length];

        NSSize sz = [tv sizeForWidth:FLT_MAX height:FLT_MAX];

        [tv setFrame:NSMakeRect(padding, padding, sz.width, sz.height)];

        return tv;
    } else if (position == 3) {
        TextView* tv = [[TextView alloc] init];
        [self addSubview:tv positioned:NSWindowBelow relativeTo:self->progressView];
        [tv setText:items length:length];

        NSSize wsz = [self frame].size;
        NSSize sz = [tv sizeForWidth:FLT_MAX height:FLT_MAX];

        double x = wsz.width - sz.width - padding;
        [tv setFrame:NSMakeRect(x, padding, sz.width, sz.height)];        
    } else if (position == 4) {   
        TextView* tv = [[TextView alloc] init];
        [self addSubview:tv positioned:NSWindowBelow relativeTo:self->progressView];
        [tv setText:items length:length];

        NSSize wsz = [self frame].size;
        NSSize sz = [tv sizeForWidth:FLT_MAX height:FLT_MAX];

        double y = (wsz.height - sz.height)/2;
        [tv setFrame:NSMakeRect(padding, y, sz.width, sz.height)];  
    } else if (position == 5) {
        TextView* tv = [[TextView alloc] init];
        [self addSubview:tv positioned:NSWindowBelow relativeTo:self->progressView];
        [tv setText:items length:length];

        NSSize wsz = [self frame].size;
        NSSize sz = [tv sizeForWidth:FLT_MAX height:FLT_MAX];

        double x = (wsz.width - sz.width)/2;
        double y = (wsz.height - sz.height)/2;
        [tv setFrame:NSMakeRect(x, y, sz.width, sz.height)];        
    } else if (position == 6) {
        TextView* tv = [[TextView alloc] init];
        [self addSubview:tv positioned:NSWindowBelow relativeTo:self->progressView];
        [tv setText:items length:length];

        NSSize wsz = [self frame].size;
        NSSize sz = [tv sizeForWidth:FLT_MAX height:FLT_MAX];

        double x = wsz.width - sz.width - padding;
        double y = (wsz.height - sz.height)/2;
        [tv setFrame:NSMakeRect(x, y, sz.width, sz.height)];        
    } else if (position == 7) {
        TextView* tv = [[TextView alloc] init];
        [self addSubview:tv positioned:NSWindowBelow relativeTo:self->progressView];
        [tv setText:items length:length];

        NSSize wsz = [self frame].size;
        NSSize sz = [tv sizeForWidth:FLT_MAX height:FLT_MAX];

        double x = padding;
        double y = wsz.height - sz.height - padding;
        [tv setFrame:NSMakeRect(x, y, sz.width, sz.height)];        
    } else if (position == 8) {
        TextView* tv = [[TextView alloc] init];
        [self addSubview:tv positioned:NSWindowBelow relativeTo:self->progressView];
        [tv setText:items length:length];

        NSSize wsz = [self frame].size;
        NSSize sz = [tv sizeForWidth:FLT_MAX height:FLT_MAX];

        double x = (wsz.width - sz.width)/2;
        double y = wsz.height - sz.height - padding;
        [tv setFrame:NSMakeRect(x, y, sz.width, sz.height)];        
    } else if (position == 9) {
        TextView* tv = [[TextView alloc] init];
        [self addSubview:tv positioned:NSWindowBelow relativeTo:self->progressView];
        [tv setText:items length:length];

        NSSize wsz = [self frame].size;
        NSSize sz = [tv sizeForWidth:FLT_MAX height:FLT_MAX];

        double x = wsz.width - sz.width - padding;
        double y = wsz.height - sz.height - padding;
        [tv setFrame:NSMakeRect(x, y, sz.width, sz.height)];        
    } else {
        TextView* tv = [[TextView alloc] init];
        [self addSubview:tv positioned:NSWindowBelow relativeTo:self->progressView];
        [tv setText:items length:length];

        NSSize wsz = [self frame].size;
        NSSize sz = [tv sizeForWidth:FLT_MAX height:FLT_MAX];

        x = x/self->originalSize.width * wsz.width;
        y = wsz.height - y/self->originalSize.height * wsz.height - sz.height;
        [tv setFrame:NSMakeRect(x, y, sz.width, sz.height)];

        return tv;
    }
}
-(void)hideText:(TextView*)tv {
    if (tv == self->textView) {
        [tv setText:NULL length:0];
    } else {
        [tv removeFromSuperview];
        [tv release];
    }
}

- (void)cursorUpdate:(NSEvent *)event {
    // setModeCursor(window, window->cursorMode);
    NSCursor* cur = self->currentCursor;
    [cur set];

    if (cur == [NSCursor arrowCursor]) {
        [self->progressView setHidden:NO];
        
        // NSRect frame = [self->textView frame];
        // frame.origin.y = 60;
        // [self->textView setFrame:frame];
        // [self->textView setNeedsDisplay:YES];
    }
    else {
        [self->progressView setHidden:YES];
        
        // NSRect frame = [self->textView frame];
        // frame.origin.y = 20;
        // [self->textView setFrame:frame];
        // [self->textView setNeedsDisplay:YES];
    }
}

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