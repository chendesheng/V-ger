//========================================================================
// GLFW 3.0 OS X - www.glfw.org
//------------------------------------------------------------------------
// Copyright (c) 2009-2010 Camilla Berglund <elmindreda@elmindreda.org>
//
// This software is provided 'as-is', without any express or implied
// warranty. In no event will the authors be held liable for any damages
// arising from the use of this software.
//
// Permission is granted to anyone to use this software for any purpose,
// including commercial applications, and to alter it and redistribute it
// freely, subject to the following restrictions:
//
// 1. The origin of this software must not be misrepresented; you must not
//    claim that you wrote the original software. If you use this software
//    in a product, an acknowledgment in the product documentation would
//    be appreciated but is not required.
//
// 2. Altered source versions must be plainly marked as such, and must not
//    be misrepresented as being the original software.
//
// 3. This notice may not be removed or altered from any source
//    distribution.
//
//========================================================================

#include "internal.h"
#include "subtitles.h"
// Needed for _NSGetProgname
#include "startupView.h"
#include <crt_externs.h>



@interface trackControl : NSView
{
@public
    NSString *leftString;
    NSString *rightString;
    CGFloat percent;
    _GLFWwindow *window;
}
@end

@implementation trackControl

-(void)drawRoundedRect:(NSRect)rect radius:(CGFloat)r{
    NSBezierPath *textViewSurround = [NSBezierPath bezierPathWithRoundedRect:rect xRadius:r yRadius:r];
    [textViewSurround fill];
}
-(void)drawRect:(NSRect)dirtyRect
{
    //    NSLog(@"draw control");
    
    CGFloat position = (dirtyRect.size.width-120)*(self->percent);
    CGFloat barHeight = 4;
    CGFloat knotHeight = 14;
    CGFloat knotWidth = 5;
    
    [[NSColor colorWithCalibratedRed:255 green:255 blue:255 alpha:0.3] setFill];
    NSRectFill(dirtyRect);
    
    CGFloat x = 8;
    if ([self->leftString length]<=5) {
        x = 22;
    }
    [self->leftString drawAtPoint:NSMakePoint(x, 18) withAttributes:@{NSFontAttributeName : [NSFont fontWithName:@"Helvetica Neue" size:12]}];
    
    [self->rightString drawAtPoint:NSMakePoint(dirtyRect.size.width-60+4, 18) withAttributes:@{NSFontAttributeName : [NSFont fontWithName:@"Helvetica Neue" size:12]}];
    
    [[NSColor colorWithCalibratedRed:0 green:0 blue:0 alpha:0.5] set];
    [self drawRoundedRect:NSMakeRect(60, (dirtyRect.size.height-barHeight)/2, dirtyRect.size.width-120, barHeight) radius:2];
    
    NSShadow* theShadow = [[NSShadow alloc] init];
    [theShadow setShadowOffset:NSMakeSize(0, 0)];
    [theShadow setShadowBlurRadius:1.0];
    
    // Use a partially transparent color for shapes that overlap.
    [theShadow setShadowColor:[[NSColor blackColor]
                               colorWithAlphaComponent:0.5]];
    
    [theShadow set];
    
    [[NSColor colorWithCalibratedRed:255 green:255 blue:255 alpha:1] setFill];
    
    [self drawRoundedRect:NSMakeRect(60, (dirtyRect.size.height-barHeight)/2, position, barHeight) radius:2];
    
    [[NSColor colorWithCalibratedRed:255 green:255 blue:255 alpha:1] setFill];
    [self drawRoundedRect:NSMakeRect(position-knotWidth/2+60, (dirtyRect.size.height-knotHeight)/2, knotWidth, knotHeight) radius:1.5];
    
    [super drawRect:dirtyRect];
}
- (void)mouseDown:(NSEvent *)event
{
    NSPoint pt = [self convertPoint:[event locationInWindow] fromView:nil];
    if (pt.x >= 60 && pt.x <= self.frame.size.width-60) {
        if (pt.y >= 10 && pt.y <= self.frame.size.height-10) {
            self->percent = (pt.x-60)/(self.frame.size.width-120);
            [self setNeedsDisplay:YES];
            
            self->window->callbacks.trackPositionChanged((GLFWwindow*)self->window, self->percent);
        }
    }
}
@end


@interface trackView : NSView
{
    CALayer *backgroundLayer;
    trackControl *control;
}

-(void)updateStatus:(NSString *)time leftTime:(NSString *)leftTime percent:(float)percent;
-(void)setWindow:(_GLFWwindow*)window;
@end

@implementation trackView

- (id)initWithFrame:(NSRect)frame
{
    self = [super initWithFrame:frame];
    if (self) {
        // Initialization code here.
        self->backgroundLayer = [CALayer layer];
        [self setLayer:self->backgroundLayer];
        [self setWantsLayer:YES];
        CIFilter *blurFilter = [CIFilter filterWithName:@"CIGaussianBlur" keysAndValues:@"inputRadius", [NSNumber numberWithFloat:20.0], nil];
        //[blurFilter setDefaults];
        
        [self->backgroundLayer setMasksToBounds:YES];
        
        [self layer].backgroundFilters = [NSArray arrayWithObject:blurFilter];
        
        self->control = [[trackControl alloc] initWithFrame:NSMakeRect(0, 0, frame.size.width, frame.size.height)];
        self->control->leftString = @"--:--:--";
        self->control->rightString = @"--:--:--";
        self->control->percent = 0;
        [self->control setAutoresizingMask:NSViewWidthSizable|NSViewHeightSizable];
        [self addSubview:self->control];
        
        self.layerContentsRedrawPolicy = NSViewLayerContentsRedrawOnSetNeedsDisplay;
    }
    
    return self;
}
-(void)setWindow:(_GLFWwindow*)window
{
    self->control->window = window;
}
-(void)setHidden:(BOOL)flag
{
    if (flag == YES)
    {
        //        [self setLayer:NULL];
        [self setFrameSize:NSMakeSize(self.frame.size.width, 0)];
    }
    else
    {
        [self setFrameSize:NSMakeSize(self.frame.size.width, 50)];
    }
    [self->control setNeedsDisplay:YES];
    [super setHidden:flag];
}
-(void)updateStatus:(NSString *)time leftTime:(NSString *)leftTime percent:(float)percent
{
    [self->control->leftString autorelease];
    self->control->leftString = time;
    [self->control->leftString retain];
    
    [self->control->rightString autorelease];
    self->control->rightString = leftTime;
    [self->control->rightString retain];
    
    self->control->percent = percent;
    
    //    NSLog(@"%@ %@ %lf", time, leftTime, percent);
    [self->control setNeedsDisplay:YES];
}
//- (void)drawRect:(NSRect)dirtyRect
//{
//    NSString *str = @"00:00:00";
//
//    NSFontManager *fontManager = [NSFontManager sharedFontManager];
//    NSFont *font = [fontManager fontWithFamily:@"Georgia"
//                                        traits:NSUnboldFontMask
//                                        weight:0
//                                          size:13];
//
//    [str drawAtPoint:NSMakePoint(10, 10) withAttributes:@{NSFontAttributeName : font, NSForegroundColorAttributeName:[NSColor blackColor]}];
//    [super drawRect:dirtyRect];
//}

- (void)mouseDragged:(NSEvent *)event
{
}

- (void)mouseUp:(NSEvent *)event
{
}

- (void)mouseMoved:(NSEvent *)event
{
}

- (void)rightMouseDown:(NSEvent *)event
{
}

- (void)rightMouseDragged:(NSEvent *)event
{
}

- (void)rightMouseUp:(NSEvent *)event
{
}

- (void)otherMouseDown:(NSEvent *)event
{
}

- (void)otherMouseDragged:(NSEvent *)event
{
}

- (void)otherMouseUp:(NSEvent *)event
{
}
@end




// Center the cursor in the view of the window
//
static void centerCursor(_GLFWwindow *window)
{
    int width, height;
    _glfwPlatformGetWindowSize(window, &width, &height);
    _glfwPlatformSetCursorPos(window, width / 2.0, height / 2.0);
}

// Update the cursor to match the specified cursor mode
//
static void setModeCursor(_GLFWwindow* window, int mode)
{
    if (mode == GLFW_CURSOR_NORMAL) {
        [[NSCursor arrowCursor] set];
        
        [window->ns.trackView setHidden:NO];
        
        NSRect frame = [window->ns.subview frame];
        frame.origin.y = 60;
        [window->ns.subview setFrameOrigin:frame.origin];
        [window->ns.subview setNeedsDisplay:YES];
    }
    else {
        [(NSCursor*) _glfw.ns.cursor set];

        [window->ns.trackView setHidden:YES];
        
        NSRect frame = [window->ns.subview frame];
        frame.origin.y = 20;
        [window->ns.subview setFrameOrigin:frame.origin];
        [window->ns.subview setNeedsDisplay:YES];
    }
}

// Enter fullscreen mode
//
static void enterFullscreenMode(_GLFWwindow* window)
{
}

// Leave fullscreen mode
//
static void leaveFullscreenMode(_GLFWwindow* window)
{
}

// Transforms the specified y-coordinate between the CG display and NS screen
// coordinate systems
//
static float transformY(float y)
{
    const float height = CGDisplayBounds(CGMainDisplayID()).size.height;
    return height - y;
}

// Returns the backing rect of the specified window
//
static NSRect convertRectToBacking(_GLFWwindow* window, NSRect contentRect)
{
#if MAC_OS_X_VERSION_MAX_ALLOWED >= 1070
    if (floor(NSAppKitVersionNumber) >= NSAppKitVersionNumber10_7)
        return [window->ns.view convertRectToBacking:contentRect];
    else
#endif /*MAC_OS_X_VERSION_MAX_ALLOWED*/
        return contentRect;
}


//------------------------------------------------------------------------
// Delegate for window related notifications
//------------------------------------------------------------------------

@interface GLFWWindowDelegate : NSObject
{
    _GLFWwindow* window;
}

- (id)initWithGlfwWindow:(_GLFWwindow *)initWndow;

@end

@implementation GLFWWindowDelegate

- (id)initWithGlfwWindow:(_GLFWwindow *)initWindow
{
    self = [super init];
    if (self != nil)
        window = initWindow;

    return self;
}

- (BOOL)windowShouldClose:(id)sender
{
    _glfwInputWindowCloseRequest(window);
    return NO;
}

- (void)windowDidResize:(NSNotification *)notification
{
    // [window->nsgl.context update];

    const NSRect contentRect = [window->ns.view frame];
    const NSRect fbRect = convertRectToBacking(window, contentRect);

    _glfwInputFramebufferSize(window, fbRect.size.width, fbRect.size.height);
    _glfwInputWindowSize(window, contentRect.size.width, contentRect.size.height);
    _glfwInputWindowDamage(window);

    if (window->cursorMode == GLFW_CURSOR_DISABLED)
        centerCursor(window);
}

- (void)windowDidMove:(NSNotification *)notification
{
    // [window->nsgl.context update];

    int x, y;
    _glfwPlatformGetWindowPos(window, &x, &y);
    _glfwInputWindowPos(window, x, y);

    if (window->cursorMode == GLFW_CURSOR_DISABLED)
        centerCursor(window);
}

- (void)windowDidMiniaturize:(NSNotification *)notification
{
    _glfwInputWindowIconify(window, GL_TRUE);
}

- (void)windowDidDeminiaturize:(NSNotification *)notification
{
    if (window->monitor)
        enterFullscreenMode(window);

    _glfwInputWindowIconify(window, GL_FALSE);
}

- (void)windowDidBecomeKey:(NSNotification *)notification
{
    _glfwInputWindowFocus(window, GL_TRUE);
    _glfwPlatformSetCursorMode(window, window->cursorMode);
}

- (void)windowDidResignKey:(NSNotification *)notification
{
    _glfwInputWindowFocus(window, GL_FALSE);
    _glfwPlatformSetCursorMode(window, GLFW_CURSOR_NORMAL);
}
- (BOOL)windowShouldZoom:(NSWindow *)sender toFrame:(NSRect)newFrame
{
    NSLog(@"windowShouldZoom:%lf,%lf", newFrame.size.width, newFrame.size.height);
    return YES;
}
- (void)windowWillEnterFullScreen:(NSNotification *)notification
{
    [window->ns.subview setFontSize:50.0];
    
    NSLog(@"windowWillEnterFullScreen");
}
- (void)windowDidEnterFullScreen:(NSNotification *)notification
{
    NSLog(@"windowDidEnterFullScreen");
}
- (void)windowWillExitFullScreen:(NSNotification *)notification
{
    [window->ns.subview setFontSize:35.0];
    NSLog(@"windowWillExitFullScreen");
}
- (void)windowDidExitFullScreen:(NSNotification *)notification
{
    NSLog(@"windowDidExitFullScreen");
}

@end


//------------------------------------------------------------------------
// Delegate for application related notifications
//------------------------------------------------------------------------

@interface GLFWApplicationDelegate : NSObject
@end

@implementation GLFWApplicationDelegate

- (NSApplicationTerminateReply)applicationShouldTerminate:(NSApplication *)sender
{
    _GLFWwindow* window;

    for (window = _glfw.windowListHead;  window;  window = window->next)
        _glfwInputWindowCloseRequest(window);

    return NSTerminateCancel;
}

- (void)applicationDidHide:(NSNotification *)notification
{
    _GLFWwindow* window;

    for (window = _glfw.windowListHead;  window;  window = window->next)
        _glfwInputWindowVisibility(window, GL_FALSE);
}

- (void)applicationDidUnhide:(NSNotification *)notification
{
    _GLFWwindow* window;

    for (window = _glfw.windowListHead;  window;  window = window->next)
    {
        if ([window->ns.object isVisible])
            _glfwInputWindowVisibility(window, GL_TRUE);
    }
}
@end

// Translates OS X key modifiers into GLFW ones
//
static int translateFlags(NSUInteger flags)
{
    int mods = 0;

    if (flags & NSShiftKeyMask)
        mods |= GLFW_MOD_SHIFT;
    if (flags & NSControlKeyMask)
        mods |= GLFW_MOD_CONTROL;
    if (flags & NSAlternateKeyMask)
        mods |= GLFW_MOD_ALT;
    if (flags & NSCommandKeyMask)
        mods |= GLFW_MOD_SUPER;

    return mods;
}

// Translates a OS X keycode to a GLFW keycode
//
static int translateKey(unsigned int key)
{
    // Keyboard symbol translation table
    // TODO: Need to find mappings for F13-F15, volume down/up/mute, and eject.
    static const unsigned int table[128] =
    {
        /* 00 */ GLFW_KEY_A,
        /* 01 */ GLFW_KEY_S,
        /* 02 */ GLFW_KEY_D,
        /* 03 */ GLFW_KEY_F,
        /* 04 */ GLFW_KEY_H,
        /* 05 */ GLFW_KEY_G,
        /* 06 */ GLFW_KEY_Z,
        /* 07 */ GLFW_KEY_X,
        /* 08 */ GLFW_KEY_C,
        /* 09 */ GLFW_KEY_V,
        /* 0a */ GLFW_KEY_GRAVE_ACCENT,
        /* 0b */ GLFW_KEY_B,
        /* 0c */ GLFW_KEY_Q,
        /* 0d */ GLFW_KEY_W,
        /* 0e */ GLFW_KEY_E,
        /* 0f */ GLFW_KEY_R,
        /* 10 */ GLFW_KEY_Y,
        /* 11 */ GLFW_KEY_T,
        /* 12 */ GLFW_KEY_1,
        /* 13 */ GLFW_KEY_2,
        /* 14 */ GLFW_KEY_3,
        /* 15 */ GLFW_KEY_4,
        /* 16 */ GLFW_KEY_6,
        /* 17 */ GLFW_KEY_5,
        /* 18 */ GLFW_KEY_EQUAL,
        /* 19 */ GLFW_KEY_9,
        /* 1a */ GLFW_KEY_7,
        /* 1b */ GLFW_KEY_MINUS,
        /* 1c */ GLFW_KEY_8,
        /* 1d */ GLFW_KEY_0,
        /* 1e */ GLFW_KEY_RIGHT_BRACKET,
        /* 1f */ GLFW_KEY_O,
        /* 20 */ GLFW_KEY_U,
        /* 21 */ GLFW_KEY_LEFT_BRACKET,
        /* 22 */ GLFW_KEY_I,
        /* 23 */ GLFW_KEY_P,
        /* 24 */ GLFW_KEY_ENTER,
        /* 25 */ GLFW_KEY_L,
        /* 26 */ GLFW_KEY_J,
        /* 27 */ GLFW_KEY_APOSTROPHE,
        /* 28 */ GLFW_KEY_K,
        /* 29 */ GLFW_KEY_SEMICOLON,
        /* 2a */ GLFW_KEY_BACKSLASH,
        /* 2b */ GLFW_KEY_COMMA,
        /* 2c */ GLFW_KEY_SLASH,
        /* 2d */ GLFW_KEY_N,
        /* 2e */ GLFW_KEY_M,
        /* 2f */ GLFW_KEY_PERIOD,
        /* 30 */ GLFW_KEY_TAB,
        /* 31 */ GLFW_KEY_SPACE,
        /* 32 */ GLFW_KEY_WORLD_1,
        /* 33 */ GLFW_KEY_BACKSPACE,
        /* 34 */ GLFW_KEY_UNKNOWN,
        /* 35 */ GLFW_KEY_ESCAPE,
        /* 36 */ GLFW_KEY_RIGHT_SUPER,
        /* 37 */ GLFW_KEY_LEFT_SUPER,
        /* 38 */ GLFW_KEY_LEFT_SHIFT,
        /* 39 */ GLFW_KEY_CAPS_LOCK,
        /* 3a */ GLFW_KEY_LEFT_ALT,
        /* 3b */ GLFW_KEY_LEFT_CONTROL,
        /* 3c */ GLFW_KEY_RIGHT_SHIFT,
        /* 3d */ GLFW_KEY_RIGHT_ALT,
        /* 3e */ GLFW_KEY_RIGHT_CONTROL,
        /* 3f */ GLFW_KEY_UNKNOWN, /* Function */
        /* 40 */ GLFW_KEY_F17,
        /* 41 */ GLFW_KEY_KP_DECIMAL,
        /* 42 */ GLFW_KEY_UNKNOWN,
        /* 43 */ GLFW_KEY_KP_MULTIPLY,
        /* 44 */ GLFW_KEY_UNKNOWN,
        /* 45 */ GLFW_KEY_KP_ADD,
        /* 46 */ GLFW_KEY_UNKNOWN,
        /* 47 */ GLFW_KEY_NUM_LOCK, /* Really KeypadClear... */
        /* 48 */ GLFW_KEY_UNKNOWN, /* VolumeUp */
        /* 49 */ GLFW_KEY_UNKNOWN, /* VolumeDown */
        /* 4a */ GLFW_KEY_UNKNOWN, /* Mute */
        /* 4b */ GLFW_KEY_KP_DIVIDE,
        /* 4c */ GLFW_KEY_KP_ENTER,
        /* 4d */ GLFW_KEY_UNKNOWN,
        /* 4e */ GLFW_KEY_KP_SUBTRACT,
        /* 4f */ GLFW_KEY_F18,
        /* 50 */ GLFW_KEY_F19,
        /* 51 */ GLFW_KEY_KP_EQUAL,
        /* 52 */ GLFW_KEY_KP_0,
        /* 53 */ GLFW_KEY_KP_1,
        /* 54 */ GLFW_KEY_KP_2,
        /* 55 */ GLFW_KEY_KP_3,
        /* 56 */ GLFW_KEY_KP_4,
        /* 57 */ GLFW_KEY_KP_5,
        /* 58 */ GLFW_KEY_KP_6,
        /* 59 */ GLFW_KEY_KP_7,
        /* 5a */ GLFW_KEY_F20,
        /* 5b */ GLFW_KEY_KP_8,
        /* 5c */ GLFW_KEY_KP_9,
        /* 5d */ GLFW_KEY_UNKNOWN,
        /* 5e */ GLFW_KEY_UNKNOWN,
        /* 5f */ GLFW_KEY_UNKNOWN,
        /* 60 */ GLFW_KEY_F5,
        /* 61 */ GLFW_KEY_F6,
        /* 62 */ GLFW_KEY_F7,
        /* 63 */ GLFW_KEY_F3,
        /* 64 */ GLFW_KEY_F8,
        /* 65 */ GLFW_KEY_F9,
        /* 66 */ GLFW_KEY_UNKNOWN,
        /* 67 */ GLFW_KEY_F11,
        /* 68 */ GLFW_KEY_UNKNOWN,
        /* 69 */ GLFW_KEY_PRINT_SCREEN,
        /* 6a */ GLFW_KEY_F16,
        /* 6b */ GLFW_KEY_F14,
        /* 6c */ GLFW_KEY_UNKNOWN,
        /* 6d */ GLFW_KEY_F10,
        /* 6e */ GLFW_KEY_UNKNOWN,
        /* 6f */ GLFW_KEY_F12,
        /* 70 */ GLFW_KEY_UNKNOWN,
        /* 71 */ GLFW_KEY_F15,
        /* 72 */ GLFW_KEY_INSERT, /* Really Help... */
        /* 73 */ GLFW_KEY_HOME,
        /* 74 */ GLFW_KEY_PAGE_UP,
        /* 75 */ GLFW_KEY_DELETE,
        /* 76 */ GLFW_KEY_F4,
        /* 77 */ GLFW_KEY_END,
        /* 78 */ GLFW_KEY_F2,
        /* 79 */ GLFW_KEY_PAGE_DOWN,
        /* 7a */ GLFW_KEY_F1,
        /* 7b */ GLFW_KEY_LEFT,
        /* 7c */ GLFW_KEY_RIGHT,
        /* 7d */ GLFW_KEY_DOWN,
        /* 7e */ GLFW_KEY_UP,
        /* 7f */ GLFW_KEY_UNKNOWN,
    };

    if (key >= 128)
        return GLFW_KEY_UNKNOWN;

    return table[key];
}


//------------------------------------------------------------------------
// Content view class for the GLFW window
//------------------------------------------------------------------------

@interface GLFWContentView : NSOpenGLView
{
    _GLFWwindow* window;
    NSTrackingArea* trackingArea;
}

- (id)initWithGlfwWindow:(_GLFWwindow *)initWindow;

@end

@implementation GLFWContentView

+ (void)initialize
{
    if (self == [GLFWContentView class])
    {
        if (_glfw.ns.cursor == nil)
        {
            NSImage* data = [[NSImage alloc] initWithSize:NSMakeSize(1, 1)];
            _glfw.ns.cursor = [[NSCursor alloc] initWithImage:data
                                                      hotSpot:NSZeroPoint];
            [data release];
        }
    }
}

- (id)initWithGlfwWindow:(_GLFWwindow *)initWindow
{
    // NSLog(@"initWithGlfwWindow");
    self = [super init];
    if (self != nil)
    {
        window = initWindow;
        trackingArea = nil;

        [self updateTrackingAreas];
    }

    NSOpenGLPixelFormatAttribute attrs[] = {
 
        // Specifying "NoRecovery" gives us a context that cannot fall back to the software renderer.  This makes the View-based context a compatible with the layer-backed context, enabling us to use the "shareContext" feature to share textures, display lists, and other OpenGL objects between the two.
        NSOpenGLPFANoRecovery, // Enable automatic use of OpenGL "share" contexts.
 
        NSOpenGLPFAColorSize, 24,
        NSOpenGLPFAAlphaSize, 8,
        NSOpenGLPFADepthSize, 16,
        NSOpenGLPFADoubleBuffer,
        NSOpenGLPFAAccelerated,
        0
    };
    // Create our pixel format.
    NSOpenGLPixelFormat* pixelFormat = [[NSOpenGLPixelFormat alloc] initWithAttributes:attrs];
    // NSRect rt = [initWindow->ns.object frame];
    // NSLog(@"%f %f", rt.size.width, rt.size.height);
    [super initWithFrame:[initWindow->ns.object frame] pixelFormat:pixelFormat];
    window->nsgl.pixelFormat = [window->ns.view pixelFormat];
    // [pixelFormat release];

    // NSLog(@"initWithGlfwWindow end");

    return self;
}
// AppKit automatically invokes the NSOpenGLView's -prepareOpenGL once when a new NSOpenGLContext becomes current.
- (void)prepareOpenGL {
    // [scene prepareOpenGL];
    // NSLog(@"prepareOpenGL");
    window->nsgl.context = [window->ns.view openGLContext];
    // window->nsgl.pixelFormat = [window->ns.view pixelFormat];
}

-(void)dealloc
{
    [trackingArea release];
    [super dealloc];
}

- (BOOL)isOpaque
{
    return YES;
}

- (BOOL)canBecomeKeyView
{
    return YES;
}

- (BOOL)acceptsFirstResponder
{
    return YES;
}

- (void)cursorUpdate:(NSEvent *)event
{
    setModeCursor(window, window->cursorMode);
}

- (void)mouseDown:(NSEvent *)event
{
    _glfwInputMouseClick(window,
                         GLFW_MOUSE_BUTTON_LEFT,
                         GLFW_PRESS,
                         translateFlags([event modifierFlags]));
}

- (void)mouseDragged:(NSEvent *)event
{
    [self mouseMoved:event];
}

- (void)mouseUp:(NSEvent *)event
{
    _glfwInputMouseClick(window,
                         GLFW_MOUSE_BUTTON_LEFT,
                         GLFW_RELEASE,
                         translateFlags([event modifierFlags]));
}

- (void)mouseMoved:(NSEvent *)event
{
    if (window->cursorMode == GLFW_CURSOR_DISABLED)
        _glfwInputCursorMotion(window, [event deltaX], [event deltaY]);
    else
    {
        const NSRect contentRect = [window->ns.view frame];
        const NSPoint p = [event locationInWindow];

        _glfwInputCursorMotion(window, p.x, contentRect.size.height - p.y);
    }
}

- (void)rightMouseDown:(NSEvent *)event
{
    _glfwInputMouseClick(window,
                         GLFW_MOUSE_BUTTON_RIGHT,
                         GLFW_PRESS,
                         translateFlags([event modifierFlags]));
}

- (void)rightMouseDragged:(NSEvent *)event
{
    [self mouseMoved:event];
}

- (void)rightMouseUp:(NSEvent *)event
{
    _glfwInputMouseClick(window,
                         GLFW_MOUSE_BUTTON_RIGHT,
                         GLFW_RELEASE,
                         translateFlags([event modifierFlags]));
}

- (void)otherMouseDown:(NSEvent *)event
{
    _glfwInputMouseClick(window,
                         (int)[event buttonNumber],
                         GLFW_PRESS,
                         translateFlags([event modifierFlags]));
}

- (void)otherMouseDragged:(NSEvent *)event
{
    [self mouseMoved:event];
}

- (void)otherMouseUp:(NSEvent *)event
{
    _glfwInputMouseClick(window,
                         (int)[event buttonNumber],
                         GLFW_RELEASE,
                         translateFlags([event modifierFlags]));
}

- (void)mouseExited:(NSEvent *)event
{
    _glfwInputCursorEnter(window, GL_FALSE);
}

- (void)mouseEntered:(NSEvent *)event
{
    _glfwInputCursorEnter(window, GL_TRUE);
}

- (void)viewDidChangeBackingProperties
{
    const NSRect contentRect = [window->ns.view frame];
    const NSRect fbRect = convertRectToBacking(window, contentRect);

    _glfwInputFramebufferSize(window, fbRect.size.width, fbRect.size.height);
}

- (void)updateTrackingAreas
{
    if (trackingArea != nil)
    {
        [self removeTrackingArea:trackingArea];
        [trackingArea release];
    }

    NSTrackingAreaOptions options = NSTrackingMouseEnteredAndExited |
                                    NSTrackingActiveInKeyWindow |
                                    NSTrackingCursorUpdate |
                                    NSTrackingInVisibleRect;

    trackingArea = [[NSTrackingArea alloc] initWithRect:[self bounds]
                                                options:options
                                                  owner:self
                                               userInfo:nil];

    [self addTrackingArea:trackingArea];
    [super updateTrackingAreas];
}

- (void)keyDown:(NSEvent *)event
{
    const int key = translateKey([event keyCode]);
    const int mods = translateFlags([event modifierFlags]);
    _glfwInputKey(window, key, [event keyCode], GLFW_PRESS, mods);

    NSString* characters = [event characters];
    NSUInteger i, length = [characters length];

    for (i = 0;  i < length;  i++)
        _glfwInputChar(window, [characters characterAtIndex:i]);

    if (key == GLFW_KEY_ESCAPE && (([window->ns.object styleMask] & NSFullScreenWindowMask) == NSFullScreenWindowMask)) {
        [window->ns.object toggleFullScreen:nil];
    }
}

- (void)flagsChanged:(NSEvent *)event
{
    int action;
    unsigned int newModifierFlags =
        [event modifierFlags] & NSDeviceIndependentModifierFlagsMask;

    if (newModifierFlags > window->ns.modifierFlags)
        action = GLFW_PRESS;
    else
        action = GLFW_RELEASE;

    window->ns.modifierFlags = newModifierFlags;

    const int key = translateKey([event keyCode]);
    const int mods = translateFlags([event modifierFlags]);
    _glfwInputKey(window, key, [event keyCode], action, mods);
}

- (void)keyUp:(NSEvent *)event
{
    const int key = translateKey([event keyCode]);
    const int mods = translateFlags([event modifierFlags]);
    _glfwInputKey(window, key, [event keyCode], GLFW_RELEASE, mods);
}

- (void)scrollWheel:(NSEvent *)event
{
    double deltaX, deltaY;

#if MAC_OS_X_VERSION_MAX_ALLOWED >= 1070
    if (floor(NSAppKitVersionNumber) >= NSAppKitVersionNumber10_7)
    {
        deltaX = [event scrollingDeltaX];
        deltaY = [event scrollingDeltaY];

        if ([event hasPreciseScrollingDeltas])
        {
            deltaX *= 0.1;
            deltaY *= 0.1;
        }
    }
    else
#endif /*MAC_OS_X_VERSION_MAX_ALLOWED*/
    {
        deltaX = [event deltaX];
        deltaY = [event deltaY];
    }

    if (fabs(deltaX) > 0.0 || fabs(deltaY) > 0.0)
        _glfwInputScroll(window, deltaX, deltaY);
}
- (void)timerTick:(NSEvent *)event
{
    _glfwTimer(window);

    // [self setNeedsDisplay:YES];
}

- (void)drawRect:(NSRect)dirtyRect
{
    // _glfwTimer(window);
//    
//	NSBezierPath* thePath = [NSBezierPath bezierPath];
//    [thePath appendBezierPathWithRoundedRect:dirtyRect xRadius:10 yRadius:10];
//    [thePath fill];

    _glfwDraw(window);
}

@end


//------------------------------------------------------------------------
// GLFW window class
//------------------------------------------------------------------------

@interface GLFWWindow : NSWindow {}
@end

@implementation GLFWWindow

- (BOOL)canBecomeKeyWindow
{
    // Required for NSBorderlessWindowMask windows
    return YES;
}
@end


//------------------------------------------------------------------------
// GLFW application class
//------------------------------------------------------------------------

@interface GLFWApplication : NSApplication
@end

@implementation GLFWApplication

// From http://cocoadev.com/index.pl?GameKeyboardHandlingAlmost
// This works around an AppKit bug, where key up events while holding
// down the command key don't get sent to the key window.
- (void)sendEvent:(NSEvent *)event
{
    if ([event type] == NSKeyUp && ([event modifierFlags] & NSCommandKeyMask))
        [[self keyWindow] sendEvent:event];
    else
        [super sendEvent:event];
}

@end

#if defined(_GLFW_USE_MENUBAR)

// Try to figure out what the calling application is called
//
static NSString* findAppName(void)
{
    size_t i;
    NSDictionary* infoDictionary = [[NSBundle mainBundle] infoDictionary];

    // Keys to search for as potential application names
    NSString* GLFWNameKeys[] =
    {
        @"CFBundleDisplayName",
        @"CFBundleName",
        @"CFBundleExecutable",
    };

    for (i = 0;  i < sizeof(GLFWNameKeys) / sizeof(GLFWNameKeys[0]);  i++)
    {
        id name = [infoDictionary objectForKey:GLFWNameKeys[i]];
        if (name &&
            [name isKindOfClass:[NSString class]] &&
            ![name isEqualToString:@""])
        {
            return name;
        }
    }

    char** progname = _NSGetProgname();
    if (progname && *progname)
        return [NSString stringWithUTF8String:*progname];

    // Really shouldn't get here
    return @"GLFW Application";
}

// Set up the menu bar (manually)
// This is nasty, nasty stuff -- calls to undocumented semi-private APIs that
// could go away at any moment, lots of stuff that really should be
// localize(d|able), etc.  Loading a nib would save us this horror, but that
// doesn't seem like a good thing to require of GLFW's clients.
//
static void createMenuBar(void)
{
    NSString* appName = findAppName();

    NSMenu* bar = [[NSMenu alloc] init];
    [NSApp setMainMenu:bar];

    NSMenuItem* appMenuItem =
        [bar addItemWithTitle:@"" action:NULL keyEquivalent:@""];
    NSMenu* appMenu = [[NSMenu alloc] init];
    [appMenuItem setSubmenu:appMenu];

    [appMenu addItemWithTitle:[NSString stringWithFormat:@"About %@", appName]
                       action:@selector(orderFrontStandardAboutPanel:)
                keyEquivalent:@""];
    [appMenu addItem:[NSMenuItem separatorItem]];
    NSMenu* servicesMenu = [[NSMenu alloc] init];
    [NSApp setServicesMenu:servicesMenu];
    [[appMenu addItemWithTitle:@"Services"
                       action:NULL
                keyEquivalent:@""] setSubmenu:servicesMenu];
    [appMenu addItem:[NSMenuItem separatorItem]];
    [appMenu addItemWithTitle:[NSString stringWithFormat:@"Hide %@", appName]
                       action:@selector(hide:)
                keyEquivalent:@"h"];
    [[appMenu addItemWithTitle:@"Hide Others"
                       action:@selector(hideOtherApplications:)
                keyEquivalent:@"h"]
        setKeyEquivalentModifierMask:NSAlternateKeyMask | NSCommandKeyMask];
    [appMenu addItemWithTitle:@"Show All"
                       action:@selector(unhideAllApplications:)
                keyEquivalent:@""];
    [appMenu addItem:[NSMenuItem separatorItem]];
    [appMenu addItemWithTitle:[NSString stringWithFormat:@"Quit %@", appName]
                       action:@selector(terminate:)
                keyEquivalent:@"q"];

    NSMenuItem* windowMenuItem =
        [bar addItemWithTitle:@"" action:NULL keyEquivalent:@""];
    NSMenu* windowMenu = [[NSMenu alloc] initWithTitle:@"Window"];
    [NSApp setWindowsMenu:windowMenu];
    [windowMenuItem setSubmenu:windowMenu];

    [windowMenu addItemWithTitle:@"Minimize"
                          action:@selector(performMiniaturize:)
                   keyEquivalent:@"m"];
    [windowMenu addItemWithTitle:@"Zoom"
                          action:@selector(performZoom:)
                   keyEquivalent:@""];
    [windowMenu addItem:[NSMenuItem separatorItem]];
    [windowMenu addItemWithTitle:@"Bring All to Front"
                          action:@selector(arrangeInFront:)
                   keyEquivalent:@""];

    // Prior to Snow Leopard, we need to use this oddly-named semi-private API
    // to get the application menu working properly.
    [NSApp performSelector:@selector(setAppleMenu:) withObject:appMenu];
}

#endif /* _GLFW_USE_MENUBAR */

// Initialize the Cocoa Application Kit
//
static GLboolean initializeAppKit(void)
{
    if (NSApp)
        return GL_TRUE;

    // Implicitly create shared NSApplication instance
    [GLFWApplication sharedApplication];

    // If we get here, the application is unbundled
    ProcessSerialNumber psn = { 0, kCurrentProcess };
    TransformProcessType(&psn, kProcessTransformToForegroundApplication);

    // Having the app in front of the terminal window is also generally
    // handy.  There is an NSApplication API to do this, but...
    SetFrontProcess(&psn);

#if defined(_GLFW_USE_MENUBAR)
    // Menu bar setup must go between sharedApplication above and
    // finishLaunching below, in order to properly emulate the behavior
    // of NSApplicationMain
    createMenuBar();
#endif

    [NSApp finishLaunching];

    return GL_TRUE;
}

// Create the Cocoa window
//
static GLboolean createWindow(_GLFWwindow* window,
                              const _GLFWwndconfig* wndconfig)
{
    // NSLog(@"createWindow");

    unsigned int styleMask = 0;

    if (wndconfig->monitor || !wndconfig->decorated)
        styleMask = NSBorderlessWindowMask;
    else
    {
        styleMask = NSTitledWindowMask | NSClosableWindowMask |
                    NSMiniaturizableWindowMask;

        if (wndconfig->resizable)
            styleMask |= NSResizableWindowMask;
    }
//    styleMask = NSBorderlessWindowMask;

    window->ns.object = [[GLFWWindow alloc]
        initWithContentRect:NSMakeRect(0, 0, wndconfig->width, wndconfig->height)
                  styleMask:styleMask
                    backing:NSBackingStoreBuffered
                      defer:NO];
    [window->ns.object setContentAspectRatio:NSMakeSize(wndconfig->width, wndconfig->height)];
    
    [window->ns.object setOpaque:NO];
    [window->ns.object setBackgroundColor:[NSColor clearColor]];
    [window->ns.object setHasShadow:YES];

    double minHeight = 500 * wndconfig->height / wndconfig->width;
    [window->ns.object setContentMinSize:NSMakeSize(500, minHeight)];
    if (window->ns.object == nil)
    {
        _glfwInputError(GLFW_PLATFORM_ERROR, "Cocoa: Failed to create window");
        return GL_FALSE;
    }

    window->ns.view = [[GLFWContentView alloc] initWithGlfwWindow:window];




// #if MAC_OS_X_VERSION_MAX_ALLOWED >= 1070
//     if (floor(NSAppKitVersionNumber) >= NSAppKitVersionNumber10_7)
//         [window->ns.view setWantsBestResolutionOpenGLSurface:YES];
// #endif /*MAC_OS_X_VERSION_MAX_ALLOWED*/

    [window->ns.object setTitle:[NSString stringWithUTF8String:wndconfig->title]];
    [window->ns.object setContentView:window->ns.view];
    [window->ns.object setDelegate:window->ns.delegate];
    [window->ns.object setAcceptsMouseMovedEvents:YES];
    [window->ns.object center];

#if MAC_OS_X_VERSION_MAX_ALLOWED >= 1070
    if (floor(NSAppKitVersionNumber) >= NSAppKitVersionNumber10_7) {
        [window->ns.object setRestorable:NO];

        if (wndconfig->resizable) {
            [window->ns.object setCollectionBehavior:NSWindowCollectionBehaviorFullScreenPrimary];
        }
    }
#endif /*MAC_OS_X_VERSION_MAX_ALLOWED*/

    NSTimer *renderTimer = [NSTimer timerWithTimeInterval:1.0/60.0 
                            target:window->ns.view
                          selector:@selector(timerTick:)
                          userInfo:nil
                           repeats:YES];

    [[NSRunLoop currentRunLoop] addTimer:renderTimer
                                forMode:NSDefaultRunLoopMode];
    [[NSRunLoop currentRunLoop] addTimer:renderTimer
                                forMode:NSEventTrackingRunLoopMode]; //Ensure timer fires during resize

    // NSLog(@"createWindow end");
    return GL_TRUE;
}


//////////////////////////////////////////////////////////////////////////
//////                       GLFW platform API                      //////
//////////////////////////////////////////////////////////////////////////

int _glfwPlatformCreateWindow(_GLFWwindow* window,
                              const _GLFWwndconfig* wndconfig,
                              const _GLFWfbconfig* fbconfig)
{
    if (!initializeAppKit())
        return GL_FALSE;

    // There can only be one application delegate, but we allocate it the
    // first time a window is created to keep all window code in this file
    if (_glfw.ns.delegate == nil)
    {
        _glfw.ns.delegate = [[GLFWApplicationDelegate alloc] init];
        if (_glfw.ns.delegate == nil)
        {
            _glfwInputError(GLFW_PLATFORM_ERROR,
                            "Cocoa: Failed to create application delegate");
            return GL_FALSE;
        }

        [NSApp setDelegate:_glfw.ns.delegate];
    }

    window->ns.delegate = [[GLFWWindowDelegate alloc] initWithGlfwWindow:window];
    if (window->ns.delegate == nil)
    {
        _glfwInputError(GLFW_PLATFORM_ERROR,
                        "Cocoa: Failed to create window delegate");
        return GL_FALSE;
    }

    // Don't use accumulation buffer support; it's not accelerated
    // Aux buffers probably aren't accelerated either

    if (!createWindow(window, wndconfig))
        return GL_FALSE;

    if (!_glfwCreateContext(window, wndconfig, fbconfig))
        return GL_FALSE;


    // [window->nsgl.context setView:window->ns.view];
    NSRect frame = [window->ns.object frame];
    //osx >= 10.5
    [window->ns.view setWantsLayer:YES];
//    
//    [[window->ns.view layer] setCornerRadius:10];
//    [[window->ns.view layer] masksToBounds];
    
    
    subtitles *subView = [[subtitles alloc] initWithFrame:NSMakeRect(0, 0, frame.size.width, 0)];
    [window->ns.view addSubview:subView];
    [subView setAutoresizingMask:NSViewWidthSizable];

    window->ns.subview = subView;
    
    
    subtitles *subWithPosition = [[subtitles alloc] initWithFrame:NSMakeRect(100, 100, 100, 0)];
//    [subWithPosition setHidden:YES];
//    SubItem item;
//    item.str = "hello";
//    item.color = 0xffffff;
//    [subWithPosition setFontSize:30];
//    [subWithPosition setText:&item length:1];
    [window->ns.view addSubview:subWithPosition];
    window->ns.subWithPosition = subWithPosition;
    
    
    trackView *track = [[trackView alloc] initWithFrame:NSMakeRect(0, 0, frame.size.width, 50.0)];
    [track setAutoresizingMask:NSViewWidthSizable];
    [window->ns.view addSubview:track];
    
    window->ns.trackView = track;
    [track setWindow:window];
    
    // [window->ns.view setAutoresizesSubviews:YES];
    // [subtitles setAutoresizingMask:NSViewWidthSizable];
    // [subtitles scaleUnitSquareToSize:NSMakeSize(1.5,1.5)];//double
    startupView *startup = [[startupView alloc] initWithFrame:NSMakeRect(0, 0, frame.size.width, frame.size.height)];
    [startup setAutoresizingMask:NSViewWidthSizable|NSViewHeightSizable];
    [window->ns.view addSubview:startup];
    
    window->ns.startupView = startup;
        
    
    NSLog(@"platform create window");
    
    return GL_TRUE;
}

void _glfwPlatformDestroyWindow(_GLFWwindow* window)
{
    [window->ns.object orderOut:nil];

    if (window->monitor)
        leaveFullscreenMode(window);

    _glfwDestroyContext(window);

    [window->ns.object setDelegate:nil];
    [window->ns.delegate release];
    window->ns.delegate = nil;

    [window->ns.view release];
    window->ns.view = nil;

    [window->ns.object close];
    window->ns.object = nil;
}

void _glfwPlatformSetWindowTitle(_GLFWwindow* window, const char *title)
{
    [window->ns.object setTitle:[NSString stringWithUTF8String:title]];
}

void _glfwPlatformGetWindowPos(_GLFWwindow* window, int* xpos, int* ypos)
{
    const NSRect contentRect =
        [window->ns.object contentRectForFrameRect:[window->ns.object frame]];

    if (xpos)
        *xpos = contentRect.origin.x;
    if (ypos)
        *ypos = transformY(contentRect.origin.y + contentRect.size.height);
}

void _glfwPlatformSetWindowPos(_GLFWwindow* window, int x, int y)
{
    const NSRect contentRect = [window->ns.view frame];
    const NSRect dummyRect = NSMakeRect(x, transformY(y + contentRect.size.height), 0, 0);
    const NSRect frameRect = [window->ns.object frameRectForContentRect:dummyRect];
    [window->ns.object setFrameOrigin:frameRect.origin];
}

void _glfwPlatformGetWindowSize(_GLFWwindow* window, int* width, int* height)
{
    const NSRect contentRect = [window->ns.view frame];

    if (width)
        *width = contentRect.size.width;
    if (height)
        *height = contentRect.size.height;
}

void _glfwPlatformSetWindowSize(_GLFWwindow* window, int width, int height)
{
    [window->ns.object setContentSize:NSMakeSize(width, height)];
}

void _glfwPlatformGetFramebufferSize(_GLFWwindow* window, int* width, int* height)
{
    const NSRect contentRect = [window->ns.view frame];
    const NSRect fbRect = convertRectToBacking(window, contentRect);

    if (width)
        *width = (int) fbRect.size.width;
    if (height)
        *height = (int) fbRect.size.height;
}

void _glfwPlatformIconifyWindow(_GLFWwindow* window)
{
    if (window->monitor)
        leaveFullscreenMode(window);

    [window->ns.object miniaturize:nil];
}

void _glfwPlatformRestoreWindow(_GLFWwindow* window)
{
    [window->ns.object deminiaturize:nil];
}

void _glfwPlatformShowWindow(_GLFWwindow* window)
{
    [window->ns.object makeKeyAndOrderFront:nil];
    _glfwInputWindowVisibility(window, GL_TRUE);

    // NSLog(@"_glfwPlatformShowWindow");
}

void _glfwPlatformHideWindow(_GLFWwindow* window)
{
    [window->ns.object orderOut:nil];
    _glfwInputWindowVisibility(window, GL_FALSE);
}

void _glfwPlatformPollEvents(void)
{
    for (;;)
    {
        NSEvent* event = [NSApp nextEventMatchingMask:NSAnyEventMask
                                            untilDate:[NSDate distantPast]
                                               inMode:NSDefaultRunLoopMode
                                              dequeue:YES];
        if (event == nil)
            break;

        [NSApp sendEvent:event];
    }

    [_glfw.ns.autoreleasePool drain];
    _glfw.ns.autoreleasePool = [[NSAutoreleasePool alloc] init];
}

void _glfwPlatformWaitEvents(void)
{
    // I wanted to pass NO to dequeue:, and rely on PollEvents to
    // dequeue and send.  For reasons not at all clear to me, passing
    // NO to dequeue: causes this method never to return.
    NSEvent *event = [NSApp nextEventMatchingMask:NSAnyEventMask
                                        untilDate:[NSDate distantFuture]
                                           inMode:NSDefaultRunLoopMode
                                          dequeue:YES];
    [NSApp sendEvent:event];

    _glfwPlatformPollEvents();
}

void _glfwPlatformSetCursorPos(_GLFWwindow* window, double x, double y)
{
    if (window->monitor)
    {
        CGDisplayMoveCursorToPoint(window->monitor->ns.displayID,
                                   CGPointMake(x, y));
    }
    else
    {
        const NSRect contentRect = [window->ns.view frame];
        const NSPoint localPoint = NSMakePoint(x, contentRect.size.height - y - 1);
        const NSPoint globalPoint = [window->ns.object convertBaseToScreen:localPoint];

        CGWarpMouseCursorPosition(CGPointMake(globalPoint.x,
                                              transformY(globalPoint.y)));
    }
}

void _glfwPlatformSetCursorMode(_GLFWwindow* window, int mode)
{
    setModeCursor(window, mode);

    if (mode == GLFW_CURSOR_DISABLED)
    {
        CGAssociateMouseAndMouseCursorPosition(false);
        centerCursor(window);
    }
    else
    {
        CGAssociateMouseAndMouseCursorPosition(true);
    }
}
void _glfwPlatformSetNeedsDisplay(_GLFWwindow* window, int b)
{
    [window->ns.view setNeedsDisplay:(b!=0)];
}

void _glfwPlatformShowText(_GLFWwindow *window, SubItem *items, int len, int withPosition, float x, float y)
{
    if (withPosition) {
        if (len == 0) {
            [window->ns.subWithPosition setHidden:YES];
        } else {
//            NSLog(@"%s %d %lf %lf", items[0].str, len, x*1280, y*720);
            [window->ns.subWithPosition setText:items length:len];
            [window->ns.subWithPosition setFrameOrigin:NSMakePoint(x*1280, y*720)];
//            [window->ns.subWithPosition setHidden:NO];
            [window->ns.subWithPosition setNeedsDisplay:YES];
        }
    } else {
//        NSLog(@"%s %d %lf %lf", items[0].str, len, x*1280, y*720);
        [window->ns.subview setText:items length:len];
    }
}

void _glfwPlatformShowLeftTime(_GLFWwindow *window, char *time, char *leftTime, float percent)
{
    [window->ns.trackView updateStatus:[NSString stringWithUTF8String:time] leftTime:[NSString stringWithUTF8String:leftTime] percent:percent];
}
void _showOrHideStartUpView(_GLFWwindow *window, int b)
{
    [window->ns.startupView setHidden:b==0];
}
//////////////////////////////////////////////////////////////////////////
//////                        GLFW native API                       //////
//////////////////////////////////////////////////////////////////////////

GLFWAPI id glfwGetCocoaWindow(GLFWwindow* handle)
{
    _GLFWwindow* window = (_GLFWwindow*) handle;
    _GLFW_REQUIRE_INIT_OR_RETURN(nil);
    return window->ns.object;
}

