//
//  VPAppDelegate.m
//  MenuTest
//
//  Created by Roy Chen on 1/26/14.
//  Copyright (c) 2014 me. All rights reserved.
//

#import "VPAppDelegate.h"

@implementation VPAppDelegate

- (void)dealloc
{
    [super dealloc];
}

- (void)applicationDidFinishLaunching:(NSNotification *)aNotification
{
    // Insert code here to initialize your application
    [NSBundle loadNibNamed:@"status" owner:self];
}

@end
