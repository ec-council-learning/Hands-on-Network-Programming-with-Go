package main

import (
	"log"

	"github.com/codered-by-ec-council/Hands-on-Network-Programming-with-Go/pkg/cmdrunner"
)

func main() {
	isisAdj := &cmdrunner.ISISAdjacencyRpcReply{Target: "labsrx", ExpectedNeighbor: "lab_srx100"}
	sr := &cmdrunner.SpecificRouteRpcReply{Target: "labsrx", ExpectedNextHop: "192.168.0.1"}
	cmds := []cmdrunner.Runner{isisAdj, sr}
	if err := cmdrunner.Stepper(cmds); err != nil {
		log.Println(err)
	}
}
