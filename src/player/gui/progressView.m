#import "progressView.h"
@implementation ProgressView

- (id)initWithFrame:(NSRect)frame {
    NSLog(@"progressView initWithFrame");
    self = [super initWithFrame:frame];
    if (self) {
        self->leftString = @"00:00:00";
        self->rightString = @"00:00:00";
        self->percent = 0;
        self->percent2 = 0;
        self->speedString = @"";
        self->paddingLeft = 0;
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

    if ([self->speedString length] > 0) {
        self->paddingLeft = 60;
    } else {
        self->paddingLeft = 0;
    }

    CGFloat position = (dirtyRect.size.width-120-self->paddingLeft)*(self->percent);
    CGFloat position2 = (dirtyRect.size.width-120-self->paddingLeft)*(self->percent2);
    
    NSSize textSize = [self->leftString sizeWithAttributes:attr];
    CGFloat textY = (dirtyRect.size.height-14)/2;
    [self->leftString drawAtPoint:NSMakePoint(60-4-textSize.width+self->paddingLeft,textY) withAttributes:attr];
    [self->rightString drawAtPoint:NSMakePoint(dirtyRect.size.width-60+4, textY) withAttributes:attr];

    if (self->paddingLeft > 0) {
        NSSize sz = [self->speedString sizeWithAttributes:attr];
        [self->speedString drawAtPoint:NSMakePoint(self->paddingLeft+60-4-textSize.width-sz.width-10, textY) withAttributes:attr];
    }
    
    [[NSColor colorWithCalibratedRed:0 green:0 blue:0 alpha:0.3] set];
    [self drawRoundedRect:NSMakeRect(60+self->paddingLeft, (dirtyRect.size.height-barHeight)/2, dirtyRect.size.width-120-self->paddingLeft, barHeight) radius:2];

    [[NSColor colorWithCalibratedRed:0 green:0 blue:0 alpha:0.3] set];
    [self drawRoundedRect:NSMakeRect(60+self->paddingLeft, (dirtyRect.size.height-barHeight)/2, position2, barHeight) radius:2];
    
    // NSShadow* theShadow = [[NSShadow alloc] init];
    // [theShadow setShadowOffset:NSMakeSize(0, 0)];
    // [theShadow setShadowBlurRadius:1.0];
    
    // Use a partially transparent color for shapes that overlap.
    // [theShadow setShadowColor:[[NSColor blackColor] colorWithAlphaComponent:0.5]];
    // [theShadow setShadowColor:nil];
    // [theShadow set];
    
    [[NSColor colorWithCalibratedRed:1 green:1 blue:1 alpha:1] setFill];
    
    [self drawRoundedRect:NSMakeRect(60+self->paddingLeft, (dirtyRect.size.height-barHeight)/2, position, barHeight) radius:2];
    
    [[NSColor colorWithCalibratedRed:1 green:1 blue:1 alpha:1] setFill];
    [self drawRoundedRect:NSMakeRect(position-knotWidth/2+60+self->paddingLeft, (dirtyRect.size.height-knotHeight)/2, knotWidth, knotHeight) radius:1.5];
    
    [super drawRect:dirtyRect];
}
- (void)mouseDown:(NSEvent *)event {
    NSPoint pt = [self convertPoint:[event locationInWindow] fromView:nil];
    NSRect bound = NSMakeRect(60+self->paddingLeft, 10, self.frame.size.width-120-self->paddingLeft, self.frame.size.height-20);
    
    if (NSPointInRect(pt, bound)) {
        self->percent = (pt.x-bound.origin.x)/bound.size.width;            
        // if ((self->percent2>0) && (self->percent > self->percent2)) {
        //     self->percent = self->percent2;
        // }
        [self setNeedsDisplay:YES];
        
        onProgressChanged((void*)[self window], 0, self->percent);
        onProgressChanged((void*)[self window], 1, self->percent);
        double lastPercent = self->percent;
            
        bool keepOn = YES;
            
        while (keepOn) {
            event = [[self window] nextEventMatchingMask: NSLeftMouseUpMask |
                            NSLeftMouseDraggedMask];
                
            switch ([event type]) {
                case NSLeftMouseDragged:
                    self->percent = [self getPercent:event bound:bound];

                    if (lastPercent != self->percent) {
                        lastPercent = self->percent;
                        [self setNeedsDisplay:YES];
                        onProgressChanged((void*)[self window], 1, self->percent);
                    }
                    break;
                case NSLeftMouseUp:
                    self->percent = [self getPercent:event bound:bound];

                    if (lastPercent != self->percent) {
                        [self setNeedsDisplay:YES];
                    }
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
- (double)getPercent:(NSEvent*)event bound:(NSRect)bound {
    NSPoint pt = [self convertPoint:[event locationInWindow] fromView:nil];
    if (pt.x < bound.origin.x) {
        pt.x = bound.origin.x;
    } else if (pt.x > bound.origin.x+bound.size.width) {
        pt.x = bound.origin.x+bound.size.width;
    }
    return (pt.x-bound.origin.x)/bound.size.width;
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