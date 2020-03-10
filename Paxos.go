package main

import (
	"encoding/json"
	"time"
)

const timeout int = 2

var ackCount int
var acceptedCount int

var noReturn bool = false

var receivedBlocks []Block

type BallotNum struct {
	Num int
	Id  int
}

type Message struct {
	Ballot   BallotNum
	Accepted bool
	block    Block
}

var lastBallot BallotNum

func isGreaterBallot(bn BallotNum) bool {
	if bn.num > lastBallot.num {
		return true
	} else if bn.num == lastBallot.num && bn.id > lastBallot.id {
		return true
	}
	return false

}

func getQuorumSize() int {
	return (len(clients) / 2) + 1
}

func beginSync() {
	ackCount = 0
	acceptedCount = 0
	receivedBlocks = nil
	var myBallot BallotNum = BallotNum{lastBallot.num + 1, getId()}
	lastBallot = myBallot
	prepareMessage := getPrepareMessage(myBallot)
	sendToClients(prepareMessage)
	time.Sleep(time.Duration(timeout) * time.Second)
	if myBallot != lastBallot {
		return
	}
	if (ackCount + 1) < getQuorumSize() {
		return
	}
	acceptMessage := getAcceptMessage(myBallot)
	sendToClients(acceptMessage)
	time.Sleep(time.Duration(timeout) * time.Second)
	if myBallot != lastBallot {
		return
	}
	if (acceptedCount + 1) < getQuorumSize() {
		return
	}
	noReturn = true
	newBlock := append(getCurrBlockChain(), receivedBlocks...)
	addBlockchain(newBlock)
	commitMessage := getCommitMessage(myBallot.num, newBlock)
	sendToClients(commitMessage)
	noReturn = false
}

func getPrepareMessage(ballot BallotNum) string {
	msg := Message{ballot, false, nil}
	serMsg, _ = json.Marshal(msg)
	return "PREPARE@" + string(serMsg)
}

func getAcceptMessage(ballot BallotNum) string {
	return "ACCEPT@" + string(ballot.num) + "@" + string(ballot.id) + "@" + rangeToString(getCurrBlockChain())
}

func getAcceptedMessage(sequenceNum int) string {
	return "ACCEPTED@" + string(sequenceNum) + "@" + rangeToString(getCurrBlockChain())
}
func getCommitMessage(block Block) string {
	ballot := BallotNum(block.sequenceNum, getId())
	msg := Message(ballot, true, block)
	serMsg, _ = json.Marshal(msg)
	return "COMMIT@" + string(ser_msg)
}
