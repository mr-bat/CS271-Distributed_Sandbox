package main

import (
	"context"
	"fmt"
	"github.com/adam-lavrik/go-imath/ix"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

var clients []*Client

func GetNumberOfClients() int {
	return len(clients)
}

func connectToClients(addrs []Addr) {
	var _clients []*Client

	for _, address := range addrs {
		if fmt.Sprintf("%v:%v", address.IP, address.Port) == getAddress() { // connecting to itself
			continue
		}
		Logger.WithField("address", address).Info("trying to connect to client")
		client := startClientMode(address)

		if client != nil {
			_clients = append(_clients, client)
			fmt.Printf("server %v at address %v\n", len(_clients), address)
		} else {
			removeServerAddr(address)
		}
	}

	clients = _clients
}

func sendToClients(message string) {
	logMessage(message, true)
	for _, client := range clients {
		Logger.WithFields(logrus.Fields{
			"clientId": client.id,
		}).Info("sending to client")

		client.Send(message)
	}
}

//nolint
func sendClient(id int, message string) {
	logMessage(message, true)
	for _, client := range clients {
		if client.id == id {
			Logger.WithFields(logrus.Fields{
				"clientId": id,
			}).Info("sending to client")

			client.Send(message)
		}
	}
}

func startClientMode(addr Addr) *Client {
	connection, error := net.DialTimeout("tcp", fmt.Sprintf("%v:%v", addr.IP, addr.Port), 2*time.Second)
	if error != nil {
		//Logger.Error(error)
		return nil
	}

	//Logger.Info("starting client...")
	Logger.WithFields(logrus.Fields{
		"server-address": fmt.Sprintf("%v:%v", addr.IP, addr.Port),
		"local-address":  getAddress(),
	}).Info("connecting to server")

	client := &Client{socket: connection}
	go client.Receive()

	return client
}

func logMessage(msg string, sending bool) {
	infoText := "handling message"
	shortInfo := "recv"
	if sending {
		infoText = "sending message"
		shortInfo = "send"
	}
	parsed := strings.Split(msg, "@")
	command := parsed[0]
	if command != "ID" {
		Logger.WithFields(logrus.Fields{
			"command": shortInfo + ":" + parsed[0],
			"message": parseMessage(parsed[1]),
			"lowest ack": lowestAck,
			"last ballot":  lastBallot,
		}).Info(infoText)
	}
}

func handleReceivedMessage(message string) {
	if !Connected {
		return
	}

	logMessage(message, false)
	parsed := strings.Split(message, "@")
	command := parsed[0]

	if command == "ID" {
		fmt.Println(parsed)
		id, _ := strconv.Atoi(parsed[1])
		addClientId(id, parsed[2])
	} else if command == "PREPARE" {
		prepareMessage := parseMessage(parsed[1])
		receivedBallot := prepareMessage.Ballot
		block := getBlock(prepareMessage.Block.SeqNum + 1)
		if !block.isEmpty() {
			commitMsg := getCommitMessage(block)
			sendClient(receivedBallot.ProcessId, commitMsg)
		} else if isGreaterBallot(receivedBallot) {
			lastBallot = receivedBallot
			sendClient(receivedBallot.ProcessId, getAckMessage(receivedBallot))
		}
		latestBallotNumber = ix.Max(latestBallotNumber, prepareMessage.Ballot.Num)
	} else if command == "ACK" {
		ackMessage := parseMessage(parsed[1])
		if ackMessage.Ballot == lastBallot {
			if ackMessage.Accepted {
				acceptedBlock = ackMessage.Block
				lowestAck = ix.Min(lowestAck, ackMessage.Block.SeqNum - 1)
			} else {
				receivedTransactions = append(receivedTransactions, ackMessage.Block.Tx...)
				lowestAck = ix.Min(lowestAck, ackMessage.Block.SeqNum)
			}
			ackCount++
		}
		latestBallotNumber = ix.Max(latestBallotNumber, ackMessage.Ballot.Num)
	} else if command == "ACCEPT" {
		acceptMessage := parseMessage(parsed[1])
		if acceptMessage.Ballot == lastBallot {
			acceptedBlock = acceptMessage.Block
			sendClient(acceptMessage.Ballot.ProcessId, getAcceptedMessage(acceptMessage.Ballot))
		}
		latestBallotNumber = ix.Max(latestBallotNumber, acceptMessage.Ballot.Num)
	} else if command == "ACCEPTED" {
		acceptedMessage := parseMessage(parsed[1])
		if acceptedMessage.Ballot == lastBallot {
			acceptedCount++
		}
		latestBallotNumber = ix.Max(latestBallotNumber, acceptedMessage.Ballot.Num)
	} else if command == "COMMIT" {
		commitMessage := parseMessage(parsed[1])
		fmt.Printf("commiting: newBlk? %v acceptedBlk %v\n", getBlock(commitMessage.Block.SeqNum).isEmpty(), acceptedBlock)
		if getBlock(commitMessage.Block.SeqNum).isEmpty() {
			if commitMessage.Block.SeqNum >= acceptedBlock.SeqNum {
				fmt.Println("accepted Blk reset")
				acceptedBlock = Block{}
			}
			if commitMessage.Block.SeqNum == 1 || !getBlock(commitMessage.Block.SeqNum - 1).isEmpty() {
				commitBlock(commitMessage.Block)
			}
		}
		latestBallotNumber = ix.Max(latestBallotNumber, commitMessage.Ballot.Num)
	}
}

func addClientId(id int, address string) {
	Logger.WithFields(logrus.Fields{
		"id":             id,
		"client-address": address,
	}).Info("identifying client")

	for _, _client := range clients {
		if _client.socket.RemoteAddr().String() == address {
			_client.id = id
			Logger.WithFields(logrus.Fields{
				"id":     id,
				"client": address,
			}).Info("identified client")
		}
	}
}

func addPurchase(from, to string, amount int) {
	Logger.WithFields(logrus.Fields{
		"from":                   from,
		"from's-initial-balance": getBalance(from),
		"to":                     to,
		"amount":                 amount,
	}).Info("current transaction")

	Logger.WithFields(logrus.Fields{
		"pending":                 pendingTx,
	}).Info("before sync")
	if getBalance(from) < amount {
		beginSync()
	}
	Logger.WithFields(logrus.Fields{
		"pending":                 pendingTx,
	}).Info("after sync")

	if getBalance(from) < amount {
		fmt.Println("INCORRECT")
	} else {
		BlockChainSemaphore.Acquire(context.Background(), 1)
		addTransaction(Transaction{
			Sender:   from,
			Receiver: to,
			Amount:   amount,
			Id:       incClock(),
		})
		fmt.Println("SUCCESS")
		BlockChainSemaphore.Release(1)
	}

	Logger.WithFields(logrus.Fields{
		"pending":                 pendingTx,
	}).Info("after add")
}

func advertiseId() {
	id := getIdFromInput()
	setId(id)
	Logger.WithField("id", getId()).Info("set id")
	sendToClients(fmt.Sprintf("ID@%d@%s", getId(), getAddress()))
}
