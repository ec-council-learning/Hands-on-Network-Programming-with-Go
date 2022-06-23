package cmdrunner

import (
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"

	"github.com/codered-by-ec-council/Hands-on-Network-Programming-with-Go/pkg/devcon"
)

var (
	user = os.Getenv("SSH_USER")
	key  = filepath.Join(os.Getenv("HOME"), ".ssh", "id_rsa")
)

type CompareError struct {
	got      string
	expected string
}

func (ce *CompareError) Error() string {
	return fmt.Sprintf("compare error - got: %q expected: %q", ce.got, ce.expected)
}

func getOutput(target, cmd string) (string, error) {
	client, err := devcon.NewClient(user, target, devcon.SetKey(key))
	if err != nil {
		return "", err
	}
	output, err := client.Run(cmd)
	if err != nil {
		return "", err
	}
	return output, nil
}

func unmarshal(bs []byte, data interface{}) error {
	return xml.Unmarshal(bs, data)
}

type ISISAdjacencyRpcReply struct {
	Target                   string
	ExpectedNeighbor         string
	IsisAdjacencyInformation struct {
		Text          string `xml:",chardata"`
		Xmlns         string `xml:"xmlns,attr"`
		Style         string `xml:"style,attr"`
		IsisAdjacency struct {
			Text           string `xml:",chardata"`
			InterfaceName  string `xml:"interface-name"`
			SystemName     string `xml:"system-name"`
			Level          string `xml:"level"`
			AdjacencyState string `xml:"adjacency-state"`
			Holdtime       string `xml:"holdtime"`
		} `xml:"isis-adjacency"`
	} `xml:"isis-adjacency-information"`
}

func (ia *ISISAdjacencyRpcReply) unmarshal(bs []byte) error {
	return xml.Unmarshal(bs, &ia)
}

func (ia *ISISAdjacencyRpcReply) Run() error {
	cmd := "show isis adjacency | display xml"
	output, err := getOutput(ia.Target, cmd)
	if err != nil {
		return err
	}
	if err := unmarshal([]byte(output), ia); err != nil {
		return err
	}
	return nil
}

func (ia *ISISAdjacencyRpcReply) Compare() error {
	gotSystemName := ia.IsisAdjacencyInformation.IsisAdjacency.SystemName
	if gotSystemName != ia.ExpectedNeighbor {
		return &CompareError{got: gotSystemName, expected: ia.ExpectedNeighbor}
	}
	return nil
}

type SpecificRouteRpcReply struct {
	Target           string
	Prefix           string
	ExpectedNextHop  string
	RouteInformation struct {
		Text       string `xml:",chardata"`
		Xmlns      string `xml:"xmlns,attr"`
		RouteTable struct {
			Text               string `xml:",chardata"`
			TableName          string `xml:"table-name"`
			DestinationCount   string `xml:"destination-count"`
			TotalRouteCount    string `xml:"total-route-count"`
			ActiveRouteCount   string `xml:"active-route-count"`
			HolddownRouteCount string `xml:"holddown-route-count"`
			HiddenRouteCount   string `xml:"hidden-route-count"`
			Rt                 struct {
				Text          string `xml:",chardata"`
				Style         string `xml:"style,attr"`
				RtDestination string `xml:"rt-destination"`
				RtEntry       struct {
					Text          string `xml:",chardata"`
					ActiveTag     string `xml:"active-tag"`
					CurrentActive string `xml:"current-active"`
					LastActive    string `xml:"last-active"`
					ProtocolName  string `xml:"protocol-name"`
					Preference    string `xml:"preference"`
					Age           struct {
						Text    string `xml:",chardata"`
						Seconds string `xml:"seconds,attr"`
					} `xml:"age"`
					Metric string `xml:"metric"`
					Nh     struct {
						Text            string `xml:",chardata"`
						SelectedNextHop string `xml:"selected-next-hop"`
						To              string `xml:"to"`
						Via             string `xml:"via"`
					} `xml:"nh"`
				} `xml:"rt-entry"`
			} `xml:"rt"`
		} `xml:"route-table"`
	} `xml:"route-information"`
}

func (sr *SpecificRouteRpcReply) Run() error {
	cmd := fmt.Sprintf("show route %s | display xml", sr.Prefix)
	output, err := getOutput(sr.Target, cmd)
	if err != nil {
		return err
	}
	if err := unmarshal([]byte(output), sr); err != nil {
		return err
	}
	return nil
}

func (sr *SpecificRouteRpcReply) Compare() error {
	gotNextHop := sr.RouteInformation.RouteTable.Rt.RtEntry.Nh.To
	if gotNextHop != sr.ExpectedNextHop {
		return &CompareError{got: gotNextHop, expected: sr.ExpectedNextHop}
	}
	return nil
}
