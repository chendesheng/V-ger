#import "windowDelegate.h"
#import "gui.h"
#import "window.h"
@implementation WindowDelegate

- (void)windowWillEnterFullScreen:(NSNotification *)notification {
	Window* w = (Window*)[notification object];

	self->savedWindowLevel = w.level;
	w.level = NSNormalWindowLevel;
	self->savedAspectRatio = w.aspectRatio; 
	w.aspectRatio = [[NSScreen mainScreen] frame].size;

	onFullScreen(WILL_ENTER_FULL_SCREEN);
}

- (void)windowDidEnterFullScreen:(NSNotification *)notification {
	onFullScreen(DID_ENTER_FULL_SCREEN);
}

- (void)windowWillExitFullScreen:(NSNotification *)notification {
	Window* w = (Window*)[notification object];

	onFullScreen(WILL_EXIT_FULL_SCREEN);
}

- (void)windowDidExitFullScreen:(NSNotification *)notification {
	Window* w = (Window*)[notification object];
	w.level = self->savedWindowLevel;

	w.aspectRatio = self->savedAspectRatio;

	onFullScreen(DID_EXIT_FULL_SCREEN);
}

- (void)windowDidResize:(NSNotification *)notification {
	//NSLog(@"DidResize");

	Window* w = (Window*)[notification object];
	[w->glView refreshTexts];
}

@end
