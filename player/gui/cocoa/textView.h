#include "gui.h"
#import <Cocoa/Cocoa.h>

//show view subtitles

@interface TextView : NSTextView {
	CGFloat _fontSize;
	CGFloat originalWindowWidth;
@public
	CGFloat x, y;
	int align;
}
- (NSSize)sizeForWidth:(float)width
				height:(float)height;

- (void)setFontSize:(CGFloat)size;

- (void)setText:(AttributedString*)items length:(int)len;

@end
