//
//  subtitles.m
//  glfwVger
//
//  Created by Roy Chen on 10/27/13.
//  Copyright (c) 2013 me. All rights reserved.
//

#include "internal.h"
#import "subtitles.h"
@implementation subtitles

int gNSStringGeometricsTypesetterBehavior = NSTypesetterLatestBehavior;

- (void)setFontSize:(CGFloat)size
{
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

- (void)setText:(SubItem*)items length:(int)len
{
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
        NSFont *font = [fontManager fontWithFamily:@"Georgia"
                                                  traits:mask
                                                  weight:0
                                                    size:self->_fontSize];
        CGFloat red = item.color&0xff0000;
        CGFloat green = item.color&0x00ff00;
        CGFloat blue = item.color&0x0000ff;
        NSColor *color = [NSColor colorWithDeviceRed:red green:green blue:blue alpha:1];
        
        NSShadow *shadow = [[NSShadow alloc] init];
        [shadow setShadowColor:color];
        [shadow setShadowBlurRadius:10];
        
        NSAttributedString *str = [[NSAttributedString alloc] initWithString:[NSString stringWithUTF8String:item.str] attributes:@{NSFontAttributeName:font,NSBackgroundColorAttributeName:[NSColor clearColor],
                                              NSForegroundColorAttributeName:color,
                                                       NSShadowAttributeName:shadow}];
        
        
        [attrStr appendAttributedString:str];
    }
    
    [self.textStorage setAttributedString:attrStr];
    
    CGFloat width = self.frame.size.width;
    CGFloat height = [self sizeForWidth:width height:FLT_MAX].height;
    
//    NSLog(@"height:%lf", height);
    [self setFrameSize:NSMakeSize(width, height)];
    [self setFrameOrigin:NSMakePoint(0, 20)];
}


- (void)mouseDown:(NSEvent *)event
{
    NSLog(@"mouseDown");
    if (self.superview != NULL) {
        NSLog(@"mouseDown1");
        [self.superview mouseDown:event];
    }
}

- (void)mouseDragged:(NSEvent *)event
{
    if (self.superview != NULL)
        [self.superview mouseDragged:event];
}

- (void)mouseUp:(NSEvent *)event
{
    if (self.superview != NULL)
        [self.superview mouseUp:event];
}

- (void)mouseMoved:(NSEvent *)event
{
    if (self.superview != NULL)
        [self.superview mouseMoved:event];
}

- (void)rightMouseDown:(NSEvent *)event
{
    if (self.superview != NULL)
        [self.superview rightMouseDown:event];
}

- (void)rightMouseDragged:(NSEvent *)event
{
    if (self.superview != NULL)
        [self.superview rightMouseDragged:event];
}

- (void)rightMouseUp:(NSEvent *)event
{
    if (self.superview != NULL)
        [self.superview rightMouseUp:event];
}

- (void)otherMouseDown:(NSEvent *)event
{
    if (self.superview != NULL)
        [self.superview otherMouseDown:event];
}

- (void)otherMouseDragged:(NSEvent *)event
{
    if (self.superview != NULL)
        [self.superview otherMouseDragged:event];
}

- (void)otherMouseUp:(NSEvent *)event
{
    if (self.superview != NULL)
        [self.superview otherMouseUp:event];
}

@end
