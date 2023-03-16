package P2P

import (
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

	localWorldView := map[string]dt.ElevData{}
	localRequestStateMatrix := map[string][dt.N_FLOORS][2]dt.OrderState{}

	worldView := dt.AllNodeInfoWithSenderIP{}
	requestStateMatrix := dt.AllNOS_WithSenderIP{}

	//set timer
	broadCastTimer := time.NewTimer(dt.BROADCAST_RATE)
	reqStateMatrixTimer := time.NewTimer(1)
	reqStateMatrixTimer.Stop()
	worldViewTimer := time.NewTimer(1)
	worldViewTimer.Stop()

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
			localRequestStateMatrix = dt.NOSSliceToMap(newNOStoNTW)
			// RSM = PP.RSM_toString(localRequestStateMatrix)
			// fmt.Printf(RSM + "/n" + WW)
		case newNodeInfoToNTW := <-nodeInfoToNTWCh:
			localWorldView = dt.NodeInfoSliceToMap(newNodeInfoToNTW)
			// WW = PP.WW_toString(localWorldView)
			// fmt.Printf(RSM + "/n" + WW)
		case newNodeOrderStates := <-receiveNodeOrderStates:
			senderData := dt.NOSSliceToMap(newNodeOrderStates.AllNOS)
			senderIP := newNodeOrderStates.SenderIP
			if localIP != senderIP && !reflect.DeepEqual(senderData[senderIP], localRequestStateMatrix[senderIP]) {
				requestStateMatrix = newNodeOrderStates
				reqStateMatrixTimer.Reset(1)
			}
		case newNodesInfo := <-receiveNodesInfo:
			senderData := dt.NodeInfoSliceToMap(newNodesInfo.AllNodeInfo)
			senderIP := newNodesInfo.SenderIP
			if localIP != senderIP && !reflect.DeepEqual(senderData[senderIP], localWorldView[senderIP]) {
				worldView = newNodesInfo
				worldViewTimer.Reset(1)
			}
		case <-broadCastTimer.C:
			transmittNodesInfo <- dt.AllNodeInfoWithSenderIP{SenderIP: localIP, AllNodeInfo: dt.NodeInfoMapToSlice(localWorldView)}
			transmittNodeOrderStates <- dt.AllNOS_WithSenderIP{SenderIP: localIP, AllNOS: dt.NOSMapToSlice(localRequestStateMatrix)}
			broadCastTimer.Reset(dt.BROADCAST_RATE)
		case <-worldViewTimer.C:
			select {
			case nodesInfoFromNTWCh <- worldView:
			default:
				worldViewTimer.Reset(1)
			}
		case <-reqStateMatrixTimer.C:
			select {
			case allNOSfromNTWCh <- requestStateMatrix:
			default:
				reqStateMatrixTimer.Reset(1)
			}
		}
	}
}
