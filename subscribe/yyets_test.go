package subscribe

import (
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"os"
	"testing"
)

func TestParse(t *testing.T) {
	if http.DefaultClient.Jar == nil {
		http.DefaultClient.Jar, _ = cookiejar.New(nil)
	}

	url := "http://www.zimuzu.tv/resource/32235"
	s, tks, err := Parse(url)
	if err != nil {
		t.Error(err)
	}

	fmt.Printf("%+v\n", s)
	fmt.Printf("%v\n", tks[0])
	fmt.Println(len(tks))
}

func TestParseSubscribeInfo(t *testing.T) {
	f, err := os.Open("fargo-info.html")
	if err != nil {
		t.Error(err)
	}
	s, err := parseSubscribeInfo(f)
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("%+v\n", s)
}

func TestYYetsFargo(t *testing.T) {
	f, err := os.Open("fargo.html")
	if err != nil {
		t.Error(err)
	}

	tasks, err := parseEpisodes(f)
	if err != nil {
		t.Error(err)
	}

	println(len(tasks))
}
func TestYYetsLogin(t *testing.T) {
	if http.DefaultClient.Jar == nil {
		http.DefaultClient.Jar, _ = cookiejar.New(nil)
	}

	YYetsLogin("", "")
}
