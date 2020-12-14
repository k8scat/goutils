// https://ding-doc.dingtalk.com/doc#/serverapi2/qf2nxq/e9d991e2
// https://github.com/JetBlink/dingtalk-notify-go-sdk
package message

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/k8scat/gotie/request"
)

const dingtalkSendAPI = "https://oapi.dingtalk.com/robot/send"

type Dingtalk struct {
	AccessToken string
	Secret      string
}

func NewDingtalk(accessToken, secret string) (client *Dingtalk, err error) {
	if accessToken == "" || secret == "" {
		err = errors.New("accessToken or secret cannot be empty")
		return
	}
	client = &Dingtalk{
		AccessToken: accessToken,
		Secret:      secret,
	}
	return
}

func (d *Dingtalk) Send(content string, atAll bool, atMobiles []string) error {
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
	payload, err := request.JSONBody(data)
	if err != nil {
		return err
	}
	currentTimestamp := strconv.FormatInt(time.Now().UnixNano()/1e6, 10)
	signStr := sign(currentTimestamp, d.Secret)
	params := &url.Values{
		"access_token": {d.AccessToken},
		"timestamp":    {currentTimestamp},
		"sign":         {signStr},
	}
	headers := map[string]string{
		"Content-Type": request.ContentTypeJSON,
	}
	client := request.NewRequest("", nil)
	resp, err := client.Post(dingtalkSendAPI, params, payload, headers, nil)
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
