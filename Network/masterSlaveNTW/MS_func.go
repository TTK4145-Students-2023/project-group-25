package masterSlaveNTW

import (
	dt "project/commonDataTypes"
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
func MS_Assigner(localIP string, P2P_IP []string) dt.MasterSlaveRole {

	localIPArr := strings.Split(localIP, ".")
	ipLocalLastByteInt, _ := strconv.Atoi(localIPArr[len(localIPArr)-1])

	var ipP2PintArray []int
	for i := 0; i < len(P2P_IP); i++ {
		ipLastByteString := strings.Split(P2P_IP[i], ".")
		ipLastByteInt, _ := strconv.Atoi(ipLastByteString[len(ipLastByteString)-1])
		ipP2PintArray = append(ipP2PintArray, ipLastByteInt)
	}

	maxIP := Max(ipP2PintArray)

	MS_role := dt.MS_Slave
	if maxIP <= ipLocalLastByteInt {
		MS_role = dt.MS_Master
	}
	return MS_role
}