#import "app.h"
#import "gui.h"

@implementation Application
- (BOOL)application:(NSApplication *)theApplication openFile:(NSString *)filename {
     NSAutoreleasePool* pool = [[NSAutoreleasePool alloc] init];

     NSLog(@"application openFile: %@",filename);
	const char *cfilename = [filename UTF8String];
	BOOL b = onOpenFile(cfilename) == 1;

     [pool drain];
     return b;
}

-(void)searchSubtitleMenuItemClick:(id)sender {
	onMenuClicked(MENU_SEARCH_SUBTITLE, 0);
}

- (BOOL)applicationShouldOpenUntitledFile:(NSApplication *)sender {
     NSOpenPanel *panel  = [NSOpenPanel openPanel];
     [panel setCanChooseDirectories:NO];
     [panel setAllowsMultipleSelection:NO];

     NSInteger type = [panel runModalForTypes:nil];
     if(type == NSOKButton){
          NSString* filename = [panel filename];
          char* cfilename = (char*)[filename UTF8String];
          onOpenFile(cfilename);
     }
     return NO;
}

-(void)openFileMenuItemClick:(id)sender {
     onOpenOpenPanel();

     NSOpenPanel *panel	= [NSOpenPanel openPanel];
     [panel setCanChooseDirectories:NO];
     [panel setAllowsMultipleSelection:NO];

     NSInteger type	= [panel runModalForTypes:nil];
     if(type == NSOKButton){
          NSString* filename = [panel filename];
          char* cfilename = (char*)[filename UTF8String];
          onCloseOpenPanel(cfilename);
     } else {
          onCloseOpenPanel("");
     	return;
     }
}
- (void)handleAppleEvent:(NSAppleEventDescriptor *)event withReplyEvent:(NSAppleEventDescriptor *)replyEvent {
     NSString *urlString = [[event paramDescriptorForKeyword:keyDirectObject] stringValue];
     const char *cstr = [[urlString substringFromIndex:13] UTF8String];
     if (onOpenFile(cstr) != 1) {
          //
     }
}
- (BOOL)applicationShouldTerminateAfterLastWindowClosed:(NSApplication*)app {
     return YES;
}
@end