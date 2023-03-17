package P2P

import (
	"fmt"
	"project/Network/Utilities/bcast"
	dt "project/commonDataTypes"
	"reflect"
	"time"
)

func P2Pntw(localIP string,
	nodeInfoToNTWCh <-chan []dt.NodeInfo,
	NOStoNTWCh <-chan []dt.NodeOrderStates,
	nodesInfoFromNTWCh chan<- dt.AllNodeInfoWithSenderIP,
	allNOSfromNTWCh chan<- dt.AllNOS_WithSenderIP,
) {
	var (
		transmittNodesInfo       = make(chan dt.AllNodeInfoWithSenderIP)
		transmittNodeOrderStates = make(chan dt.AllNOS_WithSenderIP)
		receiveNodesInfo         = make(chan dt.AllNodeInfoWithSenderIP)
		receiveNodeOrderStates   = make(chan dt.AllNOS_WithSenderIP)
	)

	localNodesInfo := map[string]dt.ElevData{}
	localNOS := map[string][dt.N_FLOORS][2]dt.OrderState{}

	NodesInfo := dt.AllNodeInfoWithSenderIP{}
	NOS := dt.AllNOS_WithSenderIP{}

	//set timer
	broadCastTimer := time.NewTimer(dt.BROADCAST_PERIOD)
	broadCastTimer.Stop()
	NOSTimer := time.NewTimer(1)
	NOSTimer.Stop()
	nodesInfoTimer := time.NewTimer(1)
	nodesInfoTimer.Stop()

	// Receive from NTW
	go bcast.Receiver(15667, receiveNodesInfo)
	go bcast.Receiver(15668, receiveNodeOrderStates)

	// Send to NTW
	go bcast.Transmitter(15667, transmittNodesInfo)
	go bcast.Transmitter(15668, transmittNodeOrderStates)

	// RSM := ""
	// WW := ""

	for {
		select {
		case newNOStoNTW := <-NOStoNTWCh:
			localNOS = dt.NOSSliceToMap(newNOStoNTW)
			//fmt.Printf("NOS TO NTW: %+v\n", localNOS)
			// RSM = PP.RSM_toString(localNOS)
			// fmt.Printf(RSM + "\n" + WW)
		case newNodeInfoToNTW := <-nodeInfoToNTWCh:
			localNodesInfo = dt.NodeInfoSliceToMap(newNodeInfoToNTW)
			broadCastTimer.Reset(1)
			//fmt.Printf("nodesInfo TO NTW: %+v\n", newNodeInfoToNTW)
			// WW = PP.WW_toString(localNodesInfo)
			// fmt.Printf(RSM + "\n" + WW)
		case newNodeOrderStates := <-receiveNodeOrderStates:
			senderData := dt.NOSSliceToMap(newNodeOrderStates.AllNOS)
			senderIP := newNodeOrderStates.SenderIP
			fmt.Printf("We get new data")

			if localIP != senderIP && !reflect.DeepEqual(senderData[senderIP], localNOS[senderIP]) {
				fmt.Printf("and we send it furtther down")
				NOS = newNodeOrderStates
				NOSTimer.Reset(1)
			}
		case newNodesInfo := <-receiveNodesInfo:
			senderData := dt.NodeInfoSliceToMap(newNodesInfo.AllNodeInfo)
			senderIP := newNodesInfo.SenderIP
			//fmt.Printf("nodesInfo fro NTW: %+v\n", senderData)
			if localIP != senderIP && !reflect.DeepEqual(senderData[senderIP], localNodesInfo[senderIP]) {
				// RSM = PP.WW_toString(localNodesInfo)
				// fmt.Printf(RSM + "\n" + WW)
				NodesInfo = newNodesInfo
				nodesInfoTimer.Reset(1)
			}
		case <-broadCastTimer.C:
			//fmt.Printf("nodesInfo to NTW: %+v\n", dt.AllNodeInfoWithSenderIP{SenderIP: localIP, AllNodeInfo: dt.NodeInfoMapToSlice(localNodesInfo)})
			transmittNodesInfo <- dt.AllNodeInfoWithSenderIP{SenderIP: localIP, AllNodeInfo: dt.NodeInfoMapToSlice(localNodesInfo)}
			transmittNodeOrderStates <- dt.AllNOS_WithSenderIP{SenderIP: localIP, AllNOS: dt.NOSMapToSlice(localNOS)}
			broadCastTimer.Reset(dt.BROADCAST_PERIOD)
		case <-nodesInfoTimer.C:
			select {
			case nodesInfoFromNTWCh <- NodesInfo:
			default:
				nodesInfoTimer.Reset(1)
			}
		case <-NOSTimer.C:
			select {
			case allNOSfromNTWCh <- NOS:
			default:
				NOSTimer.Reset(1)
			}
		}
	}
}
