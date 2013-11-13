#import "textView.h"
@implementation TextView

int gNSStringGeometricsTypesetterBehavior = NSTypesetterLatestBehavior;

- (id)initWithFrame:(NSRect)frame {
    self = [super initWithFrame:frame];
    if (self) {
        // Initialization code here.
        [self setEditable:NO];
        [self setSelectable:NO];
        [self setBackgroundColor:[NSColor clearColor]];
        [self setAlignment:NSCenterTextAlignment];
        [self setFontSize:35.0];
        self->originalWindowWidth = 1280;//fontsize 35 in 1280 pixel
    }
    
    return self;
}

- (void)setFontSize:(CGFloat)size {
    self->_fontSize = size;
}
- (NSSize)sizeForWidth:(float)width
				height:(float)height {
	NSSize answer = NSZeroSize ;
    if ([self.textStorage length] > 0) {
		// Checking for empty string is necessary since Layout Manager will give the nominal
		// height of one line if length is 0.  Our API specifies 0.0 for an empty string.
		NSSize size = NSMakeSize(width, height) ;
		NSTextContainer *textContainer = [[NSTextContainer alloc] initWithContainerSize:size] ;
		NSTextStorage *textStorage = [[NSTextStorage alloc] initWithAttributedString:self.textStorage];
		NSLayoutManager *layoutManager = [[NSLayoutManager alloc] init] ;
		[layoutManager addTextContainer:textContainer] ;
		[textStorage addLayoutManager:layoutManager] ;
		[layoutManager setHyphenationFactor:0.0] ;
		if (gNSStringGeometricsTypesetterBehavior != NSTypesetterLatestBehavior) {
			[layoutManager setTypesetterBehavior:gNSStringGeometricsTypesetterBehavior] ;
		}
		// NSLayoutManager is lazy, so we need the following kludge to force layout:
		[layoutManager glyphRangeForTextContainer:textContainer] ;
		
		answer = [layoutManager usedRectForTextContainer:textContainer].size ;
		[textStorage release] ;
		[textContainer release] ;
		[layoutManager release] ;
		
		// In case we changed it above, set typesetterBehavior back
		// to the default value.
		gNSStringGeometricsTypesetterBehavior = NSTypesetterLatestBehavior ;
	}
	
	return answer ;
}

- (CGFloat)CalcFontsize {
    CGFloat ratio = ([[self window] frame].size.width)/(self->originalWindowWidth);
    NSLog(@"radio:%lf", ratio);

    CGFloat s = self->_fontSize*ratio;

    if (s > 50) {
        s = 50;
    }

    if (s < 10) {
        s = 10;
    }

    return s;
}

- (void)setText:(SubItem*)items length:(int)len {
    NSMutableAttributedString *attrStr = [[NSMutableAttributedString alloc] init];
    for (int i=0; i < len; i++) {
        SubItem item = items[i];
//        NSFont *font = [NSFont fontWithName:@"Georgia" size:25.0];
        NSFontTraitMask mask = 0;
        if ((item.style & 1) > 0) {
            mask = mask | NSItalicFontMask;
        }
        if ((item.style & 2) > 0) {
            mask = mask | NSBoldFontMask;
        }
        NSFontManager *fontManager = [NSFontManager sharedFontManager];
        NSFont *font = [fontManager fontWithFamily:@"Palatino"
                                                  traits:mask
                                                  weight:0
                                                    size:[self CalcFontsize]];
        CGFloat red = item.color&0xff0000;
        CGFloat green = item.color&0x00ff00;
        CGFloat blue = item.color&0x0000ff;
        NSColor *color = [NSColor colorWithDeviceRed:red green:green blue:blue alpha:1];
        
        NSShadow *shadow = [[NSShadow alloc] init];
        [shadow setShadowColor:[NSColor colorWithDeviceRed:(255-red) green:(255-green) blue:(255-blue) alpha:1]];
        [shadow setShadowBlurRadius:6];
        
        NSAttributedString *str = [[NSAttributedString alloc] initWithString:[NSString stringWithUTF8String:item.str] 
            attributes:@{NSFontAttributeName:font,
              NSBackgroundColorAttributeName:[NSColor clearColor],
              NSForegroundColorAttributeName:color,
                       NSShadowAttributeName:shadow}];
//        ,
//    NSStrokeWidthAttributeName:@-4.0,
//    NSStrokeColorAttributeName:[NSColor blackColor]
        
        [attrStr appendAttributedString:str];
    }
    
    [self.textStorage setAttributedString:attrStr];
    
    CGFloat width = self.frame.size.width;
    CGFloat height = [self sizeForWidth:width height:FLT_MAX].height;
//
//    NSLog(@"height:%lf", height);
    NSPoint pt = [self frame].origin;
    [self setFrame:NSMakeRect(pt.x,pt.y,width, height)];
}


- (void)mouseDown:(NSEvent *)event {
    if (self.superview != NULL) {
        [self.superview mouseDown:event];
    }
}

- (void)mouseDragged:(NSEvent *)event {
    if (self.superview != NULL)
        [self.superview mouseDragged:event];
}

- (void)mouseUp:(NSEvent *)event {
    if (self.superview != NULL)
        [self.superview mouseUp:event];
}

- (void)mouseMoved:(NSEvent *)event {
    if (self.superview != NULL)
        [self.superview mouseMoved:event];
}

- (BOOL)canBecomeKeyView {
    return NO;
}

- (BOOL)acceptsFirstResponder {
    return NO;
}
@end
