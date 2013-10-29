//
//  subtitles.h
//  glfwVger
//
//  Created by Roy Chen on 10/27/13.
//  Copyright (c) 2013 me. All rights reserved.
//

#import <Cocoa/Cocoa.h>

@interface subtitles : NSTextView
{
 CGFloat _fontSize;
}
- (NSSize)sizeForWidth:(float)width
				height:(float)height;

- (void)setFontSize:(CGFloat)size;

- (void)setText:(SubItem*)items length:(int)len;
@end
