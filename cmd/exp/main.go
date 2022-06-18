package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/codered-by-ec-council/Hands-on-Network-Programming-with-Go/pkg/devcon"
)

func main() {
	client, err := devcon.NewClient(
		"ec2-user",
		"3.141.97.255",
		devcon.SetKey(filepath.Join(os.Getenv("CODERED"), "aws-jnpr-lab.pem")),
		devcon.SetTimeout(time.Second*5),
	)
	if err != nil {
		log.Fatal(err)
	}
	output, err := client.Run("show interfaces terse | display json")
	if err != nil {
		log.Fatal(err)
	}
	// fmt.Println(output)
	intTerse, err := parseJSON(output)
	if err != nil {
		log.Fatal(err)
	}
	for _, phy := range intTerse.InterfaceInformation {
		for _, intf := range phy.PhysicalInterface {
			fmt.Println(intf.Name[0].Data, intf.AdminStatus[0].Data, intf.OperStatus[0].Data)
		}
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

type InterfaceTerse struct {
	InterfaceInformation []struct {
		Attributes struct {
			Xmlns      string `json:"xmlns"`
			JunosStyle string `json:"junos:style"`
		} `json:"attributes"`
		PhysicalInterface []struct {
			Name []struct {
				Data string `json:"data"`
			} `json:"name"`
			AdminStatus []struct {
				Data string `json:"data"`
			} `json:"admin-status"`
			OperStatus []struct {
				Data string `json:"data"`
			} `json:"oper-status"`
			LogicalInterface []struct {
				Name []struct {
					Data string `json:"data"`
				} `json:"name"`
				AdminStatus []struct {
					Data string `json:"data"`
				} `json:"admin-status"`
				OperStatus []struct {
					Data string `json:"data"`
				} `json:"oper-status"`
				FilterInformation []struct {
				} `json:"filter-information"`
				AddressFamily []struct {
					AddressFamilyName []struct {
						Data string `json:"data"`
					} `json:"address-family-name"`
				} `json:"address-family"`
			} `json:"logical-interface,omitempty"`
		} `json:"physical-interface"`
	} `json:"interface-information"`
}

func parseJSON(in string) (InterfaceTerse, error) {
	var interfaceTerse InterfaceTerse
	if err := json.Unmarshal([]byte(in), &interfaceTerse); err != nil {
		return interfaceTerse, err
	}
	return interfaceTerse, nil
}
