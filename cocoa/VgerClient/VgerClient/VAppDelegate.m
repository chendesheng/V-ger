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
    
    [[self web] setUIDelegate:self];
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
- (WebView *)webView:(WebView *)sender createWebViewWithRequest:(NSURLRequest *)request
{
    WebView *_hiddenWebView=[[WebView alloc] init];
    [_hiddenWebView setPolicyDelegate:self];
    return _hiddenWebView;
}

- (void)webView:(WebView *)sender decidePolicyForNavigationAction:(NSDictionary *)actionInformation request:(NSURLRequest *)request frame:(WebFrame *)frame decisionListener:(id<WebPolicyDecisionListener>)listener {
    NSLog(@"%@",[[actionInformation objectForKey:WebActionOriginalURLKey] absoluteString]);
    [[NSWorkspace sharedWorkspace] openURL:[actionInformation objectForKey:WebActionOriginalURLKey]];
    [sender release];
}
@end
