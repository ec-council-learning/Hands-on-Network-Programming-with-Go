# Hands-on-Network-Programming-with-Go
Hands-on Network Programming with Go, by EC-Council


## Environment Variables

To run this project, you will need to add the following environment variables to your .env file

`SSH_USER`

`SSH_PASSWORD`

`POSTGRES_USER`

`POSTGRES_PASSWORD`

`DSN`

### For example:
```env
export POSTGRES_USER=postgres
export POSTGRES_PASSWORD=postgres
export DSN=postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@localhost:5432/device_inventory
```

## Run Locally

Clone the project

```bash
  git clone https://github.com/codered-by-ec-council/Hands-on-Network-Programming-with-Go.git
```

Go to the project directory

```bash
  cd Hands-on-Network-Programming-with-Go
```

Install dependencies

```bash
go mod download
```

Spin up the postgres container
```bash
# create the volume initially
docker-compose up
```

Start the web server

```bash
go run cmd/web/*.go
```

## Usage/Examples

### Fire up the Miminal server
```bash
cd cmd/server
go build -o custom-server
sudo -E ./custom-server
```

### Run Command via CLI
```golang
package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/codered-by-ec-council/Hands-on-Network-Programming-with-Go/pkg/devcon"
)

func main() {
	target := flag.String("target", "127.0.0.1", "target against which to run a command")
	cmd := flag.String("cmd", "", "command to run against target device")
	flag.Parse()
	client := devcon.NewClient(*target)
	output, err := client.Run(*cmd)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(output)
}
```

### Run a Change with cmdrunner
```golang
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
```