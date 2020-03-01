package main

import "time"

const timeout int = 2

var ackCount int

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
}

func getPrepareMessage(ballot BallotNum) string {
	return "PREPARE@" + string(ballot.num) + "@" + string(ballot.id)
}

func getAcceptMessage(ballot BallotNum) string {
	return "ACCEPT@" + string(ballot.num) + "@" + rangeToString(getCurrBlockChain())
}
