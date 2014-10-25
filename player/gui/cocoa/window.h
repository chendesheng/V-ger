#import <Cocoa/Cocoa.h>
#import "windowDelegate.h"
#import "gui.h"
#import "glView.h"
#import "blurView.h"
#import "titleTextView.h"
@interface Window : NSWindow {
    BlurView* bvTitleTextView;
    TitleTextView* titleTextView;
@public
	NSSize customAspectRatio;
	GLView* glView;
}

- (id)initWithWidth:(int)w height:(int)h;
- (void)makeCurrentContext;
- (void)audioMenuItemClick:(id)sender;
- (void)subtitleMenuItemClick:(id)sender;
- (void)setTitleHidden:(BOOL)b;
- (void)playPause:(id)sender;
- (void)open:(id)sender;
- (void)setRoundCorner:(bool)b;
@end
