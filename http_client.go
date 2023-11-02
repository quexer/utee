package utee

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
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

// set request header
func (p *HttpOpt) fillRequest(req *http.Request) error {
	for k, v := range p.Headers {
		k = strings.ToLower(strings.TrimSpace(k))
		if k == "" {
			return ErrEmptyHeaderName
		}

		if k == "content-type" {
			continue
		}

		req.Header.Set(cases.Title(language.English).String(k), v)
	}

	if p.BasicAuth != nil {
		req.SetBasicAuth(p.BasicAuth.Username, p.BasicAuth.Password)
	}

	return nil
}

// SetHttpClient  expose HTTP Client for further customize
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

	req, err := http.NewRequest("POST", postUrl, body)

	if err != nil {
		return nil, fmt.Errorf("[http] err %s, %s", postUrl, err)
	}

	req.Header.Set("Content-Type", contentType)

	if opt != nil {
		if err := opt.fillRequest(req); err != nil {
			return nil, err
		}
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("[http] err %s, %s", postUrl, err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("[http] status err %s, %d", postUrl, resp.StatusCode)
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("[http] read err %s, %s", postUrl, err)
	}

	return b, nil
}

func HttpGet2(getUrl string, contentType string, opt *HttpOpt) ([]byte, error) {
	httpClientThrottle.Acquire()
	defer httpClientThrottle.Release()

	var resp *http.Response

	req, err := http.NewRequest("GET", getUrl, nil)
	if err != nil {
		return nil, fmt.Errorf("[http] err %s, %s\n", getUrl, err)
	}

	contentType = strings.TrimSpace(contentType)
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}

	if opt != nil {
		if err := opt.fillRequest(req); err != nil {
			return nil, err
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

	return io.ReadAll(resp.Body)
}

// HttpPost Deprecated use HttpPost2 instead
func HttpPost(postUrl string, q url.Values, credential ...string) ([]byte, error) {
	httpClientThrottle.Acquire()
	defer httpClientThrottle.Release()

	var resp *http.Response

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

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("[http] read err %s, %s", postUrl, err)
	}

	return b, nil
}

// HttpGet Deprecated use HttpGet2 instead
func HttpGet(getUrl string, credential ...string) ([]byte, error) {
	httpClientThrottle.Acquire()
	defer httpClientThrottle.Release()

	var resp *http.Response

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

	return io.ReadAll(resp.Body)
}
