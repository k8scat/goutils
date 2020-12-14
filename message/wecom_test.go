package message

import (
	"testing"
)

func TestWechat_Send(t *testing.T) {
	key := ""
	robot, err := NewWeComGroupRobot(key)
	if err != nil {
		t.Error(err)
	}
	if err = robot.Send("hello world", nil, nil); err != nil {
		t.Error(err)
	}
}
