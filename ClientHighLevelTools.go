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
	logMessage(message, true)
	for _, client := range clients {
		//Logger.WithFields(logrus.Fields{
		//	"clientId": client.id,
		//}).Info("sending to client")

		client.Send(message)
	}
}

func sendToQuorum(bitMask int, message string) {
	logMessage(message, true)
	for i, client := range clients {
		if bitMask & (1 << uint(i)) > 0 {
			//Logger.WithFields(logrus.Fields{
			//	"clientId": client.id,
			//	"mask": bitMask & (1 << uint(i)),
			//	//"clients": clients,
			//}).Info("sending to client in quorum")

			client.Send(message)
		}
	}
}

//nolint
func sendClient(id int, message string) {
	logMessage(message, true)
	for _, client := range clients {
		if client.id == id {
			//Logger.WithFields(logrus.Fields{
			//	"clientId": id,
			//}).Info("sending to client")

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
	if Debugging {
		sz := len(msg)
		if msg[sz - 1] == '#' {
			msg = msg[:sz-1]
		}
		infoText := "handling message"
		shortInfo := "recv"
		if sending {
			infoText = "sending message"
			shortInfo = "send"
		}
		parsed := strings.Split(msg, "@")
		command := parsed[0]
		if command != "ID" {
			if len(parsed) == 3 {
				Logger.WithFields(logrus.Fields{
					"command": shortInfo + ":" + parsed[0],
					"block":   parseBlock(parsed[1]),
					"ballot":  parseBallot(parsed[2]),
				}).Info(infoText)
			} else if len(parsed) == 2 {
				Logger.WithFields(logrus.Fields{
					"command": shortInfo + ":" + parsed[0],
					"block":   parseBlock(parsed[1]),
				}).Info(infoText)
			} else {
				Logger.WithFields(logrus.Fields{
					"command": shortInfo + ":" + parsed[0],
				}).Info(infoText)
			}
		}
	}
}

var benchmarkBeganAt int64
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
	} else if command == PREPARE {
		handlePrepare(parsed)
	} else if command == ACK {
		handleAck(parsed)
	} else if command == ACCEPT {
		//panic("not possible for this test")
		handleAccept(parsed)
	} else if command == ACCEPTED {
		//panic("not possible for this test")
		handleAccepted(parsed)
	} else if command == COMMIT {
		handleCommit(parsed)
	} else if command == BENCHMARK {
		go func() {
			benchmarkBeganAt = time.Now().UnixNano()
			conflictingSender := 1
			conflictsSent := 0
			for i := 0; i < BenchmarkCnt; i++ {
				if conflictArr[i] < ConflictRatio {
					if conflictsSent == 2 {
						conflictsSent = 0
						conflictingSender++
					} else {
						conflictsSent++
					}
					sender := conflictingSender * (GetNumberOfClients() + 1) + getId()
					addPurchase(strconv.Itoa(sender), strconv.Itoa(sender), 100)
				} else {
					addPurchase("0", "0", 100)
				}
			}
		}()
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
	txn := Transaction{
		Sender:   from,
		Receiver: to,
		Amount:   amount,
		Id: incClock(),
	}
	handleNewTransaction(txn)
}

func advertiseId() {
	id := getIdFromInput()
	setId(id)
	Logger.WithField("id", getId()).Info("set id")
	sendToClients(fmt.Sprintf("ID@%d@%s", getId(), getAddress()))
}
