package main

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"net"
	"strings"
	//"time"
)

type Client struct {
	socket net.Conn
	data   chan []byte
	id int
}

var id int

func setId(_id int) {
	id = _id
}

func getId() int {
	return id
}

func (client *Client) Receive() {
	for {
		message := make([]byte, 8096000)
		length, err := client.socket.Read(message)
		if err != nil {
			_ = client.socket.Close()
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
	//time.Sleep(300 * time.Microsecond) // Simulating delays
	_, err := client.socket.Write([]byte(strings.TrimRight(message, "\n")))
	if err != nil {
		fmt.Print(err)
		panic("send failed")
	}
}