package utee

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

const (
	MAX_HTTP_CLIENT_CONCURRENT = 1000
)

var (
	HttpClientThrottle = make(chan interface{}, MAX_HTTP_CLIENT_CONCURRENT)
)

func HttpPost(postUrl string, q url.Values, credential ...string) ([]byte, error) {
	HttpClientThrottle <- nil
	defer func() {
		<-HttpClientThrottle
	}()

	var resp *http.Response
	var err error
	if len(credential) == 2 {
		req, err := http.NewRequest("POST", postUrl, strings.NewReader(q.Encode()))
		if err != nil {
			return nil, fmt.Errorf("[http] err %s, %s\n", postUrl, err)
		}
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.SetBasicAuth(credential[0], credential[1])
		resp, err = http.DefaultClient.Do(req)
	} else {
		resp, err = http.PostForm(postUrl, q)
	}

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

func HttpGet(getUrl string, credential ...string) ([]byte, error) {
	HttpClientThrottle <- nil
	defer func() {
		<-HttpClientThrottle
	}()

	var resp *http.Response
	var err error
	if len(credential) == 2 {
		req, err := http.NewRequest("GET", getUrl, nil)
		if err != nil {
			return nil, fmt.Errorf("[http] err %s, %s\n", getUrl, err)
		}
		req.SetBasicAuth(credential[0], credential[1])
		resp, err = http.DefaultClient.Do(req)
	} else {
		resp, err = http.Get(getUrl)
	}

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("[http get] status err %s, %d\n", getUrl, resp.StatusCode)
	}
	return ioutil.ReadAll(resp.Body)
}
