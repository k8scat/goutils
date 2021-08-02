package request

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type Client struct {
	Client         *http.Client
	BaseURL        string
	DefaultHeaders map[string]string
	DefaultCookies []*http.Cookie
}

func (r *Client) setRequest(req *http.Request, params *url.Values, headers map[string]string, cookies []*http.Cookie) {
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

func (c *Client) Get(endpoint string, params *url.Values, headers map[string]string, cookies []*http.Cookie) (*http.Response, error) {
	url := concatURL(c.BaseURL, endpoint)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	c.setRequest(req, params, headers, cookies)
	return c.Client.Do(req)
}

func (r *Client) Post(endpoint string, params *url.Values, body io.Reader, headers map[string]string, cookies []*http.Cookie) (resp *http.Response, err error) {
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
