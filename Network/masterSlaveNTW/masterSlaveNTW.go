package masterSlaveNTW

import (
	"fmt"
	"project/Network/Utilities/bcast"
	"project/Network/Utilities/localip"
	"project/Network/Utilities/peers"
	dt "project/commonDataTypes"
)

// datatypes
type AssignedOrders map[string][]bool

func MasterSlaveNTW(

	// Receive channels
	peerUpdate_MS chan peers.PeerUpdate,
	inputOrders_fromAss chan []byte,
	// Sending channels

	outputOrders_toAss chan []byte,
	masterOrSlave chan dt.MasterSlaveRole,
) {
	var (
		inputOrders_fromNTW = make(chan []byte)
		outputOrders_toNTW  = make(chan []byte)
	)

	localIP, _ := localip.LocalIP()

	// Receive from NTW
	go bcast.Receiver(15640, inputOrders_fromNTW)

	// Send on NTW
	go bcast.Transmitter(15640, outputOrders_toNTW)

	Local_MS_role := dt.MS_Slave

	for {
		select {
		case peerList := <-peerUpdate_MS:
			//update local MS_role
			fmt.Printf("_____PEER LIST______\n  %s\n", peerList)
			Local_MS_role = MS_Assigner(localIP, peerList.Peers)
			fmt.Printf("_____MS ROLE______\n  %s\n", Local_MS_role)
			masterOrSlave <- Local_MS_role
			fmt.Printf("_____MS IS ASSIGNED______\n  %s\n", Local_MS_role)
			//send peerlsit to other?
		case ordersFromNTW := <-inputOrders_fromNTW:
			switch Local_MS_role {
			case dt.MS_Master:
			case dt.MS_Slave:
				outputOrders_toAss <- ordersFromNTW
			}
		case ordersFromAss := <-inputOrders_fromAss:
			switch Local_MS_role {
			case dt.MS_Slave:
			case dt.MS_Master:
				outputOrders_toNTW <- ordersFromAss

			}
		}
	}
}
