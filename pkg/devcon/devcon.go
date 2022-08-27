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

// sshClient contains enough information about a device
// to remote access the device using the ssh protocol.
type sshClient struct {
	target string
	port   string
	cfg    *ssh.ClientConfig
}

// option is a functional paradigm to apply options
// to a factory function.
type option func() (func(*sshClient), error)

// failedOption deals with option error cases.
func failedOption(err error) option {
	return func() (func(*sshClient), error) {
		return nil, err
	}
}

// properOption is the success case for an option.
func properOption(setter func(*sshClient)) option {
	return func() (func(*sshClient), error) {
		return setter, nil
	}
}

// SetTimeout applies a timeout in seconds for the client connection.
func SetTimeout(timeout time.Duration) option {
	return properOption(func(c *sshClient) {
		c.cfg.Timeout = timeout
	})
}

// SetPort applies a port to the client. If not set, the default port
// of 22 is applied.
func SetPort(port string) option {
	return properOption(func(c *sshClient) {
		c.port = port
	})
}

// SetPassword applies a password to the client.
func SetPassword(pw string) option {
	return properOption(func(c *sshClient) {
		authMethod := []ssh.AuthMethod{
			ssh.Password(pw),
		}
		c.cfg.Auth = authMethod
	})
}

// SetHostKeyCallback takes in a known hosts file, usually an absolute path,
// and applies it to the client as an option to protect against
// man-in-the-middle attacks.
func SetHostKeyCallback(knownhostsFile string) option {
	hostKeyCallback, err := knownhosts.New(knownhostsFile)
	if err != nil {
		return failedOption(err)
	}
	return properOption(func(c *sshClient) {
		c.cfg.HostKeyCallback = hostKeyCallback
	})
}

// SetKey takes in private key file, usually an absolute path,
// and applies it to the client as an option.
func SetKey(keyfile string) option {
	f, err := os.Open(keyfile)
	if err != nil {
		return failedOption(err)
	}
	defer f.Close()
	bs, err := io.ReadAll(f)
	if err != nil {
		return failedOption(err)
	}
	signer, err := ssh.ParsePrivateKey(bs)
	if err != nil {
		return failedOption(err)
	}
	return properOption(func(c *sshClient) {
		authMethod := []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		}
		c.cfg.Auth = authMethod
	})
}

// setup runs through any supplied options checking for
// any errors, returning nil if there are none.
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

// NewClient is a factory function that takes in a a username, a target
// in the form of an IP or a hostname, and potentially options.
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

// Run takes in a command, opens a remote connection to a target
// device, and runs in against the device.
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

// RunAll takes in a slice of strings, normally a configuration, establishes
// an interactive session with the target, running in each command
// line by line.
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
