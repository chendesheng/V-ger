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
    setControlsVisible(w, 0, 0);
    [w setTitleHidden:NO];
}

- (void)windowWillExitFullScreen:(NSNotification *)notification {
    Window* w = (Window*)[notification object];
    setControlsVisible(w, true, true);
}

- (void)windowDidExitFullScreen:(NSNotification *)notification {
    Window* w = (Window*)[notification object];
    w.level = self->savedWindowLevel;

    w.aspectRatio = self->savedAspectRatio;
}

- (void)windowDidResize:(NSNotification *)notification {
    NSLog(@"DidResize");

    Window* w = (Window*)[notification object];
    [w->glView refreshTexts];
}

@end
