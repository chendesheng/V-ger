//
//  VAppDelegate.m
//  VgerHelper
//
//  Created by Roy Chen on 3/5/13.
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
    NSArray *arguments = [[NSProcessInfo processInfo] arguments];
    NSString *cmd = arguments[1];
    
    if ([cmd isEqualToString:@"notification"]) {
        NSString* title = arguments[3];
        NSString* message = arguments[4];
    
        [self sendNotification:title message:message];
        
        [NSTimer scheduledTimerWithTimeInterval:10.0
                                         target:self
                                       selector:@selector(timeout:)
                                       userInfo:nil
                                        repeats:NO];
    } else if ([cmd isEqualToString:@"trash"]) {
        NSInteger tag;
        
        NSString* dir = arguments[2];
        NSString* name = arguments[3];
        NSArray* array = [NSArray arrayWithObjects:name, nil];
        [[NSWorkspace sharedWorkspace] performFileOperation:NSWorkspaceRecycleOperation source:dir destination:@"" files:array tag:&tag];
        
        [NSApp terminate: nil];
    } else if ([cmd isEqualToString:@"shutdown"]) {
        [NSTimer scheduledTimerWithTimeInterval:60.0
                                         target:self
                                       selector:@selector(tick:)
                                       userInfo:nil
                                        repeats:NO];
        NSString *msg = [[NSProcessInfo processInfo] arguments][2];
        [self sendNotification:@"V'ger shutdown after 60 seconds." message:msg];
    }
}
-(void) sendNotification:(NSString*)title message:(NSString*)message
{    
    NSUserNotification *notification = [[NSUserNotification alloc] init];
    notification.title = title;
    notification.informativeText = message;
    notification.soundName = NSUserNotificationDefaultSoundName;
    notification.hasActionButton = true;
    notification.actionButtonTitle = @"Play";
    notification.deliveryRepeatInterval = nil;
    //    notification.subtitle = @"test";
    NSUserNotificationCenter *center = [NSUserNotificationCenter defaultUserNotificationCenter];
    [center setDelegate:self];
    [center deliverNotification:notification];
    // Insert code here to initialize your application
    //[NSApp terminate: nil];
    
}

- (void) timeout:(NSTimer *)timer {
    [NSApp terminate: nil];
}

- (void) tick:(NSTimer *)timer {
    [self shutdown];
}
- (void) shutdown {
    NSAppleScript* script = [[NSAppleScript alloc] initWithSource:
                                @"Tell application \"System Events\" to shut down"];
    if (script != NULL)
    {
        NSDictionary* errDict = NULL;
        // execution of the following line ends with EXC
        if (YES == [script compileAndReturnError: &errDict])
        {
            NSLog(@"compiled the script");
            [script executeAndReturnError: &errDict];
        }
        [script release];
    }
    [NSApp terminate: nil];
}
- (void) userNotificationCenter:(NSUserNotificationCenter *)center didActivateNotification:(NSUserNotification *)notification
{
    NSArray *arguments = [[NSProcessInfo processInfo] arguments];
    NSString* cmd = arguments[1];
    if (![cmd isEqualToString:@"shutdown"]) {
//        NSString *url = @"http://";
//        url = [url stringByAppendingString:arguments[2]];
//
//        [[NSWorkspace sharedWorkspace] openURL:[NSURL URLWithString:url]];
        
        [[NSWorkspace sharedWorkspace] openFile:@"/applications/V'ger.app"];
        
        [center removeDeliveredNotification: notification];
    }
    
    [NSApp terminate: nil];
}
@end
