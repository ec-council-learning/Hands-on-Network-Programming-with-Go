package devcon

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/pkg/errors"
	"golang.org/x/crypto/ssh"
)

type sshClient struct {
	target string
	cfg    *ssh.ClientConfig
}

func NewClient(target string) *sshClient {
	return &sshClient{
		target: target,
		cfg: &ssh.ClientConfig{
			User: os.Getenv("SSH_USER"),
			Auth: []ssh.AuthMethod{
				ssh.Password(os.Getenv("SSH_PASSWORD")),
			},
			HostKeyCallback: ssh.InsecureIgnoreHostKey(),
			Timeout:         time.Second * 5,
		},
	}
}

func (c *sshClient) Run(cmd string) (string, error) {
	client, err := ssh.Dial("tcp", fmt.Sprintf("%v:22", c.target), c.cfg)
	if err != nil {
		return "", errors.Wrap(err, "dial failed")
	}
	defer client.Close()
	session, err := client.NewSession()
	if err != nil {
		return "", err
	}
	defer session.Close()
	output, err := session.CombinedOutput(cmd)
	if err != nil {
		return "", err
	}
	return string(output), nil
}

func (c *sshClient) RunAll(cmds []string) (string, error) {
	client, err := ssh.Dial("tcp", fmt.Sprintf("%v:22", c.target), c.cfg)
	if err != nil {
		return "", errors.Wrap(err, "dial failed")
	}
	defer client.Close()
	session, err := client.NewSession()
	if err != nil {
		return "", err
	}
	defer session.Close()
	wc, err := session.StdinPipe()
	if err != nil {
		return "", err
	}
	defer wc.Close()
	r, err := session.StdoutPipe()
	if err != nil {
		return "", err
	}
	if err := session.Shell(); err != nil {
		return "", err
	}
	for _, cmd := range cmds {
		_, err := fmt.Fprintf(wc, "%s\n", cmd)
		if err != nil {
			return "", err
		}
	}
	bs, err := io.ReadAll(r)
	if err != nil {
		return "", err
	}
	return string(bs), nil
}
