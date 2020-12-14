package message

import "testing"

func TestDingtalk_Send(t *testing.T) {
	accessToken := ""
	secret := ""
	client, err := NewDingtalk(accessToken, secret)
	if err != nil {
		t.Error(err)
	}
	if err = client.Send("hello world", false, nil); err != nil {
		t.Error(err)
	}
}
