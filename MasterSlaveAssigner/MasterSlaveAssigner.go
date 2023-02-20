package MasterSlaveAssigner
package main

import (
    "log"
    "net"
    "strings"
	"fmt" 
)

// Get preferred outbound ip of this machine
func GetOutboundIP() net.IP {
    conn, err := net.Dial("udp", "8.8.8.8:80")
    if err != nil {
        log.Fatal(err)
    }
    defer conn.Close()

    localAddr := conn.LocalAddr().(*net.UDPAddr)

    return localAddr.IP
}

func MasterSlaveAssigner (){
	selfIPaddr = GetOutboundIP
	print.fmt(selfIPaddr)
}

func main(){
	go MasterSlaveAssigner
}

