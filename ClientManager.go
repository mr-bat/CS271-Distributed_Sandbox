package main

import (
	"bytes"
	"github.com/sirupsen/logrus"
)

type ClientManager struct {
	clients     map[*Client]bool
	mainChannel chan []byte
	register    chan *Client
	unregister  chan *Client
}

func (manager *ClientManager) Start() {
	for {
		select {
		case connection := <-manager.register:
			manager.clients[connection] = true
			Logger.WithField("client", connection.socket.LocalAddr()).Info("added new connection")
		case connection := <-manager.unregister:
			if _, ok := manager.clients[connection]; ok {
				Logger.WithField("client", connection.socket.RemoteAddr()).Info("terminated connection")
				close(connection.data)
				delete(manager.clients, connection)
			}
		}
	}
}

func (manager *ClientManager) Receive(client *Client) {
	for {
		message := make([]byte, 4096)
		length, err := client.socket.Read(message)
		if err != nil {
			Logger.WithFields(logrus.Fields{
				"data" : err,
				"client" : getAddress(),
			}).Error("received data: error")

			//manager.unregister <- client
			//client.socket.Close()
			break
		}
		if length > 0 {
			message = bytes.Trim(message, "\x00")
			Logger.WithFields(logrus.Fields{
				"data" : string(message),
				//"client" : getAddress(),
			}).Info("received data")
			handleReceivedMessage(string(message))
			//manager.mainChannel <- message
		}
	}
}

func (manager *ClientManager) Send(client *Client) {
	defer client.socket.Close()
	for {
		select {
		case message, ok := <-client.data:
			if !ok {
				return
			}
			client.socket.Write(message)
			Logger.WithFields(logrus.Fields{
				"data" : message,
				"client" : getAddress(),
			}).Info("sent data")
		}
	}
}