#import "blurView.h"

@implementation BlurView

- (id)initWithFrame:(NSRect)frame {
    self = [super initWithFrame:frame];
    if (self) {
        // Initialization code here.
        self->backgroundLayer = [CALayer layer];
        [self setLayer:self->backgroundLayer];
        [self setWantsLayer:YES];
        CIFilter *blurFilter = (CIFilter*)[CIFilter filterWithName:@"CIGaussianBlur" keysAndValues:@"inputRadius", [NSNumber numberWithFloat:20.0], nil];
        // [blurFilter setDefaults];
        
        [self->backgroundLayer setMasksToBounds:YES];
        
        [self layer].backgroundFilters = [NSArray arrayWithObject:blurFilter];
        
        self.layerContentsRedrawPolicy = NSViewLayerContentsRedrawOnSetNeedsDisplay;
    }
    
    return self;
}
-(void)setHidden:(BOOL)flag {
    if (flag == YES)
    {
        //        [self setLayer:NULL];
        [self setFrameSize:NSMakeSize(self.frame.size.width, 0)];
    }
    else
    {
        [self setFrameSize:NSMakeSize(self.frame.size.width, 30)];
    }
    [super setHidden:flag];
}

- (void)mouseDragged:(NSEvent *)event {}
- (void)mouseUp:(NSEvent *)event {}
- (void)mouseMoved:(NSEvent *)event{}
- (void)mouseDown:(NSEvent *)event{}
@end

