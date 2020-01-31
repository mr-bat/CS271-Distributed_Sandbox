package main

type Block struct {
	sender, receiver string
	amount int
}

var Blockchain []Block

func calculateBalances() map[string]int {
	balance := make(map[string]int)

	for _, block := range Blockchain {
		balance[block.receiver] += block.amount
		balance[block.sender] -= block.amount
	}

	return balance
}

func addBlock(sender, receiver string, amount int) {
	Blockchain = append(Blockchain, Block{sender: sender, receiver: receiver, amount: amount})
}

func getBalance(user string) int {
	balances := calculateBalances()

	return balances[user] + 10
}