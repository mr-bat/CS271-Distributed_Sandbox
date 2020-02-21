package main

import (
	"fmt"
	"net"
	"strconv"
	"strings"

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

func sendClient(id int, message string) {
	for _, client := range clients {
		if client.id == id {
			println("sending: " + message)
			client.Send(message)
		}
	}
}

func startClientMode(addr Addr) *Client {
	connection, error := net.Dial("tcp", fmt.Sprintf("%v:%v", addr.IP, addr.Port))
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
	} else if command == "DATA" {
		timetable := parseTt(parsed[1])
		blocks := parseRange(parsed[2])
		id, _ := strconv.Atoi(parsed[3])
		updateSelf(timetable, blocks, id)
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

	if initialBalance >= amount {
		addBlock(from, to, amount)
		fmt.Println("SUCCESS")
	} else {
		fmt.Println("INCORRECT")
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

func informClient(id int) {
	Logger.WithFields(logrus.Fields{
		"id":        id,
		"timetable": timetable,
	}).Info("informing client")

	toBeSent := pickNewBlocks(blockchain, id)
	Logger.WithFields(logrus.Fields{
		"toBeSent": toBeSent,
	}).Info("picked blocks")
	incTime()

	sendClient(id, fmt.Sprintf("DATA@%s@%s@%d", convertTtToString(), rangeToString(toBeSent), getId()))
	fmt.Printf("MESSAGE SENT TO %d\n", id)

	//updateTimetable(timetable, getId(), id)
}
