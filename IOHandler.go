package main

import (
	"bufio"
	"fmt"
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

func getCommand() string {
	reader := bufio.NewScanner(os.Stdin)

	reader.Scan()
	message := reader.Text()
	if message == "Transaction" {
		var from, to, amount int
		fmt.Println("Enter: from, to, amount")
		fmt.Scan(&from)
		fmt.Scan(&to)
		fmt.Scan(&amount)
	} else if message == "Balance" {
		var id int
		fmt.Println("Enter: id")
		fmt.Scan(&id)
	} else {
		fmt.Println("Available options:\n * Transaction\n * Balance\n")
	}
	return message
}