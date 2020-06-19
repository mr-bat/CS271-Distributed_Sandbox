package main

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"time"
)

const (
	PREPARE = "PREPARE"
	ACK = "ACK"
	ACCEPT = "ACCEPT"
	ACCEPTED = "ACCEPTED"
	COMMIT = "COMMIT"
	BENCHMARK = "BENCHMARK"
	BENCHMARK_CNT = 300
)

type Ballot struct {
	Sender int
	Quorum       int
}

var partialBlk map[int]Block
var needsSlowPath map[int]bool
var ackCnt map[int]int

func init() {
	partialBlk = make(map[int]Block, 0)
	needsSlowPath = make(map[int]bool, 0)
	ackCnt = make(map[int]int, 0)
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
	blk.Deps = unique(append(blk.Deps, calculateDependencies(blk.Tx)...))

	addBlock(blk)
	sendClient(ballot.Sender, getAckMsg(blk, ballot))
}

func handleAck(message []string) {
	blk := parseBlock(message[1])
	//ballot := parseBallot(message[2])
	Epaxos.Acquire(context.Background(), 1)
	ackCnt[blk.SeqNum]++
	_ackCnt := ackCnt[blk.SeqNum]

	if _, ok := partialBlk[blk.SeqNum]; ok {
		initialLen := len(partialBlk[blk.SeqNum].Deps)
		blk.Deps = unique(append(blk.Deps, partialBlk[blk.SeqNum].Deps...))
		if initialLen != len(blk.Deps) {
			needsSlowPath[blk.SeqNum] = true
		}
	}
	partialBlk[blk.SeqNum] = blk
	Epaxos.Release(1)

	if _ackCnt == getQuorumSize() {
		if needsSlowPath[blk.SeqNum] {
			ackCnt[-blk.SeqNum] = 1
			panic("not possible for this test")
			//updateBlock(blk) // Bad practice
			//sendToQuorum(ballot.Quorum, getAcceptMsg(partialBlk[blk.SeqNum], ballot))
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

	updateBlock(blk)
	sendClient(ballot.Sender, getAcceptedMsg(blk))
}

func handleAccepted(message []string) {
	blk := parseBlock(message[1])
	ackCnt[-blk.SeqNum]++

	if ackCnt[-blk.SeqNum] == getQuorumSize() {
		sendToClients(getCommitMsg(blk))
	}
}

var commitCnt = 0
func handleCommit(message []string) {
	blk := parseBlock(message[1])

	commitBlock(blk)
	commitCnt++
	//fmt.Printf("commited %v blocks with %v unexecuted blocks\n", commitCnt, len(unexecutedBlocks))

	if commitCnt >= BENCHMARK_CNT * (GetNumberOfClients() + 1) - 3 { // 3 is error window
		fmt.Printf("Handled %v commit messages in %v with %v unexecuted blocks\n", commitCnt, time.Now().UnixNano(), len(unexecutedBlocks))
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