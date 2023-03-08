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

	// Send to NTW
	go bcast.Transmitter(15647, allElevData_toNTW)
	go bcast.Transmitter(15648, ReqStateMatrix_toNTW)

	// Receive from dist and orderHandler comes from input channels
	// Send to dist and orderHandler is fixed by for select

	// This section is to check that the right value is written to different channels
	for {
		select {
		case ReqStateMatNTW := <-ReqStateMatrix_fromNTW:
			ReqStateMatrix_toOrderHandler <- ReqStateMatNTW
		case AllElevDataNTW := <-allElevData_fromNTW:
			allElevData_toDist <- AllElevDataNTW
		case ReqStateMatOrderHandler := <-ReqStateMatrix_fromOrderHandler:
			ReqStateMatrix_toNTW <- ReqStateMatOrderHandler
		case AllElevDataDist := <-allElevData_fromDist:
			allElevData_toNTW <- AllElevDataDist
		case p := <-peerUpdateCh:
			fmt.Printf("Peer update:\n")
			fmt.Printf("  Peers:    %q\n", p.Peers)
			fmt.Printf("  New:      %q\n", p.New)
			fmt.Printf("  Lost:     %q\n", p.Lost)
		}
	}
}
