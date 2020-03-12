package main

import (
	"context"
	"encoding/json"
	"golang.org/x/sync/semaphore"
	"reflect"
	"strconv"
)

type Transaction struct {
	Sender, Receiver string
	Amount, Id       int
}

type Block struct {
	SeqNum int
	Tx     []Transaction
}

var clock int
var BlockChainSemaphore *semaphore.Weighted
var blockchain []Block
var pendingTx []Transaction
var acceptedBlock Block

func initBlockChain() {
	clock = 0
	BlockChainSemaphore = semaphore.NewWeighted(int64(1))
	blkLength, _ := strconv.Atoi(getData("blkLength"))
	for i := 1; i <= blkLength; i++ {
		blockchain = append(blockchain, parseBlock(getData(strconv.Itoa(i))))
	}
	acceptedBlock = parseBlock(getData("accepted"))
	pendingTx = parseBlock(getData("pending")).Tx
}

func incClock() int {
	clock++
	return clock
}

func calculateBalances() map[string]int {
	balance := make(map[string]int)

	for _, tx := range pendingTx {
		balance[tx.Receiver] += tx.Amount
		balance[tx.Sender] -= tx.Amount
	}

	for _, tx := range acceptedBlock.Tx {
		balance[tx.Receiver] += tx.Amount
		balance[tx.Sender] -= tx.Amount
	}

	for _, currblockchain := range blockchain {
		for _, tx := range currblockchain.Tx {
			balance[tx.Receiver] += tx.Amount
			balance[tx.Sender] -= tx.Amount
		}
	}

	return balance
}

func (tx *Transaction) toString() string {
	res, _ := json.Marshal(tx)
	return string(res)
}

func parseTransaction(tx string) Transaction {
	var res Transaction
	if err := json.Unmarshal([]byte(tx), &res); err != nil {
		panic(err)
	}
	return res
}

func (block Block) toString() string {
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

func (block Block) isEmpty() bool {
	return block.SeqNum == 0
}

func rangeToString(txs []Transaction) string {
	res, _ := json.Marshal(txs)
	return string(res)
}

func parseRange(txs string) []Transaction {
	var res []Transaction
	if err := json.Unmarshal([]byte(txs), &res); err != nil {
		panic(err)
	}
	return res
}

func (block Block) merge(_block Block) Block {
	if block.SeqNum != _block.SeqNum {
		panic("merge: seqNumbers don't match")
	}

	var mergedTxs []Transaction
	mergedTxs = append(mergedTxs, block.Tx...)
	mergedTxs = append(mergedTxs, _block.Tx...)

	return Block{
		SeqNum: block.SeqNum,
		Tx:     mergedTxs,
	}
}

func addTransaction(tx Transaction) {
	pendingTx = append(pendingTx, tx)
	storeData("pending", Block{Tx: pendingTx}.toString())
}

func commitBlock(block Block) {
	BlockChainSemaphore.Acquire(context.Background(), 1)
	currTransaction := getCurrTransactions()
	newTransactions := make([]Transaction, 0)

	for _, tx := range currTransaction {
		shouldAdd := true
		for _, _tx := range block.Tx {
			if reflect.DeepEqual(tx, _tx) {
				shouldAdd = false
				break
			}
		}

		if shouldAdd {
			newTransactions = append(newTransactions, tx)
		}
	}

	blockchain = append(blockchain, block)
	storeData("blkLength", strconv.Itoa(len(blockchain)))
	storeData(strconv.Itoa(len(blockchain)), block.toString())


	clearCurrTransactions()
	for _, tx := range newTransactions {
		addTransaction(tx)
	}
	BlockChainSemaphore.Release(1)
}

func getCurrTransactions() []Transaction {
	return pendingTx
}

func clearCurrTransactions() {
	pendingTx = nil
	storeData("pending", Block{Tx: pendingTx}.toString())
}

func clearPersistedData() {
	blockchain = nil
	acceptedBlock = Block{}
	storeData("blkLength", strconv.Itoa(0))
	storeData("accepted", acceptedBlock.toString())

	reset()
}

func createNewBlock() Block {
	return Block{
		SeqNum: getCurrSeqNumber(),
		Tx:     pendingTx,
	}
}

func getLastBlock() Block {
	if len(blockchain) == 0 {
		return Block{SeqNum: 0, Tx: nil}
	}

	return blockchain[len(blockchain) - 1]
}

func getCurrSeqNumber() int {
	if len(blockchain) == 0 {
		return 1
	}

	return blockchain[len(blockchain) - 1].SeqNum + 1
}

func getBalance(user string) int {
	balances := calculateBalances()

	return balances[user] + 10
}

func getBlock(seqnum int) Block {
	for _, block := range blockchain {
		if block.SeqNum == seqnum {
			return block
		}
	}
	return Block{}
}
