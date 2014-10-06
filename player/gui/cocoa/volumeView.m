#import "volumeView.h"

@implementation VolumeView

-(void)drawRect:(NSRect)dirtyRect {
    NSString* title = @"Volume";
    NSDictionary *attr = @{NSFontAttributeName : [NSFont fontWithName:@"Helvetica Neue" size:13]};

    NSSize textSize = [title sizeWithAttributes:attr];
    NSRect rt = self.frame;
    // CGFloat textY = (progressHeight-13)/2;
    [title drawAtPoint:NSMakePoint((rt.size.width-textSize.width)/2, rt.size.height - textSize.height - 5) withAttributes:attr];

    [[NSColor colorWithCalibratedRed:0 green:0 blue:0 alpha:0.7] set];

    NSBezierPath* aPath = [NSBezierPath bezierPath];
    [aPath moveToPoint:NSMakePoint(33.0, 52.0)];
    [aPath lineToPoint:NSMakePoint(42.0, 52.0)];
    [aPath lineToPoint:NSMakePoint(56.0, 42.0)];
    [aPath lineToPoint:NSMakePoint(56.0, 78.0)];
    [aPath lineToPoint:NSMakePoint(42.0, 68.0)];
    [aPath lineToPoint:NSMakePoint(33.0, 68.0)];
    [aPath closePath];
    [aPath fill];

    NSBezierPath* aPath1 = [NSBezierPath bezierPath];
    [aPath1 appendBezierPathWithArcWithCenter:NSMakePoint(47, 60) radius:20 startAngle:-33 endAngle:33];
    [aPath1 setLineWidth:3.0];
    [aPath1 setLineCapStyle:NSRoundLineCapStyle];
    [aPath1 stroke];

    aPath1 = [NSBezierPath bezierPath];
    [aPath1 appendBezierPathWithArcWithCenter:NSMakePoint(47, 60) radius:30 startAngle:-33 endAngle:33];
    [aPath1 setLineWidth:3.0];
    [aPath1 setLineCapStyle:NSRoundLineCapStyle];
    [aPath1 stroke];

    aPath1 = [NSBezierPath bezierPath];
    [aPath1 appendBezierPathWithArcWithCenter:NSMakePoint(47, 60) radius:40 startAngle:-33 endAngle:33];
    [aPath1 setLineWidth:3.0];
    [aPath1 setLineCapStyle:NSRoundLineCapStyle];
    [aPath1 stroke];


    int blockWidth = 6;
    int y = 12;
    int width = (rt.size.width-2*y-1)/blockWidth*blockWidth+2;
    int x = (rt.size.width-width)/2;
    NSRectFill(NSMakeRect(x, y, width, blockWidth));

    [[NSColor whiteColor] set];
    // self->_volume = 15;
    int blocks = self->_volume;
    for (int i = 0; i < blocks;  i++) {
        NSRectFill(NSMakeRect(x+blockWidth*i+1, y+1, blockWidth-1, blockWidth-2));
    }


    [super drawRect:dirtyRect];
}


-(void)setVolume:(int)volume {
    _volume = volume;
    [self setNeedsDisplay:YES];
    return;
}
@end
