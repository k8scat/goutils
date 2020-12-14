package message

import (
	"errors"
	"io/ioutil"
	"net/url"

	"github.com/k8scat/gotie/request"
	"github.com/tidwall/gjson"
)

// 每个机器人发送的消息不能超过20条/分钟。

const sendURL = "https://qyapi.weixin.qq.com/cgi-bin/webhook/send"

type WeComGroupRobot struct {
	Key string
}

func NewWeComGroupRobot(key string) (client *WeComGroupRobot, err error) {
	if key == "" {
		err = errors.New("key cannot be empty")
		return
	}
	client = &WeComGroupRobot{
		Key: key,
	}
	return
}

func (w *WeComGroupRobot) Send(content string, mentionedUserIDList, mentionedMobileList []string) error {
	r := request.NewRequest("", nil)
	params := &url.Values{
		"key": {w.Key},
	}
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
	payload, err := request.JSONBody(data)
	if err != nil {
		return err
	}
	headers := map[string]string{
		"Content-Type": request.ContentTypeFormURLEncoded,
	}
	resp, err := r.Post(sendURL, params, payload, headers, nil)
	if err != nil {
		return err
	}
	b, err := ioutil.ReadAll(resp.Body)
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
