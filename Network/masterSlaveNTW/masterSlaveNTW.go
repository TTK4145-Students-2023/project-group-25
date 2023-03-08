package masterSlaveNTW

import (
	"Driver-go/bcast"
	"Driver-go/localip"
	"Driver-go/peers"
)

// datatypes
type MasterSlave_state int

type AssignedOrders map[string][]bool

const (
	MS_master MasterSlave_state = 0
	MS_slave  MasterSlave_state = 1
)

// channels
var (
	inputOrders_fromNTW = make(chan AssignedOrders) // input channel for reciving orders from other nodes
	inputOrders_fromAss = make(chan AssignedOrders)
	outputOrders_toNTW  = make(chan AssignedOrders) //output channel for sending order to other
	outputOrders_toAss  = make(chan AssignedOrders)
	MasterSlave         = make(chan bool) // output channel to send Master or Slave role to order assigner
)

func MasterSlaveNTW(

	// Receive channels
	inputOrders_fromNTW chan AssignedOrders,
	inputOrders_fromAsschan AssignedOrders,
	// Sending channels
	outputOrders_toNTW chan AssignedOrders,
	outputOrders_toAss chan AssignedOrders,
	masterOrSlave chan MasterSlave_state,
) {
	localIP, _ := localip.LocalIP()
	// Peer update
	peerUpdateCh := make(chan peers.PeerUpdate) // channel for receiving updates on the id of the peers that are alive on the network
	peerTxEnable := make(chan bool)             // disable/enable the transmitter after started
	go peers.Transmitter(15649, localIP, peerTxEnable)
	go peers.Receiver(15649, peerUpdateCh)

	// Receive from NTW
	go bcast.Receiver(15640, inputOrders_fromNTW)

	// Send on NTW
	go bcast.Transmitter(15640, outputOrders_toNTW)

	Local_MS_role := MS_slave

	for {
		select {
		case peerList := <-peerUpdateCh:
			//update local MS_role
			Local_MS_role = MS_Assigner(localIP, peerList.Peers)
			masterOrSlave <- Local_MS_role
			//send peerlsit to other?
		case ordersFromNTW := <-inputOrders_fromNTW:
			switch Local_MS_role {
			case MS_master:
				//if master, do nothing
			case MS_slave:
				// if slave, send orders to assigner
				outputOrders_toAss <- ordersFromNTW
			}
		case ordersFromAss := <-inputOrders_fromAss:
			switch Local_MS_role {
			case MS_master:
				//if master send orders to network
				outputOrders_toNTW <- ordersFromAss
			case MS_slave:
				//if slave, do nothing
			}
		}
	}
}
