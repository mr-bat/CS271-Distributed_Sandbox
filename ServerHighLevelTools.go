package main

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"net"
)

var PortNumber int

func getAddress() string {
	return fmt.Sprintf("%v:%v", getLocalIP(), PortNumber)
}

func startServer(port int) {
	advertiseServerAddr(PortNumber)
	go startServerMode(PortNumber)
	waitForDone()
}

func startServerMode(port int) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%v", port))
	if err != nil {
		Logger.Error(err)
		return
	}

	Logger.WithFields(logrus.Fields{
		"ip": getLocalIP(),
		"port": port,
	}).Info("starting server")
	manager := ClientManager{
		clients:     make(map[*Client]bool),
		mainChannel: make(chan []byte),
		register:    make(chan *Client),
		unregister:  make(chan *Client),
	}
	go manager.Start()
	for {
		connection, err := listener.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}
		client := &Client{socket: connection, data: make(chan []byte)}
		manager.register <- client
		go manager.Receive(client)
		go manager.Send(client)
	}
}