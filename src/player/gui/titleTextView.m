#import "titleTextView.h"
@implementation TitleTextView

- (id)initWithFrame:(NSRect)frame {
    self = [super initWithFrame:frame];
	if (self) {
		self->_title = @"";
	}

	return self;
}

- (void)drawRect:(NSRect)dirtyRect {
    [[NSColor blackColor] setFill];
    
    NSMutableParagraphStyle *style = [[NSParagraphStyle defaultParagraphStyle] mutableCopy];
	[style setAlignment:NSCenterTextAlignment];
	[style setLineBreakMode:NSLineBreakByTruncatingMiddle];

    NSDictionary *attr = @{NSFontAttributeName : [NSFont titleBarFontOfSize:13], NSParagraphStyleAttributeName: style};
    // NSSize titlesz = [self->_title sizeWithAttributes:attr];
    // [self->_title drawAtPoint:NSMakePoint((dirtyRect.size.width-titlesz.width)/2, 3) withAttributes:attr];
    [self->_title drawWithRect:NSMakeRect(82, 6, self.bounds.size.width-164, 16) options:NSStringDrawingTruncatesLastVisibleLine attributes:attr];
    [super drawRect:dirtyRect];
}
@end