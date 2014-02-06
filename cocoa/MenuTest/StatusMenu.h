//
//  StatusMenu.h
//  V'gerMenu
//
//  Created by Roy Chen on 1/26/14.
//  Copyright (c) 2014 me. All rights reserved.
//

#import <Foundation/Foundation.h>

@interface StatusMenu : NSObject

-(IBAction)sleepAfterFinishClick:(id)sender;


@property (assign) IBOutlet NSMenuItem *menuItemSleepAfterFinish;
@property (assign) IBOutlet NSMenu *menu;


@end
