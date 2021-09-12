package request

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type Request struct {
	HttpClient    *http.Client
	BackOffClient *BackOffClient

	BaseURL        string
	DefaultHeaders map[string]string
	DefaultCookies []*http.Cookie
}

func (r *Request) Get(endpoint string, params *url.Values, headers map[string]string, cookies []*http.Cookie) (*http.Response, error) {
	url := concatURL(r.BaseURL, endpoint)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	r.setRequest(req, params, headers, cookies)
	if r.BackOffClient != nil {
		return r.BackOffClient.Do(r.HttpClient, req)
	}
	return r.HttpClient.Do(req)
}

func (r *Request) Post(endpoint string, params *url.Values, body io.Reader, headers map[string]string, cookies []*http.Cookie) (*http.Response, error) {
	url := fmt.Sprintf("%s%s", r.BaseURL, endpoint)
	req, err := http.NewRequest(http.MethodPost, url, body)
	if err != nil {
		return nil, err
	}
	r.setRequest(req, params, headers, cookies)
	if r.BackOffClient != nil {
		return r.BackOffClient.Do(r.HttpClient, req)
	}
	return r.HttpClient.Do(req)
}

func (r *Request) setRequest(req *http.Request, params *url.Values, headers map[string]string, cookies []*http.Cookie) {
	for k, v := range r.DefaultHeaders {
		req.Header.Set(k, v)
	}
	for _, c := range r.DefaultCookies {
		req.AddCookie(c)
	}

	if params != nil {
		req.URL.RawQuery = params.Encode()
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	for _, c := range cookies {
		req.AddCookie(c)
	}
}
