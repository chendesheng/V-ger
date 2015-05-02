#import <Cocoa/Cocoa.h>
#import "gui.h"
@interface ProgressView : NSView {
	NSString *_titleString;
	NSString *_leftString;
	NSString *_rightString;
	CGFloat _percent;
	CGFloat _percent2;
	NSString *_speedString;
	CGFloat _paddingLeft;
	NSTrackingArea* _trackingArea;
}
-(void)updatePorgressInfo:(NSString*)leftStr rightString:(NSString*)rightStr percent:(CGFloat)p;
-(void)updateBufferInfo:(NSString*)speed bufferPercent:(CGFloat)p;
@end
