package masterSlaveNTW

import (
	"strconv"
	"strings"
)

// var (
// 	IPAddrP2PRx = make(chan peers.PeerUpdate) // input channel to recieve IP adresses from P2P NTW
// 	MasterSlave = make(chan bool)             // output channel to send Master or Slave role to order assigner
// )

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
func MS_Assigner(localIP string, P2P_IP []string) MasterSlave_state {

	localIPArr := strings.Split(localIP, ".")
	ipLocalLastByteInt, _ := strconv.Atoi(localIPArr[len(localIPArr)-1])

	var ipP2PintArray []int
	for i := 0; i < len(P2P_IP); i++ {
		ipLastByteString := strings.Split(P2P_IP[i], ".")
		ipLastByteInt, _ := strconv.Atoi(ipLastByteString[len(ipLastByteString)-1])
		ipP2PintArray = append(ipP2PintArray, ipLastByteInt)
	}

	maxIP := Max(ipP2PintArray)

	MS_role := MS_slave
	if maxIP <= ipLocalLastByteInt {
		MS_role = MS_master
	}
	return MS_role
}

// // assign master or slave to order assigner
// func MasterSlaveAssigner() {
// 	MasterSlave := make(chan bool)
// 	localIP, _ := localip.LocalIP()
// 	for {
// 		select {
// 		case P2P_IP := <-IPAddrP2PRx:
// 			i := CompIP(localIP, P2P_IP.Peers)
// 			if i == 1 {
// 				MasterSlave <- true
// 				fmt.Println("Master")
// 			} else {
// 				MasterSlave <- false
// 				fmt.Println("Slave")
// 			}
// 		}
// 	}
// }
