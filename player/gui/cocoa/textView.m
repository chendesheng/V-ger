#import "textView.h"
#include <stdlib.h>
#include <string.h>

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
        [self setFontSize:25.0];
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
        [self updateFontSize];

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
                
                answer.height += 6;
	}
	
	return answer ;
}

- (CGFloat)CalcFontsize {
    CGFloat ratio = ([[self window] frame].size.width)/(self->originalWindowWidth);
    // NSLog(@"radio:%lf", ratio);

    CGFloat s = self->_fontSize*ratio;

    if (s > 45) {
        s = 45;
    }

    if (s < 10) {
        s = 10;
    }

    return s;
}

- (void)updateFontSize {
    if ([self.textStorage length] == 0) {
        return;
    }

    NSFontManager* fontManager = [NSFontManager sharedFontManager];
    NSTextStorage* textStorage = [self textStorage];
    [textStorage beginEditing];
    [textStorage enumerateAttribute:NSFontAttributeName
                            inRange:NSMakeRange(0, [textStorage length])
                            options:0
                         usingBlock:^(id value,
                                    NSRange range,
                                    BOOL * stop) {
        NSFont * font = value;
        font = [fontManager convertFont:font
                                toSize:[self CalcFontsize]];
        if (font != nil) {
            [textStorage removeAttribute:NSFontAttributeName
                                   range:range];
            [textStorage addAttribute:NSFontAttributeName
                                value:font
                                range:range];
        }
    }];
    [textStorage endEditing];
}

- (void)setText:(AttributedString*)items length:(int)len {
    NSMutableAttributedString *attrStr = [[NSMutableAttributedString alloc] init];
    for (int i=0; i < len; i++) {
        AttributedString item = items[i];
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
        CGFloat red = ((item.color&0xff0000) >> 16)/255.0;
        CGFloat green = ((item.color&0x00ff00) >> 8)/255.0;
        CGFloat blue = (item.color&0x0000ff)/255.0;
        // NSLog(@"color:%f,%f,%f", red, green, blue);
        NSColor *color = [NSColor colorWithDeviceRed:red green:green blue:blue alpha:1];
        
        NSShadow *shadow = [[NSShadow alloc] init];
        [shadow setShadowColor:[NSColor blackColor]];
        [shadow setShadowBlurRadius:6];
       
        NSMutableParagraphStyle *paragrapStyle = NSMutableParagraphStyle.new;
        paragrapStyle.alignment = kCTTextAlignmentCenter;

        NSAttributedString *str = [[NSAttributedString alloc] initWithString:[NSString stringWithUTF8String:item.str] 
            attributes:@{
                         NSFontAttributeName:font,
              NSBackgroundColorAttributeName:[NSColor clearColor],
              NSForegroundColorAttributeName:color,
                       NSShadowAttributeName:shadow,
               NSParagraphStyleAttributeName:paragrapStyle
            }];
        
        [attrStr appendAttributedString:str];
    }
    
    [self.textStorage setAttributedString:attrStr];
}

- (BOOL)canBecomeKeyView {
    return NO;
}

- (BOOL)acceptsFirstResponder {
    return NO;
}
- (NSView*)hitTest:(NSPoint)aPoint
{
    return nil;
}
@end
