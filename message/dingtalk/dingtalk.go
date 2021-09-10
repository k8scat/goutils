// https://ding-doc.dingtalk.com/doc#/serverapi2/qf2nxq/e9d991e2
// https://github.com/JetBlink/dingtalk-notify-go-sdk
package dingtalk

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/k8scat/goutils/request"
)

const (
	RobotAPI = "https://oapi.dingtalk.com/robot/send"

	RequestTimeout = time.Second * time.Duration(10)
)

type Robot struct {
	AccessToken string
	Secret      string
}

func (r *Robot) Send(content string, atAll bool, atMobiles []string) error {
	at := map[string]interface{}{
		"isAtAll": atAll,
	}
	if atMobiles != nil && len(atMobiles) > 0 {
		at["atMobiles"] = atMobiles
	}
	data := map[string]interface{}{
		"msgtype": "text",
		"text": map[string]string{
			"content": content,
		},
		"at": at,
	}
	b, err := json.Marshal(data)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, RobotAPI, bytes.NewReader(b))
	if err != nil {
		return err
	}

	currentTimestamp := strconv.FormatInt(time.Now().UnixNano()/1e6, 10)
	signStr := sign(currentTimestamp, r.Secret)
	q := req.URL.Query()
	q.Add("access_token", r.AccessToken)
	q.Add("timestamp", currentTimestamp)
	q.Add("sign", signStr)
	req.URL.RawQuery = q.Encode()

	req.Header.Add("Content-Type", "application/json")
	client := &http.Client{
		Timeout: RequestTimeout,
	}
	backoff := &request.BackOffClient{
		BackOff: request.DefaultBackOff,
		Notify:  request.DefaultNotify,
	}
	resp, err := backoff.Do(client, req)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return request.Error(resp)
	}
	return nil
}

func sign(t string, secret string) string {
	s := fmt.Sprintf("%s\n%s", t, secret)
	hmac256 := hmac.New(sha256.New, []byte(secret))
	hmac256.Write([]byte(s))
	data := hmac256.Sum(nil)
	return base64.StdEncoding.EncodeToString(data)
}
