package main

import (
	"github.com/sirupsen/logrus"
	"net"
	"strings"
	"time"
)

type Client struct {
	socket net.Conn
	data   chan []byte
}

func (client *Client) Receive() {
	for {
		message := make([]byte, 4096)
		length, err := client.socket.Read(message)
		if err != nil {
			client.socket.Close()
			break
		}
		if length > 0 {
			Logger.WithFields(logrus.Fields{
				"data" : message,
				"client" : client.socket.RemoteAddr(),
			}).Info("received data")
		}
	}
}

func (client *Client) Send(message string) {
	time.Sleep(5 * time.Second)
	client.socket.Write([]byte(strings.TrimRight(message, "\n")))
}