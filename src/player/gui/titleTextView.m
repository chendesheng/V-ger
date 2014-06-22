#import "titleTextView.h"
@implementation TitleTextView

- (id)initWithFrame:(NSRect)frame {
    self = [super initWithFrame:frame];
	if (self) {
		self->titleString = @"";
	}

	return self;
}

- (void)drawRect:(NSRect)dirtyRect {
    // [[NSColor colorWithCalibratedRed:1 green:1 blue:1 alpha:0.9] setFill];
    // NSRectFill(dirtyRect);

    // [[NSColor colorWithCalibratedRed:1 green:1 blue:1 alpha:0.5] setFill];
    // NSRectFill(NSMakeRect(0, 0, 4, 22));

    // [[NSColor colorWithCalibratedRed:1 green:1 blue:1 alpha:0.5] setFill];
    // NSRectFill(NSMakeRect(dirtyRect.size.width-4, 0, 4, 22));


	// NSLog(@"titleTextView draw:%@", self->titleString);
    [[NSColor colorWithCalibratedRed:0 green:0 blue:0 alpha:1] setFill];
    NSDictionary *attr = @{NSFontAttributeName : [NSFont titleBarFontOfSize:13]};
    NSSize titlesz = [self->titleString sizeWithAttributes:attr];
    [self->titleString drawAtPoint:NSMakePoint((dirtyRect.size.width-titlesz.width)/2, 3) withAttributes:attr];

    [super drawRect:dirtyRect];
}
- (void)setTitle:(NSString*)title {
	self->titleString = title;

    NSDictionary *attr = @{NSFontAttributeName : [NSFont titleBarFontOfSize:13]};
    NSSize titlesz = [self->titleString sizeWithAttributes:attr];

    NSRect rt = self.bounds;
    rt.size.width = titlesz.width;
    [self setFrame:rt];
}
@end