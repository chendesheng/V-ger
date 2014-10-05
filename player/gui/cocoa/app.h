#import <Cocoa/Cocoa.h>

@interface Application : NSObject {
}
-(BOOL)application:(NSApplication *)theApplication openFile:(NSString *)filename;
-(void)handleAppleEvent:(NSAppleEventDescriptor *)event withReplyEvent:(NSAppleEventDescriptor *)replyEvent;
-(void)openFileMenuItemClick:(id)sender;
-(void)searchSubtitleMenuItemClick:(id)sender;
-(void)timerTick:(NSEvent *)event;
@end