//
//  trackView.m
//  glfwVger
//
//  Created by Roy Chen on 10/27/13.
//  Copyright (c) 2013 me. All rights reserved.
//

#import "trackView.h"

@implementation trackView

- (id)initWithFrame:(NSRect)frame
{
    self = [super initWithFrame:frame];
    if (self) {
        // Initialization code here.
        self->backgroundLayer = [CALayer layer];
        [self setLayer:self->backgroundLayer];
        [self setWantsLayer:YES];
        CIFilter *blurFilter = [CIFilter filterWithName:@"CIGaussianBlur" keysAndValues:@"inputRadius", [NSNumber numberWithFloat:20.0], nil];
        //[blurFilter setDefaults];

        [self->backgroundLayer setMasksToBounds:YES];
        
        [self layer].backgroundFilters = [NSArray arrayWithObject:blurFilter];
        
        self->control = [[trackControl alloc] initWithFrame:NSMakeRect(0, 0, frame.size.width, frame.size.height)];
        self->control->leftString = @"--:--:--";
        self->control->rightString = @"--:--:--";
        self->control->percent = 0;
        [self->control setAutoresizingMask:NSViewWidthSizable|NSViewHeightSizable];
        [self addSubview:self->control];
        
        self.layerContentsRedrawPolicy = NSViewLayerContentsRedrawOnSetNeedsDisplay;
    }
    
    return self;
}
-(void)setHidden:(BOOL)flag
{
    if (flag == YES)
    {
//        [self setLayer:NULL];
        [self setFrameSize:NSMakeSize(self.frame.size.width, 0)];
    }
    else
    {
        [self setFrameSize:NSMakeSize(self.frame.size.width, 50)];
    }
    [self->control setNeedsDisplay:YES];
    [super setHidden:flag];
}
-(void)updateStatus:(NSString *)time leftTime:(NSString *)leftTime percent:(float)percent
{
    [self->control->leftString autorelease];
    self->control->leftString = time;
    [self->control->leftString retain];
    
    [self->control->rightString autorelease];
    self->control->rightString = leftTime;
    [self->control->rightString retain];
    
    self->control->percent = percent;
    
//    NSLog(@"%@ %@ %lf", time, leftTime, percent);
    [self->control setNeedsDisplay:YES];
}
//- (void)drawRect:(NSRect)dirtyRect
//{
//    NSString *str = @"00:00:00";
//    
//    NSFontManager *fontManager = [NSFontManager sharedFontManager];
//    NSFont *font = [fontManager fontWithFamily:@"Georgia"
//                                        traits:NSUnboldFontMask
//                                        weight:0
//                                          size:13];
//
//    [str drawAtPoint:NSMakePoint(10, 10) withAttributes:@{NSFontAttributeName : font, NSForegroundColorAttributeName:[NSColor blackColor]}];
//    [super drawRect:dirtyRect];
//}

@end
