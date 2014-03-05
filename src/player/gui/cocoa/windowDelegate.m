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
    
    Window* w = (Window*)(self->window);
    self->savedAspectRatio = w->customAspectRatio;
    w->customAspectRatio = frame.size;

    // NSLog(@"windowWillEnterFullScreen");

}
- (void)windowDidEnterFullScreen:(NSNotification *)notification
{
    onFullscreenChanged((void*)[notification object], 1);
}
- (void)windowWillExitFullScreen:(NSNotification *)notification
{ 
    // NSLog(@"windowWillExitFullScreen");
}
- (void)updateRoundCorner {
    Window* w = (Window*)(self->window);
    NSView* rv = [[w contentView] superview];
    [rv setWantsLayer:YES];
    rv.layer.cornerRadius=4.3;
    rv.layer.masksToBounds=YES;
    [rv setNeedsDisplay:YES];
}

- (void)windowDidExitFullScreen:(NSNotification *)notification
{
    [self updateRoundCorner];

    Window* w = (Window*)(self->window);
    w->customAspectRatio = self->savedAspectRatio;

    onFullscreenChanged((void*)[notification object], 0);
}
- (NSSize)windowWillResize:(NSWindow *)sender toSize:(NSSize)frameSize {
    // NSLog(@"windowWillResize");
    [self updateRoundCorner];
 
	NSRect r;
	Window* w = (Window*)sender;

	r = NSMakeRect([w frame].origin.x, [w frame].origin.y,
	frameSize.width, frameSize.height);

	NSSize aspectRatio = w->customAspectRatio;
	// r = [w contentRectForFrameRect:r];
	r.size.height = r.size.width * aspectRatio.height / aspectRatio.width;
	// r = [w frameRectForContentRect:r];

    [w->titlebarView setFrame:NSMakeRect(0, r.size.height-30, r.size.width, 30)];
	return r.size;
}
-(void)windowDidEndLiveResize:(NSNotification *)notification {
    [self updateRoundCorner];
}
@end