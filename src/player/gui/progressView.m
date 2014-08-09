#import "progressView.h"
@implementation ProgressView

- (id)initWithFrame:(NSRect)frame {
    self = [super initWithFrame:frame];
    if (self) {
        self->leftString = @"00:00:00";
        self->rightString = @"00:00:00";
        self->percent = 0;
        self->percent2 = 0;
        self->speedString = @"";
        self->paddingLeft = 0;
        self->titleString = @"";
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
    CGFloat barHeight = 2;
    CGFloat knotHeight = 12;
    CGFloat knotWidth = 3;

    CGFloat progressHeight = 22;

    NSDictionary *attr = @{NSFontAttributeName : [NSFont fontWithName:@"Helvetica Neue" size:11]};
    NSDictionary *attrLarge = @{NSFontAttributeName : [NSFont fontWithName:@"Helvetica Neue" size:13]};

    CGFloat stringWidth = 60;

    if ([self->speedString length] > 0) {
        self->paddingLeft = 50;
    } else {
        self->paddingLeft = 0;
    }

    CGFloat position = (dirtyRect.size.width-2*stringWidth-self->paddingLeft)*(self->percent);
    CGFloat position2 = (dirtyRect.size.width-2*stringWidth-self->paddingLeft)*(self->percent2);
    
    NSSize textSize = [self->leftString sizeWithAttributes:attr];
    CGFloat textY = (progressHeight-13)/2;
    [self->leftString drawAtPoint:NSMakePoint(stringWidth-4-textSize.width+self->paddingLeft,textY) withAttributes:attr];
    [self->rightString drawAtPoint:NSMakePoint(dirtyRect.size.width-stringWidth+4, textY) withAttributes:attr];

    NSSize titlesz = [self->titleString sizeWithAttributes:attrLarge];
    [self->titleString drawAtPoint:NSMakePoint((dirtyRect.size.width-titlesz.width)/2, textY+knotHeight+3) withAttributes:attrLarge];

    if (self->paddingLeft > 0) {
        NSSize sz = [self->speedString sizeWithAttributes:attr];
        [self->speedString drawAtPoint:NSMakePoint(self->paddingLeft+stringWidth-4-textSize.width-sz.width-10, textY) withAttributes:attr];
    }
    
    [[NSColor colorWithCalibratedRed:0 green:0 blue:0 alpha:0.3] set];
    [self drawRoundedRect:NSMakeRect(stringWidth+self->paddingLeft, (progressHeight-barHeight)/2, dirtyRect.size.width-2*stringWidth-self->paddingLeft, barHeight) radius:2];
    [self drawRoundedRect:NSMakeRect(stringWidth+self->paddingLeft, (progressHeight-barHeight)/2, position2, barHeight) radius:2];

    
    [[NSColor whiteColor] setFill];    
    [self drawRoundedRect:NSMakeRect(stringWidth+self->paddingLeft, (progressHeight-barHeight)/2, position, barHeight) radius:2];
    [self drawRoundedRect:NSMakeRect(position-knotWidth/2+stringWidth+self->paddingLeft, (progressHeight-knotHeight)/2, knotWidth, knotHeight) radius:1.5];
    
    [super drawRect:dirtyRect];
}
- (void)mouseDown:(NSEvent *)event {
    CGFloat stringWidth = 60;
    NSPoint pt = [self convertPoint:[event locationInWindow] fromView:nil];
    NSRect bound = NSMakeRect(stringWidth+self->paddingLeft, 4, self.frame.size.width-2*stringWidth-self->paddingLeft, 22-8);
    
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
    [v setHidden:b];
}
- (void)mouseDragged:(NSEvent *)event{}
- (void)mouseUp:(NSEvent *)event{}
- (void)mouseMoved:(NSEvent *)event{}
@end