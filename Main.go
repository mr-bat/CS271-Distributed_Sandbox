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

func handleCommand(command Command) {
	if command.cType == UnknownCode {
	} else if command.cType == TransactionCode {
		if addTransaction(command.from, command.to, command.amount) {
			fmt.Println("SUCCESS")
		} else {
			fmt.Println("INCORRECT")
		}
	} else if command.cType == BalanceCode {
		fmt.Println("User balance:", getUserBalance(command.id))
	} else {
		fmt.Println("Unknown Command")
	}
}

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	PortNumber = 7180 + rand.Intn(100)
	startServer(PortNumber)

	addrs := getClientAddrs()
	connectToClients(addrs)

	fmt.Println("Please enter your command: ")
	for {
		command := getCommand()
		AcquireLock()
		println("acquired")
		handleCommand(command)
		println("handled")
		ReleaseLock()
		println("released")
	}
}