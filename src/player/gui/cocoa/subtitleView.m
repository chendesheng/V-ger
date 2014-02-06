#import "subtitleView.h"
@implementation SubtitleView

- (id)initWithFrame:(NSRect)frame {
    NSLog(@"popupView initWithFrame");
    self = [super initWithFrame:frame];
    if (self != nil) {
     //    [self setBorderType:NSNoBorder];
    	// [self setDrawsBackground:NO];

  //   	NSButton *button = [[NSButton alloc] initWithFrame:NSMakeRect(0,0,50,30)]; 
		// [button setTitle:@"Click me!"]; 
		// [self addSubview:button]; 


    	NSSearchField *search = [[NSSearchField alloc] initWithFrame:NSMakeRect(0,0,frame.size.width,30)]; 
		// [search setTitle:@"Click me!"]; 
		[self addSubview:search]; 
    	[search setAutoresizingMask:NSViewWidthSizable];

    	NSTableView *tv = [[NSTableView alloc] initWithFrame:NSMakeRect(0,30,frame.size.width,frame.size.height-30)];
    	[self addSubview:tv];
    	[tv setAutoresizingMask:NSViewWidthSizable];

    }
    return self;
}
@end