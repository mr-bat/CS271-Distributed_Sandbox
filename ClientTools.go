package main

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"net"
)

var clients []*Client

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
			fmt.Printf("server %v at address %v", len(_clients), address)
		} else {
			removeServerAddr(address)
		}
	}

	clients =  _clients
}

func sendToClients(message string) {
	for _, client := range clients {
		client.Send(message)
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
		"local-address": getAddress(),
	}).Info("connecting to server")

	client := &Client{socket: connection}
	go client.Receive()

	return client
}