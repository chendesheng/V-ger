#import "windowDelegate.h"
#import "gui.h"
#import "window.h"
@implementation WindowDelegate


- (void)windowWillClose:(id)sender
{
    [NSApp terminate:nil];
}

- (void)windowWillEnterFullScreen:(NSNotification *)notification
{
    NSScreen *mainScreen = [NSScreen mainScreen];
    NSRect frame = [mainScreen frame];
    // NSLog(@"%lf %lf", frame.size.width, frame.size.height);
    
    Window* w = (Window*)[notification object];
    self->savedAspectRatio = w->customAspectRatio;
    w->customAspectRatio = frame.size;

    // NSLog(@"windowWillEnterFullScreen");

}
- (void)windowDidEnterFullScreen:(NSNotification *)notification
{
    Window* w = (Window*)[notification object];
    [w->glView hideCursor];
    [w->glView hideProgress];

    onFullscreenChanged((void*)[notification object], 1);
}
// - (void)windowWillExitFullScreen:(NSNotification *)notification
// { 
//     // NSLog(@"windowWillExitFullScreen");
// }

- (void)windowDidExitFullScreen:(NSNotification *)notification
{
    Window* w = (Window*)[notification object];

    [w updateRoundCorner];

    [w->glView hideCursor];
    [w->glView hideProgress];

    w->customAspectRatio = self->savedAspectRatio;

    onFullscreenChanged((void*)[notification object], 0);
}
- (NSSize)windowWillResize:(NSWindow *)sender toSize:(NSSize)frameSize {
    // NSLog(@"windowWillResize");
    Window* w = (Window*)sender;
    [w updateRoundCorner];
 
	NSRect r = NSMakeRect([w frame].origin.x, [w frame].origin.y, frameSize.width, frameSize.height);

	NSSize aspectRatio = w->customAspectRatio;
	// r = [w contentRectForFrameRect:r];
	r.size.height = r.size.width * aspectRatio.height / aspectRatio.width;
	// r = [w frameRectForContentRect:r];

	return r.size;
}
-(void)windowDidEndLiveResize:(NSNotification *)notification {
    NSLog(@"windowDidEndLiveResize");
    Window* w = (Window*)[notification object];
    [w updateRoundCorner];
    
    onFullscreenChanged((void*)[notification object], 0);
}
- (NSRect)windowWillUseStandardFrame:(NSWindow *)window defaultFrame:(NSRect)newFrame {
    // NSLog(@"windowWillUseStandardFrame:%lf,%lf,%lf,%lf", newFrame.origin.x, newFrame.origin.y, newFrame.size.width, newFrame.size.height);

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
-(void)windowDidResignKey:(NSNotification *)notification {
    Window* w = (Window*)[notification object];

    [w->glView hideCursor];
    [w->glView hideProgress];
}
@end