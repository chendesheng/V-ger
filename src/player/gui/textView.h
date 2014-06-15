#include "gui.h"
#import <Cocoa/Cocoa.h>

//show view subtitles

@interface TextView : NSTextView {
	CGFloat _fontSize;
	CGFloat originalWindowWidth;
}
- (NSSize)sizeForWidth:(float)width
				height:(float)height;

- (void)setFontSize:(CGFloat)size;

- (void)setText:(SubItem*)items length:(int)len;
@end
