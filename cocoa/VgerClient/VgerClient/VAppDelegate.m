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

- (void)refreshClick:(id)sender {
    NSLog(@"refresh clicked");
    [[[self web] mainFrame] reload];
}

- (void)applicationDidFinishLaunching:(NSNotification *)aNotification
{
    NSURLCache *sharedCache = [[NSURLCache alloc] initWithMemoryCapacity:0 diskCapacity:0 diskPath:nil];
    [NSURLCache setSharedURLCache:sharedCache];
    [sharedCache release];
    sharedCache = nil;

    [[[self web] preferences] setDefaultFontSize:16];
    [[self web] setFrameLoadDelegate:self];
    
    [[self web] setMainFrameURL:@"http://127.0.0.1:9527"];
}
- (BOOL)applicationShouldTerminateAfterLastWindowClosed:(NSApplication *)theApplication {
    return YES;
}

- (void)webView:(WebView *)sender didFailProvisionalLoadWithError:(NSError *)error forFrame:(WebFrame *)frame {
    [[self web] setMainFrameURL:@"http://127.0.0.1:9527"];
}
- (void)applicationWillTerminate:(NSNotification *)notification {
    [_web stringByEvaluatingJavaScriptFromString: @"onbeforeunload()" ];    
}
//- (void)webView:(WebView *)sender didFinishLoadForFrame:(WebFrame *)frame {
//    NSScrollView *mainScrollView = sender.mainFrame.frameView.documentView.enclosingScrollView;
//    [mainScrollView setVerticalScrollElasticity:NSScrollElasticityNone];
//    [mainScrollView setHorizontalScrollElasticity:NSScrollElasticityNone];
//}
@end
