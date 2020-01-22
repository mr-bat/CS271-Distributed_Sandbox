package main

import (
	"bufio"
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
)

func waitForDone() {
	fmt.Println("Awaiting \"done\"")
	reader := bufio.NewScanner(os.Stdin)
	message := "await"
	for message != "done" {
		reader.Scan()
		message = reader.Text()
	}
}

func getCommand() Command {
	reader := bufio.NewScanner(os.Stdin)

	reader.Scan()
	message := reader.Text()
	if message == "Transaction" {
		var from, to, amount int
		fmt.Println("Enter: from, to, amount")
		fmt.Scan(&from)
		fmt.Scan(&to)
		fmt.Scan(&amount)

		Logger.WithFields(logrus.Fields{
			"from": from,
			"to": to,
			"amount": amount,
		}).Info("received transaction command")
		return Command{cType: TransactionCode, from: from, to: to, amount: amount}
	} else if message == "Balance" {
		var id int
		fmt.Println("Enter: id")
		fmt.Scan(&id)

		Logger.WithFields(logrus.Fields{
			"id": id,
		}).Info("received balance command")
		return Command{cType: BalanceCode, id: id}
	} else {
		fmt.Println("Available options:\n * Transaction\n * Balance\n")
		return Command{cType: UnknownCode}
	}
}