#import <Cocoa/Cocoa.h>
@interface TitleTextView : NSView {
@public
	NSString* titleString;
}
-(void)setTitle:(NSString*)title;
@end