package main

import (
	"bufio"
	"encoding/xml"
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
	output, err := client.Run("show isis adjacency | display xml")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(output)
	// isisAdj, err := parseXML(output)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println(isisAdj.IsisAdjacencyInformation.IsisAdjacency.SystemName)
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

type ISISAdjInfoRPCReply struct {
	// XMLName                  xml.Name `xml:"ISISAdjInfoRPCReply"`
	Junos                    string `xml:"junos,attr"`
	IsisAdjacencyInformation struct {
		Xmlns         string `xml:"xmlns,attr"`
		Style         string `xml:"style,attr"`
		IsisAdjacency struct {
			InterfaceName  string `xml:"interface-name"`
			SystemName     string `xml:"system-name"`
			Level          string `xml:"level"`
			AdjacencyState string `xml:"adjacency-state"`
			Holdtime       string `xml:"holdtime"`
		} `xml:"isis-adjacency"`
	} `xml:"isis-adjacency-information"`
}

func parseXML(in string) (ISISAdjInfoRPCReply, error) {
	var isisAdjInfo ISISAdjInfoRPCReply
	if err := xml.Unmarshal([]byte(in), &isisAdjInfo); err != nil {
		return isisAdjInfo, err
	}
	return isisAdjInfo, nil
}
