package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

const (
	UnknownCode     = iota
	TransactionCode = iota
	BalanceCode     = iota
	ResetDataCode	= iota
)

type Command struct {
	cType      int
	from, to   string
	amount, id int
}

func handleCommand(command Command) {
	if command.cType == UnknownCode {
	} else if command.cType == TransactionCode {
		addPurchase(command.from, command.to, command.amount)
	} else if command.cType == BalanceCode {
		fmt.Println("User balance:", getBalance(strconv.Itoa(command.id)))
	} else if command.cType == ResetDataCode {
		clearCurrTransactions()
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
	advertiseId()
	lastBallot = Ballot{0, getId()}

	for {
		fmt.Println("Please enter your command: ")
		command := getCommand()
		handleCommand(command)
	}
}
