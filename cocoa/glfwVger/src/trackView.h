//
//  trackView.h
//  glfwVger
//
//  Created by Roy Chen on 10/27/13.
//  Copyright (c) 2013 me. All rights reserved.
//

#import <Cocoa/Cocoa.h>
#import "trackControl.h"
@interface trackView : NSView
{
    CALayer *backgroundLayer;
    trackControl *control;
}

-(void)updateStatus:(NSString *)time leftTime:(NSString *)leftTime percent:(float)percent;
@end
