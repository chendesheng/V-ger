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

	Window* w = (Window*)[notification object];
	[w setTitleHidden:NO];
}

- (void)windowWillExitFullScreen:(NSNotification *)notification {
	onFullScreen(WILL_EXIT_FULL_SCREEN);

	Window* w = (Window*)[notification object];

	[w setTitleHidden:YES];
	setControlsVisible(w, 0, 0);
}

- (void)windowDidExitFullScreen:(NSNotification *)notification {
	Window* w = (Window*)[notification object];
	w.level = self->savedWindowLevel;

	w.aspectRatio = self->savedAspectRatio;

	onFullScreen(DID_EXIT_FULL_SCREEN);

	//This is required or the progress view will broken, the progress view must show up right after exit from full screen. FIXME: I don't known why. 
	setControlsVisible(w, 1, 1);
}

- (void)windowDidResize:(NSNotification *)notification {
	//NSLog(@"DidResize");

	Window* w = (Window*)[notification object];
	[w->glView refreshTexts];
}

@end
