#import <Cocoa/Cocoa.h>

@interface Application : NSObject <NSFileManagerDelegate> {
}

@property (strong, nonatomic) IBOutlet NSMenu *mainMenu;

-(BOOL)application:(NSApplication *)theApplication openFile:(NSString *)filename;
-(void)handleAppleEvent:(NSAppleEventDescriptor *)event withReplyEvent:(NSAppleEventDescriptor *)replyEvent;
@end
