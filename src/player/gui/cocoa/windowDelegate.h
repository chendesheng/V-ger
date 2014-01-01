#import <Cocoa/Cocoa.h>
@interface WindowDelegate : NSObject<NSWindowDelegate>
{
@public
	NSWindow* window;
	NSSize savedAspectRatio;
}

@end
