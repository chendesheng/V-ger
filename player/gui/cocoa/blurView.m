#import "blurView.h"
@implementation BlurView

// wrap view make view blur background
-(BlurView*) initWithView:(NSView*)v frame:(NSRect)bounds {
    BlurView* bv = [[BlurView alloc] initWithFrame:bounds];
    bv.appearance = [NSAppearance appearanceNamed:NSAppearanceNameVibrantLight];
    bv.state = NSVisualEffectStateActive;
    bv.blendingMode = NSVisualEffectBlendingModeWithinWindow;
    v.frame = NSMakeRect(0, 0, bounds.size.width, bounds.size.height);
    [bv addSubview:v];
    [v setAutoresizingMask:NSViewWidthSizable|NSViewHeightSizable];
    return bv;
}

-(void) setCornerRadius:(CGFloat)r {
        self.wantsLayer = YES;
        self.layer.masksToBounds = YES;
        self.layer.cornerRadius = r;
}


-(NSView*) getWrappedView {
    if (self.subviews.count > 0) {
        return [self.subviews objectAtIndex:0];
    } else {
        return nil;
    }
}

@end





