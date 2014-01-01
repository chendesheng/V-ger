#import "progressView.h"
@implementation ProgressView

- (id)initWithFrame:(NSRect)frame {
    NSLog(@"progressView initWithFrame");
    self = [super initWithFrame:frame];
    if (self) {
        self->leftString = @"--:--:--";
        self->rightString = @"--:--:--";
        self->percent = 0;
        self->percent2 = 0;
    }
    
    return self;
}
-(BOOL)mouseDownCanMoveWindow {
    return NO;
}

-(void)drawRoundedRect:(NSRect)rect radius:(CGFloat)r {
    NSBezierPath *textViewSurround = [NSBezierPath bezierPathWithRoundedRect:rect xRadius:r yRadius:r];
    [textViewSurround fill];
}
-(void)drawRect:(NSRect)dirtyRect {
    CGFloat position = (dirtyRect.size.width-120)*(self->percent);
    CGFloat position2 = (dirtyRect.size.width-120)*(self->percent2);
    // CGFloat position3 = position2;
    // if (position3 > 65) {
    //     position3 -= 5;
    // }
    CGFloat barHeight = 4;
    CGFloat knotHeight = 14;
    CGFloat knotWidth = 5;
    
    [[NSColor colorWithCalibratedRed:1 green:1 blue:1 alpha:0.3] setFill];
    NSRectFill(dirtyRect);
    
    NSDictionary *attr = @{NSFontAttributeName : [NSFont fontWithName:@"Helvetica Neue" size:12]};
    NSSize textSize = [self->leftString sizeWithAttributes:attr];
    CGFloat textY = (dirtyRect.size.height-14)/2;
    [self->leftString drawAtPoint:NSMakePoint(60-4-textSize.width,textY) withAttributes:attr];
    [self->rightString drawAtPoint:NSMakePoint(dirtyRect.size.width-60+4, textY) withAttributes:@{NSFontAttributeName : [NSFont fontWithName:@"Helvetica Neue" size:12]}];
    
    [[NSColor colorWithCalibratedRed:0 green:0 blue:0 alpha:0.3] set];
    [self drawRoundedRect:NSMakeRect(60, (dirtyRect.size.height-barHeight)/2, dirtyRect.size.width-120, barHeight) radius:2];

    [[NSColor colorWithCalibratedRed:0 green:0 blue:0 alpha:0.3] set];
    [self drawRoundedRect:NSMakeRect(60, (dirtyRect.size.height-barHeight)/2, position2, barHeight) radius:2];
    
    // NSShadow* theShadow = [[NSShadow alloc] init];
    // [theShadow setShadowOffset:NSMakeSize(0, 0)];
    // [theShadow setShadowBlurRadius:1.0];
    
    // Use a partially transparent color for shapes that overlap.
    // [theShadow setShadowColor:[[NSColor blackColor] colorWithAlphaComponent:0.5]];
    // [theShadow setShadowColor:nil];
    // [theShadow set];
    
    [[NSColor colorWithCalibratedRed:1 green:1 blue:1 alpha:1] setFill];
    
    [self drawRoundedRect:NSMakeRect(60, (dirtyRect.size.height-barHeight)/2, position, barHeight) radius:2];
    
    [[NSColor colorWithCalibratedRed:1 green:1 blue:1 alpha:1] setFill];
    [self drawRoundedRect:NSMakeRect(position-knotWidth/2+60, (dirtyRect.size.height-knotHeight)/2, knotWidth, knotHeight) radius:1.5];
    
    [super drawRect:dirtyRect];
}
- (void)mouseDown:(NSEvent *)event {
    NSPoint pt = [self convertPoint:[event locationInWindow] fromView:nil];
    NSRect bound = NSMakeRect(60, 10, self.frame.size.width-120, self.frame.size.height-20);
    
    if (NSPointInRect(pt, bound)) {
        self->percent = (pt.x-bound.origin.x)/bound.size.width;
        [self setNeedsDisplay:YES];
            
        // self->window->callbacks.trackPositionChanged((GLFWwindow*)self->window, self->percent, 0);
        onProgressChanged((void*)[self window], 0, self->percent);
            
        bool keepOn = YES;
            
        while (keepOn) {
            event = [[self window] nextEventMatchingMask: NSLeftMouseUpMask |
                            NSLeftMouseDraggedMask];
                
            switch ([event type]) {
                case NSLeftMouseDragged:
                    pt = [self convertPoint:[event locationInWindow] fromView:nil];
                    if (pt.x < bound.origin.x) {
                        pt.x = bound.origin.x;
                    } else if (pt.x > bound.origin.x+bound.size.width) {
                        pt.x = bound.origin.x+bound.size.width;
                    }
                    self->percent = (pt.x-bound.origin.x)/bound.size.width;
                    if ((self->percent2>0) && (self->percent > self->percent2)) {
                        self->percent = self->percent2;
                    }
                    [self setNeedsDisplay:YES];
                    onProgressChanged((void*)[self window], 1, self->percent);
                    break;
                case NSLeftMouseUp:
                    [self setNeedsDisplay:YES];
                    onProgressChanged((void*)[self window], 2, self->percent);
                    keepOn = NO;
                    break;
                default:
                    /* Ignore any other kind of event. */
                    break;
            }
        }
    }
}
- (void)setHidden:(BOOL)b {
    NSView *v = self.superview;
    if (b) {
        [v setFrameSize:NSMakeSize(v.frame.size.width, 0)];
    } else {
        [v setFrameSize:NSMakeSize(v.frame.size.width, 30)];
    }
    [v setHidden:b];
}
- (void)mouseDragged:(NSEvent *)event{}
- (void)mouseUp:(NSEvent *)event{}
- (void)mouseMoved:(NSEvent *)event{}
@end