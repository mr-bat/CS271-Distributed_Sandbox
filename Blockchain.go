package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/looplab/tarjan"
	"golang.org/x/sync/semaphore"
)

type Transaction struct {
	Sender, Receiver string
	Amount, Id       int
}

type Block struct {
	SeqNum int
	Deps   []int
	Tx     []Transaction
}

const (
	BlkUnknown    = iota
	BlkReceived   = iota
	BlkUnexecuted = iota
	BlkCommitted  = iota
)

var clock int
var seqNum int
var seenBlks map[int]int
var depCnt map[int]int
var dependants map[int][]int
var readyToExecute []int
var BlockChainSemaphore *semaphore.Weighted
var Epaxos *semaphore.Weighted
var blockchain []Block
var unexecutedBlocks []Block

func initBlockChain() {
	clock = -1
	seqNum = -1
	seenBlks = make(map[int]int)
	depCnt = make(map[int]int)
	dependants = make(map[int][]int)
	readyToExecute = make([]int, 0)
	BlockChainSemaphore = semaphore.NewWeighted(int64(1))
	Epaxos = semaphore.NewWeighted(int64(1))
}

func incClock() int {
	clock++
	return clock * (GetNumberOfClients() + 1) + getId()
}

func incSeqNum() int {
	seqNum++
	return seqNum * (GetNumberOfClients() + 1) + getId()
}

func contains(s []int, e int) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func unique(intSlice []int) []int {
	keys := make(map[int]bool)
	list := []int{}
	for _, entry := range intSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

func calculateBalances() map[string]int {
	balance := make(map[string]int)

	for _, currblockchain := range blockchain {
		for _, tx := range currblockchain.Tx {
			balance[tx.Receiver] += tx.Amount
			balance[tx.Sender] -= tx.Amount
		}
	}

	return balance
}

func hasConflict(txs []Transaction, participant string) bool {
	for _, transaction := range txs {
		if transaction.Receiver == participant || transaction.Sender == participant {
			return true
		}
	}

	return false
}

func calculateDependencies(txs []Transaction) []int {
	deps := make([]int, 0)

	for _, tx := range txs {
		if tx.Sender == "0" {
			continue
		}
		for i := len(unexecutedBlocks) - 1; i > -1; i-- {
			if hasConflict(unexecutedBlocks[i].Tx, tx.Sender) {
				deps = append(deps, unexecutedBlocks[i].SeqNum)
				break
			}
		}

		//for i := len(unexecutedBlocks) - 1; i > -1; i-- {
		//	if hasConflict(unexecutedBlocks[i].Tx, tx.Receiver) {
		//		deps = append(deps, unexecutedBlocks[i].SeqNum)
		//		break
		//	}
		//}
	}

	return unique(deps)
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

func fillDependants(block Block) {
	deps := 0
	for _, dep := range block.Deps {
		if seenBlks[dep] != BlkCommitted {
			deps++
			dependants[dep] = append(dependants[dep], block.SeqNum)
		}
	}
	depCnt[block.SeqNum] = deps
}

var nonTrivialComponents int
func tryExecuteTarjan() {
	graph := make(map[interface{}][]interface{})
	for _, blk := range unexecutedBlocks {
		adj := make([]interface{}, len(blk.Deps))
		for i, v := range blk.Deps {
			adj[i] = v
		}
		graph[blk.SeqNum] = adj
	}
	scc := tarjan.Connections(graph)

	i := 0
	for ; i < len(scc); i++ {
		for _, _v := range scc[i] {
			if v, ok := _v.(int); ok {
				if seenBlks[v] == BlkUnknown || seenBlks[v] == BlkReceived {
					goto removeExecuted
				}
			} else {
				panic(fmt.Sprintf("%v is not integer!", _v))
			}
		}

		if len(scc[i]) > 1 {
			//println(len(scc[i]))
			nonTrivialComponents++
		}
		for _, _v := range scc[i] {
			if v, ok := _v.(int); ok {
				if seenBlks[v] == BlkUnexecuted {
					seenBlks[v] = BlkCommitted
					blk := getBlock(v)
					blockchain = append(blockchain, blk)
					//printBlock(blk)
				}
			} else {
				panic(fmt.Sprintf("%v is not integer!", _v))
			}
		}
	}
removeExecuted:
	newUnexecuted := make([]Block, 0)
	for ; i < len(scc); i++ {
		for _, _v := range scc[i] {
			if v, ok := _v.(int); ok {
				if seenBlks[v] != BlkUnknown && seenBlks[v] != BlkCommitted {
					newUnexecuted = append(newUnexecuted, getBlock(v))
				}
			} else {
				panic(fmt.Sprintf("%v is not integer!", _v))
			}
		}
	}
	unexecutedBlocks = newUnexecuted
}

func tryExecuteDAG() {
	updated := false
	newReadyToExecute := make([]int, 0)
	for _, v := range readyToExecute {
		if seenBlks[v] == BlkUnexecuted {
			blockchain = append(blockchain, getBlock(v))
			seenBlks[v] = BlkCommitted
			for _, u := range dependants[v] {
				depCnt[u]--
				if depCnt[u] == 0 {
					newReadyToExecute = append(newReadyToExecute, u)
					updated = true
				}
			}
		} else { // seenBlks[v] is never BlkCommitted here
			newReadyToExecute = append(newReadyToExecute, v)
		}
	}

	readyToExecute = newReadyToExecute
	if updated {
		tryExecuteDAG()
	}
}

func commitBlock(block Block) {
	BlockChainSemaphore.Acquire(context.Background(), 1)
	fillDependants(block)
	if depCnt[block.SeqNum] == 0 {
		readyToExecute = append(readyToExecute, block.SeqNum)
	}
	if seenBlks[block.SeqNum] == BlkUnknown {
		unexecutedBlocks = append(unexecutedBlocks, block) // addBlock will get stuck on semaphore
	}
	seenBlks[block.SeqNum] = BlkUnexecuted
	tryExecuteTarjan()
	BlockChainSemaphore.Release(1)
}

func printUnexecBlks() {
	fmt.Print("unexecuted:{ ")
	for i := 0; i < len(unexecutedBlocks); i++ {
		seqNum := unexecutedBlocks[i].SeqNum
		fmt.Printf("(%v:%v) ", seqNum, seenBlks[seqNum])
	}
	fmt.Print("}\n")
}

func printBlock(block Block) {
	fmt.Printf("executed %v with deps: %v\n", block.SeqNum, block.Deps)
}

func getBalance(user string) int {
	balances := calculateBalances()

	return balances[user] + 10
}

func createAndAddBlock(txns []Transaction) Block {
	BlockChainSemaphore.Acquire(context.Background(), 1)
	defer	BlockChainSemaphore.Release(1)

	blk := Block{
		SeqNum: incSeqNum(),
		Deps:   calculateDependencies(txns),
		Tx:     txns,
	}
	unexecutedBlocks = append(unexecutedBlocks, blk)
	seenBlks[blk.SeqNum] = BlkReceived
	//printUnexecBlks()

	return blk
}

func addBlock(block Block) {
	BlockChainSemaphore.Acquire(context.Background(), 1)
	defer	BlockChainSemaphore.Release(1)

	unexecutedBlocks = append(unexecutedBlocks, block)
	seenBlks[block.SeqNum] = BlkReceived
	//printUnexecBlks()
}

func getBlock(seqnum int) Block {
	for _, block := range unexecutedBlocks {
		if block.SeqNum == seqnum {
			return block
		}
	}

	panic(fmt.Sprintf("block %v not found getBlock", seqnum))
}

func updateBlock(block Block) {
	BlockChainSemaphore.Acquire(context.Background(), 1)
	defer	BlockChainSemaphore.Release(1)

	for i := 0; i < len(unexecutedBlocks); i++ {
		if unexecutedBlocks[i].SeqNum == block.SeqNum {
			unexecutedBlocks[i] = block
			return
		}
	}

	panic(fmt.Sprintf("block %v not found updateBlock", block.SeqNum))
}
