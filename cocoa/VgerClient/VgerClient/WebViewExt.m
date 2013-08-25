//
//  V.m
//  VgerClient
//
//  Created by Roy Chen on 8/25/13.
//  Copyright (c) 2013 me. All rights reserved.
//

#import <WebKit/WebKit.h>

@interface WebView (WebViewExt)

- (BOOL)performKeyEquivalent:(NSEvent *)theEvent;

@end

@implementation WebView (WebViewExt)

- (BOOL)performKeyEquivalent:(NSEvent *)theEvent
{
    if ( [theEvent modifierFlags] & NSCommandKeyMask)
    {
        NSString *chars = [theEvent charactersIgnoringModifiers];
        
        if ([chars isEqualToString:@"x"])
        {
            [self cut:self];
            return YES;
        }
        
        if ([chars isEqualToString:@"c"])
        {
            [self copy:self];
            return YES;
        }
        
        if ([chars isEqualToString:@"v"])
        {
            [self paste:self];
            return YES;
        }
        
        if ([chars isEqualToString:@"z"])
        {
            [[self undoManager] undo];
            return YES;
        }
        
        if ([chars isEqualToString:@"y"])
        {
            [[self undoManager] redo];
            return YES;
        }
    }
    
    return [super performKeyEquivalent:theEvent];
}

@end
