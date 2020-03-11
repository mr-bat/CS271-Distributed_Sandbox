package main

import (
	"bufio"
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"strconv"
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

func getIdFromInput() int {
	fmt.Println("What is your id?")
	reader := bufio.NewScanner(os.Stdin)
	reader.Scan()
	message := reader.Text()
	id, _ := strconv.Atoi(message)

	return id
}

func getCommand() Command {
	reader := bufio.NewScanner(os.Stdin)

	reader.Scan()
	message := reader.Text()
	if message == "Transaction" {
		var to, amount int
		fmt.Println("Enter: to, amount")
		fmt.Scan(&to)
		fmt.Scan(&amount)

		Logger.WithFields(logrus.Fields{
			"from": getId(),
			"to": to,
			"amount": amount,
		}).Info("received transaction command")
		return Command{cType: TransactionCode, from: strconv.Itoa(getId()), to: strconv.Itoa(to), amount: amount}
	} else if message == "Balance" {
		var id int
		fmt.Println("Enter: id")
		fmt.Scan(&id)

		Logger.WithFields(logrus.Fields{
			"id": id,
		}).Info("received balance command")
		return Command{cType: BalanceCode, id: id}
	} else if message == "Reset" {
		return Command{cType: ResetDataCode}
	} else {
		fmt.Println("Available options:\n * Transaction\n * Balance\n * Reset")
		return Command{cType: UnknownCode}
	}
}