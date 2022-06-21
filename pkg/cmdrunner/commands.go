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
	client, err := devcon.NewClient(user, ia.Target, devcon.SetKey(key))
	if err != nil {
		return err
	}
	output, err := client.Run(cmd)
	if err != nil {
		return err
	}
	if err := ia.unmarshal([]byte(output)); err != nil {
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
