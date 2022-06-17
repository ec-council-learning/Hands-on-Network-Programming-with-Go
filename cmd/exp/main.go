package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/codered-by-ec-council/Hands-on-Network-Programming-with-Go/pkg/devcon"
)

func main() {
	target := flag.String("target", "127.0.0.1", "target against which to run a command")
	cmdFile := flag.String("cmdfile", "", "command filename")
	flag.Parse()
	// set hostKey
	knownhosts := filepath.Join(os.Getenv("HOME"), ".ssh", "known_hosts")
	client := devcon.NewClient(
		os.Getenv("SSH_USER"),
		*target,
		devcon.SetPassword(os.Getenv("SSH_PASSWORD")),
		devcon.SetHostKeyCallback(knownhosts),
	)
	f, err := os.Open(*cmdFile)
	if err != nil {
		log.Fatal(err)
	}
	bs, err := io.ReadAll(f)
	if err != nil {
		log.Fatal(err)
	}
	cmds := strings.Split(string(bs), "\n")
	output, err := client.RunAll(cmds)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(output)
}
