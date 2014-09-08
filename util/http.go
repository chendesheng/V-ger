package util

import (
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"time"
)

func SetCookie(name, value, domainUrl string) {
	u, _ := url.Parse(domainUrl)
	cookie := http.Cookie{
		Name:    name,
		Value:   value,
		Domain:  u.Host,
		Expires: time.Now().AddDate(100, 0, 0),
	}
	cookies := []*http.Cookie{&cookie}

	if http.DefaultClient.Jar == nil {
		jar, _ := cookiejar.New(nil)
		http.DefaultClient.Jar = jar
	}

	http.DefaultClient.Jar.SetCookies(u, cookies)
}
