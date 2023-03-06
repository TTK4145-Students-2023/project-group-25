package P2Pntw

import (
	"Driver-go/bcast"
	"Driver-go/localip"
	"Driver-go/peers"
	"time"
)

// type AllElevState struct {
// 	Behaviour   string `json:"behaviour"`
// 	Floor       int    `json:"floor"`
// 	Direction   string `json:"direction"`
// 	CabRequests []bool `json:"cabRequests"`
// }

// Defining data type to send on P2P network
type P2Pmsg struct {
	IPAddrFrom   string
	HallRequests [][2]bool               `json:"hallRequests"` //Btn types? Up/down
	States       map[string]HRAElevState `json:"states"`
}

type HRAElevState struct {
	Behavior    string `json:"behaviour"`
	Floor       int    `json:"floor"`
	Direction   string `json:"direction"`
	CabRequests []bool `json:"cabRequests"`
}

type HRAInput struct {
	HallRequests [][2]bool               `json:"hallRequests"`
	States       map[string]HRAElevState `json:"states"`
}

type OrderMatrix [][2]bool

func P2Pntw(ordrFromDist chan HRAInput) {
	localIP, _ := localip.LocalIP() // get local IP

	peerUpdateCh := make(chan peers.PeerUpdate) // channel for receiving updates on the id of the peers that are alive on the network
	peerTxEnable := make(chan bool)             // disable/enable the transmitter after started

	go peers.Transmitter(15647, localIP, peerTxEnable)
	go peers.Receiver(15647, peerUpdateCh)

	msgTx := make(chan P2Pmsg) // channels for sending and receiving custom data types
	msgRx := make(chan P2Pmsg) // channels for sending and receiving custom data types

	// channels for sending info to order distributor
	ordrMatIncoming := make(chan [][2]bool)
	AllElevStatesIncoming := make(chan map[string]HRAElevState)
	IPaddrIncoming := make(chan string)

	go bcast.Transmitter(16569, msgTx)
	go bcast.Receiver(16569, msgRx)

	for {
		select {
		case orderFromDist := <-ordrFromDist: // when we get a new update from orderDistributor
			go func() {
				outputMsg := P2Pmsg{
					IPAddrFrom:   localIP,
					HallRequests: orderFromDist.HallRequests,
					States:       orderFromDist.States,
				}
				for {
					msgTx <- outputMsg                 // putting outputMsg onto msgTx channel
					time.Sleep(500 * time.Millisecond) // choose pushrate onto network
				}
			}()
		case incomingMsg := <-msgRx:
			ordrMatIncoming <- incomingMsg.HallRequests
			AllElevStatesIncoming <- incomingMsg.States
			IPaddrIncoming <- incomingMsg.IPAddrFrom
		}
	}
}
