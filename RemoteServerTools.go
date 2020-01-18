package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net"
	"net/http"
	"os"
)

const ServerAddr = "http://178.128.139.251:8123/"

func getLocalIP() net.IP {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil
	}
	for _, i := range ifaces {
		addrs, err := i.Addrs()
		if err != nil {
			return nil
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip != nil && ip.To4() != nil && ip.IsLoopback() == false {
				return ip
			}
		}
	}

	return nil
}

func advertiseServerAddr(port int) {
	var jsonStr = []byte(fmt.Sprintf(`{ "ip": "%v", "port": "%v" }`, getLocalIP(), port))
	res, err := http.Post(ServerAddr, "application/json", bytes.NewBuffer(jsonStr))
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()
	//fmt.Printf("added %v:%v to server\n", getLocalIP(), port)
	Logger.WithFields(logrus.Fields{
		"ip": getLocalIP(),
		"port": port,
	}).Info("posted addr to remote-server")
}

func removeServerAddr(ip string, port int) {
	var jsonStr = []byte(fmt.Sprintf(`{ "ip": "%v", "port": "%v" }`, ip, port))
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodPut, ServerAddr, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		panic(err)
	}
	res, err := client.Do(req)
	defer res.Body.Close()
	if err != nil {
		panic(err)
	}

	Logger.WithFields(logrus.Fields{
		"ip": ip,
		"port": port,
	}).Info("removed addr from remote-server")
	//fmt.Printf("removed %v:%v from server\n", ip, port)
}

func getClientAddrs() []Addr {
	response, err := http.Get(ServerAddr)
	if err != nil {
		fmt.Printf("%s", err)
		os.Exit(1)
	} else {
		defer response.Body.Close()
		contents, err := ioutil.ReadAll(response.Body)
		if err != nil {
			fmt.Printf("%s", err)
			os.Exit(1)
		}

		var addrs []Addr
		json.Unmarshal([]byte(string(contents)), &addrs)

		Logger.WithField("addrs", addrs).Info("received addrs from remote-server")
		//fmt.Printf("received addrs: %v", addrs)
		return addrs
	}
	return nil
}
