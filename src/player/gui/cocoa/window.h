#import <Cocoa/Cocoa.h>
#import "windowDelegate.h"
#import "gui.h"
#import "glView.h"

@interface Window : NSWindow {
@public
	NSSize customAspectRatio;
	GLView* glView;
	NSView* titlebarView;
}

- (id)initWithTitle:(NSString*)title width:(int)w height:(int)h;
- (void)setContentViewNeedsDisplay:(BOOL)b;
- (void)timerTick:(NSEvent *)event;
- (void)makeCurrentContext;
- (void)updateRoundCorner;
@end
