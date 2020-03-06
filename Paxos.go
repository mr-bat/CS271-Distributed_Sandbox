package main

import "time"

const timeout int = 2

var ackCount int
var acceptedCount int

var noReturn bool = false

var receivedBlocks []Block

type BallotNum struct {
	num int
	id  int
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
	return "PREPARE@" + string(ballot.num) + "@" + string(ballot.id)
}

func getAcceptMessage(ballot BallotNum) string {
	return "ACCEPT@" + string(ballot.num) + "@" + string(ballot.id) + "@" + rangeToString(getCurrBlockChain())
}

func getAcceptedMessage(sequenceNum int) string {
	return "ACCEPTED@" + string(sequenceNum) + "@" + rangeToString(getCurrBlockChain())
}
func getCommitMessage(sequenceNum int, blocks []Block) string {
	return "COMMIT@" + string(sequenceNum) + "@" + rangeToString(blocks)
}
