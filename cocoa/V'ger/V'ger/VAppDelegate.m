//
//  VAppDelegate.m
//  V'ger
//
//  Created by Roy Chen on 3/4/13.
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
//    [[NSWorkspace alloc] launchApplication:@"~/vgerproj/bin/vger"];  
//    NSTask *task = [[[NSTask alloc] init] autorelease];
//    [task setLaunchPath:@"~/vgerproj/bin/vger"];
////    NSPipe *outputPipe = [NSPipe pipe];
////    [task setStandardOutput:outputPipe];
////    [[NSNotificationCenter defaultCenter] addObserver:self selector:@selector(readCompleted:) name:NSFileHandleReadToEndOfFileCompletionNotification object:[outputPipe fileHandleForReading]];
////    [[outputPipe fileHandleForReading] readToEndOfFileInBackgroundAndNotify];
//    [task launch];
//    [task waitUntilExit];
//    if([task isRunning]) {
//        NSLog(@"Task is running");
//        // Insert code here to initialize your application
//        [[_web mainFrame] loadRequest:[NSURLRequest requestWithURL:[NSURL URLWithString:@"http://localhost:3824"]]];
//    } else {
//        NSLog(@"end");
//    }
    NSTask *task;
    task = [[NSTask alloc] init];
    [task setLaunchPath: @"~/vgerproj/bin/vger"];
    
    NSArray *arguments;
    arguments = [NSArray arrayWithObjects: @"-l", @"-a", @"-t", nil];
    [task setArguments: arguments];
    
    NSPipe *pipe;
    pipe = [NSPipe pipe];
    [task setStandardOutput: pipe];
    
    NSFileHandle *file;
    file = [pipe fileHandleForReading];
    
    [task launch];
    
    NSData *data;
    data = [file readDataToEndOfFile];
    
    NSString *string;
    string = [[NSString alloc] initWithData: data
                                   encoding: NSUTF8StringEncoding];
    NSLog (@"got\n%@", string);
}

- (void)readCompleted:(NSNotification *)notification {
    NSLog(@"Read data: %@", [[notification userInfo] objectForKey:NSFileHandleNotificationDataItem]);
    [[NSNotificationCenter defaultCenter] removeObserver:self name:NSFileHandleReadToEndOfFileCompletionNotification object:[notification object]];
}
@end
