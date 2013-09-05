package thunder

import (
	"crypto/md5"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"
	"util"
)

func singleMd5(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	return fmt.Sprintf("%x", h.Sum(nil))
}
func doubleMD5(p string) string {
	return singleMd5(singleMd5(p))
}

var isLogined = false

func Login() error {
	if isLogined {
		return nil
	}

	config := util.ReadAllConfigs()
	user := config["thunder-user"]
	password := config["thunder-password"]

	_, err := sendGet("http://login.xunlei.com/check",
		&url.Values{
			"u": {user},
		})
	if err != nil {
		return err
	}

	verifyCode := strings.Split(getCookieValue("check_result"), ":")[1]
	passwordMd5 := singleMd5(doubleMD5(password) + strings.ToUpper(verifyCode))

	_, err = sendPost("http://login.xunlei.com/sec2login/", nil,
		&url.Values{
			"login_enable": {"1"},
			"login_hour":   {"720"},
			"verifycode":   {verifyCode},
			"u":            {user},
			"p":            {passwordMd5},
		})
	if err != nil {
		return err
	}

	//gdriveid
	_, err = sendGet("http://dynamic.lixian.vip.xunlei.com/login?from=0", &url.Values{})
	if err != nil {
		return err
	}

	html, err := sendGet("http://dynamic.cloud.vip.xunlei.com/user_task",
		&url.Values{
			"userid": {getCookieValue("userid")},
			"st":     {"4"},
		})
	if err != nil {
		return err
	}

	gdriveidReg := regexp.MustCompile(`input type="hidden" id="cok" value="([^"]+)"`)
	matches := gdriveidReg.FindStringSubmatch(html)
	if matches == nil {
		return fmt.Errorf("Can't find gdriveid.")
	}

	gdriveid := matches[1]

	log.Print("gdriveid: ", gdriveid)

	cookie := http.Cookie{
		Name:    "gdriveid",
		Value:   gdriveid,
		Domain:  "xunlei.com",
		Expires: time.Now().AddDate(100, 0, 0),
	}
	cookies := []*http.Cookie{&cookie}
	url, _ := url.Parse("http://vip.lixian.xunlei.com")
	http.DefaultClient.Jar.SetCookies(url, cookies)

	isLogined = true

	log.Println("Thunder login success.")
	return nil
}
