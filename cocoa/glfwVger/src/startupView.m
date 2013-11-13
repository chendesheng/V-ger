//
//  StartupView.m
//  glfwVger
//
//  Created by Roy Chen on 10/27/13.
//  Copyright (c) 2013 me. All rights reserved.
//

#import "StartupView.h"

@implementation StartupView

- (id)initWithFrame:(NSRect)frame
{
    self = [super initWithFrame:frame];
    if (self) {
        // Initialization code here.
    }
    
    return self;
}

- (void)drawRect:(NSRect)dirtyRect
{
    [[NSColor blackColor] setFill];
    NSRectFill(dirtyRect);
    [super drawRect:dirtyRect];
}

@end
