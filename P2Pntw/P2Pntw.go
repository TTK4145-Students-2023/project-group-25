package P2P

import (
	"Driver-go/bcast"
	"Driver-go/localip"
	"Driver-go/peers"
	"fmt"
)

// Datatype to and from DistElevData
type AllElevData map[string]ElevData // ID_1 : elev_strct

type AllElevData_withID struct { // denne sendes på kanal til elevDataDist
	ID      string
	AllData AllElevData
}

type ElevData struct {
	Behavior    string `json:"behaviour"`
	Floor       int    `json:"floor"`
	Direction   string `json:"direction"`
	CabRequests []bool `json:"cabRequests"`
}

// Data type to and from OrderStateHandler
type RequestStateMatrix map[string]singleNode_requestStates // denne sendes på kanal til orderstatehandler
type requestState int
type singleNode_requestStates [][2]requestState

// // input channels
// var (
// 	allElevData_fromDist            = make(chan AllElevData_withID)
// 	ReqStateMatrix_fromOrderHandler = make(chan RequestStateMatrix)

// 	// We make channels for receiving our custom data types
// 	allElevData_fromNTW    = make(chan AllElevData)
// 	ReqStateMatrix_fromNTW = make(chan RequestStateMatrix)
// )

// // output channels
// var (
// 	allElevData_toDist            = make(chan AllElevData_withID)
// 	ReqStateMatrix_toOrderHandler = make(chan RequestStateMatrix)

// 	// We make channels for sending our custom data types
// 	allElevData_toNTW    = make(chan AllElevData)
// 	ReqStateMatrix_toNTW = make(chan RequestStateMatrix)
// )

func P2Pntw(
	// Receive channels
	allElevData_fromDist chan AllElevData_withID,
	ReqStateMatrix_fromOrderHandler chan RequestStateMatrix,
	allElevData_fromNTW chan AllElevData_withID,
	ReqStateMatrix_fromNTW chan RequestStateMatrix,

	// Sending channels
	allElevData_toDist chan AllElevData_withID,
	ReqStateMatrix_toOrderHandler chan RequestStateMatrix,
	allElevData_toNTW chan AllElevData_withID,
	ReqStateMatrix_toNTW chan RequestStateMatrix,
) {
	localIP, _ := localip.LocalIP()
	peerUpdateCh := make(chan peers.PeerUpdate) // channel for receiving updates on the id of the peers that are alive on the network
	peerTxEnable := make(chan bool)             // disable/enable the transmitter after started

	// Peer update
	go peers.Transmitter(15649, localIP, peerTxEnable)
	go peers.Receiver(15649, peerUpdateCh)

	// Receive from NTW
	go bcast.Receiver(15647, allElevData_fromNTW)
	go bcast.Receiver(15648, ReqStateMatrix_fromNTW)

	// Receive from dist and orderHandler
	go bcast.Receiver(15647, allElevData_fromDist)
	go bcast.Receiver(15648, ReqStateMatrix_fromOrderHandler)

	// Send to NTW
	go bcast.Transmitter(15647, allElevData_toNTW)
	go bcast.Transmitter(15648, ReqStateMatrix_toNTW)

	// Send to dist and orderHandler
	go bcast.Transmitter(15647, allElevData_toDist)
	go bcast.Transmitter(15648, ReqStateMatrix_toOrderHandler)

	// This section is to check that the right value is written to different channels
	for {
		select {
		case val := <-ReqStateMatrix_fromNTW:
			fmt.Println(val)
		case testMsg := <-allElevData_fromNTW:
			fmt.Println(testMsg)
		case p := <-peerUpdateCh:
			fmt.Printf("Peer update:\n")
			fmt.Printf("  Peers:    %q\n", p.Peers)
			fmt.Printf("  New:      %q\n", p.New)
			fmt.Printf("  Lost:     %q\n", p.Lost)
		}
	}
}
