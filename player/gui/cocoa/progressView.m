#import "progressView.h"
@implementation ProgressView

- (id)initWithFrame:(NSRect)frame {
    self = [super initWithFrame:frame];
    if (self) {
        _leftString = @"00:00:00";
        _rightString = @"00:00:00";
        _percent = 0;
        _percent2 = 0;
        _speedString = @"";
        _paddingLeft = 0;
        _titleString = @"";
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
    CGFloat knotHeight = 6;
    CGFloat knotWidth = 6;

    CGFloat progressHeight = 22;

    NSDictionary *attr = @{NSFontAttributeName : [NSFont fontWithName:@"Helvetica Neue" size:11]};
    NSDictionary *attrLarge = @{NSFontAttributeName : [NSFont fontWithName:@"Helvetica Neue" size:13]};

    CGFloat stringWidth = 60;

    if ([_speedString length] > 0) {
        _paddingLeft = 53;
    } else {
        _paddingLeft = 0;
    }

    CGFloat position = (dirtyRect.size.width-2*stringWidth-_paddingLeft)*(_percent);
    CGFloat position2 = (dirtyRect.size.width-2*stringWidth-_paddingLeft)*(_percent2);
    
    NSSize textSize = [_leftString sizeWithAttributes:attr];
    CGFloat textY = (progressHeight-13)/2;
    [_leftString drawAtPoint:NSMakePoint(stringWidth-4-textSize.width+_paddingLeft,textY) withAttributes:attr];
    [_rightString drawAtPoint:NSMakePoint(dirtyRect.size.width-stringWidth+4, textY) withAttributes:attr];

    NSSize titlesz = [_titleString sizeWithAttributes:attrLarge];
    [_titleString drawAtPoint:NSMakePoint((dirtyRect.size.width-titlesz.width)/2, textY+knotHeight+3) withAttributes:attrLarge];

    if (_paddingLeft > 0) {
        NSSize sz = [_speedString sizeWithAttributes:attr];
        [_speedString drawAtPoint:NSMakePoint(_paddingLeft+stringWidth-4-textSize.width-sz.width-8, textY) withAttributes:attr];
    }
    
    [[NSColor colorWithCalibratedRed:0 green:0 blue:0 alpha:0.3] set];
    [self drawRoundedRect:NSMakeRect(stringWidth+_paddingLeft, (progressHeight-barHeight)/2, dirtyRect.size.width-2*stringWidth-_paddingLeft, barHeight) radius:2];
    [self drawRoundedRect:NSMakeRect(stringWidth+_paddingLeft, (progressHeight-barHeight)/2, position2, barHeight) radius:2];

    
    [[NSColor whiteColor] setFill];    
    [self drawRoundedRect:NSMakeRect(stringWidth+_paddingLeft, (progressHeight-barHeight)/2, position, barHeight) radius:2];
    [self drawRoundedRect:NSMakeRect(position-knotWidth/2+stringWidth+_paddingLeft, (progressHeight-knotHeight)/2, knotWidth, knotHeight) radius:5];
    
    [super drawRect:dirtyRect];
}
- (void)mouseDown:(NSEvent *)event {
    if (_leftString == _rightString) {
        return;
    }

    CGFloat stringWidth = 60;
    NSPoint pt = [self convertPoint:[event locationInWindow] fromView:nil];
    NSRect bound = NSMakeRect(stringWidth+_paddingLeft, 4, self.frame.size.width-2*stringWidth-_paddingLeft, 22-8);
    
    if (NSPointInRect(pt, bound)) {
        _percent = (pt.x-bound.origin.x)/bound.size.width;            
        // if ((percent2>0) && (percent > percent2)) {
        //     percent = percent2;
        // }
        [self setNeedsDisplay:YES];
        
        onPlaybackChange(0, _percent);
        onPlaybackChange(1, _percent);
        double lastPercent = _percent;
            
        bool keepOn = YES;
            
        while (keepOn) {
            event = [[self window] nextEventMatchingMask: NSLeftMouseUpMask |
                            NSLeftMouseDraggedMask];
                
            switch ([event type]) {
                case NSLeftMouseDragged:
                    _percent = [self getPercent:event bound:bound];

                    if (lastPercent != _percent) {
                        lastPercent = _percent;
                        [self setNeedsDisplay:YES];
                        onPlaybackChange(1, _percent);
                    }
                    break;
                case NSLeftMouseUp:
                    _percent = [self getPercent:event bound:bound];

                    if (lastPercent != _percent) {
                        [self setNeedsDisplay:YES];
                    }
                    onPlaybackChange(2, _percent);
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
- (void)updatePorgressInfo:(NSString*)leftString rightString:(NSString*)rightString percent:(CGFloat)percent {
    _leftString = leftString;
    _rightString = rightString;
    _percent = percent;
    [self setNeedsDisplay:YES];
}
-(void)updateBufferInfo:(NSString*)speed bufferPercent:(CGFloat)percent {
    _speedString = speed;
    _percent2 = percent;
    [self setNeedsDisplay:YES];
}
@end
