#import <Cocoa/Cocoa.h>
#import "YRKSpinningProgressIndicator.h"
#import "blurView.h"

@interface StartupView : NSView {
	NSProgressIndicator* _progressIndicator;
	BlurView* _blurView;
}
@end
