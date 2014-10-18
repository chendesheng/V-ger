#import "titleTextView.h"
@implementation TitleTextView

- (id)initWithFrame:(NSRect)frame {
    self = [super initWithFrame:frame];
	if (self) {
		self->_title = @"";
	}

	return self;
}

- (CGFloat)widthOfString:(NSString *)string {
     NSDictionary *attributes = [NSDictionary dictionaryWithObjectsAndKeys:[NSFont titleBarFontOfSize:13], NSFontAttributeName, nil];
     NSAttributedString *str = [[[NSAttributedString alloc] initWithString:string attributes:attributes] autorelease];
     return str.size.width;
 }

- (void)drawRect:(NSRect)dirtyRect {
    [[NSColor blackColor] setFill];

    CGFloat width = [self widthOfString:self->_title];
    // NSLog(@"width:%f", width);
    // NSLog(@"b width:%f", self.bounds.size.width-100);

    NSMutableParagraphStyle *style = [[NSParagraphStyle defaultParagraphStyle] mutableCopy];
	[style setLineBreakMode:NSLineBreakByTruncatingTail];

	CGFloat drawWidth;
    if (self.bounds.size.width-100 > width) {
		[style setAlignment:NSCenterTextAlignment];
		drawWidth = self.bounds.size.width-132;
    } else {
		drawWidth = self.bounds.size.width-86;
    }


    NSDictionary *attr = @{NSFontAttributeName : [NSFont titleBarFontOfSize:13], NSParagraphStyleAttributeName: style};

    [self->_title drawWithRect:NSMakeRect(66, 6, drawWidth, 16) options:NSStringDrawingTruncatesLastVisibleLine attributes:attr];
    [super drawRect:dirtyRect];
}
@end