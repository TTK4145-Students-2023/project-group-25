package main

import (
	"fmt"
	"log"
	"net"
	"encoding/binary"
	"math/big"
	"bytes"
)

// Get preferred outbound ip of this machine
func getLocalIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP
}

// input channel 
var (
	IPAddr_P2P = make(chan string)
)

// output channel (true = master, false = slave)
var(
	MasterSlave = make(chan bool)
)

// assign IP address from P2P network to a string 
type IPAddr_NTW string {
	IPAddr_NTW  	string `json:"IPAddr"`
}

func StringtoInt(string) Int64 {
	IPv4Int := big.NewInt(0)
	IPv4Int.SetBytes(string.To4())
	return IPv4Int.Int64()
}

func MasterSlaveAssigner(IPAddr_P2P <- chan IPAddr_NTW) {
	var localIPaddr string = getLocalIP().String()
	var localIPaddrBin Int64 = StringtoInt(localAddr) 

	fmt.Printf(localIPaddrBin)
}
