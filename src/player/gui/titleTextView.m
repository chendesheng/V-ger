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
    NSDictionary *attr = @{NSFontAttributeName : [NSFont titleBarFontOfSize:13]};
    NSSize titlesz = [self->_title sizeWithAttributes:attr];
    [self->_title drawAtPoint:NSMakePoint((dirtyRect.size.width-titlesz.width)/2, 3) withAttributes:attr];

    [super drawRect:dirtyRect];
}
@end