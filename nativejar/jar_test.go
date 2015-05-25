package nativejar

import (
	"fmt"
	"net/http"
	"net/url"
	"testing"
	"time"
)

func TestCookies(t *testing.T) {
	u, _ := url.Parse("http://www.zimuzu.tv")

	jar, _ := New()
	fmt.Printf("%#v\n", jar.Cookies(u))
}

func TestSetCookie(t *testing.T) {
	u, _ := url.Parse("http://www.zimuzu.tv")
	jar, _ := New()

	cookies := make([]*http.Cookie, 1, 1)
	cookies[0] = &http.Cookie{"goname123", "govalue", "/", ".zimuzu.tv", time.Now().UTC().Add(time.Hour * 1000), "", 0, false, false, "", nil}
	jar.SetCookies(u, cookies)

	fmt.Printf("%#v\n", jar.Cookies(u))
}
