package MasterSlaveAssigner

import (
	"net"
	"strings"
)

// Get preferred outbound ip of this machine
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

// input channel
var (
	IPAddr_P2P = make(chan string)
)

// output channel (true = master, false = slave)
var (
	MasterSlave = make(chan bool)
)

// assign IP address from P2P network to a string
type IPAddr_NTW struct {
	IPAddr_NTW string `json:"IPAddr"`
}

// func StringtoInt(IPAddr_string string) int {
// 	IPint, err := strconv.Atoi(IPAddr_string)
// 	if err != nil {
// 		panic(err)
// 	}
// 	return IPint
// }

// func MasterSlaveAssigner(IPAddr_P2P <-chan IPAddr_NTW) {
// 	var localIPaddr string = getLocalIP().String()
// 	var localIPaddrInt int64 = StringtoInt(localIPAddr)

// 	fmt.Println(localIPaddrInt)
// }
