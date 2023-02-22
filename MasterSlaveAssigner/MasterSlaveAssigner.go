package MasterSlaveAssigner

import (
	"net"
	"strconv"
	"strings"
	// "Network-go/network/peers"
)

// type of message on P2P network
type msgP2P struct {
	//elevData - some format
	IPAddrP2P string `json:"IPAddr"`
}

// Get preferred outbound ip of this machine
// work as long as you have internet connection
var localIP string

func LocalIP() (string, error) {
	if localIP == "" {
		conn, err := net.DialTCP("tcp4", nil, &net.TCPAddr{IP: []byte{8, 8, 8, 8}, Port: 53})
		if err != nil {
			return "", err
		}
		defer conn.Close()
		localIP = strings.Split(conn.LocalAddr().String(), ":")[0]
	}
	return localIP, nil
}

// returning -1 if ip1 > ip2, 0 if ip1 == ip2 and 1 if ip1 < ip2
// this function only checks for which last byte is the biggest (can be developed)
func CompIP(localIP, P2P_IP string) int {
	localIPArr := strings.Split(localIP, ".")
	P2P_IPArr := strings.Split(P2P_IP, ".")
	ip1LastByte, _ := strconv.Atoi(localIPArr[len(localIPArr)-1])
	ip2LastByte, _ := strconv.Atoi(P2P_IPArr[len(P2P_IPArr)-1])
	i := 0
	if ip1LastByte > ip2LastByte {
		i = -1
	} else if ip1LastByte == ip2LastByte {
		i = 0
	} else if ip1LastByte < ip2LastByte {
		i = 1
	}
	return i
}

var (
	IPAddrP2PRx = make(chan string) // input channel to recieve IP adress from P2P NTW
	MasterSlave = make(chan bool)   // output channel to send Master or Slave role to order assigner
)

// assign master or slave to prder assigner
// dont know if this works yet...
func MasterSlaveAssigner() bool {
	localIP, _ := LocalIP() // dont know where to assign this...
	for {
		select {
		case P2P_IP := <-IPAddrP2PRx:
			i := CompIP(localIP, P2P_IP)
			if i == -1 {
				MasterSlave <- true
			} else {
				MasterSlave <- false
			}
		}
	}
}
