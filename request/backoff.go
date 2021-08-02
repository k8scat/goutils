package request

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/cenkalti/backoff/v4"
)

var DefaultBackOff = &BackOff{
	BackOff: &backoff.ExponentialBackOff{
		InitialInterval:     200 * time.Millisecond,
		RandomizationFactor: backoff.DefaultRandomizationFactor,
		Multiplier:          backoff.DefaultMultiplier,
		MaxInterval:         5 * time.Second,
		MaxElapsedTime:      10 * time.Second,
		Clock:               backoff.SystemClock,
	},
}

type BackOff struct {
	BackOff backoff.BackOff
	Notify  func(error, time.Duration)
}

func (b *BackOff) Do(client *http.Client, req *http.Request) (*http.Response, error) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return nil, fmt.Errorf("Read request body failed: %+v", err)
	}
	req.Body.Close()
	req.Body = io.NopCloser(bytes.NewReader(body))

	var resp *http.Response
	op := func() (err error) {
		req.Body = ioutil.NopCloser(bytes.NewReader(body))
		resp, err = client.Do(req)
		if err != nil {
			err = fmt.Errorf("Request err: %+v", err)
		}
		return err
	}
	return resp, backoff.RetryNotify(op, b.BackOff, b.Notify)
}
