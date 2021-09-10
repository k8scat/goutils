package request

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/cenkalti/backoff/v4"
)

var (
	DefaultBackOff = &backoff.ExponentialBackOff{
		InitialInterval:     200 * time.Millisecond,
		RandomizationFactor: backoff.DefaultRandomizationFactor,
		Multiplier:          backoff.DefaultMultiplier,
		MaxInterval:         5 * time.Second,
		MaxElapsedTime:      10 * time.Second,
		Clock:               backoff.SystemClock,
	}

	DefaultNotify = func(err error, t time.Duration) {
		log.Printf("BackOff err: %+v, retry duration: %d ms", err, t.Milliseconds())
	}
)

type BackOffClient struct {
	BackOff backoff.BackOff
	Notify  func(error, time.Duration)
}

func (c *BackOffClient) Do(client *http.Client, req *http.Request) (*http.Response, error) {
	hasBody := req.Body != nil
	var body []byte
	var err error
	if hasBody {
		body, err = ioutil.ReadAll(req.Body)
		if err != nil {
			return nil, fmt.Errorf("Read request body failed: %+v", err)
		}
		req.Body.Close()
		req.Body = io.NopCloser(bytes.NewReader(body))
	}

	var resp *http.Response
	op := func() (err error) {
		if hasBody {
			req.Body = ioutil.NopCloser(bytes.NewReader(body))
		}
		resp, err = client.Do(req)
		if err != nil {
			err = fmt.Errorf("Request err: %+v", err)
		}
		return err
	}
	return resp, backoff.RetryNotify(op, c.BackOff, c.Notify)
}
