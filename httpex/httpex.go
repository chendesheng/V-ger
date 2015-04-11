package httpex

import (
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
)

func CancelRequest(req *http.Request) {
	http.DefaultTransport.(*http.Transport).CancelRequest(req)
}

func Get(url string, quit chan struct{}) (*http.Response, error) {
	finish := make(chan error)
	defer close(finish)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	go func() {
		select {
		case <-quit:
			CancelRequest(req)
		case <-finish:
		}
	}()

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func GetStringResp(url string, params *url.Values, quit chan struct{}) (string, error) {
	if params != nil {
		url = url + "?" + params.Encode()
	}
	resp, err := Get(url, quit)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	return readBody(resp.Body)
}

func readBody(body io.ReadCloser) (string, error) {
	bytes, err := ioutil.ReadAll(body)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

func PostFormRespString(url string, params *url.Values, data *url.Values) (string, error) {
	if params != nil {
		url = url + "?" + params.Encode()
	}
	resp, err := http.PostForm(url, *data)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	return readBody(resp.Body)
}
