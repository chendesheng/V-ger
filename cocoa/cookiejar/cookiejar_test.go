package cookiejar

import (
	"net/http"
	"net/url"
	"testing"
	"time"
)

func TestGet(t *testing.T) {
	jar := SafariCookieJar{}
	u, _ := url.Parse("http://lixian.vip.xunlei.com/")
	cookies := jar.Cookies(u)
	for _, c := range cookies {
		println(c.String())
	}
}

func TestSet(t *testing.T) {
	jar := SafariCookieJar{}
	u, _ := url.Parse("http://login.xunlei.com/check?u=123456")
	jar.SetCookies(u, []*http.Cookie{&http.Cookie{
		Name:    "VERIFY_KEY",
		Value:   "1AEFA9F6716D9623D9C152F839A04E61",
		Path:    "/",
		Domain:  "xunlei.com",
		Expires: time.Now().Add(time.Hour),
	}})

	u1, _ := url.Parse("http://lixian.vip.xunlei.com/")
	cookies := jar.Cookies(u1)
	finded := false
	for _, c := range cookies {
		if c.Name == "VERIFY_KEY" && c.Value == "1AEFA9F6716D9623D9C152F839A04E61" {
			finded = true
		}
	}

	if !finded {
		t.Errorf("not find cookie")
	}
}
