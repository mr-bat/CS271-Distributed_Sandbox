package main

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"time"
)

const (
	PREPARE       = "PREPARE"
	ACK           = "ACK"
	ACCEPT        = "ACCEPT"
	ACCEPTED      = "ACCEPTED"
	COMMIT        = "COMMIT"
	BENCHMARK     = "BENCHMARK"
	BenchmarkCnt  = 300
	ConflictRatio = 100
)

type Ballot struct {
	Sender int
	Quorum       int
}

var partialBlk map[int]Block
var needsSlowPath map[int]bool
var ackCnt map[int]int
var acptCnt map[int]int
var conflictArr []int

func init() {
	partialBlk = make(map[int]Block, 0)
	needsSlowPath = make(map[int]bool, 0)
	ackCnt = make(map[int]int, 0)
	acptCnt = make(map[int]int, 0)
	conflictArr = make([]int, 0, BenchmarkCnt)
	for i := 0; i < BenchmarkCnt; i++ {
		conflictArr = append(conflictArr, rand.Int()%100)
	}
}

func (ballot *Ballot) toString() string {
	res, _ := json.Marshal(ballot)
	return string(res)
}

func parseBallot(ballot string) Ballot {
	var res Ballot
	if err := json.Unmarshal([]byte(ballot), &res); err != nil {
		panic(err)
	}
	return res
}

func getQuorumSize() int {
	return (len(clients) / 2) + 1
}

func chooseQuorum() int {
	rand.Seed(time.Now().UnixNano())
	p := rand.Perm(GetNumberOfClients())
	bitMask := 0
	for _, r := range p[:(getQuorumSize()-1)] {
		bitMask += 1 << uint(r)
	}

	return bitMask
}

func handleNewTransaction(txn Transaction) {
	ballot := Ballot{
		Sender: getId(),
		Quorum: chooseQuorum(),
	}
	blk := createAndAddBlock([]Transaction{txn})
	Epaxos.Acquire(context.Background(), 1)
	ackCnt[blk.SeqNum] = 1
	Epaxos.Release(1)
	sendToQuorum(ballot.Quorum, getPrepareMsg(blk, ballot))
}

func handlePrepare(message []string) {
	blk := parseBlock(message[1])
	ballot := parseBallot(message[2])
	blk = calculateDependencies(blk)

	addBlock(blk)
	sendClient(ballot.Sender, getAckMsg(blk, ballot))
}

func divideDeps(block Block) Block {
	cerDeps := make([]int, 0)
	potDeps := make([]int, 0)
	for _, dep := range block.CerDeps {
		if dep < block.SeqNum {
			cerDeps = append(cerDeps, dep)
		} else {
			potDeps = append(potDeps, dep)
		}
	}
	block.CerDeps = cerDeps
	block.PotDeps = potDeps

	return block
}

func handleAck(message []string) {
	blk := parseBlock(message[1])
	ballot := parseBallot(message[2])
	Epaxos.Acquire(context.Background(), 1)
	ackCnt[blk.SeqNum]++
	_ackCnt := ackCnt[blk.SeqNum]

	if _, ok := partialBlk[blk.SeqNum]; ok {
		blk.CerDeps = unique(append(blk.CerDeps, partialBlk[blk.SeqNum].CerDeps...))
		if Debugging {
			fmt.Print("Updating partial block\n")
		}
	}
	partialBlk[blk.SeqNum] = blk
	Epaxos.Release(1)

	if _ackCnt == getQuorumSize() {
		blk = divideDeps(blk)

		if len(blk.PotDeps) != 0 {
			msg := getAcceptMsg(blk, ballot)
			sendToQuorum(ballot.Quorum, getAcceptMsg(blk, ballot))

			sz := len(msg)
			if msg[sz - 1] == '#' {
				msg = msg[:sz-1]
			}
			handleReceivedMessage(msg)
		} else {
			msg := getCommitMsg(blk)
			sendToClients(msg)

			sz := len(msg)
			if msg[sz - 1] == '#' {
				msg = msg[:sz-1]
			}
			handleReceivedMessage(msg)
		}
	}
}

func handleAccept(message []string) {
	blk := parseBlock(message[1])
	ballot := parseBallot(message[2])
	Epaxos.Acquire(context.Background(), 1)
	for i := len(unexecutedBlocks) - 1; i > -1; i-- {
		if unexecutedBlocks[i].SeqNum == blk.SeqNum {
			break
		}
		if contains(unexecutedBlocks[i].CerDeps, blk.SeqNum) {
			blk.ToBeDel = append(blk.ToBeDel, unexecutedBlocks[i].SeqNum)
		}
	} // Remove committeds?
	Epaxos.Release(1)

	updateBlock(blk)
	if ballot.Sender != getId() {
		sendClient(ballot.Sender, getAcceptedMsg(blk))
	} else {
		msg := getAcceptedMsg(blk)
		sz := len(msg)
		if msg[sz - 1] == '#' {
			msg = msg[:sz-1]
		}
		handleReceivedMessage(msg)
	}
}

func handleAccepted(message []string) {
	blk := parseBlock(message[1])
	Epaxos.Acquire(context.Background(), 1)
	acptCnt[blk.SeqNum]++
	blk.ToBeDel = unique(append(blk.ToBeDel, partialBlk[blk.SeqNum].ToBeDel...))
	partialBlk[blk.SeqNum] = blk
	_acptCnt := acptCnt[blk.SeqNum]
	Epaxos.Release(1)

	if _acptCnt == getQuorumSize() {
		for _, dep := range blk.PotDeps {
			if !contains(blk.ToBeDel, dep) {
				blk.CerDeps = append(blk.CerDeps, dep)
			}
		}

		msg := getCommitMsg(blk)
		sendToClients(msg)
		sz := len(msg)
		if msg[sz - 1] == '#' {
			msg = msg[:sz-1]
		}
		handleReceivedMessage(msg)
	}
}

var commitCnt = 0
func handleCommit(message []string) {
	blk := parseBlock(message[1])

	commitBlock(blk)
	commitCnt++
	if Debugging {
		fmt.Printf("commited %v blocks with %v unexecuted blocks\n", commitCnt, len(unexecutedBlocks))
	}

	if commitCnt >= BenchmarkCnt* (GetNumberOfClients() + 1) {
		fmt.Printf("Handled %v commit messages in %vms with %v unexecuted blocks with %v non-trivial components\n",
			commitCnt,
			(time.Now().UnixNano() - benchmarkBeganAt) / 1e6,
			len(unexecutedBlocks),
			nonTrivialComponents,
		)
		if commitCnt > BenchmarkCnt* (GetNumberOfClients() + 1) {
			panic("WTH?!")
		}
	}
}

func getPrepareMsg(block Block, ballot Ballot) string {
	return PREPARE + "@" + block.toString() + "@" + ballot.toString() + "#"
}

func getAckMsg(block Block, ballot Ballot) string {
	return ACK + "@" + block.toString() + "@" + ballot.toString() + "#"
}

func getAcceptMsg(block Block, ballot Ballot) string {
	return ACCEPT + "@" + block.toString() + "@" + ballot.toString() + "#"
}

func getAcceptedMsg(block Block) string {
	return ACCEPTED + "@" + block.toString() + "#"
}

func getCommitMsg(block Block) string {
	return COMMIT + "@" + block.toString() + "#"
}

func beginBenchmark() {
	sendToClients(BENCHMARK)
	handleReceivedMessage(BENCHMARK)
}