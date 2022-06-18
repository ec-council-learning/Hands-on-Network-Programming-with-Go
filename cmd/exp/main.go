package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/codered-by-ec-council/Hands-on-Network-Programming-with-Go/pkg/devcon"
)

func main() {
	client, err := devcon.NewClient(
		os.Getenv("SSH_USER"),
		"labsrx",
		devcon.SetPassword(os.Getenv("SSH_PASSWORD")),
		devcon.SetTimeout(time.Second*5),
	)
	if err != nil {
		log.Fatal(err)
	}
	output, err := client.Run("show version")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(output)
}
