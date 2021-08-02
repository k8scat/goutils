package dingtalk

import "testing"

func TestDingtalk_Send(t *testing.T) {
	client := &Robot{
		AccessToken: "",
		Secret:      "",
	}
	err := client.Send("hello world", false, nil)
	if err != nil {
		t.Error(err)
	}
}
