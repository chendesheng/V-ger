//
//  VAppDelegate.h
//  V'ger
//
//  Created by Roy Chen on 3/4/13.
//  Copyright (c) 2013 me. All rights reserved.
//

#import <Cocoa/Cocoa.h>
#import <WebKit/WebKit.h>

@interface VAppDelegate : NSObject <NSApplicationDelegate>

@property (assign) IBOutlet NSWindow *window;
@property (assign) IBOutlet WebView *web;

@end
