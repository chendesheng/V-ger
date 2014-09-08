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
	"vger/httpex"
	"vger/util"
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

var UserName string
var Password string

func Login(quit chan struct{}) error {
	config := util.ReadAllConfigs()
	user := config["thunder-user"]
	password := config["thunder-password"]

	gdriveid, _, err := Login2(util.ReadConfig("gdriveid"), user, password, quit)
	if err == nil {
		util.SaveConfig("gdriveid", gdriveid)
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
		return "", "", fmt.Errorf("Login faild")
	}

	args := strings.Split(result, ":")
	if len(args) < 2 {
		return "", "", fmt.Errorf("Login faild")
	}
	verifyCode := args[1]
	passwordMd5 := singleMd5(doubleMD5(password) + strings.ToUpper(verifyCode))

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
	http.DefaultClient.Jar.SetCookies(url, cookies)

	// http.DefaultClient.Get("http://lixian.vip.xunlei.com")
}
