package utee

import (
	"crypto/tls"
	"fmt"
	"golang.org/x/net/http2"
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
	http2Client        = &http.Client{
		Transport: &http2.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}},
	}
)

func HttpPost(postUrl string, q url.Values, credential ...string) ([]byte, error) {
	return httpPost(1, postUrl, q, credential...)
}

func Http2Post(postUrl string, q url.Values, credential ...string) ([]byte, error) {
	return httpPost(2, postUrl, q, credential...)
}

func httpPost(v int, postUrl string, q url.Values, credential ...string) ([]byte, error) {
	HttpClientThrottle <- nil
	defer func() {
		<-HttpClientThrottle
	}()

	var resp *http.Response
	var err error
	req, err := http.NewRequest("POST", postUrl, strings.NewReader(q.Encode()))
	if err != nil {
		return nil, fmt.Errorf("[http] err %s, %s", postUrl, err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if len(credential) == 2 {
		req.SetBasicAuth(credential[0], credential[1])
	}

	var client *http.Client
	if v == 1 {
		client = http.DefaultClient
	} else {
		client = http2Client
	}

	resp, err = client.Do(req)

	if err != nil {
		return nil, fmt.Errorf("[http] err %s, %s", postUrl, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("[http] status err %s, %d", postUrl, resp.StatusCode)
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("[http] read err %s, %s", postUrl, err)
	}
	return b, nil
}

func HttpGet(getUrl string, credential ...string) ([]byte, error) {
	return httpGet(1, getUrl, credential...)
}

func Http2Get(getUrl string, credential ...string) ([]byte, error) {
	return httpGet(2, getUrl, credential...)
}

func httpGet(v int, getUrl string, credential ...string) ([]byte, error) {
	HttpClientThrottle <- nil
	defer func() {
		<-HttpClientThrottle
	}()

	var resp *http.Response
	var err error
	req, err := http.NewRequest("GET", getUrl, nil)
	if err != nil {
		return nil, fmt.Errorf("[http] err %s, %s\n", getUrl, err)
	}
	if len(credential) == 2 {
		req.SetBasicAuth(credential[0], credential[1])
	}

	var client *http.Client
	if v == 1 {
		client = http.DefaultClient
	} else {
		client = http2Client
	}

	resp, err = client.Do(req)

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("[http get] status err %s, %d\n", getUrl, resp.StatusCode)
	}
	return ioutil.ReadAll(resp.Body)
}
