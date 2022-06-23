package main

import (
	"log"

	"github.com/codered-by-ec-council/Hands-on-Network-Programming-with-Go/pkg/cmdrunner"
)

func main() {
	// isisAdj := cmdrunner.ISISAdjacencyRpcReply{Target: "labsrx", ExpectedNeighbor: "lab_srx10"}
	// if err := isisAdj.Run(); err != nil {
	// 	log.Fatal(err)
	// }
	// err := isisAdj.Compare()
	// if err != nil {
	// 	log.Println(err)
	// }
	sr := cmdrunner.SpecificRouteRpcReply{Target: "labsrx", ExpectedNextHop: "192.168.0.1"}
	if err := sr.Run(); err != nil {
		log.Fatal(err)
	}
	if err := sr.Compare(); err != nil {
		log.Println(err)
	}
}
