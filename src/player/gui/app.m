#import "app.h"
#import "gui.h"

@implementation Application
- (BOOL)application:(NSApplication *)theApplication openFile:(NSString *)filename {
     NSLog(@"application openFile: %@",filename);
	const char *cfilename = [filename UTF8String];
	return onOpenFile(cfilename) == 1;
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
- (void)handleAppleEvent:(NSAppleEventDescriptor *)event withReplyEvent:(NSAppleEventDescriptor *)replyEvent {
     NSString *urlString = [[event paramDescriptorForKeyword:keyDirectObject] stringValue];
     const char *cstr = [[urlString substringFromIndex:13] UTF8String];
     if (onOpenFile(cstr) != 1) {
          //
     }
}
@end