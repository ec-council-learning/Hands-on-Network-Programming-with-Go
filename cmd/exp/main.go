package main

import (
	"log"

	"github.com/codered-by-ec-council/Hands-on-Network-Programming-with-Go/pkg/cmdrunner"
)

func main() {
	isisAdj := cmdrunner.ISISAdjacencyRpcReply{Target: "labsrx", ExpectedNeighbor: "lab_srx10"}
	if err := isisAdj.Run(); err != nil {
		log.Fatal(err)
	}
	err := isisAdj.Compare()
	if err != nil {
		log.Println(err)
	}
}
