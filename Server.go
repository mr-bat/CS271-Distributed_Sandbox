package main

import (
	"bufio"
	"fmt"
	"github.com/sirupsen/logrus"
	"math/rand"
	"net"
	"os"
	"strings"
)

type ClientManager struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
}

type Client struct {
	socket net.Conn
	data   chan []byte
}

type Addr struct {
	IP string
	Port string
}

func (manager *ClientManager) start() {
	for {
		select {
		case connection := <-manager.register:
			manager.clients[connection] = true
			Logger.WithField("client", connection.socket.RemoteAddr()).Info("added new connection")
		case connection := <-manager.unregister:
			if _, ok := manager.clients[connection]; ok {
				Logger.WithField("client", connection.socket.RemoteAddr()).Info("terminated connection")
				close(connection.data)
				delete(manager.clients, connection)
			}
		case message := <-manager.broadcast:
			for connection := range manager.clients {
				select {
				case connection.data <- message:
				default:
					close(connection.data)
					delete(manager.clients, connection)
				}
			}
		}
	}
}

func (manager *ClientManager) receive(client *Client) {
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
			manager.broadcast <- message
		}
	}
}

func (client *Client) receive() {
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

func (manager *ClientManager) send(client *Client) {
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

func getAddress() string {
	return string(getLocalIP()) + string(portNumber)
}

var portNumber int
func startServerMode(port int) {
	portNumber = port
	listener, err := net.Listen("tcp", ":" + string(port))
	if err != nil {
		Logger.Error(err)
		return
	}

	Logger.WithFields(logrus.Fields{
		"ip": getLocalIP(),
		"port": port,
	}).Info("starting server")
	manager := ClientManager{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
	go manager.start()
	advertiseServerAddr(port)
	for {
		connection, err := listener.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}
		client := &Client{socket: connection, data: make(chan []byte)}
		manager.register <- client
		go manager.receive(client)
		go manager.send(client)
	}
}

func startClientMode(addr Addr) {
	fmt.Println("Starting client...")
	connection, error := net.Dial("tcp", fmt.Sprintf("%v:%v", addr.IP, addr.Port))
	if error != nil {
		fmt.Println(error)
	}

	Logger.WithFields(logrus.Fields{
		"server-address": addr,
		"local-address": getAddress(),
	}).Info("connecting to server")

	client := &Client{socket: connection}
	go client.receive()
	for {
		reader := bufio.NewReader(os.Stdin)
		message, _ := reader.ReadString('\n')
		connection.Write([]byte(strings.TrimRight(message, "\n")))
	}
}

func main() {
	var addrs []Addr = getClientAddrs()

	startServerMode(7180 + rand.Intn(100))

	for _, address := range addrs {
		go startClientMode(address)
	}
}