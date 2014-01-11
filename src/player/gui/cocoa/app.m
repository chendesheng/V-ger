#import "app.h"
#import "gui.h"

@implementation Application
- (BOOL)application:(NSApplication *)theApplication openFile:(NSString *)filename {
	const char *cfilename = [filename UTF8String];
	return onOpenFile(cfilename) == 1;
}
- (void)applicationWillTerminate:(NSNotification *)aNotification {
	NSLog(@"applicationWillTerminate");
	onWillTerminate();
}
@end