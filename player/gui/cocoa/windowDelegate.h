#import <Cocoa/Cocoa.h>
@interface WindowDelegate : NSObject<NSWindowDelegate>
{
	NSSize savedAspectRatio;
	NSInteger savedWindowLevel;
}

@end
