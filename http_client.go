package utee

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	MAX_HTTP_CLIENT_CONCURRENT = 1000
)

var (
	httpClientThrottle = NewThrottle(MAX_HTTP_CLIENT_CONCURRENT)
	insecureClient     = &http.Client{
		Timeout:   15 * time.Second,
		Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}},
	}
)

func HttpPost(postUrl string, q url.Values, credential ...string) ([]byte, error) {
	httpClientThrottle.Acquire()
	defer httpClientThrottle.Release()

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

	resp, err = insecureClient.Do(req)

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
	httpClientThrottle.Acquire()
	defer httpClientThrottle.Release()

	var resp *http.Response
	var err error
	req, err := http.NewRequest("GET", getUrl, nil)
	if err != nil {
		return nil, fmt.Errorf("[http] err %s, %s\n", getUrl, err)
	}
	if len(credential) == 2 {
		req.SetBasicAuth(credential[0], credential[1])
	}

	resp, err = insecureClient.Do(req)

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("[http get] status err %s, %d\n", getUrl, resp.StatusCode)
	}
	return ioutil.ReadAll(resp.Body)
}
