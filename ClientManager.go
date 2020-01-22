package main

import (
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
				Logger.WithField("client", connection.socket.LocalAddr()).Info("terminated connection")
				close(connection.data)
				delete(manager.clients, connection)
			}
		case message := <-manager.mainChannel:
			Logger.WithFields(logrus.Fields{
				"data" : message,
				"client" : getAddress(),
			}).Info("main channel received")
			//for connection := range manager.clients {
			//	select {
			//	case connection.data <- message:
			//	default:
			//		close(connection.data)
			//		delete(manager.clients, connection)
			//	}
			//}
		}
	}
}

func (manager *ClientManager) Receive(client *Client) {
	for {
		message := make([]byte, 4096)
		length, err := client.socket.Read(message)
		if err != nil {
			manager.unregister <- client
			client.socket.Close()
			break
		}
		if length > 0 {
			Logger.WithFields(logrus.Fields{
				"data" : message,
				"client" : getAddress(),
			}).Info("received data")
			manager.mainChannel <- message
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