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
- (BOOL)performKeyEquivalent:(NSEvent *)theEvent
{
    if ( [theEvent modifierFlags] & NSCommandKeyMask)
    {
        NSString *chars = [theEvent charactersIgnoringModifiers];
        
        if ([chars isEqualToString:@"x"])
        {
            [_web cut:_web];
            return YES;
        }
        
        if ([chars isEqualToString:@"c"])
        {
            [_web copy:_web];
            return YES;
        }
        
        if ([chars isEqualToString:@"v"])
        {
            [_web paste:_web];
            return YES;
        }
        
    }
    
    return [_web performKeyEquivalent:theEvent];
}

//- (void)webView:(WebView *)sender didFinishLoadForFrame:(WebFrame *)frame {
//    NSScrollView *mainScrollView = sender.mainFrame.frameView.documentView.enclosingScrollView;
//    [mainScrollView setVerticalScrollElasticity:NSScrollElasticityNone];
//    [mainScrollView setHorizontalScrollElasticity:NSScrollElasticityNone];
//}
@end
