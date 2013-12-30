#import <Cocoa/Cocoa.h>

@interface Application : NSObject {
}
- (BOOL)application:(NSApplication *)theApplication openFile:(NSString *)filename;
@end