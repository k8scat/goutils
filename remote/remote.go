package remote

import (
	"fmt"
	"github.com/bramvdbogaerde/go-scp"
	"github.com/bramvdbogaerde/go-scp/auth"
	"golang.org/x/crypto/ssh"
	"log"
	"os"
	"path"
	"regexp"
	"time"
)

const (
	DefaultSSHPort     = "22"
	DefaultSSHUser     = "root"
	DefaultPermissions = "0777"
	TmpDir             = "/tmp"

	PatternIP   = "^\\d{1,3}\\.\\d{1,3}\\.\\d{1,3}\\.\\d{1,3}$"
	PatternPort = "^\\d{1,5}$"
)

func NewSSHClient(ip, port, user, password string) (client *ssh.Client, err error) {
	if err = asset(PatternIP, ip, "ip"); err != nil {
		return
	}
	if err = asset(PatternPort, port, "port"); err != nil {
		return
	}
	if err = asset("", user, "user"); err != nil {
		return
	}
	if err = asset("", password, "password"); err != nil {
		return
	}

	config := &ssh.ClientConfig{
		User:            user,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         time.Second * time.Duration(10),
		Auth:            []ssh.AuthMethod{ssh.Password(password)},
	}
	addr := fmt.Sprintf("%s:%s", ip, port)
	client, err = ssh.Dial("tcp", addr, config)
	return
}

func NewScpClient(ip, port, user, password string) (client scp.Client, err error) {
	if err = asset(PatternIP, ip, "ip"); err != nil {
		return
	}
	if err = asset(PatternPort, port, "port"); err != nil {
		return
	}
	if err = asset("", user, "user"); err != nil {
		return
	}
	if err = asset("", password, "password"); err != nil {
		return
	}

	config, _ := auth.PasswordKey(user, password, ssh.InsecureIgnoreHostKey())
	addr := fmt.Sprintf("%s:%s", ip, port)
	return scp.NewClient(addr, &config), nil
}

func NewScpClientBySSH(sshClient *ssh.Client) (client scp.Client, err error) {
	return scp.NewClientBySSH(sshClient)
}

func RunCommandRemote(client *ssh.Client, command string) error {
	session, err := client.NewSession()
	if err != nil {
		return err
	}
	return session.Run(command)
}

func Scp(client scp.Client, localFile, remotePath, permissions string) error {
	err := client.Connect()
	if err != nil {
		return err
	}
	f, err := os.Open(localFile)
	if err != nil {
		return err
	}
	defer f.Close()
	return client.CopyFile(f, remotePath, permissions)
}

func RunScriptRemote(sshClient *ssh.Client, scriptFile, permissions string, needClean bool) error {
	f, err := os.Stat(scriptFile)
	if err != nil {
		return err
	}
	if f.IsDir() {
		return fmt.Errorf("%s is not a regular file", scriptFile)
	}

	scpClient, err := NewScpClientBySSH(sshClient)
	if err != nil {
		return err
	}
	defer scpClient.Close()

	createTime := time.Now().Unix()
	filename := fmt.Sprintf("%d_%s", createTime, path.Base(scriptFile))
	remotePath := path.Join(TmpDir, filename)
	err = Scp(scpClient, scriptFile, remotePath, permissions)
	if err != nil {
		return err
	}

	remoteScript := path.Join(TmpDir, filename)
	command := fmt.Sprintf("/bin/sh %s;", remoteScript)
	if needClean {
		command = fmt.Sprintf("%s rm %s;", command, remotePath)
	}
	return RunCommandRemote(sshClient, command)
}

func asset(pattern, s, name string) (err error) {
	if s == "" {
		err = fmt.Errorf("%s cannot be empty", name)
		return err
	}
	if pattern != "" {
		log.Println(pattern)
		var matched bool
		matched, err = regexp.MatchString(pattern, s)
		if err != nil {
			return err
		}
		if !matched {
			err = fmt.Errorf("wrong %s", name)
			return
		}
	}
	return nil
}
