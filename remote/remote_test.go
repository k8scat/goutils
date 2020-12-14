package remote

import "testing"

func TestRunCommandRemote(t *testing.T) {
	ip := ""
	port := DefaultSSHPort
	user := DefaultSSHUser
	password := ""
	sshClient, err := NewSSHClient(ip, port, user, password)
	if err != nil {
		t.Error(err)
	}
	if err = RunCommandRemote(sshClient, "cd /tmp; echo \"hello world\" >> test.txt"); err != nil {
		t.Error(err)
	}
}

func TestRunScriptRemote(t *testing.T) {
	ip := ""
	port := DefaultSSHPort
	user := DefaultSSHUser
	password := ""
	scriptFile := "/tmp/test.sh"
	sshClient, err := NewSSHClient(ip, port, user, password)
	if err != nil {
		t.Error(err)
	}
	if err := RunScriptRemote(sshClient, scriptFile, DefaultPermissions, false); err != nil {
		t.Error(err)
	}
}
