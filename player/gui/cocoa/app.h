#import <Cocoa/Cocoa.h>
#import "openURL.h"

@interface Application : NSObject <NSFileManagerDelegate> {
    NSWindowController* winOpenURL;
}

@property (strong, nonatomic) IBOutlet NSMenu *mainMenu;

-(BOOL)application:(NSApplication *)theApplication openFile:(NSString *)filename;
-(void)handleAppleEvent:(NSAppleEventDescriptor *)event withReplyEvent:(NSAppleEventDescriptor *)replyEvent;
-(void)openURL:(id)sender;
-(void)open:(id)sender;
@end
