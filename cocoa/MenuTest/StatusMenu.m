//
//  StatusMenu.m
//  V'gerMenu
//
//  Created by Roy Chen on 1/26/14.
//  Copyright (c) 2014 me. All rights reserved.
//

#import "StatusMenu.h"

@implementation StatusMenu
-(IBAction)sleepAfterFinishClick:(id)sender {
    
}
- (void)awakeFromNib
{
    NSLog(@"awakeFromNib");
//    NSMenu* menu = [[NSMenu alloc] initWithTitle:@"foo"];
    NSMenuItem* item = [[NSApp mainMenu] addItemWithTitle:@"foo"
                                                   action: nil
                                            keyEquivalent: @""];
    [item setSubmenu:[self menu]];

}

@end
