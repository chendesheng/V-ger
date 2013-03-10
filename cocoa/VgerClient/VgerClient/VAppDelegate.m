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
    [[self web] setPolicyDelegate:self];
    [[self web] setUIDelegate:self];
}
- (BOOL)applicationShouldTerminateAfterLastWindowClosed:(NSApplication *)theApplication {
    return YES;
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
- (void)webView:(WebView *)sender runJavaScriptAlertPanelWithMessage:(NSString *)message initiatedByFrame:(WebFrame *)frame {
    NSAlert *alert = [[NSAlert alloc] init];
    [alert addButtonWithTitle:@"OK"];
    [alert setInformativeText:message];
    [alert setMessageText:@"V'ger problem"];
    [alert runModal];
    [alert release];
}
@end
