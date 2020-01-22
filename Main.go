package main

import (
	"fmt"
	"math/rand"
	"time"
)

const (
	UnknownCode = iota
	TransactionCode = iota
	BalanceCode = iota
)

type Command struct {
	cType int
	from, to, amount int
	id int
}

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	PortNumber = 7180 + rand.Intn(100)
	startServer(PortNumber)

	addrs := getClientAddrs()
	connectToClients(addrs)

	sendToClients(getAddress())
	for {
		command := getCommand()
		fmt.Println(command)
	}

	//flagMode := flag.String("initiator", "silent", "Start in initiator or silent mode")
	//flag.Parse()
	//if strings.ToLower(*flagMode) == "initiator" {
	//
	//}
}