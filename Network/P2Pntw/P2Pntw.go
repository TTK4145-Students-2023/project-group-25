package P2P

import (
	"project/Network/Utilities/bcast"
	dt "project/commonDataTypes"
	"reflect"
	"time"
)

func P2Pntw(localIP string,
	nodeInfoToNTWCh <-chan []dt.NodeInfo,
	nodesInfoFromNTWCh chan<- dt.AllNodeInfoWithSenderIP,
) {
	var (
		transmittNodesInfo = make(chan dt.AllNodeInfoWithSenderIP)

		receiveNodesInfo = make(chan dt.AllNodeInfoWithSenderIP)
	)

	localNodesInfo := map[string]dt.ElevData{}

	NodesInfo := dt.AllNodeInfoWithSenderIP{}

	//set timer
	NOSTimer := time.NewTimer(1)
	NOSTimer.Stop()
	nodesInfoTimer := time.NewTimer(1)
	nodesInfoTimer.Stop()
	broadCastTimer := time.NewTimer(1)
	broadCastTimer.Stop()

	// Receive from NTW
	go bcast.Receiver(15667, receiveNodesInfo)
	// Send to NTW
	go bcast.Transmitter(15667, transmittNodesInfo)

	for {
		select {
		case newNodeInfoToNTW := <-nodeInfoToNTWCh:
			localNodesInfo = dt.NodeInfoSliceToMap(newNodeInfoToNTW)
			broadCastTimer.Reset(1)
			//fmt.Printf("nodesInfo TO NTW: %+v\n", newNodeInfoToNTW)
			// WW = PP.WW_toString(localNodesInfo)
			// fmt.Printf(RSM + "\n" + WW)

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
			broadCastTimer.Reset(dt.BROADCAST_PERIOD)
		case <-nodesInfoTimer.C:
			select {
			case nodesInfoFromNTWCh <- NodesInfo:
			default:
				nodesInfoTimer.Reset(1)
			}
		}
	}
}
