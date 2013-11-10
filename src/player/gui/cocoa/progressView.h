#import <Cocoa/Cocoa.h>
#import "gui.h"
@interface ProgressView : NSView {
@public
    NSString *leftString;
    NSString *rightString;
    CGFloat percent;
}
@end