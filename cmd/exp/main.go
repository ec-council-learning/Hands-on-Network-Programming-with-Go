package main

import (
	"fmt"
	"os"

	"github.com/codered-by-ec-council/Hands-on-Network-Programming-with-Go/pkg/devcon"
)

func main() {
	e7 := "10.224.132.17"
	cl, err := devcon.NewClient(
		"chh973", e7,
		devcon.SetPassword(os.Getenv("FTR_PASSWORD")),
		devcon.SetKeyExchange("diffie-hellman-group1-sha1"),
		devcon.SetCipher("3des-cbc"),
	)
	if err != nil {
		panic(err)
	}
	alarms, err := cl.RunAll([]string{"show alarm", "exit"})
	if err != nil {
		panic(err)
	}
	fmt.Println(alarms)
	// isisAdj := &cmdrunner.ISISAdjacencyRpcReply{Target: "labsrx", ExpectedNeighbor: "lab_srx100"}
	// sr := &cmdrunner.SpecificRouteRpcReply{Target: "labsrx", ExpectedNextHop: "192.168.0.1"}
	// cmds := []cmdrunner.Runner{isisAdj, sr}
	// if err := cmdrunner.Stepper(cmds); err != nil {
	// 	log.Println(err)
	// }
}
