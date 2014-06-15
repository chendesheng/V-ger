#import <Cocoa/Cocoa.h>
#import "blurView.h"

@interface StartupView : NSView {
	NSProgressIndicator* _progressIndicator;
	BlurView* _blurView;
}
@end
