#import "windowDelegate.h"
#import "gui.h"
#import "window.h"
@implementation WindowDelegate


- (void)windowWillEnterFullScreen:(NSNotification *)notification {
    NSScreen *mainScreen = [NSScreen mainScreen];
    NSRect frame = [mainScreen frame];
    
    Window* w = (Window*)[notification object];
    self->savedAspectRatio = w->customAspectRatio;
    w->customAspectRatio = frame.size;

    self->savedWindowLevel = w.level;
    w.level = NSNormalWindowLevel;

    setControlsVisible(w, 0, 0);
    [w setRoundCorner:NO];  //hide round corner in full screen
}
- (void)windowDidEnterFullScreen:(NSNotification *)notification {
    Window* w = (Window*)[notification object];
    [w->glView refreshTexts];
}

- (void)windowWillExitFullScreen:(NSNotification *)notification {
    Window* w = (Window*)[notification object];
    setControlsVisible(w, true, true);
    [w setRoundCorner:YES];
}
- (void)windowDidExitFullScreen:(NSNotification *)notification {
    Window* w = (Window*)[notification object];

    w->customAspectRatio = self->savedAspectRatio;

    [w->glView refreshTexts];

    w.level = self->savedWindowLevel;
}
- (void)windowDidResize:(NSNotification *)notification {
    Window* w = (Window*)[notification object];
    [w->glView refreshTexts];
}
- (NSSize)windowWillResize:(NSWindow *)sender toSize:(NSSize)frameSize {
    Window* w = (Window*)sender;
 
	NSRect r = NSMakeRect([w frame].origin.x, [w frame].origin.y, frameSize.width, frameSize.height);
	NSSize aspectRatio = w->customAspectRatio;
	r.size.height = r.size.width * aspectRatio.height / aspectRatio.width;

    setControlsVisible(w, false, false);
	return r.size;
}

-(void)windowDidEndLiveResize:(NSNotification *)notification {
    Window* w = (Window*)[notification object];
    [w->glView refreshTexts];
}
- (NSRect)windowWillUseStandardFrame:(NSWindow *)window defaultFrame:(NSRect)newFrame {
    Window* w = (Window*)window;
    NSRect r = newFrame;
    NSSize aspectRatio = w->customAspectRatio;
    double d = aspectRatio.width/aspectRatio.height;
    if (r.size.width > r.size.height*d) {
        r.size.width = r.size.height * d;
    } else {
        r.size.height = r.size.width / d;
    }
    r.origin.x += (newFrame.size.width - r.size.width)/2;
    r.origin.y += (newFrame.size.height - r.size.height)/2;
    return r;
}
//lost focus
// -(void)windowDidResignKey:(NSNotification *)notification {
//     Window* w = (Window*)[notification object];
//     setControlsVisible(w, 0, 0);
// }
@end
