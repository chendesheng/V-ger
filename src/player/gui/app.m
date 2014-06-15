#import "app.h"
#import "gui.h"

@implementation Application
- (BOOL)application:(NSApplication *)theApplication openFile:(NSString *)filename {
     NSLog(@"doOpen filename = %@",filename);
	const char *cfilename = [filename UTF8String];
	return onOpenFile(cfilename) == 1;
}
- (void)applicationWillTerminate:(NSNotification *)aNotification {
	NSLog(@"applicationWillTerminate");
	onWillTerminate();
}

-(void)searchSubtitleMenuItemClick:(id)sender {
	onSearchSubtitleMenuItemClick();
}
-(void)openFileMenuItemClick:(id)sender {
     // onOpenOpenPanel();

     // NSOpenPanel *panel	= [NSOpenPanel openPanel];
     // [panel setCanChooseDirectories:NO];
     // [panel setAllowsMultipleSelection:NO];

     // NSInteger type	= [panel runModalForTypes:nil];
     // if(type == NSOKButton){
     //      NSString* filename = [panel filename];
     //      char* cfilename = (char*)[filename UTF8String];
     //      onOpenFile(cfilename);

     //      onCloseOpenPanel(cfilename);
     // } else {
     //      onCloseOpenPanel("");
     // 	return;
     // }
}
@end