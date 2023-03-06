package P2P

import (
	"Driver-go/bcast"
	"Driver-go/localip"
	"Driver-go/peers"
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

// input channels
var (
	allElevData_fromDist            = make(chan AllElevData_withID)
	ReqStateMatrix_fromOrderHandler = make(chan RequestStateMatrix)

	// We make channels for receiving our custom data types
	allElevData_fromNTW    = make(chan AllElevData)
	ReqStateMatrix_fromNTW = make(chan RequestStateMatrix)
)

// output channels
var (
	allElevData_toDist            = make(chan AllElevData_withID)
	ReqStateMatrix_toOrderHandler = make(chan RequestStateMatrix)

	// We make channels for receiving our custom data types
	allElevData_toNTW    = make(chan AllElevData)
	ReqStateMatrix_toNTW = make(chan RequestStateMatrix)
)

func P2Pntw(
	allElevData_fromDist <-chan AllElevData_withID,
	ReqStateMatrix_fromOrderHandler <-chan RequestStateMatrix,
	allElevData_fromNTW <-chan AllElevData_withID,
	ReqStateMatrix_fromNTW <-chan RequestStateMatrix,

	allElevData_toNTW chan<- AllElevData_withID,
	ReqStateMatrix_toNTW chan<- RequestStateMatrix,
	allElevData_toDist chan<- AllElevData_withID,
	ReqStateMatrix_toOrderHandler chan<- RequestStateMatrix,
) {
	localIP, _ := localip.LocalIP()
	peerUpdateCh := make(chan peers.PeerUpdate) // channel for receiving updates on the id of the peers that are alive on the network
	peerTxEnable := make(chan bool)             // disable/enable the transmitter after started

	go peers.Transmitter(15647, localIP, peerTxEnable)
	go peers.Receiver(15647, peerUpdateCh)

	go bcast.Receiver(15647, allElevData_fromNTW)
	go bcast.Receiver(15647, ReqStateMatrix_fromNTW)

	for {
		select {
		case ElevDataDistTx := <-allElevData_fromDist:
			go bcast.Transmitter(15647, allElevData_fromDist)
			allElevData_toNTW <- ElevDataDistTx

		case ReqStateMatrixHandlerTx := <-ReqStateMatrix_fromOrderHandler:
			go bcast.Transmitter(15647, ReqStateMatrix_fromOrderHandler)
			ReqStateMatrix_toNTW <- ReqStateMatrixHandlerTx

		case allElevDataNtwRx := <-allElevData_fromNTW:
			allElevData_toDist <- allElevDataNtwRx

		case ReqStateMatrixRx := <-ReqStateMatrix_fromNTW:
			ReqStateMatrix_toOrderHandler <- ReqStateMatrixRx
		}
	}
}

// func P2PntwTx(ordrFromDist chan P2Pmsg) {
// 	localIP, _ := localip.LocalIP() // get local IP

// 	peerUpdateCh := make(chan peers.PeerUpdate) // channel for receiving updates on the id of the peers that are alive on the network
// 	peerTxEnable := make(chan bool)             // disable/enable the transmitter after started

// 	go peers.Transmitter(15647, localIP, peerTxEnable)
// 	go peers.Receiver(15647, peerUpdateCh)

// 	go bcast.Transmitter(16569, ordrFromDist) // broadcasting msg on port

// }

// func P2PntwRx(){
// 	msgRx := make(chan P2Pmsg) // channels for receiving custom data types
// 	go bcast.Receiver(16569, msgRx) // receive msg on port

// 	peerUpdateCh := make(chan peers.PeerUpdate) // channel for receiving updates on the id of the peers that are alive on the network
// 	go peers.Receiver(15647, peerUpdateCh)
// }

// func P2Pntw(ordrFromDist chan P2Pmsg) {

// 	localIP, _ := localip.LocalIP() // get local IP

// 	peerUpdateCh := make(chan peers.PeerUpdate) // channel for receiving updates on the id of the peers that are alive on the network
// 	peerTxEnable := make(chan bool)             // disable/enable the transmitter after started

// 	go peers.Transmitter(15647, localIP, peerTxEnable)
// 	go peers.Receiver(15647, peerUpdateCh)

// 	go bcast.Transmitter(16569, ordrFromDist) // broadcasting msg on port

// 	msgRx := make(chan P2Pmsg) // channels for receiving custom data types
// 	go bcast.Receiver(16569, msgRx) // receive msg on port

// 	go peers.Receiver(15647, peerUpdateCh)

// 	for {
// 		select {
// 		case orderFromDist := <-ordrFromDist: // when we get a new update from orderDistributor
// 			outputMsg := P2Pmsg{
// 				IPAddrFrom:   localIP,
// 				HallRequests: orderFromDist.HallRequests,
// 				States:       orderFromDist.States,
// 			}
// 			fmt.Printf("Going to distribute this msg on P2P ntw: %#v\n", outputMsg)
// 			for {
// 				msgTx <- outputMsg // putting outputMsg onto msgTx channel
// 				//time.Sleep(500 * time.Millisecond) // choose pushrate onto network
// 			}
// 		case msgFromP2P := <- msgRx: // when we get a msg from P2P ntw

// 		}
// 	}

// }

// UNUSED allows unused variables to be included in Go programs
