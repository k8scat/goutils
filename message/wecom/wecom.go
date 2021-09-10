package wecom

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/k8scat/goutils/request"
	"github.com/tidwall/gjson"
)

const (
	// 每个机器人发送的消息不能超过20条/分钟。
	WebhookAPI = "https://qyapi.weixin.qq.com/cgi-bin/webhook/send"

	RequestTimeout = time.Second * time.Duration(10)
)

type GroupRobot struct {
	Key string
}

func (r *GroupRobot) Send(content string, mentionedUserIDList, mentionedMobileList []string) error {
	text := map[string]interface{}{
		"content": content,
	}
	if mentionedUserIDList != nil {
		text["mentioned_list"] = mentionedUserIDList
	}
	if mentionedMobileList != nil {
		text["mentioned_mobile_list"] = mentionedMobileList
	}
	data := map[string]interface{}{
		"msgtype": "text",
		"text":    text,
	}
	b, err := json.Marshal(data)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, WebhookAPI, bytes.NewReader(b))
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	q := req.URL.Query()
	q.Add("key", r.Key)
	req.URL.RawQuery = q.Encode()

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

	b, err = request.ReadBody(resp)
	if err != nil {
		return err
	}
	respContent := string(b)
	errCode := gjson.Get(respContent, "errcode").Int()
	errMsg := gjson.Get(respContent, "errmsg").String()
	if errCode != 0 {
		return errors.New(errMsg)
	}
	return nil
}
