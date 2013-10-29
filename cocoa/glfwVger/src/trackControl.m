//
//  trackControl.m
//  glfwVger
//
//  Created by Roy Chen on 10/28/13.
//  Copyright (c) 2013 me. All rights reserved.
//

#import "trackControl.h"

@implementation trackControl

-(void)drawRoundedRect:(NSRect)rect radius:(CGFloat)r{
    NSBezierPath *textViewSurround = [NSBezierPath bezierPathWithRoundedRect:rect xRadius:r yRadius:r];
    [textViewSurround fill];
}
-(void)drawRect:(NSRect)dirtyRect
{
//    NSLog(@"draw control");
    
    CGFloat position = (dirtyRect.size.width-120)*(self->percent);
    CGFloat barHeight = 4;
    CGFloat knotHeight = 14;
    CGFloat knotWidth = 5;
    
    [[NSColor colorWithCalibratedRed:255 green:255 blue:255 alpha:0.3] setFill];
    NSRectFill(dirtyRect);
    
    CGFloat x = 8;
    if ([self->leftString length]<=5) {
        x = 22;
    }
    [self->leftString drawAtPoint:NSMakePoint(x, 18) withAttributes:@{NSFontAttributeName : [NSFont fontWithName:@"Helvetica Neue" size:12]}];
    
    [self->rightString drawAtPoint:NSMakePoint(dirtyRect.size.width-60+4, 18) withAttributes:@{NSFontAttributeName : [NSFont fontWithName:@"Helvetica Neue" size:12]}];
    
    [[NSColor colorWithCalibratedRed:0 green:0 blue:0 alpha:0.5] set];
    [self drawRoundedRect:NSMakeRect(60, (dirtyRect.size.height-barHeight)/2, dirtyRect.size.width-120, barHeight) radius:2];
    
    NSShadow* theShadow = [[NSShadow alloc] init];
    [theShadow setShadowOffset:NSMakeSize(0, 0)];
    [theShadow setShadowBlurRadius:1.0];
    
    // Use a partially transparent color for shapes that overlap.
    [theShadow setShadowColor:[[NSColor blackColor]
                               colorWithAlphaComponent:0.5]];
    
    [theShadow set];
    
    [[NSColor colorWithCalibratedRed:255 green:255 blue:255 alpha:1] setFill];
    
    [self drawRoundedRect:NSMakeRect(60, (dirtyRect.size.height-barHeight)/2, position, barHeight) radius:2];
    
    [[NSColor colorWithCalibratedRed:255 green:255 blue:255 alpha:1] setFill];
    [self drawRoundedRect:NSMakeRect(position-knotWidth/2+60, (dirtyRect.size.height-knotHeight)/2, knotWidth, knotHeight) radius:1.5];
    
    [super drawRect:dirtyRect];
}
@end
