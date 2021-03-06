//
//  RMBlurredView.m
//
//  Created by Raffael Hannemann on 08.10.13.
//  Copyright (c) 2013 Raffael Hannemann. All rights reserved.
//

#import "blurView.h"

#define kRMBlurredViewDefaultTintColor [NSColor colorWithCalibratedWhite:1.0 alpha:0.5]
#define kRMBlurredViewDefaultSaturationFactor 2.0
#define kRMBlurredViewDefaultBlurRadius 20.0

@implementation BlurView

- (id)initWithFrame:(NSRect)frame
{
    self = [super initWithFrame:frame];
    if (self) {
        self->_height = frame.size.height;

        [self setUp];
    }
    return self;
}

- (id)initWithCoder:(NSCoder *)coder
{
    self = [super initWithCoder:coder];
    if (self) {
        [self setUp];
    }
    return self;
}

- (void) setTintColor:(NSColor *)tintColor {
    _tintColor = tintColor;
    
    // Since we need a CGColor reference, store it for the drawing of the layer.
    if (_tintColor) {
        [self.layer setBackgroundColor:_tintColor.CGColor];
    }
    
    // Trigger a re-drawing of the layer
    [self.layer setNeedsDisplay];
}

- (void) setBlurRadius:(float)blurRadius {
    // Setting the blur radius requires a resetting of the filters
    _blurRadius = blurRadius;
    [self resetFilters];
}

- (void) setBlurRadiusNumber:(NSNumber *)blurRadiusNumber {
    [self setBlurRadius:blurRadiusNumber.floatValue];
}

- (void) setSaturationFactor:(float)saturationFactor {
    // Setting the saturation factor also requires a resetting of the filters
    _saturationFactor = saturationFactor;
    [self resetFilters];
}

- (void) setSaturationFactorNumber:(NSNumber *)saturationFactorNumber {
    [self setSaturationFactor:saturationFactorNumber.floatValue];
}

- (void) setUp {
    // Instantiate a new CALayer and set it as the NSView's layer (layer-hosting)
    _hostedLayer = [CALayer layer];
    [self setWantsLayer:YES];
    [self setLayer:_hostedLayer];
    
    // Set up the default parameters
    _blurRadius = kRMBlurredViewDefaultBlurRadius;
    _saturationFactor = kRMBlurredViewDefaultSaturationFactor;
    [self setTintColor:kRMBlurredViewDefaultTintColor];
    
    // It's important to set the layer to mask to its bounds, otherwise the whole parent view might get blurred
    [self.layer setMasksToBounds:YES];

#ifdef __MAC_10_9
    // To apply CIFilters on OS X 10.9, we need to set the property accordingly:
    if ([self respondsToSelector:@selector(setLayerUsesCoreImageFilters:)]) {
        BOOL flag = YES;
        NSInvocation *inv = [NSInvocation invocationWithMethodSignature:[self methodSignatureForSelector:@selector(setLayerUsesCoreImageFilters:)]];
        [inv setSelector:@selector(setLayerUsesCoreImageFilters:)];
        [inv setArgument:&flag atIndex:2];
        [inv invokeWithTarget:self];
    }
#endif

    // Set the layer to redraw itself once it's size is changed
    [self.layer setNeedsDisplayOnBoundsChange:YES];
    
    // Initially create the filter instances
    [self resetFilters];
}

- (void) resetFilters {
    
    // To get a higher color saturation, we create a ColorControls filter
    _saturationFilter = [CIFilter filterWithName:@"CIColorControls"];
    [_saturationFilter setDefaults];
    [_saturationFilter setValue:[NSNumber numberWithFloat:_saturationFactor] forKey:@"inputSaturation"];
    
    // Next, we create the blur filter
    _blurFilter = [CIFilter filterWithName:@"CIGaussianBlur"];
    [_blurFilter setDefaults];
    [_blurFilter setValue:[NSNumber numberWithFloat:_blurRadius] forKey:@"inputRadius"];
    
    // Now we apply the two filters as the layer's background filters
    [self.layer setBackgroundFilters:@[_saturationFilter, _blurFilter]];
    
    // ... and trigger a refresh
    [self.layer setNeedsDisplay];
}
-(void)setHidden:(BOOL)flag {
    if (flag == YES) {
        [self setFrameSize:NSMakeSize(self.frame.size.width, 0)];
    }
    else {
        [self setFrameSize:NSMakeSize(self.frame.size.width, self->_height)];
    }
    
    [super setHidden:flag];
}

// wrap view make view blur background
-(BlurView*) initWithView:(NSView*)v frame:(NSRect)bounds {
    BlurView* bv = [[BlurView alloc] initWithFrame:bounds];
    v.frame = NSMakeRect(0, 0, bounds.size.width, bounds.size.height);
    [bv addSubview:v];
    [v setAutoresizingMask:NSViewWidthSizable|NSViewHeightSizable];
    return bv;
}

-(NSView*) getWrappedView {
    if (self.subviews.count > 0) {
        return [self.subviews objectAtIndex:0];
    } else {
        return nil;
    }
}
-(void) setCornerRadius:(CGFloat)r {
        self.wantsLayer = YES;
        self.layer.masksToBounds = YES;
        self.layer.cornerRadius = r;
}
@end





