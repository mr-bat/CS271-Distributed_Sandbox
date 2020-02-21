package main

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

type Block struct {
	sender, receiver string
	amount     		 int
}

var blockchain [][]Block

func init()  {
	blockchain = append(blockchain, make([]Block, 0))
}

func calculateBalances() map[string]int {
	balance := make(map[string]int)

	for _, currblockchain := range blockchain {
		for _, block := range currblockchain {
			balance[block.receiver] += block.amount
			balance[block.sender] -= block.amount
		}
	}

	return balance
}

func (block *Block) toString() string {
	return fmt.Sprintf("%s&%s&%d\n", block.sender, block.receiver, block.amount)
}

func parseBlock(block string) Block {
	parsed := strings.Split(block, "&")
	amount, _ := strconv.Atoi(parsed[2])

	return Block{sender: parsed[0], receiver: parsed[1], amount: amount}
}

func rangeToString(blocks []Block) string {
	result := ""

	for _, block := range blocks {
		result += block.toString()
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
	blockchain[len(blockchain) - 1] = append(blockchain[len(blockchain) - 1],
		Block{sender: sender, receiver: receiver, amount: amount})
}

func addBlockchain(blocks []Block) {
	blockchain[len(blockchain) - 1] = append(blockchain[len(blockchain) - 1], blocks...)
}

func commitBlockchain(blocks []Block) {
	sampleBlock := getCurrBlockChain()[0]

	containsCurrBlock := false
	for _, block := range blocks {
		if reflect.DeepEqual(block, sampleBlock) {
			containsCurrBlock = true
			break
		}
	}

	var newBlock []Block
	if !containsCurrBlock {
		newBlock = getCurrBlockChain()
	}

	clearCurrBlockChain()
	addBlockchain(blocks)
	blockchain = append(blockchain, newBlock)
}

func getCurrBlockChain() []Block {
	return blockchain[len(blockchain) - 1]
}

func clearCurrBlockChain() {
	blockchain[len(blockchain) - 1] = nil
}

func getCurrSeqNumber() int {
	return len(blockchain)
}

func getBalance(user string) int {
	balances := calculateBalances()

	return balances[user] + 10
}
