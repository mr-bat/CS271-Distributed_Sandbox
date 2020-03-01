package main

import (
	"fmt"
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
	for _, client := range clients {
		println("sending: " + message)
		client.Send(message)
	}
}

//nolint
func sendClient(id int, message string) {
	for _, client := range clients {
		if client.id == id {
			println("sending: " + message)
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

func handleReceivedMessage(message string) {
	parsed := strings.Split(message, "@")
	command := parsed[0]

	if command == "ID" {
		id, _ := strconv.Atoi(parsed[1])
		addClientId(id, parsed[2])
	} else if command == "PREPARE" {
		if noReturn {
			return
		}
		var receivedBallot BallotNum
		receivedBallot.num, _ = strconv.Atoi(parsed[1])
		receivedBallot.id, _ = strconv.Atoi(parsed[2])
		if isGreaterBallot(receivedBallot) {
			lastBallot = receivedBallot
			ackMessage := "ACK@" + string(receivedBallot.num) + "@" + string(receivedBallot.id)
			sendClient(receivedBallot.id, ackMessage)
		}
	} else if command == "ACK" {
		ackCount++
	} else if command == "ACCEPT" {
		var acceptBallot BallotNum
		acceptBallot.num, _ = strconv.Atoi(parsed[1])
		acceptBallot.id, _ = strconv.Atoi(parsed[2])
		if (acceptBallot == lastBallot) || (isGreaterBallot(acceptBallot)) {
			noReturn = true
			lastBallot = acceptBallot
			acceptedMessage := getAcceptedMessage(acceptBallot.num)
			sendToClients(acceptedMessage)
		}
	} else if command == "ACCEPTED" {
		//sequenceNum, _ := strconv.Atoi(parsed[1])
		receivedBlocks = append(receivedBlocks, parseRange(parsed[2])...)
		acceptedCount++

	} else if command == "COMMIT" {
		addBlockchain(parseRange(parsed[2]))
		noReturn = false
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

func updateSelf(_timetable [][]int, blocks []Block, informerId int) {
	newBlocks := pickNewBlocks(blocks, getId())
	Logger.WithFields(logrus.Fields{
		"received-timetable": _timetable,
		"received-blocks":    blocks,
		"blockchain":         blockchain,
		"filtered-blocks":    newBlocks,
		"informer-id":        informerId,
	}).Info("updating self")

	updateTimetable(_timetable, informerId, getId())
	addBlockchain(newBlocks)

	Logger.WithFields(logrus.Fields{
		"updated-timetable":  timetable,
		"updated-blockchain": blockchain,
	}).Info("updated self")
}

func addTransaction(from, to string, amount int) {
	initialBalance := getBalance(from)
	Logger.WithFields(logrus.Fields{
		"from":                   from,
		"from's-initial-balance": getBalance(from),
		"to":                     to,
		"amount":                 amount,
	}).Info("current transaction")

	if initialBalance < amount {
		beginSync()
	}

	if initialBalance < amount {
		fmt.Println("INCORRECT")
	} else {
		addBlock(from, to, amount)
		fmt.Println("SUCCESS")
	}
	Logger.WithFields(logrus.Fields{
		"timetable":          timetable,
		"from's-new-balance": getBalance(from),
	}).Info("updated timetable")
}

func advertiseId() {
	id := getIdFromInput()
	setId(id)
	Logger.WithField("id", getId()).Info("set id")
	sendToClients(fmt.Sprintf("ID@%d@%s", getId(), getAddress()))
}
