package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/codered-by-ec-council/Hands-on-Network-Programming-with-Go/pkg/devcon"
)

func main() {
	target := flag.String("target", "127.0.0.1", "target against which to run a command")
	cmd := flag.String("cmd", "", "command to run against target device")
	flag.Parse()
	client := devcon.NewClient(*target)
	output, err := client.Run(*cmd)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(output)
}

// func run(target, cmd string) (string, error) {
// 	cfg := &ssh.ClientConfig{
// 		User: os.Getenv("SSH_USER"),
// 		Auth: []ssh.AuthMethod{
// 			ssh.Password(os.Getenv("SSH_PASSWORD")),
// 		},
// 		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
// 		Timeout:         time.Second * 5,
// 	}
// 	client, err := ssh.Dial("tcp", fmt.Sprintf("%v:22", target), cfg)
// 	if err != nil {
// 		return "", errors.Wrap(err, "dial failed")
// 	}
// 	defer client.Close()
// 	session, err := client.NewSession()
// 	if err != nil {
// 		return "", err
// 	}
// 	defer session.Close()
// 	output, err := session.CombinedOutput(cmd)
// 	if err != nil {
// 		return "", err
// 	}
// 	return string(output), nil
// }
