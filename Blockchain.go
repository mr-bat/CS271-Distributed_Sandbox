package main

import (
	"fmt"
	"strconv"
	"strings"
)

type Block struct {
	sender, receiver string
	amount, time int
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

func (block *Block) toString() string {
	return fmt.Sprintf("%s&%s&%d&%d\n", block.sender, block.receiver, block.amount, block.time)
}

func parseBlock(block string) Block{
	parsed := strings.Split(block, "&")
	amount, _ := strconv.Atoi(parsed[2])
	time, _ := strconv.Atoi(parsed[3])

	return Block{sender: parsed[0], receiver: parsed[1], amount: amount, time: time}
}

func rangeToString(start int) string {
	result := ""

	for i := start; i < len(Blockchain); i++ {
		result += Blockchain[i].toString()
	}

	return result
}

func parseRange(blocks string) []Block {
	parsedBlocks := strings.Split(blocks, "\n")
	var createdBlocks []Block

	for _, block := range parsedBlocks {
		if len(block) > 0 {
			createdBlocks = append(createdBlocks, parseBlock(block))
		}
	}

	return createdBlocks
}

func addBlock(sender, receiver string, amount int) {
	Blockchain = append(Blockchain, Block{sender: sender, receiver: receiver, amount: amount, time: incTime()})
}

func getBalance(user string) int {
	balances := calculateBalances()

	return balances[user] + 10
}