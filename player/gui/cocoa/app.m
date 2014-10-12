#import "app.h"
#import "gui.h"

@implementation Application
- (BOOL)application:(NSApplication *)theApplication openFile:(NSString *)filename {
     @autoreleasepool {
          NSLog(@"application openFile: %@",filename);

          const char *cfilename = [filename UTF8String];
          BOOL b = onOpenFile(cfilename) == 1;
 
          return b;
     }
}

- (BOOL)applicationShouldOpenUntitledFile:(NSApplication *)sender {
     NSOpenPanel *panel  = [NSOpenPanel openPanel];
     [panel setCanChooseDirectories:NO];
     [panel setAllowsMultipleSelection:NO];

     NSInteger type = [panel runModal];
     if(type == NSOKButton){
          NSString* filename = [[panel URL] path];
          char* cfilename = (char*)[filename UTF8String];
          onOpenFile(cfilename);
     }
     return NO;
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