package thunder

import (
	"crypto/md5"
	"fmt"
	"net/url"
	"strings"
)

func singleMd5(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	return fmt.Sprintf("%x", h.Sum(nil))
}
func doubleMD5(p string) string {
	return singleMd5(singleMd5(p))
}

func Login(user string, password string) {

	sendGet("http://login.xunlei.com/check",
		&url.Values{
			"u": {user},
		})

	verifyCode := strings.Split(getCookieValue("check_result"), ":")[1]
	passwordMd5 := singleMd5(doubleMD5(password) + strings.ToUpper(verifyCode))

	sendPost("http://login.xunlei.com/sec2login/", nil,
		&url.Values{
			"login_enable": {"1"},
			"login_hour":   {"720"},
			"verifycode":   {verifyCode},
			"u":            {user},
			"p":            {passwordMd5},
		})
}
