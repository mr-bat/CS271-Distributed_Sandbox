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
	ConnectCode		= iota
	DisconnectCode  = iota
	PrintCode		= iota
)

type Command struct {
	cType      int
	from, to   string
	amount, id int
}

var Connected = true
func handleCommand(command Command) {
	if command.cType == ConnectCode {
		Connected = true
		return
	}
	if !Connected {
		return
	}

	if command.cType == UnknownCode {
	} else if command.cType == TransactionCode {
		addPurchase(command.from, command.to, command.amount)
	} else if command.cType == BalanceCode {
		fmt.Println("User balance:", getBalance(strconv.Itoa(command.id)))
	} else if command.cType == PrintCode {
		fmt.Println("Printing blockchain")
		for i, block := range blockchain {
			fmt.Printf("Blk %v: %v\n", i + 1, block)
		}
		fmt.Printf("Accepted blk: %v\n", acceptedBlock)
		fmt.Printf("Pending Txs: %v\n", pendingTx)
	} else if command.cType == ResetDataCode {
		//clearCurrTransactions()
		clearPersistedData()
	}  else if command.cType == DisconnectCode {
		Connected = false
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
	initBlockChain()
	lastBallot = Ballot{0, getId()}

	for {
		fmt.Println("Please enter your command: ")
		command := getCommand()
		handleCommand(command)
	}
}
