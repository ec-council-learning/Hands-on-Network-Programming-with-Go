package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
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
	output, err := client.Run("show isis adjacency")
	if err != nil {
		log.Fatal(err)
	}
	isisAdj := parseRegex(output)
	for k, v := range isisAdj {
		fmt.Printf("%v: %v\n", k, v)
	}
}

func parseLineSplit(in string) map[string]string {
	var keys []string
	isisAdj := make(map[string]string)
	s := bufio.NewScanner(strings.NewReader(in))
	for s.Scan() {
		ln := s.Text()
		if strings.Contains(ln, "Interface") {
			fields := strings.Fields(ln)
			for i := 0; i < 4; i++ {
				keys = append(keys, fields[i])
			}
		}
		fields := strings.Fields(ln)
		for i := 0; i < 4; i++ {
			isisAdj[keys[i]] = fields[i]
		}
	}
	return isisAdj
}

func parseRegex(in string) map[string]string {
	isisAdj := make(map[string]string)
	pat := regexp.MustCompile(`(?P<intf>[fg]e-\d\/\d/\d\.\d+)\s{2,}(?P<hostname>[a-zA-Z_0-9]+)\s{2,}(?P<level>\d)\s{2,}(?P<state>[a-zA-Z]+)`)
	match := pat.FindStringSubmatch(in)
	for idx, name := range pat.SubexpNames() {
		if idx > 0 && idx <= len(match) {
			isisAdj[name] = match[idx]
		}
	}
	return isisAdj
}
