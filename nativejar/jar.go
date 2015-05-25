package nativejar

// #cgo CFLAGS: -x objective-c
// #cgo LDFLAGS: -framework Cocoa -framework OpenGL -framework QuartzCore
//#include "cookie.h"
//#include "stdlib.h"
import "C"

import (
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"sync"
	"unsafe"
)

type NativeJar struct {
	sync.Mutex
	jar *cookiejar.Jar
}

func New() (NativeJar, error) {
	jar, _ := cookiejar.New(nil)
	return NativeJar{sync.Mutex{}, jar}, nil
}

func (nj NativeJar) SetCookies(u *url.URL, cookies []*http.Cookie) {
	nj.Lock()
	defer nj.Unlock()

	cstr := C.CString(u.String())
	defer C.free(unsafe.Pointer(cstr))

	var ccookies []C.Cookie
	for _, c := range cookies {
		cc := C.Cookie{}

		cc.name = C.CString(c.Name)
		defer C.free(unsafe.Pointer(cc.name))

		cc.value = C.CString(c.Value)
		defer C.free(unsafe.Pointer(cc.value))

		cc.path = C.CString(c.Path)
		defer C.free(unsafe.Pointer(cc.path))

		cc.domain = C.CString(c.Domain)
		defer C.free(unsafe.Pointer(cc.domain))

		cc.expires = C.long(c.Expires.Unix())

		cc.secure = 0
		if c.Secure {
			cc.secure = 1
		}
		cc.httpOnly = 0
		if c.HttpOnly {
			cc.httpOnly = 1
		}

		ccookies = append(ccookies, cc)
	}

	C.setCookies(cstr, &ccookies[0], C.int(len(ccookies)))
	nj.jar.SetCookies(u, cookies)
}

func (nj NativeJar) Cookies(u *url.URL) []*http.Cookie {
	nj.Lock()
	defer nj.Unlock()

	cstr := C.CString(u.String())
	defer C.free(unsafe.Pointer(cstr))

	length := C.int(0)
	ptr := C.getCookies(cstr, &length)
	if length == 0 {
		return nil
	}

	ccookies := (*[1 << 30]C.Cookie)(unsafe.Pointer(ptr))[:length:length]
	defer C.free(unsafe.Pointer(&ccookies[0]))

	cookies := make([]*http.Cookie, len(ccookies))
	for i, cc := range ccookies {
		cookies[i] = &http.Cookie{
			Name:     C.GoString(cc.name),
			Value:    C.GoString(cc.value),
			Domain:   C.GoString(cc.domain),
			HttpOnly: cc.httpOnly != 0,
			Secure:   cc.secure != 0,
		}
	}

	res := nj.jar.Cookies(u)
	for _, c := range cookies {
		exists := false
		for _, cc := range res {
			if cc.Name == c.Name {
				exists = true
			}
		}

		if !exists {
			res = append(res, c)
		}
	}

	//log.Printf("%+v", cookies)
	return res
}
