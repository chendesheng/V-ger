#import <Cocoa/Cocoa.h>
#import "windowDelegate.h"
#import "gui.h"

@interface Window : NSWindow {}

- (id)initWithTitle:(NSString*)title width:(int)w height:(int)h;
- (void)setContentViewNeedsDisplay:(BOOL)b;
- (void)timerTick:(NSEvent *)event;
- (void)makeCurrentContext;
@end
