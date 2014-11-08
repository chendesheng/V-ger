#import <Cocoa/Cocoa.h>
@interface WindowDelegate : NSObject<NSWindowDelegate>
{
	NSInteger savedWindowLevel;
    NSSize savedAspectRatio;
}

@end
