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

type Addr struct {
	IP   string
	Port string
}

func (this *Addr) String() string{
	return this.IP + ":" + this.Port
}

const (
	//ServerAddr = "http://178.128.139.251:8123/"
	ServerAddr = "https://locshare-167709.appspot.com/"
)

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
			if ip != nil && ip.To4() != nil && !ip.IsLoopback() {
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
	Logger.WithFields(logrus.Fields{
		"ip": getLocalIP(),
		"port": port,
	}).Info("posted addr to remote-server")
}

func removeServerAddr(addr Addr) {
	ip := addr.IP
	port := addr.Port

	var jsonStr = []byte(fmt.Sprintf(`{ "ip": "%v", "port": "%v" }`, ip, port))
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodPut, ServerAddr, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		panic(err)
	}
	res, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	Logger.WithFields(logrus.Fields{
		"ip": ip,
		"port": port,
	}).Info("removed addr from remote-server")
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
		_ = json.Unmarshal([]byte(string(contents)), &addrs)

		Logger.WithField("addrs", addrs).Info("received addrs from remote-server")
		//fmt.Printf("received addrs: %v", addrs)
		return addrs
	}
	return nil
}