#import <Cocoa/Cocoa.h>
@interface WindowDelegate : NSObject
{
@public
	NSWindow* window;
	NSSize savedAspectRatio;
}

@end
