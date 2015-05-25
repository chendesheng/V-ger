#include "cookie.h"

#import <Cocoa/Cocoa.h>
#include "stdio.h"

Cookie* getCookies(char* url, int* len) {
        @autoreleasepool {
		NSURL* u = [NSURL URLWithString:[NSString stringWithUTF8String:url]];
		NSArray* cookies = [[NSHTTPCookieStorage sharedHTTPCookieStorage] cookiesForURL:u];
		Cookie* rets = (Cookie*)malloc(sizeof(Cookie)*cookies.count);
		for (int i = 0; i < cookies.count; i++) {
			NSHTTPCookie* c = [cookies objectAtIndex:i];
			//NSLog(@"get cookie: %@", c);
			Cookie* cc = &rets[i];
			cc->name = [c.name UTF8String];
			cc->value = [c.value UTF8String];
			cc->path = [c.path UTF8String];
			cc->domain = [c.domain UTF8String];
			cc->httpOnly = c.HTTPOnly;
			cc->secure = c.secure;
			cc->expires = [c.expiresDate timeIntervalSince1970];
		}
		
		*len = (int)cookies.count;

//		for (int i = 0; i < cookies.count; i++) {
//			printf("get cookie: %s=%s\n", rets[i].name, rets[i].value);
//		}
		return rets;
        }
}

void setCookies(char* url, Cookie* cookies, int len) {
	@autoreleasepool {
		for (int i = 0; i < len; i++) {
			Cookie* c = &cookies[i];
			//printf("set cookie: %s=%s\n", c->name, c->value);
			//printf("%s\n", c->domain);
			//printf("%ld\n", c->expires);
			
			NSMutableDictionary *dic =[@{
				NSHTTPCookieOriginURL: [NSString stringWithUTF8String:url],
				NSHTTPCookieName: [NSString stringWithUTF8String:c->name],
				NSHTTPCookieValue: [NSString stringWithUTF8String:c->value],
				NSHTTPCookiePath: [NSString stringWithUTF8String:c->path],
				NSHTTPCookieDomain: [NSString stringWithUTF8String:c->domain],
				NSHTTPCookieExpires: [NSDate dateWithTimeIntervalSince1970:c->expires].description
			} mutableCopy];
			if (c->secure == 1) {
				[dic setObject:@"TRUE" forKey:NSHTTPCookieSecure];
			}
			//if (c->expires <= 0) {
			//	[dic setObject:@"TRUE" forKey:NSHTTPCookieDiscard];
			//}

			
			NSHTTPCookie* nsc = [NSHTTPCookie cookieWithProperties: dic];
			//NSLog(@"set cookie: %@", nsc);
					
			[[NSHTTPCookieStorage sharedHTTPCookieStorage] setCookie:nsc];
		}
	}
}
