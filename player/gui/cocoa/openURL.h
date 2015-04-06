#import <Cocoa/Cocoa.h>
#import "gui.h"

@interface OpenURL : NSWindowController

- (IBAction)openClick:(id)sender;
- (IBAction)cancelClick:(id)sender;

@property (strong, nonatomic) IBOutlet NSTextField *txtURL;
@property (strong, nonatomic) IBOutlet NSButton *btnOpen;

@end
