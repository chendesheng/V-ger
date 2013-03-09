//
//  main.m
//  VgerHelper
//
//  Created by Roy Chen on 3/5/13.
//  Copyright (c) 2013 me. All rights reserved.
//

#import <Cocoa/Cocoa.h>
#import "VAppDelegate.h"
int main(int argc, char *argv[])
{
    if (argc == 2 || argc == 1) {
        [[NSWorkspace sharedWorkspace] openFile:@"/applications/V'ger.app"];
        return 0;
    }
    
    NSAutoreleasePool* pool = [[NSAutoreleasePool alloc] init];
    
    [pool release];
    
    VAppDelegate * appDelegate = [[[VAppDelegate alloc]init]autorelease];
    
    NSApplication * application = [NSApplication sharedApplication];
    [application setDelegate:appDelegate];
    
    [application run];
    
    [pool drain];
}
