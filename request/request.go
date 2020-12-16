package request

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
)

const (
	ContentTypeJSON           = "application/json"
	ContentTypeFormURLEncoded = "application/x-www-form-urlencoded"
)

type Request struct {
	BaseURL        string
	DefaultHeaders map[string]string
	DefaultCookies []*http.Cookie
	Client         *http.Client
}

func NewRequest(baseURL string, defaultHeaders map[string]string, defaultCookies []*http.Cookie) *Request {
	return &Request{
		BaseURL:        baseURL,
		DefaultHeaders: defaultHeaders,
		DefaultCookies: defaultCookies,
		Client:         http.DefaultClient,
	}
}

func (r *Request) setRequest(req *http.Request, params *url.Values, headers map[string]string, cookies []*http.Cookie) *http.Request {
	for k, v := range r.DefaultHeaders {
		req.Header.Set(k, v)
	}
	for _, c := range r.DefaultCookies {
		req.AddCookie(c)
	}

	if params != nil {
		req.URL.RawQuery = params.Encode()
	}
	if headers != nil {
		for k, v := range headers {
			req.Header.Set(k, v)
		}
	}
	for _, c := range cookies {
		req.AddCookie(c)
	}
	return req
}

func (r *Request) Get(endpoint string, params *url.Values, headers map[string]string, cookies []*http.Cookie) (resp *http.Response, err error) {
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

func (r *Request) Post(endpoint string, params *url.Values, body io.Reader, headers map[string]string, cookies []*http.Cookie) (resp *http.Response, err error) {
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

func CreateJSONPayload(data map[string]interface{}) (body *bytes.Buffer, err error) {
	if data == nil {
		return
	}
	body = &bytes.Buffer{}
	err = json.NewEncoder(body).Encode(data)
	return
}

func CreateFormPayload(textFields, fileFields map[string]string) (*bytes.Buffer, error) {
	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	defer writer.Close()
	for key, value := range textFields {
		if err := writer.WriteField(key, value); err != nil {
			return nil, err
		}
	}
	for key, value := range fileFields {
		file, err := os.Open(value)
		if err != nil {
			return nil, err
		}
		fw, err := writer.CreateFormFile(key, filepath.Base(value))
		if _, err := io.Copy(fw, file); err != nil {
			return nil, err
		}
		if err := file.Close(); err != nil {
			return nil, err
		}
	}
	return payload, nil
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
