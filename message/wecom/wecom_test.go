package wecom

import (
	"testing"
)

func TestWechat_Send(t *testing.T) {
	robot := &GroupRobot{
		Key: "",
	}
	err := robot.Send("hello world", nil, nil)
	if err != nil {
		t.Error(err)
	}
}
