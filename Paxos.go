package main

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"reflect"
	"time"
)

const timeout int = 3

var ackCount int
var acceptedCount int
var lowestAck int


var receivedTransactions []Transaction

type Ballot struct {
	Num       int
	ProcessId int
}

type Message struct {
	Ballot   Ballot
	Accepted bool
	Block    Block
}

var latestBallotNumber int
var lastBallot Ballot

func init() {
	latestBallotNumber = 0
}

func reset() {
	latestBallotNumber = 0
}

func isGreaterBallot(bn Ballot) bool {
	if bn.Num > lastBallot.Num {
		return true
	} else if bn.Num == lastBallot.Num && bn.ProcessId > lastBallot.ProcessId {
		return true
	}
	return false
}

func (ballot Ballot) toString() string {
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

func (msg Message) toString() string {
	res, _ := json.Marshal(msg)
	return string(res)
}

func parseMessage(tx string) Message {
	var res Message
	if err := json.Unmarshal([]byte(tx), &res); err != nil {
		panic(err)
	}
	return res
}

func getQuorumSize() int {
	return (len(clients) / 2) + 1
}

func beginSync() {
	var lastCommitedBlock Block
	var commitingAcceptedBlock bool

	beginProtocol:

	ackCount = 0
	acceptedCount = 0
	receivedTransactions = nil
	lastCommitedBlock = getLastBlock()
	commitingAcceptedBlock = false
	lowestAck = getCurrSeqNumber()
	latestBallotNumber++
	var myBallot = Ballot{latestBallotNumber, getId()}
	lastBallot = myBallot
	Logger.WithFields(logrus.Fields{
		"lowest ack": lowestAck,
	}).Info("begin protocol")

	prepareMessage := getPrepareMessage(myBallot)
	sendToClients(prepareMessage)
	time.Sleep(time.Duration(timeout) * time.Second)
	if (ackCount + 1) < getQuorumSize() || !reflect.DeepEqual(lastCommitedBlock, getLastBlock()) {
		goto beginProtocol
	}
	if lowestAck < getCurrSeqNumber() {
		for lowestAck < getCurrSeqNumber() {
			sendToClients(getCommitMessage(getBlock(lowestAck)))
			lowestAck++
		}
		goto beginProtocol
	}

	if acceptedBlock.isEmpty() {
		newBlock := createNewBlock()
		newBlock.Tx = append(newBlock.Tx, receivedTransactions...)
		receivedTransactions = nil
		acceptedBlock = newBlock
		sendToClients(getAcceptMessage(myBallot, acceptedBlock))
	} else {
		commitingAcceptedBlock = true
		sendToClients(getAcceptMessage(myBallot, acceptedBlock))
	}
	time.Sleep(time.Duration(timeout) * time.Second)
	if (acceptedCount + 1) < getQuorumSize() {
		goto beginProtocol
	}

	commitBlock(acceptedBlock)
	sendToClients(getCommitMessage(acceptedBlock))
	acceptedBlock = Block{}
	if commitingAcceptedBlock {
		goto beginProtocol
	}
}

func getPrepareMessage(ballot Ballot) string {
	return "PREPARE@" + Message{ballot, false, getLastBlock()}.toString()
}

func getAckMessage(ballot Ballot) string {
	var currBlk Block
	if acceptedBlock.isEmpty() {
		currBlk = createNewBlock()
	} else {
		currBlk = acceptedBlock
	}
	fmt.Printf("currBlock %v\n", currBlk)
	fmt.Println("ACK@" + Message{
		Ballot:   ballot,
		Accepted: !acceptedBlock.isEmpty(),
		Block:    currBlk,
	}.toString())

	return "ACK@" + Message{
		Ballot:   ballot,
		Accepted: !acceptedBlock.isEmpty(),
		Block:    currBlk,
	}.toString()
}

func getAcceptMessage(ballot Ballot, block Block) string {
	return "ACCEPT@" + Message{Ballot: ballot, Accepted: false, Block: block}.toString()
}

func getAcceptedMessage(ballot Ballot) string {
	return "ACCEPTED@" + Message{Ballot: ballot, Accepted: false, Block: Block{}}.toString()
}

func getCommitMessage(block Block) string {
	return "COMMIT@" + Message{Ballot{latestBallotNumber, 0}, true, block}.toString()
}
