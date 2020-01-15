package main

import (
	"bufio"
	"fmt"
	"os"
)

const WorkerCount = 3


func work(id int)  {

}

func main()  {
	fmt.Printf("Number of workers %v\n", WorkerCount)
	fmt.Println("Input id of the client and then the type of request (for \"balance\" 1, for \"transfer\" 2)")
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}

	if scanner.Err() != nil {
		// handle error.
	}
}