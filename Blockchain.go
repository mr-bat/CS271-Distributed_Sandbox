package main

import (
	"encoding/json"
	"reflect"
	"strings"
)

type Transaction struct {
	Sender, Receiver string
	Amount           int
}

type Block struct {
	SeqNum int
	Tx     []Transaction
}

var blockchain []Block
var pendingTx []Transaction

/*
func init() {
	blockchain = append(blockchain, make([]Block, 0))
}
*/
func calculateBalances() map[string]int {
	balance := make(map[string]int)

	for _, tx := range pendingTx {
		balance[tx.Receiver] += tx.Amount
		balance[tx.Sender] -= tx.Amount
	}

	for _, currblockchain := range blockchain {
		for _, block := range currblockchain.Tx {
			balance[block.Receiver] += block.Amount
			balance[block.Sender] -= block.Amount
		}
	}

	return balance
}

func (tx *Transaction) toString() string {
	res, _ := json.Marshal(tx)
	return string(res)
}

func (block *Block) toString() string {
	res, _ := json.Marshal(block)
	return string(res)
}

func parseBlock(block string) Block {
	var res Block
	if err := json.Unmarshal([]byte(block), &res); err != nil {
		panic(err)
	}
	return res
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

/*
func addBlock(sender, receiver string, amount int) {
	blockchain[len(blockchain)-1] = append(blockchain[len(blockchain)-1],
		Block{sender: sender, receiver: receiver, amount: amount})
}
*/
func addBlock(block Block) {
	blockchain = append(blockchain, block)
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
	return blockchain[len(blockchain)-1]
}

func clearCurrBlockChain() {
	blockchain[len(blockchain)-1] = nil
}

func getCurrSeqNumber() int {
	return len(blockchain)
}

func getBalance(user string) int {
	balances := calculateBalances()

	return balances[user] + 10
}

func getBlock(seqnum int) *Block {
	for _, block := range blockchain {
		if block.SeqNum == seqnum {
			return &block
		}
	}
	return nil
}
