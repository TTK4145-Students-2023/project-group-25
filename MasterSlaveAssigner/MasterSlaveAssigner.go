package MasterSlaveAssigner

import (
	"Driver-go/peers"
	"fmt"
	"net"
	"strconv"
	"strings"
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

var (
	IPAddrP2PRx = make(chan peers.PeerUpdate) // input channel to recieve IP adresses from P2P NTW
	MasterSlave = make(chan bool)             // output channel to send Master or Slave role to order assigner
)

func Max(array []int) int {
	var max int = array[0]
	for _, value := range array {
		if max < value {
			max = value
		}
	}
	return max
}

// this function only checks which last byte is the biggest (can be developed)
func CompIP(localIP string, P2P_IP []string) int {

	if len(P2P_IP) == 0 { // if array is empty i = 1 --> only node on NTW == master
		return 1
	}

	localIPArr := strings.Split(localIP, ".")
	ipLocalLastByteInt, _ := strconv.Atoi(localIPArr[len(localIPArr)-1])

	var ipP2PintArray []int
	for i := 0; i < len(P2P_IP); i++ {
		ipLastByteString := strings.Split(P2P_IP[i], ".")
		ipLastByteInt, _ := strconv.Atoi(ipLastByteString[len(ipLastByteString)-1])
		ipP2PintArray = append(ipP2PintArray, ipLastByteInt)
	}
	for i := 0; i < len(ipP2PintArray); i++ {
		fmt.Println(ipP2PintArray[i])
	}

	maxIP := Max(ipP2PintArray)

	i := 0
	if maxIP <= ipLocalLastByteInt {
		i = 1
	} else {
		i = 0
	}
	return i
}

// assign master or slave to order assigner
func MasterSlaveAssigner() bool {
	localIP, _ := LocalIP()
	for {
		select {
		case P2P_IP := <-IPAddrP2PRx:
			i := CompIP(localIP, P2P_IP.Peers)
			if i == 1 {
				MasterSlave <- true
				fmt.Println("Master")
			} else {
				MasterSlave <- false
				fmt.Println("Slave")
			}
		}
	}
}
