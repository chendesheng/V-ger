#import <Cocoa/Cocoa.h>

@interface BlurView : NSVisualEffectView {
}

-(BlurView*) initWithView:(NSView*)v frame:(NSRect)bounds;
-(void) setCornerRadius:(CGFloat)r;

@end
