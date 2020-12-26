package utee

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	MAX_HTTP_CLIENT_CONCURRENT = 1000

	ContentTypeForm = "application/x-www-form-urlencoded"
	ContentTypeJson = "application/json; charset=utf-8"
)

var (
	httpClientThrottle = NewThrottle(MAX_HTTP_CLIENT_CONCURRENT)
	client             = &http.Client{
		Timeout: 15 * time.Second,
	}
	ErrEmptyHeaderName = errors.New("header name must not be empty")
)

type BasicAuth struct {
	Username string
	Password string
}

type HttpOpt struct {
	Headers   map[string]string
	BasicAuth *BasicAuth
}

// HTTP Client expose for further customize
func SetHttpClient(hc *http.Client) {
	client = hc
}

func HttpPost2(postUrl string, contentType string, body io.Reader, opt *HttpOpt) ([]byte, error) {
	httpClientThrottle.Acquire()
	defer httpClientThrottle.Release()

	contentType = strings.TrimSpace(contentType)
	if contentType == "" {
		contentType = ContentTypeForm
	}

	var resp *http.Response
	var err error
	req, err := http.NewRequest("POST", postUrl, body)
	if err != nil {
		return nil, fmt.Errorf("[http] err %s, %s", postUrl, err)
	}
	req.Header.Set("Content-Type", contentType)

	if opt != nil {
		for k, v := range opt.Headers {
			k := strings.ToLower(strings.TrimSpace(k))
			if k == "" {
				return nil, ErrEmptyHeaderName
			}
			if k == "content-type" {
				continue
			}
			req.Header.Set(strings.Title(k), v)
		}

		if opt.BasicAuth != nil {
			req.SetBasicAuth(opt.BasicAuth.Username, opt.BasicAuth.Password)
		}
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

func HttpGet2(getUrl string, contentType string, opt *HttpOpt) ([]byte, error) {
	httpClientThrottle.Acquire()
	defer httpClientThrottle.Release()

	var resp *http.Response
	var err error
	req, err := http.NewRequest("GET", getUrl, nil)
	if err != nil {
		return nil, fmt.Errorf("[http] err %s, %s\n", getUrl, err)
	}

	contentType = strings.TrimSpace(contentType)
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}

	if opt != nil {
		for k, v := range opt.Headers {
			k := strings.ToLower(strings.TrimSpace(k))
			if k == "" {
				return nil, ErrEmptyHeaderName
			}
			if k == "content-type" {
				continue
			}
			req.Header.Set(strings.Title(k), v)
		}

		if opt.BasicAuth != nil {
			req.SetBasicAuth(opt.BasicAuth.Username, opt.BasicAuth.Password)
		}
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
