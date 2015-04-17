package mpush

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

const (
	MAX_HTTP_CLIENT_CONCURRENT = 1000
)

var (
	httpClientThrottle = make(chan interface{}, MAX_HTTP_CLIENT_CONCURRENT)
)

func HttpPost(postUrl string, q url.Values) ([]byte, error) {
	httpClientThrottle <- nil
	defer func() {
		<-httpClientThrottle
	}()

	resp, err := http.PostForm(postUrl, q)
	if err != nil {
		return nil, fmt.Errorf("[http] err %s, %s\n", postUrl, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("[http] status err %s, %d\n", postUrl, resp.StatusCode)
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("[http] read err %s, %s\n", postUrl, err)
	}
	return b, nil
}

func HttpGet(url string) ([]byte, error) {
	httpClientThrottle <- nil
	defer func() {
		<-httpClientThrottle
	}()
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("[http get] status err %s, %d\n", url, resp.StatusCode)
	}
	return ioutil.ReadAll(resp.Body)
}
