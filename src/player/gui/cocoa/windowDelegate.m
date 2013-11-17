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
    NSLog(@"%lf %lf", frame.size.width, frame.size.height);
    
    Window* w = (Window*)(self->window);
    self->savedAspectRatio = w->customAspectRatio;
    w->customAspectRatio = frame.size;

    NSLog(@"windowWillEnterFullScreen");

    onFullscreenChanged((void*)[notification object], 1);
}
- (void)windowDidEnterFullScreen:(NSNotification *)notification
{
    NSLog(@"windowDidEnterFullScreen");
}
- (void)windowWillExitFullScreen:(NSNotification *)notification
{ 
    onFullscreenChanged((void*)[notification object], 0);
    NSLog(@"windowWillExitFullScreen");
}
- (void)windowDidExitFullScreen:(NSNotification *)notification
{
    NSLog(@"windowDidExitFullScreen");

    Window* w = (Window*)(self->window);
    w->customAspectRatio = self->savedAspectRatio;

}
- (NSSize)windowWillResize:(NSWindow *)sender toSize:(NSSize)frameSize {
    NSLog(@"windowWillResize");

	NSRect r;
	Window* w = (Window*)sender;

	r = NSMakeRect([w frame].origin.x, [w frame].origin.y,
	frameSize.width, frameSize.height);

	NSSize aspectRatio = w->customAspectRatio;
	r = [w contentRectForFrameRect:r];
	r.size.height = r.size.width * aspectRatio.height / aspectRatio.width;
	r = [w frameRectForContentRect:r];
	return r.size;
}
@end