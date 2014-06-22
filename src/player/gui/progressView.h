#import <Cocoa/Cocoa.h>
#import "gui.h"
@interface ProgressView : NSView {
@public
	NSString *titleString;
    NSString *leftString;
    NSString *rightString;
    CGFloat percent;
    CGFloat percent2;
    NSString *speedString;
    CGFloat paddingLeft;
}
@end