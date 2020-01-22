package main

import (
	"math/rand"
	"time"
)

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	PortNumber = 7180 + rand.Intn(100)
	startServer(PortNumber)

	addrs := getClientAddrs()
	connectToClients(addrs)

	//flagMode := flag.String("initiator", "silent", "Start in initiator or silent mode")
	//flag.Parse()
	//if strings.ToLower(*flagMode) == "initiator" {
	//
	//}
}