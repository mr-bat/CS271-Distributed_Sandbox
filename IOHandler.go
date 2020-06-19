package main

import (
	"bufio"
	"fmt"
	"github.com/manifoldco/promptui"
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
	fmt.Println("What is your id? (one base)")
	reader := bufio.NewScanner(os.Stdin)
	reader.Scan()
	message := reader.Text()
	id, _ := strconv.Atoi(message)

	return id
}

func getInput() string{
	fmt.Printf("number of participants: %v\n\n\n\n\n", GetNumberOfClients() + 1)
	prompt := promptui.Select{
		Label: "Select one",
		Items: []string{"Benchmark", "Transaction", "Balance", "Print", "Reset", "Connect", "Disconnect"},
	}

	_, result, err := prompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v", err)
		panic(err)
	}

	return result
}

func getCommand() Command {
	message := getInput()
	if message == "Transaction" {
		var to, amount int
		fmt.Println("Enter: to, amount")
		fmt.Scan(&to)
		fmt.Scan(&amount)

		Logger.WithFields(logrus.Fields{
			"from":   getId(),
			"to":     to,
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
	} else if message == "Benchmark" {
		return Command{cType: BenchmarkCode}
	} else if message == "Print" {
		return Command{cType: PrintCode}
	} else if message == "Reset" {
		return Command{cType: ResetDataCode}
	} else if message == "Connect" {
		return Command{cType: ConnectCode}
	} else if message == "Disconnect" {
		return Command{cType: DisconnectCode}
	} else {
		return Command{cType: UnknownCode}
	}
}