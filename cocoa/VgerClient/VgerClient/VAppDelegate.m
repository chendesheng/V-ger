//
//  VAppDelegate.m
//  VgerClient
//
//  Created by Roy Chen on 3/9/13.
//  Copyright (c) 2013 me. All rights reserved.
//

#import "VAppDelegate.h"

@implementation VAppDelegate

- (void)dealloc
{
    [super dealloc];
}

- (void)applicationDidFinishLaunching:(NSNotification *)aNotification
{
    // Insert code here to initialize your application
    [[WebPreferences standardPreferences] setCacheModel:WebCacheModelDocumentViewer];
    [[WebPreferences standardPreferences] setUsesPageCache:NO];
    
    NSURLCache *sharedCache = [[NSURLCache alloc] initWithMemoryCapacity:0 diskCapacity:0 diskPath:nil];
    [NSURLCache setSharedURLCache:sharedCache];
    [sharedCache release];
    sharedCache = nil;

    [[[self web] preferences] setDefaultFontSize:16];

    [[[self web] mainFrame] loadRequest:[NSURLRequest requestWithURL:[NSURL URLWithString:@"http://127.0.0.1:9527"]]];

}
- (BOOL)applicationShouldTerminateAfterLastWindowClosed:(NSApplication *)theApplication {
    return YES;
}
@end
