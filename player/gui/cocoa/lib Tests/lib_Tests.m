//
//  lib_Tests.m
//  lib Tests
//
//  Created by Roy Chen on 10/5/14.
//  Copyright (c) 2014 vp. All rights reserved.
//

#import <XCTest/XCTest.h>
#import "gui.h"

@interface lib_Tests : XCTestCase

@end

@implementation lib_Tests

- (void)setUp
{
    [super setUp];
    // Put setup code here. This method is called before the invocation of each test method in the class.
}

- (void)tearDown
{
    // Put teardown code here. This method is called after the invocation of each test method in the class.
    [super tearDown];
}

- (void)testExample
{
    initialize();
    void* w = newWindow("test", 400, 400);
    showWindow(w);
    pollEvents();
}

@end
