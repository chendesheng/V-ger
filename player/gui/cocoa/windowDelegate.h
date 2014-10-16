#import <Cocoa/Cocoa.h>
@interface WindowDelegate : NSObject<NSWindowDelegate>
{
	NSSize savedAspectRatio;
	int savedWindowLevel;
}

@end
