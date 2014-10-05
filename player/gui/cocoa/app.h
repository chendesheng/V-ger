#import <Cocoa/Cocoa.h>

@interface Application : NSObject {
}

@property (strong, nonatomic) IBOutlet NSMenu *mainMenu;

-(BOOL)application:(NSApplication *)theApplication openFile:(NSString *)filename;
-(void)handleAppleEvent:(NSAppleEventDescriptor *)event withReplyEvent:(NSAppleEventDescriptor *)replyEvent;
-(void)open:(id)sender;
-(void)searchSubtitle:(id)sender;
-(void)timerTick:(NSEvent *)event;
@end