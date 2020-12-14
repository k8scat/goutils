package request

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
)

const (
	ContentTypeJSON           = "application/json"
	ContentTypeFormURLEncoded = "application/x-www-form-urlencoded"
)

type Request struct {
	BaseURL        string
	DefaultHeaders map[string]string
	Client         *http.Client
}

func NewRequest(baseURL string, defaultHeaders map[string]string) *Request {
	return &Request{
		BaseURL:        baseURL,
		DefaultHeaders: defaultHeaders,
		Client:         http.DefaultClient,
	}
}

func (r *Request) setRequest(req *http.Request, params *url.Values, headers map[string]string, cookies map[string]string) *http.Request {
	if r.DefaultHeaders != nil {
		for k, v := range r.DefaultHeaders {
			req.Header.Set(k, v)
		}
	}
	if params != nil {
		req.URL.RawQuery = params.Encode()
	}
	if headers != nil {
		for k, v := range headers {
			req.Header.Set(k, v)
		}
	}
	if cookies != nil {
		for k, v := range cookies {
			req.AddCookie(&http.Cookie{
				Name:  k,
				Value: v,
			})
		}
	}
	return req
}

func (r *Request) Get(endpoint string, params *url.Values, headers map[string]string, cookies map[string]string) (resp *http.Response, err error) {
	url := fmt.Sprintf("%s%s", r.BaseURL, endpoint)
	var req *http.Request
	req, err = http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return
	}
	r.setRequest(req, params, headers, cookies)
	resp, err = r.Client.Do(req)
	return
}

func (r *Request) Post(endpoint string, params *url.Values, body io.Reader, headers map[string]string, cookies map[string]string) (resp *http.Response, err error) {
	url := fmt.Sprintf("%s%s", r.BaseURL, endpoint)
	var req *http.Request
	req, err = http.NewRequest(http.MethodPost, url, body)
	if err != nil {
		return
	}
	r.setRequest(req, params, headers, cookies)
	resp, err = r.Client.Do(req)
	return
}

func Error(resp *http.Response) error {
	content, _ := ReadBody(resp)
	return fmt.Errorf("%s %s %d: %s", resp.Request.Method, resp.Request.URL.String(), resp.StatusCode, content)
}

func RandomUA() string {
	userAgentList := []string{
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/85.0.4183.102 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/85.0.4183.102 Safari/537.36 Edg/85.0.564.60",
	}
	return userAgentList[rand.Intn(len(userAgentList))]
}

func JSONBody(data map[string]interface{}) (body *bytes.Buffer, err error) {
	if data == nil {
		err = errors.New("data cannot be nil")
		return
	}
	body = &bytes.Buffer{}
	err = json.NewEncoder(body).Encode(data)
	return
}

func FormBody(values *url.Values) (body *strings.Reader, err error) {
	if values == nil {
		err = errors.New("values cannot be nil")
		return
	}
	body = strings.NewReader(values.Encode())
	return
}

func ReadBody(resp *http.Response) (string, error) {
	if resp.Body == nil {
		return "", nil
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	resp.Body = ioutil.NopCloser(bytes.NewBuffer(b))
	return string(b), nil
}
