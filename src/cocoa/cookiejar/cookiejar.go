package cookiejar

import (
	// "fmt"
	"github.com/mkrautz/objc"
	. "github.com/mkrautz/objc/Foundation"
	"net/http"
	"net/url"
	// "strings"
	// "time"
	"sync"
)

type SafariCookieJar struct {
	sync.RWMutex
}

func (jar *SafariCookieJar) SetCookies(u *url.URL, cookies []*http.Cookie) {
	jar.Lock()
	defer jar.Unlock()
	// println("set cookie:", u.String())

	pool := NewNSAutoreleasePool()
	defer pool.Drain()
	cs := NSSharedHTTPCookieStorage()
	for _, c := range cookies {
		// println(c.String())
		if len(c.Domain) > 0 && c.Domain[0] != '.' {
			c.Domain = "." + c.Domain
		}
		cs.SetCookieForURL(NewCookie(c, u.String()), URLWithString(u.String()), NSURL{objc.NilObject()})
	}
}

func (jar *SafariCookieJar) Cookies(u *url.URL) []*http.Cookie {
	jar.RLock()
	defer jar.RUnlock()

	pool := NewNSAutoreleasePool()
	defer pool.Drain()

	cs := NSSharedHTTPCookieStorage()
	// println(cs.String())

	cookies := make([]*http.Cookie, 0)
	for _, c := range cs.CookiesForURL(URLWithString(u.String())) {
		val := c.Value()

		// expires := time.Time{}
		// if !c.IsSessionOnly() {
		// 	expires = c.ExpiresDate()
		// }
		cookies = append(cookies, &http.Cookie{
			Name:     c.Name(),
			Value:    val,
			Domain:   c.Domain(),
			Path:     c.Path(),
			Expires:  c.ExpiresDate(),
			Secure:   c.IsSecure(),
			HttpOnly: c.IsHttpOnly(),
		})
	}

	// if strings.Contains(u.String(), "xunlei") {
	// 	println("get cookie:", u.String())
	// 	fmt.Printf("get cookies: %v\n", cookies)
	// }
	return cookies

}
