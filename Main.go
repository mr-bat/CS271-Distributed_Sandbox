package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"strconv"
	"syscall"
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
	BenchmarkCode = iota
)

type Command struct {
	cType      int
	from, to   string
	amount, id int
}
const Debugging = false
var yield = make(chan bool)
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
			fmt.Printf("Blk %v: %v\n", i+1, block)
		}
	} else if command.cType == ResetDataCode {
		//clearCurrTransactions()
	} else if command.cType == BenchmarkCode {
		beginBenchmark()
		//<- yield
	} else if command.cType == DisconnectCode {
		Connected = false
	} else {
		fmt.Println("Unknown Command")
	}
}

func main() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)
	go func() {
		_ = <- sigChan
		clearData()
		for _, addr := range getClientAddrs() {
			removeServerAddr(addr)
		}
		os.Exit(0)
	}()

	rand.Seed(time.Now().UTC().UnixNano())
	PortNumber = 7180 + rand.Intn(500)
	startServer(PortNumber)

	addrs := getClientAddrs()
	connectToClients(addrs)
	advertiseId()
	initBlockChain()

	for {
		fmt.Println("Please enter your command: ")
		command := getCommand()
		handleCommand(command)
	}
}
