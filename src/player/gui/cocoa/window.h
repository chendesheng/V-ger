#import <Cocoa/Cocoa.h>
#import "windowDelegate.h"
#import "gui.h"

@interface Window : NSWindow {
@public
	NSSize customAspectRatio;
}

- (id)initWithTitle:(NSString*)title width:(int)w height:(int)h;
- (void)setContentViewNeedsDisplay:(BOOL)b;
- (void)timerTick:(NSEvent *)event;
- (void)makeCurrentContext;
@end
