package devcon

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/pkg/errors"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/knownhosts"
)

type sshClient struct {
	target string
	port   string
	cfg    *ssh.ClientConfig
}

type option func() (func(*sshClient), error)

func failedOption(err error) option {
	return func() (func(*sshClient), error) {
		return nil, err
	}
}

func properOption(setter func(*sshClient)) option {
	return func() (func(*sshClient), error) {
		return setter, nil
	}
}

func SetPort(port string) option {
	return properOption(func(c *sshClient) {
		c.port = port
	})
}

func SetPassword(pw string) option {
	return properOption(func(c *sshClient) {
		authMethod := []ssh.AuthMethod{
			ssh.Password(pw),
		}
		c.cfg.Auth = authMethod
	})
}

func SetHostKeyCallback(knownhostsFile string) option {
	hostKeyCallback, err := knownhosts.New(knownhostsFile)
	if err != nil {
		failedOption(err)
	}
	return properOption(func(c *sshClient) {
		c.cfg.HostKeyCallback = hostKeyCallback
	})
}

func SetKey(keyfile string) option {
	f, err := os.Open(keyfile)
	if err != nil {
		failedOption(err)
	}
	defer f.Close()
	bs, err := io.ReadAll(f)
	if err != nil {
		failedOption(err)
	}
	signer, err := ssh.ParsePrivateKey(bs)
	if err != nil {
		failedOption(err)
	}
	return properOption(func(c *sshClient) {
		authMethod := []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		}
		c.cfg.Auth = authMethod
	})
}

func (c *sshClient) setup(opts ...option) error {
	if c == nil {
		return nil
	}
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		setter, err := opt()
		if err != nil {
			return err
		}
		if setter != nil {
			setter(c)
		}
	}
	return nil
}

func NewClient(user, target string, opts ...option) (*sshClient, error) {
	defaultPort := "22"
	client := &sshClient{
		port:   defaultPort,
		target: target,
		cfg: &ssh.ClientConfig{
			User:            user,
			HostKeyCallback: ssh.InsecureIgnoreHostKey(),
			Timeout:         time.Second * 5,
		},
	}
	if err := client.setup(opts...); err != nil {
		return client, err
	}
	return client, nil
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
