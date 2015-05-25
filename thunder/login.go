package thunder

import (
	"crypto/md5"
	"fmt"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"regexp"
	"strings"
	"time"
	"vger/httpex"
	"vger/util"
)

func singleMd5(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	return fmt.Sprintf("%x", h.Sum(nil))
}

var isLogined = false

var UserName string
var Password string
var Gdriveid string

func Login(quit chan struct{}) error {
	Gdriveid, _, err := Login2(Gdriveid, UserName, Password, quit)
	if err == nil {
		util.SaveConfig("gdriveid", Gdriveid)
	}
	return err
}

func Login2(gdriveid string, user, password string, quit chan struct{}) (string, string, error) {
	setCookie("gdriveid", gdriveid)

	if isLogined {
		return gdriveid, getCookieValue("userid"), nil
	}

	_, err := httpex.GetStringResp("http://login.xunlei.com/check",
		&url.Values{
			"u": {user},
		}, quit)
	if err != nil {
		return "", "", err
	}

	result := getCookieValue("check_result")
	if len(result) == 0 {
		return "", "", fmt.Errorf("Login failed")
	}
	log.Println("check_result:", result)

	args := strings.Split(result, ":")
	if len(args) < 2 {
		return "", "", fmt.Errorf("Login failed")
	}
	verifyCode := args[1]

	log.Println("args:", args[0], args[1])

	passwordMd5 := singleMd5(password + strings.ToUpper(verifyCode))

	_, err = httpex.PostFormRespString("http://login.xunlei.com/sec2login/", nil,
		&url.Values{
			"login_enable": {"1"},
			"login_hour":   {"720"},
			"verifycode":   {verifyCode},
			"u":            {user},
			"p":            {passwordMd5},
		})
	if err != nil {
		return "", "", err
	}

	//gdriveid
	_, err = httpex.GetStringResp("http://dynamic.lixian.vip.xunlei.com/login?from=0", &url.Values{}, quit)
	if err != nil {
		return "", "", err
	}

	userid := getCookieValue("userid")
	log.Print("userid:", userid)

	html, err := httpex.GetStringResp("http://dynamic.cloud.vip.xunlei.com/user_task",
		&url.Values{
			"userid": {userid},
			"st":     {"4"},
		}, quit)
	if err != nil {
		return "", "", err
	}

	// println(html)

	gdriveidReg := regexp.MustCompile(`input type="hidden" id="cok" value="([^"]+)"`)
	matches := gdriveidReg.FindStringSubmatch(html)
	if matches == nil {
		return "", "", fmt.Errorf("Can't find gdriveid.")
	}

	gdriveid = matches[1]

	log.Print("gdriveid: ", gdriveid)

	setCookie("gdriveid", gdriveid)

	isLogined = true

	log.Println("Thunder login success.")
	return gdriveid, userid, nil
}

func GetUserId() string {
	return getCookieValue("userid")
}

func setCookie(name, value string) {
	cookie := http.Cookie{
		Name:    name,
		Value:   value,
		Domain:  "xunlei.com",
		Expires: time.Now().AddDate(100, 0, 0),
	}
	cookies := []*http.Cookie{&cookie}
	url, _ := url.Parse("http://lixian.vip.xunlei.com")
	if http.DefaultClient.Jar == nil {
		http.DefaultClient.Jar, _ = cookiejar.New(nil)
	}
	http.DefaultClient.Jar.SetCookies(url, cookies)

	// http.DefaultClient.Get("http://lixian.vip.xunlei.com")
}
